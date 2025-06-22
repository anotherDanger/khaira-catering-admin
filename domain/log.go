package domain

type HitSource struct {
	Entity    string `json:"entity"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type Hit struct {
	Index  string    `json:"_index"`
	Source HitSource `json:"_source"`
}

type HitsWrapper struct {
	Hits []Hit `json:"hits"`
}

type SearchResponse struct {
	Hits HitsWrapper `json:"hits"`
}
