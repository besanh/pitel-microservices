package model

type SearchReponse struct {
	ScrollId string `json:"_scroll_id"`
	Took     int    `json:"took"`
	TimedOut bool   `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore any `json:"max_score"`
		Hits     []struct {
			Index   string   `json:"_index"`
			Type    string   `json:"_type"`
			ID      string   `json:"_id"`
			Score   any      `json:"_score"`
			Routing string   `json:"_routing"`
			Source  any      `json:"_source"`
			Sort    []string `json:"sort"`
			Fields  struct {
				NotesList []any `json:"notes_list"`
			} `json:"fields"`
		} `json:"hits"`
	} `json:"hits"`
}
