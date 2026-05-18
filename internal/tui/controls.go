package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/Abhiram86/echotune/internal/models"
	"github.com/Abhiram86/echotune/internal/platform"
)

type AppSession interface {
	CurrentSong() models.Download
	Pause() error
	Next(ctx context.Context) error
	Previous(ctx context.Context) error
	AddToPlaylist(ctx context.Context, storage *models.Storage, title string) error
	RemoveFromPlaylist(ctx context.Context, storage *models.Storage, title string) error
}

type DownloadFunc func(ctx context.Context, storage *models.Storage, song models.SearchResult, mgr *models.DownloadManager) error

func getMpvProperty(socketPath string, name string) (float64, error) {
	conn, err := platform.DialIPC(socketPath)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	cmd := map[string]any{
		"command": []any{"get_property", name},
	}
	data, err := json.Marshal(cmd)
	if err != nil {
		return 0, err
	}

	if _, err := conn.Write(append(data, '\n')); err != nil {
		return 0, err
	}

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	var resp struct {
		Data  *float64 `json:"data"`
		Error string   `json:"error"`
	}
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		return 0, err
	}

	if resp.Error != "success" {
		return 0, fmt.Errorf("mpv error: %s", resp.Error)
	}

	if resp.Data == nil {
		return 0, nil
	}

	return *resp.Data, nil
}

func Controls(
	ctx context.Context,
	app AppSession,
	player *models.Player,
	storage *models.Storage,
	extraControls []string,
	downloadFunc DownloadFunc,
) error {
	enabled := map[string]bool{
		"q": true,
		"p": true,
		"f": true,
		"b": true,
	}

	for _, ctrl := range extraControls {
		enabled[ctrl] = true
	}

	m := playerModel{
		app:          app,
		player:       player,
		storage:      storage,
		enabled:      enabled,
		ctx:          ctx,
		downloadFunc: downloadFunc,
		songTitle:    app.CurrentSong().Title,
		status:       player.GetStatus(),
		downloadMgr:  &models.DownloadManager{},
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	fm := finalModel.(playerModel)
	return fm.exitErr
}

type mode int

const (
	normalMode mode = iota
	inputMode
)

type tickMsg struct{}
type progressMsg struct {
	pos      float64
	duration float64
}
type playerDoneMsg struct{}

type playerModel struct {
	app          AppSession
	player       *models.Player
	storage      *models.Storage
	enabled      map[string]bool
	ctx          context.Context
	downloadFunc DownloadFunc

	songTitle string
	status    models.PlayerStatus
	pos       float64
	duration  float64

	mode          mode
	pendingAction string
	pendingInput  string
	quitting      bool

	downloadMgr *models.DownloadManager
	exitErr     error
}

func (m playerModel) Init() tea.Cmd {
	return tea.Batch(
		tea.Tick(time.Second, func(time.Time) tea.Msg {
			return tickMsg{}
		}),
		func() tea.Msg {
			<-m.player.Done
			return playerDoneMsg{}
		},
	)
}

func (m playerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		return m, tea.Batch(
			tea.Tick(time.Second, func(time.Time) tea.Msg {
				return tickMsg{}
			}),
			func() tea.Msg {
				pos, _ := getMpvProperty(m.player.SocketPath, "time-pos")
				dur, _ := getMpvProperty(m.player.SocketPath, "duration")
				return progressMsg{pos: pos, duration: dur}
			},
		)

	case progressMsg:
		m.pos = msg.pos
		m.duration = msg.duration
		m.status = m.player.GetStatus()
		return m, nil

	case playerDoneMsg:
		m.quitting = true
		return m, tea.Quit

	case tea.KeyMsg:
		if m.mode == inputMode {
			return m.handleInputModeKey(msg)
		}
		return m.handleNormalModeKey(msg)
	}

	return m, nil
}

func formatDuration(seconds float64) string {
	secs := int(seconds)
	mins := secs / 60
	secs = secs % 60
	return fmt.Sprintf("%d:%02d", mins, secs)
}

func (m playerModel) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}
	var v tea.View
	if m.mode == inputMode {
		v = tea.NewView(m.inputView())
	} else {
		v = tea.NewView(m.normalView())
	}
	return v
}

