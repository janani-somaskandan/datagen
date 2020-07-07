package operations

import(
	"testing"
	"../utils"
	"../config"
	"time"
)

func TestComputeRangeMap(t *testing.T){
	probMap := make(map[string]float64)
	probMap["E1"]= 0.0001
	probMap["E2"]= 0.5
	probMap["E3"]= 0.4
	probMap["E4"]= 0.0009
	probMap["E5"]= 0.009
	probMap["E6"]= 0.06
	probMap["E7"]= 0.029
	probMap["E8"]= 0.001
	resultRangeMap, resultMultiplier  := ComputeRangeMap(probMap)
	if(resultMultiplier != 10000){
		t.Errorf("Expected: 10000 Result: %v", resultMultiplier)
	}
	if(len(resultRangeMap.Keys) != 8){
		t.Errorf("Expected: 8 Result: %v", len(resultRangeMap.Keys))
	}
	for i := 0; i < len(resultRangeMap.Keys); i++ {
		dataRange := resultRangeMap.Keys[i].U - resultRangeMap.Keys[i].L + 1
		if(resultRangeMap.Values[i] == "E1"){
			if(dataRange != 1){
				t.Errorf("Expected: 1 Result: %v", dataRange)
			}
		}
		if(resultRangeMap.Values[i] == "E2"){
			if(dataRange != 5000){
				t.Errorf("Expected: 5000 Result: %v", dataRange)
			}
		}
		if(resultRangeMap.Values[i] == "E3"){
			if(dataRange != 4000){
				t.Errorf("Expected: 4000 Result: %v", dataRange)
			}
		}
		if(resultRangeMap.Values[i] == "E4"){
			if(dataRange != 9){
				t.Errorf("Expected: 9 Result: %v", dataRange)
			}
		}
		if(resultRangeMap.Values[i] == "E5"){
			if(dataRange != 90){
				t.Errorf("Expected: 90 Result: %v", dataRange)
			}
		}
		if(resultRangeMap.Values[i] == "E6"){
			if(dataRange != 600){
				t.Errorf("Expected: 600 Result: %v", dataRange)
			}
		}
		if(resultRangeMap.Values[i] == "E7"){
			if(dataRange != 290){
				t.Errorf("Expected: 290 Result: %v", dataRange)
			}
		}
		if(resultRangeMap.Values[i] == "E8"){
			if(dataRange != 10){
				t.Errorf("Expected: 10 Result: %v", dataRange)
			}
		}
	}
}

func TestGetRandomEventWithCorrelation(t *testing.T){
	lastKnownGoodState := "E1"
	seedEvents := []string{"E1"}
	var probMap SegmentProbMap
	probMap.EventCorrelationMultiplier = make(map[string]int)
	probMap.EventCorrelationMultiplier["E1"] = 10
	probMap.EventCorrelationMultiplier["E2"] = 10
	probMap.EventCorrelationMultiplier["E3"] = 10
	probMap.EventCorrelationRangeMap = make(map[string]utils.RangeMap)
	probRangeMapE1 := utils.RangeMap{}
	probRangeMapE1.Keys = []utils.Range{utils.Range{0,3}, utils.Range{4,9}}
	probRangeMapE1.Values = []string{"E2", "E3"}
	probMap.EventCorrelationRangeMap["E1"] = probRangeMapE1
	probRangeMapE2 := utils.RangeMap{}
	probRangeMapE2.Keys = []utils.Range{utils.Range{0,9}}
	probRangeMapE2.Values = []string{"E3"}
	probMap.EventCorrelationRangeMap["E2"] = probRangeMapE2
	probRangeMapE3 := utils.RangeMap{}
	probRangeMapE3.Keys = []utils.Range{utils.Range{0,1},utils.Range{1,9}}
	probRangeMapE3.Values = []string{"E3","E1"}
	probMap.EventCorrelationRangeMap["E3"] = probRangeMapE3
	
	result1 := GetRandomEventWithCorrelation(&lastKnownGoodState, seedEvents, probMap)
	if(!(result1 == "E2" || result1 == "E3")){
		t.Errorf("Expected: E2 || E3 Result: %v", result1)
	}
	if(!(lastKnownGoodState == "E2" || lastKnownGoodState == "E3")){
		t.Errorf("Expected: E2 || E3 Result: %v", lastKnownGoodState)
	}
	if(!(result1 == lastKnownGoodState)){
		t.Errorf("Expected to have same result and lastKnownGoodState")
	}

	lastKnownGoodState = ""
	result2 := GetRandomEventWithCorrelation(&lastKnownGoodState, seedEvents, probMap)
	if(!(result2 == lastKnownGoodState)){
		t.Errorf("Expected to have same result and lastKnownGoodState")
	}
	if(!(result2 == "E1")){
		t.Errorf("Expected: E1 Result: %v", result2)
	}

	lastKnownGoodState = "E2"
	result3 := GetRandomEventWithCorrelation(&lastKnownGoodState, seedEvents, probMap)
	if(!(result3 == lastKnownGoodState)){
		t.Errorf("Expected to have same result and lastKnownGoodState")
	}
	if(!(result3 == "E3")){
		t.Errorf("Expected: E3 Result: %v", result2)
	}
}

func TestGetRandomEvent(t *testing.T){
	var probMap SegmentProbMap
	probMap.eventMultiplier = 10
	eventsRangeMap := utils.RangeMap{}
	eventsRangeMap.Keys = []utils.Range{utils.Range{0,3}, utils.Range{4,9}}
	eventsRangeMap.Values = []string{"E4", "E5"}
	probMap.eventProbRangeMap = eventsRangeMap

	result1 := GetRandomEvent(probMap)
	if(!(result1 == "E4" || result1 == "E5")){
		t.Errorf("Expected: E4 || E5 Result: %v", result1)
	}
}

