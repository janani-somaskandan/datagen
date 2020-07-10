package operations

import(	
	"math/rand"
	Log "../utils/Log"
	"fmt"
	"reflect"
	"../config"
)

func SeedUserOrNot(probMap ProbMap)bool{
	return DecideYesOrNo(probMap.yesOrNoProbMap.SeedNewUser, "Seed New User")
}

func AddCustomEventAttributeOrNot(probMap ProbMap)bool {
	return DecideYesOrNo(probMap.yesOrNoProbMap.AddCustomEventAttribute, "Add custom Event Attribute")
}

func AddCustomUserAttributeOrNot(probMap ProbMap)bool {
	return DecideYesOrNo(probMap.yesOrNoProbMap.AddCustomUserAttribute, "Add Custom User Attribute")
}

func DecideYesOrNo(rangeMap RangeMapMultiplierTuple, tag string)bool{
	r := rand.Intn(rangeMap.multiplier)
	yesOrNo, state := rangeMap.probRangeMap.Get(r)
	if(state == false){
		Log.Error.Fatal(fmt.Sprintf("Tag: %s Key not found %v",tag, r))
	}
	if(yesOrNo == "Yes") {
		return true
	}
	return false
}

func GetRandomSegment()string{
	segmentKeys := reflect.ValueOf(config.ConfigV2.User_segments).MapKeys()
	seg := (segmentKeys[rand.Intn(len(segmentKeys))].Interface()).(string)
	return seg
}

func GetRandomActivity(probMap SegmentProbMap) string {
	return GetRandomValueWithProbablity(probMap.activityProbMap, "Activity")
}

func GetRandomEvent(probMap SegmentProbMap) string {
	return GetRandomValueWithProbablity(probMap.eventProbMap, "Event")
}

func GetRandomEventWithCorrelation(lastKnownGoodState *string, seedEvents []string, probMap SegmentProbMap) (string) {
	if *lastKnownGoodState == "" {
		*lastKnownGoodState = seedEvents[rand.Intn(len(seedEvents))]
		return *lastKnownGoodState
	}
	*lastKnownGoodState = GetRandomValueWithProbablity(
		probMap.EventCorrelationProbMap[*lastKnownGoodState], fmt.Sprintf("EventWithCorrelation:%s",*lastKnownGoodState))
	return *lastKnownGoodState
}

func GetRandomValueWithProbablity(rangeMap RangeMapMultiplierTuple, tag string) string {
	r := rand.Intn(rangeMap.multiplier)
	value, state := rangeMap.probRangeMap.Get(r)
	if(state == false){
		Log.Error.Fatal(fmt.Sprintf("Tag: %s, RangeMap: Key not found %v", tag, r))
	}
	return value
}