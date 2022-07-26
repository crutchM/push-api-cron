package models

type InputGroup struct {
	AppId    int    `json:"app_id" db:"app_id"`
	Name     string `json:"name" db:"name"`
	SendRate int    `json:"send_rate" db:"send_rate"`
}

type OutputGroup struct {
	Id       int    `json:"id" db:"id"`
	AppId    int    `json:"app_id" db:"app_id"`
	Name     string `json:"name" db:"name"`
	SendRate int    `json:"send_rate" db:"send_rate"`
}
