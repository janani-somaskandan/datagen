package operations

import(
	"../utils"
	"time"
	"../registration"
	"../config"
	"sync"
	Log "../utils/Log"
	"math/rand"
	"strconv"
)

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
	probMap.yesOrNoProbMap = YesOrNoProbablityMap{ 
		ComputeYesOrNoProbablityMap(config.ConfigV2.New_user_probablity),
		ComputeYesOrNoProbablityMap(config.ConfigV2.Custom_event_attribute_probablity),
		ComputeYesOrNoProbablityMap(config.ConfigV2.Custom_user_attribute_probablity)}
	Log.Debug.Printf("RangeMaps %v", probMap)

	// Generate events per USER SEGMENT
	// segmentStatus variable is used to check if all the segments are done executing
	segmentStatus := make(map[string]bool)
	for item,element := range config.ConfigV2.User_segments {
		segmentWg.Add(1)
		segmentStatus[item] = false
		go OperateOnSegment(
			&segmentWg, 
			probMap,
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
				probMap,
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

func OperateOnSegment(segmentWg *sync.WaitGroup, probMap ProbMap, segmentName string, segment config.UserSegmentV2, segmentProbMap SegmentProbMap, userRangeStart int, userRangeEnd int, segmentStatus map[string]bool){

	defer segmentWg.Done()
	var wg sync.WaitGroup
	segmentProbMap.UserToUserAttributeMap = make(map[string]map[string]string)
	segmentProbMap.EventToEventAttributeMap = PreloadEventAttributes(probMap, segment, segmentProbMap)
	Log.Debug.Printf("Main: Operating on %s with User Range %v - %v", segmentName , userRangeStart ,userRangeEnd)
	//Generating events per user in the segment
	for i := userRangeStart; i<= userRangeEnd; i++ {
		wg.Add(1)
		userId := config.ConfigV2.User_id_prefix+strconv.Itoa(i)
		segmentProbMap.UserToUserAttributeMap[userId] = make(map[string]string)
		segmentProbMap.UserToUserAttributeMap[userId] = GetUserAttributes(probMap, segmentProbMap, segment)
		go GenerateEvents(
			&wg,
			probMap,
			segment,
			(int)(config.ConfigV2.Activity_time_in_seconds / segment.Activity_ticker_in_seconds), 
			userId,
			segmentProbMap)
	}
	
	Log.Debug.Printf("Main: Waiting for %s to finish for user Range %v - %v", segmentName , userRangeStart , userRangeEnd)
	wg.Wait()
	segmentStatus[segmentName] = true
	Log.Debug.Printf("Main: %s Completed for user Range %v - %v", segmentName, userRangeStart ,userRangeEnd)
}

func GenerateEvents(wg *sync.WaitGroup,probMap ProbMap, segmentConfig config.UserSegmentV2, activityDuration int, userId string, segmentProbMap SegmentProbMap) {
	
	defer wg.Done()
	rand.Seed(time.Now().UTC().UnixNano())
	var lastKnownGoodState string
	
	// Setting attributes in output
	userAttributes := SetUserAttributes(segmentProbMap, segmentConfig, userId)

	Log.Debug.Printf("Starting %s for duration %v", userId, activityDuration)
    for i := 0; i < activityDuration; i++ {
		
		activity := GetRandomActivity(segmentProbMap)
		// TODO: Janani Have enums for these
		if activity == "DoSomething" {
			event := GetRandomEvent(segmentProbMap)

			if event == "EventCorrelation" {

				event = GetRandomEventWithCorrelation(
					&lastKnownGoodState, 
					segmentConfig.Event_probablity_map.Correlation_matrix.Seed_events, 
					segmentProbMap)

				if(utils.Contains(segmentConfig.Event_probablity_map.Correlation_matrix.Exit_events,event)){
					Log.Debug.Printf("User %s Exit events: %s", userId, event)
					break;
				}
			}
			eventAttributes := SetEventAttributes(segmentProbMap, segmentConfig, event)

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

func PreloadEventAttributes(probMap ProbMap, segmentConfig config.UserSegmentV2, segmentProbMap SegmentProbMap)(map[string]map[string]string){
	eventToEventAttributes := make(map[string]map[string]string)
	for item, _ := range segmentConfig.Event_probablity_map.Correlation_matrix.Events{
		eventToEventAttributes[item] = GetEventAttributes(
			probMap,
			segmentProbMap,
			segmentConfig,
			item)
	}
	for item, _ := range segmentConfig.Event_probablity_map.Independent_events{
		eventToEventAttributes[item] = GetEventAttributes(
			probMap,
			segmentProbMap,
			segmentConfig,
			item)
	}
	return eventToEventAttributes
}
func GetUserAttributes(probMap ProbMap, segmentProbMap SegmentProbMap, segmentConfig config.UserSegmentV2) map[string]string{
	userAttr := make(map[string]string)
	userAttr = PickAttributes(
		segmentConfig.User_attributes.Default,
		segmentProbMap.defaultUserAttrProbMap)
	if(AddCustomUserAttributeOrNot(probMap)){
		utils.AppendMaps(userAttr, PickAttributes(
			segmentConfig.User_attributes.Custom,
			segmentProbMap.customUserAttrProbMap))
	}
	return userAttr
}

func SetUserAttributes(segmentProbMap SegmentProbMap, segmentConfig config.UserSegmentV2, userId string) map[string]string{
	if(segmentConfig.Set_attributes == true){
		return segmentProbMap.UserToUserAttributeMap[userId]
	}
	return nil
}

func GetEventAttributes(probMap ProbMap, segmentProbMap SegmentProbMap, segmentConfig config.UserSegmentV2,eventName string) map[string]string{
	eventAttr := make(map[string]string)
	eventAttr = segmentConfig.Event_attributes.Predefined[eventName]
	utils.AppendMaps(eventAttr, PickAttributes(
		segmentConfig.Event_attributes.Default,
		segmentProbMap.defaultEventAttrProbMap))
	if(AddCustomUserAttributeOrNot(probMap)){
		utils.AppendMaps(eventAttr, PickAttributes(
			segmentConfig.Event_attributes.Custom,
			segmentProbMap.defaultEventAttrProbMap))
	}
	return eventAttr
}

func SetEventAttributes(segmentProbMap SegmentProbMap, segmentConfig config.UserSegmentV2, event string) map[string]string{
	if(segmentConfig.Set_attributes == true){
		return segmentProbMap.EventToEventAttributeMap[event]
	}
	return nil
}
