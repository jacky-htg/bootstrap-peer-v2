package pkg

// Request adalah model untuk request dari client.
type Request struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}
