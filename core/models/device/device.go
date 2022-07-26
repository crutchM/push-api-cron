package device

type Device struct {
	Id        string `json:"id" db:"id"`
	PushToken string `json:"push_token" db:"push_token"`
}
