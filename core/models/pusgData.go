package models

type PushContent struct {
	GroupId          int    `json:"group_id"`
	ClientTransferId int64  `json:"client_transfer_id"`
	Tag              string `json:"tag"`
}

type PushData struct {

	//TODO:согласовать контент пушей
}

type Message struct {
}

type Android struct {
	Silent  bool    `json:"silent"`
	Content Content `json:"content"`
}

type Content struct {
	Data    string `json:"data"`
	Urgency string `json:"urgency"`
}
