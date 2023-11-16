package response

type Response struct {
	Status bool   `json:"status"`
	Data   any    `json:"data,omitempty"`
	Total  int64  `json:"total,omitempty"`
	Error  string `json:"error,omitempty"`
}
