package model

import "time"

type Comment struct {
	Id           uint64
	PostId       uint64
	Comment      string
	Author       string
	CreationDate time.Time
}

type Post struct {
	Id           uint64
	Title        string
	Content      string
	CreationDate time.Time
}
