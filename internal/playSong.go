package internal

import (
	"os/exec"
	"time"

	"github.com/Abhiram86/echotune/internal/models"
)

func PlaySong(p *models.Player, url string) error {
	p.SocketPath = "/tmp/echotune.sock"

	p.Cmd = exec.Command(
		"mpv",
		"--no-video",
		"--input-ipc-server="+p.SocketPath,
		"--ytdl-format=bestaudio[ext=webm]/bestaudio",
		url,
	)

	err := p.Cmd.Start()
	if err != nil {
		return err
	}

	p.Done = make(chan bool)

	go func() {
		err := p.Cmd.Wait()
		if err != nil {
			panic(err)
		}

		p.Done <- true
		p.Status = models.Stopped
	}()

	time.Sleep(300 * time.Millisecond)
	p.Status = models.Playing

	return nil
}
