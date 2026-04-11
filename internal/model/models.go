package model

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type StorageRecord struct {
	ShortedURL  string `json:"shorted_url"`
	OriginalURL string `json:"original_url"`
}