func (m playerModel) normalView() string {
	var b strings.Builder

	statusIcon := "▶"
	statusText := "Playing"
	switch m.status {
	case models.Paused:
		statusIcon = "⏸"
		statusText = "Paused"
	case models.Stopped:
		statusIcon = "■"
		statusText = "Stopped"
	}

	m.downloadMgr.Mu.Lock()
	isDownloading := m.downloadMgr.IsDownloading
	m.downloadMgr.Mu.Unlock()

	if isDownloading {
		statusText += " (Downloading...)"
	}

	fmt.Fprintf(&b, "\n  %s Now Playing: %s\n", statusIcon, m.songTitle)

	barWidth := 30
	filled := 0
	if m.duration > 0 {
		filled = int((m.pos / m.duration) * float64(barWidth))
	}
	if filled > barWidth {
		filled = barWidth
	} else if filled < 0 {
		filled = 0
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	current := formatDuration(m.pos)
	total := formatDuration(m.duration)
	fmt.Fprintf(&b, "  ● %s  [%s]  %s / %s\n", statusText, bar, current, total)

	fmt.Fprintf(&b, "\n  [q] quit  [space] pause  [j/←] -5s  [k/→] +5s")
	if m.enabled["d"] {
		fmt.Fprintf(&b, "  [d] dl")
	}
	if m.enabled["n"] {
		fmt.Fprintf(&b, "  [n/↑] next")
	}
	if m.enabled["z"] {
		fmt.Fprintf(&b, "  [b/↓] prev")
	}
	if m.enabled["a"] {
		fmt.Fprintf(&b, "  [a] +list")
	}
	if m.enabled["x"] {
		fmt.Fprintf(&b, "  [x] -list")
	}
	fmt.Fprintf(&b, "\n")

	return b.String()
}

func (m playerModel) inputView() string {
	return fmt.Sprintf("\n  Enter %s: %s\n  (press Enter to confirm, Esc to cancel)\n", m.pendingAction, m.pendingInput)
}

func (m playerModel) handleInputModeKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if len(m.pendingInput) < 1 {
			return m, nil
		}
		switch m.pendingAction {
		case "playlist title to add to":
			m.app.AddToPlaylist(m.ctx, m.storage, m.pendingInput)
		case "playlist title to remove from":
			m.app.RemoveFromPlaylist(m.ctx, m.storage, m.pendingInput)
		}
		m.mode = normalMode
		m.pendingAction = ""
		m.pendingInput = ""
		return m, nil

	case "esc":
		m.mode = normalMode
		m.pendingAction = ""
		m.pendingInput = ""
		return m, nil

	case "backspace":
		if len(m.pendingInput) > 0 {
			m.pendingInput = m.pendingInput[:len(m.pendingInput)-1]
		}
		return m, nil

	default:
		if len(msg.String()) == 1 && msg.String()[0] >= 32 && msg.String()[0] <= 126 {
			m.pendingInput += msg.String()
		}
		return m, nil
	}
}

func (m playerModel) handleNormalModeKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		m.player.Stop()
		m.exitErr = fmt.Errorf("interrupted, user quit")
		m.quitting = true
		return m, tea.Quit

	case " ", "space", "p":
		m.app.Pause()
		m.status = m.player.GetStatus()
		return m, nil

	case "k", "right":
		m.player.Seek(5)
		return m, nil

	case "j", "left":
		m.player.Seek(-5)
		return m, nil

	case "n", "up":
		if m.enabled["n"] {
			m.app.Next(m.ctx)
			return m, func() tea.Msg {
				<-m.player.Done
				return playerDoneMsg{}
			}
		}

	case "b", "down":
		if m.enabled["z"] {
			m.app.Previous(m.ctx)
			return m, func() tea.Msg {
				<-m.player.Done
				return playerDoneMsg{}
			}
		}

	case "d":
		if m.enabled["d"] {
			m.downloadFunc(m.ctx, m.storage, m.app.CurrentSong().Metadata, m.downloadMgr)
		}

	case "a":
		if m.enabled["a"] {
			m.mode = inputMode
			m.pendingAction = "playlist title to add to"
			m.pendingInput = ""
			return m, nil
		}

	case "x":
		if m.enabled["x"] {
			m.mode = inputMode
			m.pendingAction = "playlist title to remove from"
			m.pendingInput = ""
			return m, nil
		}
	}

	return m, nil
}