func TestGetRandomActivity(t *testing.T){
	var probMap SegmentProbMap
	probMap.activityMultiplier = 10
	activityRangeMap := utils.RangeMap{}
	activityRangeMap.Keys = []utils.Range{utils.Range{0,3}, utils.Range{4,9}}
	activityRangeMap.Values = []string{"A1", "A2"}
	probMap.activityProbRangeMap = activityRangeMap

	result1 := GetRandomActivity(probMap)
	if(!(result1 == "A1" || result1 == "A2")){
		t.Errorf("Expected: A1 || A2 Result: %v", result1)
	}
}

func TestSetEventAttributes(t *testing.T){
	var userSegment config.UserSegmentV2
	userSegment.Set_attributes = true
	userSegment.Event_attributes = make(map[string]map[string]string)
	attributes := make(map[string]string)
	attributes["Category"] = "C1"
	attributes["Type"] = "T1"
	userSegment.Event_attributes["E1"] = attributes

	result1 := SetEventAttributes(userSegment,"E1")
	if(result1 == nil){
		t.Errorf("Expected: NotNull Result: %v", result1)
	}

	result2 := SetEventAttributes(userSegment,"E2")
	if(!(result2 == nil)){
		t.Errorf("Expected: Null Result: %v", result2)
	}

	userSegment.Set_attributes = false
	result3 := SetEventAttributes(userSegment,"E1")
	if(!(result3 == nil)){
		t.Errorf("Expected: Null Result: %v", result3)
	}
}

func TestSetUserAttributes(t *testing.T){
	var userSegment config.UserSegmentV2
	userSegment.Set_attributes = true
	userSegment.User_attributes = make(map[string]map[string]string)
	attributes := make(map[string]string)
	attributes["Gender"] = "Male"
	attributes["Age"] = "18-25"
	userSegment.User_attributes["U1"] = attributes

	result1 := SetUserAttributes(userSegment,"U1")
	if(result1 == nil){
		t.Errorf("Expected: NotNull Result: %v", result1)
	}

	result2 := SetUserAttributes(userSegment,"U2")
	if(!(result2 == nil)){
		t.Errorf("Expected: Null Result: %v", result2)
	}

	userSegment.Set_attributes = false
	result3 := SetUserAttributes(userSegment,"U1")
	if(!(result3 == nil)){
		t.Errorf("Expected: Null Result: %v", result3)
	}
}

func TestFormatOutput(t *testing.T){
	var userSegment config.UserSegmentV2
	userSegment.Activity_ticker_in_seconds = 1
	userSegment.Start_Time = time.Date(
		2009, 11, 17, 20, 34, 58, 651387237, time.UTC) 
	result1 := FormatOutput(userSegment, "U1", "E1", 1, nil, nil)
	output1 := "{\"user_id\":\"U1\",\"event_name\":\"E1\",\"timestamp\":1258490099,\"user_properties\":null,\"event_properties\":null}"
	if(result1 != output1){
		t.Errorf("Expected %v Result %v", output1, result1)
	}

	attr := make(map[string]string)
	attr["A1"] = "U1"
	result2 := FormatOutput(userSegment, "U1", "E1", 1, attr, attr)
	output2 := "{\"user_id\":\"U1\",\"event_name\":\"E1\",\"timestamp\":1258490099,\"user_properties\":{\"A1\":\"U1\"},\"event_properties\":{\"A1\":\"U1\"}}"
	if(result2 != output2){
		t.Errorf("Expected %v Result %v", output2, result2)
	}
}

func TestIsAllSegmentsDone(t *testing.T){
	segmentStatus := make(map[string]bool)
	segmentStatus["E1"] = true
	segmentStatus["E2"] = false
	segmentStatus["E3"] = true
	result1 := IsAllSegmentsDone(segmentStatus)
	if(result1 == true){
		t.Errorf("Expected false. Result %v",result1)
	}
	segmentStatus["E2"] = true
	result2 := IsAllSegmentsDone(segmentStatus)
	if(result2 == false){
		t.Errorf("Expected true. Result %v",result2)
	}
}

func TestGetRandomSegment(t *testing.T){
	config.ConfigV2.User_segments = make(map[string]config.UserSegmentV2)
	config.ConfigV2.User_segments["Segment1"] = config.UserSegmentV2{}
	result1 := GetRandomSegment()
	if(result1 != "Segment1"){
		t.Errorf("Expected Segment1. Result %v", result1)
	}
	config.ConfigV2.User_segments["Segment2"] = config.UserSegmentV2{}
	result2 := GetRandomSegment()
	if(!(result2 == "Segment1" || result2 == "Segment2")){
		t.Errorf("Expected Segment1 || segment2. Result %v", result2)
	}
}

func TestCreateNewUserProbMap(t *testing.T){
	config.ConfigV2.New_user_probablity = 0.2
	resultMap, resultMultiplier := CreateNewUserProbMap()
	if(resultMultiplier != 10){
		t.Errorf("Expected 10. Result %v", resultMultiplier)
	}
	for item,element := range resultMap.Values {
		if( element == "Insert"){
			if(resultMap.Keys[item].U-resultMap.Keys[item].L+1 != 2){
				t.Errorf("Expected 2. Result %v", (resultMap.Keys[item].U-resultMap.Keys[item].L+1))
			}
		}
		if( element == "NoInsert"){
			if(resultMap.Keys[item].U-resultMap.Keys[item].L+1 != 8){
				t.Errorf("Expected 8. Result %v", (resultMap.Keys[item].U-resultMap.Keys[item].L+1))
			}
		}
	}
}