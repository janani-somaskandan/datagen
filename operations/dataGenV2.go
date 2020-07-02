package operations

import(
	"../utils"
	"math/rand"
	"time"
	"strconv"
	"../registration"
	"../config"
	"fmt"
	"sync"
	"reflect"
	"encoding/json"
)

func OperateV2(){

	//Declaring WaitGroup for SegmentLevel and newUser Concurrency
	var segmentWg sync.WaitGroup
	var newUserWg sync.WaitGroup

	// Calculating user count level indexing across segments
	var userCounter int = 1
	userIndex := make(map[string]int)
	for item, element := range config.ConfigV2.User_segments {
		userIndex[item] = userCounter
		userCounter = userCounter + element.Number_of_users 
	}

	// Pre-Computing Probablity RangeMap for ActivityProb, EventProbablity, EventCorrelation
	segmentProbMap := make(map[string]ProbMap)
	for item, element := range config.ConfigV2.User_segments {
		segmentProbMap[item] = PreComputeRangeMap(element)
	}

	// Operating per user segment
	segmentStatus := make(map[string]bool)
	for item,element := range config.ConfigV2.User_segments {
		segmentWg.Add(1)
		segmentStatus[item] = false
		go OperateOnSegment(&segmentWg, item, element, segmentProbMap[item], userIndex[item], userIndex[item] + element.Number_of_users -1, segmentStatus)
	}

	fmt.Println("Main: Waiting for Tasks to finish")

	allSegmentsDone := false
	newUserSegmentStatus := make(map[string]bool)
	CreateNewUserProbMap()
	
	for i := userCounter; allSegmentsDone == false; i++ {

		time.Sleep(time.Duration(config.ConfigV2.New_user_poll_time) * time.Second)
		if(GetRandomNewUserInsertStatus() == true) {
			
			seg := GetRandomSegment()
			fmt.Println("Getting User ", i ," to the system with Segment ", seg)
			newUserWg.Add(1)
			go OperateOnSegment(&newUserWg,seg,config.ConfigV2.User_segments[seg],segmentProbMap[seg],i,i,newUserSegmentStatus)
			allSegmentsDone = CheckIfAllUsersDone(segmentStatus)
				
		}
	}

	fmt.Println("Exit")
	newUserWg.Wait()
	segmentWg.Wait()
	fmt.Println("Main: Completed")
}

var newUserRangeMap utils.RangeMap
var newUserMultiplier int
func CreateNewUserProbMap(){
	
	newUserProbablityMap := make(map[string]float64)
	newUserProbablityMap["Insert"] = config.ConfigV2.New_user_probablity
	newUserProbablityMap["NoInsert"] = (1.0 - config.ConfigV2.New_user_probablity)
	newUserRangeMap, newUserMultiplier = ComputeRangeMap(newUserProbablityMap)
}

func GetRandomNewUserInsertStatus()bool{
	newUserInsert, _ := newUserRangeMap.Get(rand.Intn(newUserMultiplier))
	if(newUserInsert == "Insert") {
		return true
	}
	return false
}

func GetRandomSegment()string{
	segmentKeys := reflect.ValueOf(config.ConfigV2.User_segments).MapKeys()
	seg := (segmentKeys[rand.Intn(len(segmentKeys))].Interface()).(string)
	return seg
}
func CheckIfAllUsersDone(segmentStatus map[string]bool) bool {

	allSegmentsDone := true
	for _,element := range segmentStatus {
		if element == false {
			allSegmentsDone = false
			break
		}
	}
	return allSegmentsDone
}

func OperateOnSegment(segmentWg *sync.WaitGroup, segmentName string, segment config.UserSegmentV2,probMap ProbMap,userRangeStart int, userRangeEnd int, segmentStatus map[string]bool){

	defer segmentWg.Done()
	var wg sync.WaitGroup

	fmt.Println("Main: Operating on ", segmentName ," with User Range ", userRangeStart , "-" ,userRangeEnd)
	//Generating events per user in the segment
	for i := userRangeStart; i<= userRangeEnd; i++ {
		wg.Add(1)
		go GenerateEvents(
			&wg,
			segment,
			(int)(config.ConfigV2.Activity_time_in_seconds / segment.Activity_ticker_in_seconds), 
			config.ConfigV2.User_id_prefix+strconv.Itoa(i),
			probMap)
	}
	
	fmt.Println("Main: Waiting for ", segmentName ," to finish for user Range ", userRangeStart , "-" ,userRangeEnd)
	wg.Wait()
	segmentStatus[segmentName] = true
	fmt.Println("Main: ", segmentName ," Completed for user Range ", userRangeStart , "-" ,userRangeEnd)
}

type ProbMap struct {
	EventCorrelationRangeMap map[string]utils.RangeMap
	EventCorrelationMultiplier map[string]int
	eventProbRangeMap utils.RangeMap
	eventMultiplier int
	activityProbRangeMap utils.RangeMap
	activityMultiplier int
}

