package emails

type APIResponse struct {
	Took     int    `json:"took"`
	TimedOut bool   `json:"timed_out"`
	Shards   Shards `json:"_shards"`
	Hits     Hits   `json:"hits"`
}

type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type Hits struct {
	Total    Total   `json:"total"`
	MaxScore float64 `json:"max_score"`
	Hits     []Hit   `json:"hits"`
}

type Total struct {
	Value int `json:"value"`
}

type Hit struct {
	Index     string  `json:"_index"`
	Type      string  `json:"_type"`
	ID        string  `json:"_id"`
	Score     float64 `json:"_score"`
	Timestamp string  `json:"@timestamp"`
	Source    Email   `json:"_source"`
}

type Email struct {
	ID         int    `json:"ID"`
	Message_ID string `json:"Message-ID"`
	Date       string `json:"Date"`
	From       string `json:"from"`
	To         string `json:"to"`
	Subject    string `json:"subject"`
	Body       string `json:"Body"`
}
