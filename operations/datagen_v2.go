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
	"math"
	Log "../utils/Log"
)
var globalTimer bool
func OperateV2(){

	//Declaring WaitGroup for SegmentLevel and newUser Concurrency
	var segmentWg sync.WaitGroup
	var newUserWg sync.WaitGroup
	var globalTimerWg sync.WaitGroup
	/*Calculating USERNAME indexing across segments
	Ex: Segment1 has 10 users and Segment2 has 5 users
	Segment1 will have users named U1,U2...U10 and 
	Segment2 will have U11... U15
	New seeded users will have name from U16*/
	var userCounter int = 1
	userIndex := make(map[string]int)
	for item, element := range config.ConfigV2.User_segments {
		userIndex[item] = userCounter
		userCounter = userCounter + element.Number_of_users 
	}
	Log.Debug.Printf("UserIndex Map %v", userIndex)
	/* Pre-Computing the following probablityRangeMaps per segment
		1. Activity
		2. Event
		3. Event Correlation
		4. New User seed probablity
	*/
	var probMap ProbMap
	probMap.segmentProbMap = make(map[string]SegmentProbMap)
	for item, element := range config.ConfigV2.User_segments {
		probMap.segmentProbMap[item] = PreComputeRangeMap(element)
	}
	probMap.newUserProbMap = CreateNewUserProbMap()
	Log.Debug.Printf("RangeMaps %v", probMap)

	// Generate events per USER SEGMENT
	// segmentStatus variable is used to check if all the segments are done executing
	segmentStatus := make(map[string]bool)
	for item,element := range config.ConfigV2.User_segments {
		segmentWg.Add(1)
		segmentStatus[item] = false
		go OperateOnSegment(
			&segmentWg, 
			item, 
			element, 
			probMap.segmentProbMap[item], 
			userIndex[item], 
			userIndex[item] + element.Number_of_users -1, 
			segmentStatus)
	}

	Log.Debug.Printf("Main: Waiting for All Segments to finish")

	allSegmentsDone := false
	//newUserSegmentStatus is used to check if the new users seeded into the system are done executing
	newUserSegmentStatus := make(map[string]bool)
	
	// Seeding new users based on the seed probablity till the pre-defined segments executes
	i := userCounter
	globalTimer = false
	globalTimerWg.Add(1)
	go WaitForNSeconds(&globalTimerWg, config.ConfigV2.Activity_time_in_seconds)
	for (allSegmentsDone == false && IsRealTime() == true) || (IsRealTime() == true && globalTimer == false) {

		WaitIfRealTime(config.ConfigV2.New_user_poll_time)
		if(SeedUserOrNot(probMap) == true) {
			
			seg := GetRandomSegment()
			end := i+config.ConfigV2.Per_tick_new_user_seed_count-1
			Log.Debug.Printf("Getting User %v - %v to the system with Segment %s", i ,end, seg)
			newUserWg.Add(1)
			go OperateOnSegment(
				&newUserWg,
				seg,config.ConfigV2.User_segments[seg],
				probMap.segmentProbMap[seg],
				i,
				end,
				newUserSegmentStatus)
			i = end + 1
			allSegmentsDone = IsAllSegmentsDone(segmentStatus)
				
		}
	}
	segmentWg.Wait()
	Log.Debug.Printf("All Segments - Done !!!")
	newUserWg.Wait()
	Log.Debug.Printf("New Users - Done !!!")
	globalTimerWg.Wait()
	Log.Debug.Printf("Global Timer - Exit !!!")
	Log.Debug.Printf("Main - Done !!!")
}

type RangeMapMultiplierTuple struct {
	probRangeMap utils.RangeMap
	multiplier int
}

type SegmentProbMap struct {
	EventCorrelationProbMap map[string]RangeMapMultiplierTuple
	eventProbMap RangeMapMultiplierTuple
	activityProbMap RangeMapMultiplierTuple
}

type ProbMap struct {
	newUserProbMap RangeMapMultiplierTuple
	segmentProbMap map[string]SegmentProbMap
}

func WaitForNSeconds(wg *sync.WaitGroup, duration int){
	defer wg.Done()
	Log.Debug.Printf("Waiting for Total Activity Time")
	WaitIfRealTime(duration)
	globalTimer = true
}

func CreateNewUserProbMap()(RangeMapMultiplierTuple){
	
	newUserProbablityMap := make(map[string]float64)
	newUserProbablityMap["Insert"] = config.ConfigV2.New_user_probablity
	newUserProbablityMap["NoInsert"] = (1.0 - config.ConfigV2.New_user_probablity)
	return ComputeRangeMap(newUserProbablityMap)
}

