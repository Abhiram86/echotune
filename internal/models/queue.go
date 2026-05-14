package models

type Queue struct {
	Songs        []Download
	CurrentIndex int
}
