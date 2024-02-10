package s9y

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type S9Y struct {
	db     *sqlx.DB
	prefix string
}

type Entry struct {
	ID                  int64          `db:"id"`
	Title               string         `db:"title"`
	Timestamp           int64          `db:"timestamp"`
	Body                string         `db:"body"`
	Extended            string         `db:"extended"`
	Author              string         `db:"author"`
	IsDraft             bool           `db:"isdraft"`
	ExFlag              bool           `db:"exflag"`
	LastModified        int64          `db:"last_modified"`
	Permalink           string         `db:"permalink"`
	CategoryName        sql.NullString `db:"category_name"`
	CategoryDescription sql.NullString `db:"category_description"`
	CategoryLink        sql.NullString `db:"category_link"`
	Tags                []EntryTag
}

type EntryTag struct {
	EntryID int64  `db:"entryid"`
	Tag     string `db:"tag"`
}
