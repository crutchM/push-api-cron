package models

import (
	"math/rand"
	"time"
)

type Push struct {
	PushBatchRequest PushBatchRequest `json:"push_batch_request"`
}

type PushBatchRequest struct {
	GroupId          int     `json:"group_id"`
	ClientTransferId int64   `json:"client_transfer_id"`
	Tag              string  `json:"tag"`
	Batch            []Batch `json:"batch"`
}

func NewPushBatchRequest(groupId int, batch []Batch) PushBatchRequest {
	rand.Seed(time.Now().UnixNano())
	return PushBatchRequest{
		GroupId:          groupId,
		ClientTransferId: rand.Int63(),
		Tag:              "tag",
		Batch:            batch}
}

type Batch struct {
	Messages Messages `json:"messages"`
	Device   []Device `json:"devices"`
}

type Messages struct {
	Android Android `json:"android"`
}

type Android struct {
	Silent  bool    `json:"silent"`
	Content Content `json:"content"`
}

type Content struct {
	Title            string `json:"title,omitempty"`
	Text             string `json:"text,omitempty"`
	Icon             string `json:"icon,omitempty"`
	Image            string `json:"image,omitempty"`
	Banner           string `json:"banner,omitempty"`
	Data             string `json:"data,omitempty"`
	ChannelId        string `json:"channel_id,omitempty"`
	Priority         int    `json:"priority,omitempty"`
	CollapseKey      int    `json:"collapse_key,omitempty"`
	Vibration        []int  `json:"vibration,omitempty"`
	LedInterval      int    `json:"led_interval,omitempty"`
	LedPauseInterval int    `json:"led_pause_interval,omitempty"`
	TimeToAlive      int    `json:"time_to_alive,omitempty"`
	Visibility       string `json:"visibility,omitempty"`
	Urgency          string `json:"urgency,omitempty"`
}

type Device struct {
	IdType   string   `json:"id_type"`
	IdValues []string `json:"id_values"`
}
