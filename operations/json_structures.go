package operations

type UserDataOutput struct {
	UserId string `json:"user_id"`
	UserAttributes map[string]string `json:"user_properties"`
}

type EventOutput struct {
	UserId string `json:"user_id"`
	Event string `json:"event_name"`
	Timestamp int `json:"timestamp"`
	UserAttributes map[string]string `json:"user_properties"`
	EventAttributes map[string]string `json:"event_properties"`
}