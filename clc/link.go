package clc

type Link struct {
	Rel   string   `json:"rel,omitempty"`
	Href  string   `json:"href,omitempty"`
	ID    string   `json:"id,omitempty"`
	Verbs []string `json:"verbs,omitempty"`
}

type Links []Link

func (l Links) GetID(rel string) (bool, string) {
	for _, v := range l {
		if v.Rel == rel {
			return true, v.ID
		}
	}
	return false, ""
}
