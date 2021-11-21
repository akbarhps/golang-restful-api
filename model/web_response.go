package model

type WebResponse struct {
	Code   int                 `json:"code"`
	Status string              `json:"status"`
	Errors []map[string]string `json:"errors"`
	Data   interface{}         `json:"data"`
}
