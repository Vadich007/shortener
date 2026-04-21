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

type BatchResponse struct {
	Records []BatchRecordResponse
}

type BatchRecordResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortedURL    string `json:"shorted_url"`
}

type BatchRequest struct {
	Records []BatchRecordRequest
}

type BatchRecordRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}
