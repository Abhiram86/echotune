package internal

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/Abhiram86/echotune/internal/models"
	"github.com/Abhiram86/echotune/internal/platform"
)

type PlaybackSession struct {
	Player    *models.Player
	Queue     *models.Queue
	Playlists *models.Playlists
}

func NewPlaybackSession(storage *models.Storage, songs []models.Download) *PlaybackSession {
	paths, err := platform.NewAppPaths()
	if err != nil {
		panic(err)
	}
	return &PlaybackSession{
		Player: &models.Player{
			SocketPath: paths.SocketFile,
		},
		Queue: &models.Queue{
			Songs:        songs,
			CurrentIndex: 0,
		},
		Playlists: &storage.Playlists,
	}
}

func (app *PlaybackSession) CurrentSong() models.Download {
	return app.Queue.Songs[app.Queue.CurrentIndex]
}

func (app *PlaybackSession) Play(ctx context.Context) error {
	app.Player.Song = app.CurrentSong().Metadata
	fmt.Printf("Playing song %s\n", app.CurrentSong().Title)
	if app.CurrentSong().Path == "__SEARCHED__" {
		return app.Player.PlaySong(ctx, models.Playable{URL: app.CurrentSong().Metadata.URL})
	}
	return app.Player.PlaySong(ctx, models.Playable{URL: app.CurrentSong().Path})
}

func (app *PlaybackSession) PlayALL(ctx context.Context, storage *models.Storage, additional ...Control) models.PlayerStatus {
	for app.Queue.CurrentIndex >= 0 && app.Queue.CurrentIndex < len(app.Queue.Songs) {
		storage.History.Add(app.CurrentSong().Metadata)
		err := app.Play(ctx)
		if err != nil {
			fmt.Printf("playback error: %v\n", err)
			return models.Stopped
		}

		reader := bufio.NewReader(os.Stdin)

		err = Controls(ctx, app, storage, reader, additional...)
		if err != nil {
			return models.Stopped
		}

		app.Queue.CurrentIndex++
	}
	return app.Player.GetStatus()
}

func (app *PlaybackSession) Stop() error {
	return app.Player.Stop()
}

func (app *PlaybackSession) Pause() error {
	return app.Player.TogglePlay()
}

func (app *PlaybackSession) Seek(second int) error {
	return app.Player.Seek(second)
}

func (app *PlaybackSession) Next(ctx context.Context) error {
	if len(app.Queue.Songs) > 1 {
		// Stop will trigger the Done channel, unblocking the PlayALL loop,
		// which will automatically increment CurrentIndex and play the next song.
		app.Player.Stop()
	}
	return nil
}

func (app *PlaybackSession) Previous(ctx context.Context) error {
	if len(app.Queue.Songs) > 1 {
		// Since PlayALL will increment CurrentIndex by 1 after Stop,
		// we subtract 2 to end up at CurrentIndex - 1.
		app.Queue.CurrentIndex -= 2
		if app.Queue.CurrentIndex < -1 {
			app.Queue.CurrentIndex = len(app.Queue.Songs) - 2
		}
		app.Player.Stop()
	}
	return nil
}

func (app *PlaybackSession) AddToPlaylist(ctx context.Context, storage *models.Storage, playlistTitle string) error {
	song := app.CurrentSong()
	_, exists := storage.Playlists.Get(playlistTitle)
	if !exists {
		app.Playlists.AddPlayList(models.Playlist{
			Title: playlistTitle,
			Songs: make(map[string]models.Download, 0),
		})
		return app.Playlists.AddSong(playlistTitle, song)
	}
	return storage.Playlists.AddSong(playlistTitle, song)
}

func (app *PlaybackSession) RemoveFromPlaylist(ctx context.Context, storage *models.Storage, playlistTitle string) error {
	song := app.CurrentSong()
	_, exists := storage.Playlists.Get(playlistTitle)
	if !exists {
		return fmt.Errorf("playlist '%s' does not exist", playlistTitle)
	}
	return app.Playlists.RemoveSong(playlistTitle, song)
}
