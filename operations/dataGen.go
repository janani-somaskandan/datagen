package operations

import(
	"../utils"
	"math/rand"
	"time"
	"strconv"
	"../registration"
	"../config"
	"fmt"
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


func GenerateEvents(activityDuration int, userId string, eventProbRangeMap utils.RangeMap,
	activityProbRangeMap utils.RangeMap,eventMultiplier int, activityMultiplier int) {
	
	rand.Seed(time.Now().UTC().UnixNano())

    for i := 0; i < activityDuration; i++ {
		activity := rand.Intn(activityMultiplier + 1)
		activityName, _ := activityProbRangeMap.Get(activity)
		fmt.Println(activityName)
		// TODO Have enums for these
		if activityName == "DoSomething" {
			event := rand.Intn(eventMultiplier + 1)
			eventName, ok := eventProbRangeMap.Get(event)
			if ok == true {
				op := 
					userId +
					"," +
					eventName +
					"," +
					config.Config.Start_Time.Add(time.Second * time.Duration(i * config.Config.Activity_ticker_in_seconds)).String()
				 registration.WriterInstance.Write(op)
				 if(config.Config.Real_Time == true){
					time.Sleep(time.Duration(config.Config.Activity_ticker_in_seconds) * time.Second)
				 }
			}
		}	
		if activityName == "Exit" {	
			break
		}
	}
}

func Operate(){

	fmt.Println(config.Config)
	eventProbRangeMap, eventMultiplier := PreComputeRangeMap(config.Config.Event_probablity_map)
	activityProbRangeMap, activityMultiplier := PreComputeRangeMap(config.Config.Activity_probablity_map)

	//TODO Add logic to bring users back
	for i := 1; i<= config.Config.Number_of_users; i++ {
		go GenerateEvents(
			(int)(config.Config.Activity_time_in_seconds / config.Config.Activity_ticker_in_seconds), 
			"U"+strconv.Itoa(i),
			eventProbRangeMap,
			activityProbRangeMap,
			eventMultiplier,
			activityMultiplier)
	}
	
	//TODO Wait for events instead of having a standard wait
	time.Sleep(30 * time.Second)
}


