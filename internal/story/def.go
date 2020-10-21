package story

import "github.com/grafov/kiwi"

var log = kiwi.Fork().With("pkg", "story")

type Story struct {
	Category string `json:"category" db:"category"`
	Title    string `json:"title" db:"title"`
	Type     string `json:"type" db:"type"`
	URL      string `json:"url,omitempty"`
	Text     string `json:"text,omitempty"`
	Body     string `db:"body"` // url or text stored to body
}
