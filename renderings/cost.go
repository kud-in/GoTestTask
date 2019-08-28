package renderings

type Cost struct {
	Code string  `json:"code"`
	Rate float32 `json:"rate"`
}

type CostResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Payload []Cost `json:"payload"`
}
