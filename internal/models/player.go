package models

import (
	"fmt"
	"net"
	"os/exec"
)

type PlayerStatus int

const (
	Paused PlayerStatus = iota
	Playing
	Stopped
)

type Player struct {
	Cmd        *exec.Cmd
	SocketPath string
	Status     PlayerStatus
	Done       chan bool
	Song       SearchResult
}

func (p *Player) sendCommand(command string) error {
	conn, err := net.Dial("unix", p.SocketPath)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(command + "\n"))
	return err
}

func (p *Player) TogglePlay() error {
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
