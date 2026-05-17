package models

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/Abhiram86/echotune/internal/platform"
)

type PlayerStatus int

const (
	Paused PlayerStatus = iota
	Playing
	Stopped
)

type Player struct {
	mu         sync.Mutex
	Cmd        *exec.Cmd
	SocketPath string
	Status     PlayerStatus
	Done       chan struct{}
	Song       SearchResult
}

type Playable struct {
	URL string
}

func (p *Player) GetStatus() PlayerStatus {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Status
}

func (p *Player) SetStatus(status PlayerStatus) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Status = status
}

func (p *Player) sendCommand(args ...any) error {
	conn, err := platform.DialIPC(p.SocketPath)
	if err != nil {
		return err
	}
	defer conn.Close()

	cmd := map[string]any{
		"command": args,
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	_, err = conn.Write(append(data, '\n'))
	return err
}

func waitForSocket(ctx context.Context, path string, timeout time.Duration) error {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	timeoutTimer := time.NewTimer(timeout)
	defer timeoutTimer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-timeoutTimer.C:
			return fmt.Errorf("timed out waiting for mpv IPC socket")

		case <-ticker.C:
			conn, err := platform.DialIPC(path)
			if err == nil {
				conn.Close()
				return nil
			}
		}
	}
}

func (p *Player) PlaySong(ctx context.Context, song Playable) error {
	if p.SocketPath == "" {
		return fmt.Errorf("socket path not set")
	}

	// Remove stale unix socket
	if runtime.GOOS != "windows" {
		_ = os.Remove(p.SocketPath)
	}

	args := []string{
		"--no-video",
		"--input-ipc-server=" + p.SocketPath,
		"--ytdl-format=bestaudio[ext=webm]/bestaudio",
		"--cache=yes",
		"--cache-secs=30",
		"--demuxer-max-bytes=50M",
		"--ytdl-raw-options=js-runtimes=node",
		song.URL,
	}

	p.Cmd = exec.CommandContext(ctx, "mpv", args...)

	// Optional debug
	// p.Cmd.Stdout = os.Stdout
	// p.Cmd.Stderr = os.Stderr

	if err := p.Cmd.Start(); err != nil {
		return err
	}

	p.Done = make(chan struct{})

	go func() {
		err := p.Cmd.Wait()

		fmt.Printf(
			"\n[DEBUG] mpv exited. Wait err: %v, ctx.Err: %v\n",
			err,
			ctx.Err(),
		)

		p.SetStatus(Stopped)
		close(p.Done)
	}()

	// Wait a bit for IPC socket/pipe to exist
	if err := waitForSocket(ctx, p.SocketPath, 2*time.Second); err != nil {
		return err
	}

	p.SetStatus(Playing)

	return nil
}

func (p *Player) TogglePlay() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.Status == Stopped {
		return fmt.Errorf("player stopped, audio ended")
	}

	if p.Status == Playing {
		p.Status = Paused
		return p.sendCommand("cycle", "pause")
	}

	p.Status = Playing
	return p.sendCommand("cycle", "pause")
}

func (p *Player) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.Status == Stopped {
		return nil
	}

	err := p.sendCommand("quit")
	p.Status = Stopped
	return err
}

func (p *Player) Seek(second int) error {
	return p.sendCommand("seek", second)
}