func SeedUserOrNot(probMap ProbMap)bool{
	r := rand.Intn(probMap.newUserProbMap.multiplier)
	newUserInsert, state := probMap.newUserProbMap.probRangeMap.Get(r)
	if(state == false){
		Log.Error.Fatal(fmt.Sprintf("NewUserRangeMap: Key not found %v", r))
	}
	Log.Debug.Printf("Seed New User - %s", newUserInsert)
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

func IsAllSegmentsDone(segmentStatus map[string]bool) bool {

	allSegmentsDone := true
	for _,element := range segmentStatus {
		if element == false {
			allSegmentsDone = false
			break
		}
	}
	return allSegmentsDone
}

func OperateOnSegment(segmentWg *sync.WaitGroup, segmentName string, segment config.UserSegmentV2, probMap SegmentProbMap, userRangeStart int, userRangeEnd int, segmentStatus map[string]bool){

	defer segmentWg.Done()
	var wg sync.WaitGroup

	Log.Debug.Printf("Main: Operating on %s with User Range %v - %v", segmentName , userRangeStart ,userRangeEnd)
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
	
	Log.Debug.Printf("Main: Waiting for %s to finish for user Range %v - %v", segmentName , userRangeStart , userRangeEnd)
	wg.Wait()
	segmentStatus[segmentName] = true
	Log.Debug.Printf("Main: %s Completed for user Range %v - %v", segmentName, userRangeStart ,userRangeEnd)
}

func PreComputeRangeMap(segment config.UserSegmentV2) (SegmentProbMap) {

	var probMap SegmentProbMap
	probMap.EventCorrelationProbMap = make(map[string]RangeMapMultiplierTuple)
	for item, element := range segment.Event_probablity_map.Correlation_matrix.Events {
		probMap.EventCorrelationProbMap[item] = ComputeRangeMap(element)
	}

	events := make(map[string]float64)
	sum := 0.0
	if segment.Event_probablity_map.Independent_events != nil {
		for item, element := range segment.Event_probablity_map.Independent_events {
			sum += element
			events[item] = element
		}
	}

	events["EventCorrelation"] = (1.0 - sum)
	probMap.eventProbMap = ComputeRangeMap(events)
	probMap.activityProbMap = ComputeRangeMap(segment.Activity_probablity_map)

	return probMap
}

func GenerateEvents(wg *sync.WaitGroup, segmentConfig config.UserSegmentV2, activityDuration int, userId string, probMap SegmentProbMap) {
	
	defer wg.Done()
	rand.Seed(time.Now().UTC().UnixNano())
	var lastKnownGoodState string
	
	// Setting attributes in output
	userAttributes := SetUserAttributes(segmentConfig, userId)

	Log.Debug.Printf("Starting %s for duration %v", userId, activityDuration)
    for i := 0; i < activityDuration; i++ {
		
		activity := GetRandomActivity(probMap)
		// TODO: Janani Have enums for these
		if activity == "DoSomething" {
			event := GetRandomEvent(probMap)

			if event == "EventCorrelation" {
				event = GetRandomEventWithCorrelation(
					&lastKnownGoodState, 
					segmentConfig.Event_probablity_map.Correlation_matrix.Seed_events, 
					probMap)
				if(utils.Contains(segmentConfig.Event_probablity_map.Correlation_matrix.Exit_events,event)){
					Log.Debug.Printf("User %s Exit events: %s", userId, event)
					break;
				}
			}
			eventAttributes := SetEventAttributes(segmentConfig, event)

			op := FormatOutput(segmentConfig, userId, event, i, userAttributes, eventAttributes)

			registration.WriterInstance.Write(op)
			WaitIfRealTime(segmentConfig.Activity_ticker_in_seconds)
			
		}
		if(activity == "Exit"){
			Log.Debug.Printf("Exit %s", userId)
			break;
		}	
	}
	Log.Debug.Printf("Done %s", userId)
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
	op.Timestamp, _ = strconv.Atoi(fmt.Sprintf("%v", segmentConfig.Start_Time.Add(time.Second * time.Duration(eventCounter * segmentConfig.Activity_ticker_in_seconds)).Unix()))
	op.UserAttributes = userAttributes
	op.EventAttributes = eventAttributes
	e, _ := json.Marshal(&op)
	return string(e)
}

func WaitIfRealTime(duration int) {
	if(IsRealTime() == true){
		time.Sleep(time.Duration(duration) * time.Second)
	}
}

func IsRealTime() bool {
	if(config.ConfigV2.Real_Time == true){
		return true
	}
	return false
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

func GetRandomActivity(probMap SegmentProbMap) string {
	activity := rand.Intn(probMap.activityProbMap.multiplier)
	activityName, state := probMap.activityProbMap.probRangeMap.Get(activity)
	if(state == false){
		Log.Error.Fatal(fmt.Sprintf("ActivityProbablityRangeMap: Key not found %v", activity))
	}
	return activityName
}

func GetRandomEvent(probMap SegmentProbMap) string {
	event := rand.Intn(probMap.eventProbMap.multiplier)
	eventName, state := probMap.eventProbMap.probRangeMap.Get(event)
	if(state == false){
		Log.Error.Fatal(fmt.Sprintf("EventProbablityRangeMap: Key not found %v", event))
	}
	return eventName
}

func GetRandomEventWithCorrelation(lastKnownGoodState *string, seedEvents []string, probMap SegmentProbMap) (string) {
	if *lastKnownGoodState == "" {
		*lastKnownGoodState = seedEvents[rand.Intn(len(seedEvents))]
		return *lastKnownGoodState
	}
    
	event := rand.Intn(probMap.EventCorrelationProbMap[*lastKnownGoodState].multiplier)
	eventName, state := probMap.EventCorrelationProbMap[*lastKnownGoodState].probRangeMap.Get(event)
	if(state == false){
		Log.Error.Fatal(fmt.Sprintf("EventProbablityRangeMapWithCorrelation: Key not found %v", event))
	}
	*lastKnownGoodState = eventName
	return eventName
}

func ComputeRangeMap(probMap map[string]float64) (RangeMapMultiplierTuple) {

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
		rangeEnd := int(math.Round(element * multiplier))
		probRangeMap.Keys = append(probRangeMap.Keys,utils.Range{ start, start+rangeEnd-1 })
		probRangeMap.Values = append(probRangeMap.Values, item)
		start = start + rangeEnd
	}

	return RangeMapMultiplierTuple{ probRangeMap, int(multiplier)}
}