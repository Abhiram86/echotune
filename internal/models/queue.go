package models

type Queue struct {
	Songs        []SearchResult
	CurrentIndex int
	Shuffle      bool
}
