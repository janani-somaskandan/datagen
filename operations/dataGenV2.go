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
)

func OperateV2(){

	//Declaring WaitGroup for SegmentLevel Concurrency
	var segmentWg sync.WaitGroup

	// Calculating user count level indexing across segments
	var userCounter int = 1
	userIndex := make(map[string]int)
	for item, element := range config.ConfigV2.User_segments {
		userIndex[item] = userCounter
		userCounter = userCounter + element.Number_of_users 
	}

	// Operating per user segment
	for item,element := range config.ConfigV2.User_segments {
		segmentWg.Add(1)
		go OperateOnSegment(&segmentWg, item, element, userIndex[item])
	}

	fmt.Println("Main: Waiting for Tasks to finish")
	segmentWg.Wait()
	fmt.Println("Main: Completed")
}

func OperateOnSegment(segmentWg *sync.WaitGroup, segmentName string, segment config.UserSegmentV2, userRangeStart int){

	defer segmentWg.Done()
	var wg sync.WaitGroup

	fmt.Println("Main: Operating on ", segmentName ," with User Range ", userRangeStart , "-" ,userRangeStart + segment.Number_of_users - 1)
	// Pre-Computing Probablity RangeMap for ActivityProb, EventProbablity, EventCorrelation
	probMap := PreComputeRangeMap(segment)

	//TODO: janani Add logic to bring users back
	//Generating events per user in the segment
	for i := userRangeStart; i<= userRangeStart + segment.Number_of_users - 1; i++ {
		wg.Add(1)
		go GenerateEvents(
			&wg,
			segment,
			(int)(config.ConfigV2.Activity_time_in_seconds / segment.Activity_ticker_in_seconds), 
			config.ConfigV2.User_id_prefix+strconv.Itoa(i),
			probMap)
	}
	
	fmt.Println("Main: Waiting for ", segmentName ," to finish")
	wg.Wait()
	fmt.Println("Main: ", segmentName ," Completed")
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
	for _, element := range segment.Event_probablity_map.Independent_events {
		sum += element
	}
	events =  segment.Event_probablity_map.Independent_events
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

			op := 
				userId +
				"," +
				event +
				"," +
				segmentConfig.Start_Time.Add(time.Second * time.Duration(i * segmentConfig.Activity_ticker_in_seconds)).String() +
				userAttributes +
				eventAttributes

			registration.WriterInstance.Write(op)
			WaitIfRealTime(config.ConfigV2.Real_Time, segmentConfig.Activity_ticker_in_seconds)
			
		}	
	}
	fmt.Println("Finished ", userId)
}

func WaitIfRealTime(realTime bool, duration int) {
	if(realTime == true){
		time.Sleep(time.Duration(duration) * time.Second)
	}
}

func SetUserAttributes(segmentConfig config.UserSegmentV2, userId string) string{
	var userAttributes string
	if(segmentConfig.Set_attributes == true){
		attr := segmentConfig.User_attributes[userId]
		if(attr != ""){
			userAttributes = "," + attr
		}
	}
	return userAttributes
}

func SetEventAttributes(segmentConfig config.UserSegmentV2,eventName string) string{
	var eventAttributes string
	if(segmentConfig.Set_attributes == true){
		attr := segmentConfig.Event_attributes[eventName]
		if(attr != ""){
			eventAttributes = "," + attr
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


	for item,element := range probMap {
		probMap[item] = element * multiplier
	}		

	start := 0
	probRangeMap := utils.RangeMap{}
	for item,element := range probMap {
		probRangeMap.Keys = append(probRangeMap.Keys,utils.Range{ start, start+int(element)-1 })
		probRangeMap.Values = append(probRangeMap.Values, item)
		start = start + int(element)
	}

	return probRangeMap, int(multiplier)
}