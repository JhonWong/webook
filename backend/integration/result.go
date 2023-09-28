package integration

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
	Data T      `json:"data"`
}