func PreComputeRangeMap(segment config.UserSegmentV2) (ProbMap) {

	var probMap ProbMap
	// Pre-Computing Probablity RangeMap for ActivityProb, EventProbablity, EventCorrelation
	probMap.EventCorrelationRangeMap = make(map[string]utils.RangeMap)
	probMap.EventCorrelationMultiplier = make(map[string]int)
	for item, element := range segment.Event_probablity_map.Correlation_matrix.Events {
		probMap.EventCorrelationRangeMap[item], probMap.EventCorrelationMultiplier[item] = ComputeRangeMap(element)
	}

	events := make(map[string]float64)
	sum := 0.0
	for item, element := range segment.Event_probablity_map.Independent_events {
		sum += element
		events[item] = element
	}

	events["EventCorrelation"] = (1.0 - sum)
	probMap.eventProbRangeMap, probMap.eventMultiplier = ComputeRangeMap(events)
	probMap.activityProbRangeMap, probMap.activityMultiplier = ComputeRangeMap(segment.Activity_probablity_map)

	return probMap
}

func GenerateEvents(wg *sync.WaitGroup, segmentConfig config.UserSegmentV2, activityDuration int, userId string, probMap ProbMap) {
	
	defer wg.Done()
	rand.Seed(time.Now().UTC().UnixNano())
	var lastKnownGoodState string
	
	// Setting attributes in output
	userAttributes := SetUserAttributes(segmentConfig, userId)

	fmt.Println("Starting ", userId, "for duration ", activityDuration)
    for i := 0; i < activityDuration; i++ {
		
		activity := GetRandomActivity(probMap)
		// TODO: Janani Have enums for these
		if activity == "DoSomething" {
			event := GetRandomEvent(probMap)

			if event == "EventCorrelation" {
				event = GetRandomEventWithCorrelation(&lastKnownGoodState, segmentConfig.Event_probablity_map.Correlation_matrix.Seed_events, probMap)
			}
			eventAttributes := SetEventAttributes(segmentConfig, event)

			op := FormatOutput(segmentConfig, userId, event, i, userAttributes, eventAttributes)

			registration.WriterInstance.Write(op)
			WaitIfRealTime(config.ConfigV2.Real_Time, segmentConfig.Activity_ticker_in_seconds)
			
		}	
	}
	fmt.Println("Finished ", userId)
}

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
	op.Timestamp, _ = strconv.Atoi(fmt.Sprintf("%v", segmentConfig.Start_Time.Add(time.Second * time.Duration(eventCounter * segmentConfig.Activity_ticker_in_seconds)).UnixNano()))
	op.UserAttributes = userAttributes
	op.EventAttributes = eventAttributes
	e, _ := json.Marshal(&op)
	return string(e)
}

func WaitIfRealTime(realTime bool, duration int) {
	if(realTime == true){
		time.Sleep(time.Duration(duration) * time.Second)
	}
}

func SetUserAttributes(segmentConfig config.UserSegmentV2, userId string) map[string]string{
	var userAttributes map[string]string
	if(segmentConfig.Set_attributes == true){
		attr := segmentConfig.User_attributes[userId]
		if(attr != nil){
			userAttributes = attr
		}
	}
	return userAttributes
}

func SetEventAttributes(segmentConfig config.UserSegmentV2,eventName string) map[string]string{
	var eventAttributes map[string]string
	if(segmentConfig.Set_attributes == true){
		attr := segmentConfig.Event_attributes[eventName]
		if(attr != nil){
			eventAttributes = attr
		}
	}
	return eventAttributes
}

func GetRandomActivity(probMap ProbMap) string {
	activity := rand.Intn(probMap.activityMultiplier)
	activityName, _ := probMap.activityProbRangeMap.Get(activity)
	return activityName
}

func GetRandomEvent(probMap ProbMap) string {
	event := rand.Intn(probMap.eventMultiplier)
	eventName, _ := probMap.eventProbRangeMap.Get(event)
	return eventName
}

func GetRandomEventWithCorrelation(lastKnownGoodState *string, seedEvents []string, probMap ProbMap) (string) {
	if *lastKnownGoodState == "" {
		*lastKnownGoodState = seedEvents[rand.Intn(len(seedEvents))]
		return *lastKnownGoodState
	}
    
	event := rand.Intn(probMap.EventCorrelationMultiplier[*lastKnownGoodState])
	eventName, _ := probMap.EventCorrelationRangeMap[*lastKnownGoodState].Get(event)
	*lastKnownGoodState = eventName
	return eventName
}

func ComputeRangeMap(probMap map[string]float64) (utils.RangeMap, int) {

	// Assuming sum of the probablities of elements is 1
	// TODO janani: yet to calculate relative probablities if the sum is not 1
	min := 1.0
		
	//TODO call this from util once you find a way to iterate values
	for _, element := range probMap {
		if element < min && element != 0{
			min = element
		}
	}

	multiplier := 1.0
	temp := 0.0
	for temp < 1.0 {
		multiplier = multiplier * 10.0
		temp = min * multiplier;
	}

	start := 0
	probRangeMap := utils.RangeMap{}
	for item,element := range probMap {
		probRangeMap.Keys = append(probRangeMap.Keys,utils.Range{ start, start+int(element * multiplier)-1 })
		probRangeMap.Values = append(probRangeMap.Values, item)
		start = start + int(element * multiplier)
	}

	return probRangeMap, int(multiplier)
}