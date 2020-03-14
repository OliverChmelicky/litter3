package trash

import (
	"github.com/aodin/aspect/postgis"
)

type TrashModel struct {
	tableName     struct{} `pg:"trash"`
	Id            string
	Cleaned       bool
	Size          size
	Accessibility accessibility
	Gps           postgis.Point
	Description   string
	Finder        string
	Created       int64
}

type size string
type accessibility string

type Collection struct {
	tableName struct{} `pg:"trash"`
	Id        string
	Created   int64
}
