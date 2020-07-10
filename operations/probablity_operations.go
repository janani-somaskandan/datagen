package operations

import(
	"../utils"
	"../config"
	"math"
)

type RangeMapMultiplierTuple struct {
	probRangeMap utils.RangeMap
	multiplier int
}

type SegmentProbMap struct {
	EventCorrelationProbMap map[string]RangeMapMultiplierTuple
	eventProbMap RangeMapMultiplierTuple
	activityProbMap RangeMapMultiplierTuple
	defaultUserAttrProbMap AttributeProbMap
	customUserAttrProbMap AttributeProbMap
	defaultEventAttrProbMap AttributeProbMap
	customEventAttrProbMap AttributeProbMap
	UserToUserAttributeMap map[string]map[string]string
	EventToEventAttributeMap map[string]map[string]string
}

type YesOrNoProbablityMap struct {
	SeedNewUser RangeMapMultiplierTuple
	AddCustomEventAttribute RangeMapMultiplierTuple
	AddCustomUserAttribute RangeMapMultiplierTuple
}
type ProbMap struct {
	yesOrNoProbMap YesOrNoProbablityMap
	segmentProbMap map[string]SegmentProbMap
}

func ComputeYesOrNoProbablityMap(trueProb float64)(RangeMapMultiplierTuple)  {
	probMap := make(map[string]float64)
	probMap["Yes"] = trueProb
	probMap["No"] = (1.0 - trueProb)
	return ComputeRangeMap(probMap)
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

	probMap.defaultUserAttrProbMap, probMap.customUserAttrProbMap = 
		PreComputeUserAttributeProbMap(segment.User_attributes)
	probMap.defaultEventAttrProbMap, probMap.customEventAttrProbMap = 
		PreComputeUserAttributeProbMap(segment.User_attributes)
	return probMap
}

func PreComputeUserAttributeProbMap(userAttributes config.UserAttributes)(AttributeProbMap, AttributeProbMap) {
	SortAttributeMap(userAttributes.Default)
	SortAttributeMap(userAttributes.Custom)
	d := GenerateProbablityMapForAttributes(
		userAttributes.Default)
	c := GenerateProbablityMapForAttributes(
			userAttributes.Custom)
	return d, c
}

func PreComputeEventAttributeProbMap(eventAttributes config.EventAttributes)(AttributeProbMap, AttributeProbMap) {
	SortAttributeMap(eventAttributes.Default)
	SortAttributeMap(eventAttributes.Custom)
	d := GenerateProbablityMapForAttributes(
		eventAttributes.Default)
	c := GenerateProbablityMapForAttributes(
		eventAttributes.Custom)
	return d, c
}