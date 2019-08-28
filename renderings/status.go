package renderings

type IntervalAverage struct {
	Code  string  `json:"code"`
	Last  float32 `json:"last"`
	Day   float32 `json:"day"`
	Week  float32 `json:"week"`
	Month float32 `json:"month"`
}

type StatusResponse struct {
	Message string            `json:"message"`
	Code    int               `json:"code"`
	Payload []IntervalAverage `json:"payload"`
}
