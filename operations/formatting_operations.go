package operations

/*
This is to format the output in required format
*/

import(
	"encoding/json"
)

func FormatOutput(timeStamp int, userId string, event string, userAttributes map[string]string, eventAttributes map[string]string) (string){
	var op EventOutput 
	op.UserId = userId
	op.Event = event
	op.Timestamp = timeStamp
	op.UserAttributes = userAttributes
	op.EventAttributes = eventAttributes
	e, _ := json.Marshal(&op)
	return string(e)
}

func FormatUserData(userId string, attributes map[string]string)string{
	var op UserDataOutput
	op.UserId = userId
	op.UserAttributes = attributes
	e, _ := json.Marshal(&op)
	return string(e)
}