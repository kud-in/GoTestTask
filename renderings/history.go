package renderings

type HistoryItem struct {
	Time string  `json:"time"`
	Rate float32 `json:"rate"`
}

type HistoryResponse struct {
	Message string        `json:"message"`
	Code    int           `json:"code"`
	Payload []HistoryItem `json:"payload"`
}
