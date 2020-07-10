package operations

import(
	"encoding/json"
	"../config"
	"fmt"
	"strconv"
	"time"
)

func FormatOutput(segmentConfig config.UserSegmentV2, userId string, event string, eventCounter int, userAttributes map[string]string, eventAttributes map[string]string) (string){

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
	op.Timestamp, _ = strconv.Atoi(fmt.Sprintf("%v", segmentConfig.Start_Time.Add(time.Second * time.Duration(eventCounter * segmentConfig.Activity_ticker_in_seconds)).Unix()))
	op.UserAttributes = userAttributes
	op.EventAttributes = eventAttributes
	e, _ := json.Marshal(&op)
	return string(e)
}