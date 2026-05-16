package models

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"sync"
	"time"
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
	Done       chan bool
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

func (p *Player) sendCommand(command string) error {
	conn, err := net.DialTimeout("unix", p.SocketPath, 10*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(command + "\n"))
	return err
}

func (p *Player) PlaySong(ctx context.Context, song Playable) error {
	p.SocketPath = "/tmp/echotune.sock"

	p.Cmd = exec.CommandContext(ctx,
		"mpv",
		"--no-video",
		"--input-ipc-server="+p.SocketPath,
		"--ytdl-format=bestaudio[ext=webm]/bestaudio",
		"--cache=yes",
		"--cache-secs=30",
		"--demuxer-max-bytes=50M",
		"--ytdl-raw-options=js-runtimes=node",
		song.URL,
	)

	// p.Cmd.Stderr = os.Stderr

	err := p.Cmd.Start()
	if err != nil {
		return err
	}

	p.Done = make(chan bool)

	go func() {
		err := p.Cmd.Wait()
		fmt.Printf("\n[DEBUG] mpv exited. Wait err: %v, ctx.Err: %v\n", err, ctx.Err())

		if err != nil && ctx.Err() == nil {
			fmt.Printf("[DEBUG] Killing mpv process...\n")
			p.Cmd.Process.Kill()
		}

		fmt.Printf("[DEBUG] Sending done signal...\n")
		p.Done <- true
		p.SetStatus(Stopped)
	}()

	// time.Sleep(300 * time.Millisecond)
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
		return p.sendCommand("cycle pause")
	}

	p.Status = Playing
	return p.sendCommand("cycle pause")
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
	return p.sendCommand(fmt.Sprintf("seek %d", second))
}
