package clc

type Link struct {
	Rel   string   `json:"rel,omitempty"`
	Href  string   `json:"href,omitempty"`
	ID    string   `json:"id,omitempty"`
	Verbs []string `json:"verbs,omitempty"`
}
