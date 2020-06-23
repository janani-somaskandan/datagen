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

func PreComputeRangeMap(probMap map[string]float64) (utils.RangeMap, int) {

// Assuming sum of the probablities of elements is 1
// TODO janani: yet to calculate relative probablities if the sum is not 1
	min := 1.0
		
	//TODO call this from util once you find a way to iterate values
	for _, element := range probMap {
		if element < min {
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
		probRangeMap.Keys = append(probRangeMap.Keys,utils.Range{ start+1, start+int(element) })
		probRangeMap.Values = append(probRangeMap.Values, item)
		start = start + int(element)
	}

	return probRangeMap, int(multiplier)
}


func GenerateEvents(wg *sync.WaitGroup, segmentConfig config.UserSegment, activityDuration int, userId string, eventProbRangeMap utils.RangeMap,
	activityProbRangeMap utils.RangeMap,eventMultiplier int, activityMultiplier int) {
	
	defer wg.Done()
	rand.Seed(time.Now().UTC().UnixNano())
	
	var userAttributes string
	if(segmentConfig.Set_attributes == true){
		attr := segmentConfig.User_attributes[userId]
		if(attr != ""){
			userAttributes = "," + attr
		}
	}

	fmt.Println("Starting ", userId)
    for i := 0; i < activityDuration; i++ {
		activity := rand.Intn(activityMultiplier + 1)
		activityName, _ := activityProbRangeMap.Get(activity)
		// TODO Have enums for these
		if activityName == "DoSomething" {
			event := rand.Intn(eventMultiplier + 1)
			eventName, ok := eventProbRangeMap.Get(event)

			var eventAttributes string
			if(segmentConfig.Set_attributes == true){
				attr := segmentConfig.Event_attributes[eventName]
				if(attr != ""){
					eventAttributes = "," + attr
				}
			}

			if ok == true {
				op := 
					userId +
					"," +
					eventName +
					"," +
					segmentConfig.Start_Time.Add(time.Second * time.Duration(i * segmentConfig.Activity_ticker_in_seconds)).String() +
					userAttributes +
					eventAttributes

				 registration.WriterInstance.Write(op)
				 if(config.Config.Real_Time == true){
					time.Sleep(time.Duration(segmentConfig.Activity_ticker_in_seconds) * time.Second)
				 }
			}
		}	
		if activityName == "Exit" {	
			fmt.Println("Exit ", userId)
			break
		}
	}
	fmt.Println("Finished ", userId)
}

func OperatePerSegment(segmentWg *sync.WaitGroup, segmentName string, segment config.UserSegment, userRangeStart int){

	defer segmentWg.Done()
	var wg sync.WaitGroup
	eventProbRangeMap, eventMultiplier := PreComputeRangeMap(segment.Event_probablity_map)
	activityProbRangeMap, activityMultiplier := PreComputeRangeMap(segment.Activity_probablity_map)

	//TODO Add logic to bring users back
	for i := userRangeStart; i<= userRangeStart + segment.Number_of_users - 1; i++ {
		wg.Add(1)
		go GenerateEvents(
			&wg,
			segment,
			(int)(config.Config.Activity_time_in_seconds / segment.Activity_ticker_in_seconds), 
			config.Config.User_id_prefix+strconv.Itoa(i),
			eventProbRangeMap,
			activityProbRangeMap,
			eventMultiplier,
			activityMultiplier)
	}
	
	fmt.Println("Main: Waiting for ", segmentName ," to finish")
	wg.Wait()
	fmt.Println("Main: ", segmentName ," Completed")
}

func Operate(){

	var segmentWg sync.WaitGroup
	fmt.Println(config.Config)
	var userCounter int = 1
	userIndex := make(map[string]int)
	for item, element := range config.Config.User_segments {
		userIndex[item] = userCounter
		userCounter = userCounter + element.Number_of_users
	}

	for item,element := range config.Config.User_segments {
		segmentWg.Add(1)
		go OperatePerSegment(&segmentWg, item, element, userIndex[item])
	}
	fmt.Println("Main: Waiting for Tasks to finish")
	segmentWg.Wait()
	fmt.Println("Main: Completed")
}


