package models

type SearchList struct {
	Query   string         `json:"query"`
	Results []SearchResult `json:"results"`
}

type SearchResult struct {
	ID         string  `json:"id"`
	Title      string  `json:"title"`
	URL        string  `json:"url"`
	Duration   float32 `json:"duration"`
	ViewCount  int     `json:"view_count"`
	LikeCount  int     `json:"like_count"`
	Uploader   string  `json:"uploader"`
	Channel    string  `json:"channel"`
	UploadDate string  `json:"upload_date"`
}
