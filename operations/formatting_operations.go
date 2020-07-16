package operations

/*
This is to format the output in required format
*/

import(
	"encoding/json"
)

func FormatOutput(timeStamp int, userId string, event string, userAttributes map[string]string, eventAttributes map[string]string) (string){

	type output struct {
		UserId string `json:"user_id"`
		Event string `json:"event_name"`
		Timestamp int `json:"timestamp"`
		UserAttributes map[string]string `json:"user_properties"`
		EventAttributes map[string]string `json:"event_properties"`
	}

	var op output 
	op.UserId = userId
	op.Event = event
	op.Timestamp = timeStamp
	op.UserAttributes = userAttributes
	op.EventAttributes = eventAttributes
	e, _ := json.Marshal(&op)
	return string(e)
}