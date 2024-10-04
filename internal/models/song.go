package models

import "time"

type Song struct {
	ID          int       `db:"id" json:"id"`
	GroupName   string    `db:"group_name" json:"group"`
	SongName    string    `db:"song_name" json:"song"`
	ReleaseDate time.Time `db:"release_date" json:"releaseDate"`
	Text        string    `db:"text" json:"text"`
	Link        string    `db:"link" json:"link"`
}
