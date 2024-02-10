package s9y_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/till/s9y-hugo/s9y"
)

// test to confirm the logic from s9y works
func TestUrl(t *testing.T) {
	fixtures := []struct {
		ID       int64
		Title    string
		Expected string
	}{
		{
			ID:       1,
			Title:    "PEAR & Plesk",
			Expected: "1-PEAR-Plesk",
		},
		{
			ID:       2,
			Title:    "Prometheus: relabel your scrape_config",
			Expected: "2-Prometheus-relabel-your-scrape_config",
		},
		{
			ID:       3,
			Title:    "node.js & socket.io fun",
			Expected: "3-node.js-socket.io-fun",
		},
	}

	for _, f := range fixtures {
		t.Run(f.Title, func(t *testing.T) {
			e := s9y.Entry{ID: f.ID, Title: f.Title}
			assert.Equal(t, f.Expected, e.URL())
		})
	}
}
