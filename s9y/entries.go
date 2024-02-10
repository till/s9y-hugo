package s9y

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// create a s9y 'object'
func New(db *sqlx.DB, prefix string) *S9Y {
	return &S9Y{
		db:     db,
		prefix: prefix,
	}
}

func (s *S9Y) Entries(blogUrl string) ([]Entry, error) {
	q := fmt.Sprintf(`
		SELECT
			id, title,
			timestamp, last_modified,
			body, extended, isdraft, exflag
			author,
			p.permalink as permalink,
			category_name, category_description,
			pc.permalink as category_link
		FROM %sentries e
		LEFT JOIN %sentrycat ec
			ON (e.ID = ec.entryid)
		LEFT JOIN %scategory c
			ON (ec.categoryid = c.categoryid)
		LEFT JOIN %spermalinks p
			ON (e.ID = p.entry_id AND p.type = 'entry')
		LEFT JOIN %spermalinks pc
			ON (c.categoryid = pc.entry_id AND pc.type = 'category')
		ORDER BY id ASC`,
		s.prefix, s.prefix, s.prefix, s.prefix, s.prefix)
	rows, err := s.db.Queryx(q)
	if err != nil {
		return nil, fmt.Errorf("Entries (%s): %v", q, err)
	}
	defer rows.Close()

	var entries []Entry

	for rows.Next() {
		var entry Entry
		if err := rows.StructScan(&entry); err != nil {
			return nil, fmt.Errorf("Entries: %v", err)
		}

		entry.Permalink = fmt.Sprintf("%s/%s", blogUrl, entry.Permalink)

		entry.Tags, err = s.tags(entry.ID)
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func (s *S9Y) tags(e int64) ([]EntryTag, error) {
	q := fmt.Sprintf(`SELECT * FROM %sentrytags WHERE entryid = ?`, s.prefix)
	rows, err := s.db.Queryx(q, e)
	if err != nil {
		return nil, fmt.Errorf("tags: %v", err)
	}
	defer rows.Close()

	var tags []EntryTag
	for rows.Next() {
		var tag EntryTag
		if err := rows.StructScan(&tag); err != nil {
			return nil, fmt.Errorf("tags: %v", err)
		}

		tags = append(tags, tag)
	}
	return tags, nil
}
