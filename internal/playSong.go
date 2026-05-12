package internal

import (
	"context"
	"os/exec"
	"time"

	"github.com/Abhiram86/echotune/internal/models"
)

func PlaySong(ctx context.Context, p *models.Player, song models.Playable) error {
	p.SocketPath = "/tmp/echotune.sock"

	p.Cmd = exec.CommandContext(ctx,
		"mpv",
		"--no-video",
		"--input-ipc-server="+p.SocketPath,
		"--ytdl-format=bestaudio[ext=webm]/bestaudio",
		song.URL,
	)

	err := p.Cmd.Start()
	if err != nil {
		return err
	}

	p.Done = make(chan bool)

	go func() {
		err := p.Cmd.Wait()
		if err != nil && ctx.Err() == nil {
			p.Cmd.Process.Kill()
		}
		p.Done <- true
		p.SetStatus(models.Stopped)
	}()

	time.Sleep(300 * time.Millisecond)
	p.SetStatus(models.Playing)

	return nil
}
