package models

type InputGroup struct {
	AppId    int    `json:"app_id"`
	Name     string `json:"name"`
	SendRate string `json:"send_rate"`
}

type OutputGroup struct {
	Id       int    `json:"id"`
	AppId    int    `json:"app_id"`
	Name     string `json:"name"`
	SendRate string `json:"send_rate"`
}
