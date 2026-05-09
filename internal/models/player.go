package models

import (
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
	conn, err := net.DialTimeout("unix", p.SocketPath, 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(command + "\n"))
	return err
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

func (p *Player) Seek(second int) error {
	return p.sendCommand(fmt.Sprintf("seek %d", second))
}
