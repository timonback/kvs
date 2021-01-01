package model

import "time"

type Path string

type Item struct {
	Content string
	Time    time.Time
}
