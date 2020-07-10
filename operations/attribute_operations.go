package operations

/*
File with all attribute level operations
*/

import(
	"sort"
	"../config"
	"fmt"
	Log "../utils/Log"
)


type AttributeProbMap struct{
	Attributes_Order1 map[string]RangeMapMultiplierTuple
	Attributes_OrderOver1 map[string]map[string]RangeMapMultiplierTuple
}

func Convert1(input map[string]interface{}) map[string]float64{
	op := make(map[string]float64)
	for k, v := range input {
		op[k] = v.(float64)
	}
	return op
}

func Convert2(input map[interface{}]interface{}) map[string]float64{
	op := make(map[string]float64)
	for k, v := range input {
		op[k.(string)] = v.(float64)
	}
	return op
}

func SortAttributeMap(attributes []config.AttributeData){
	Log.Debug.Printf("Attributes before sorting: %v", attributes)
	sort.Slice(attributes, func(i,j int) bool {
		return attributes[i].Order_Level < attributes[j].Order_Level
	})
	Log.Debug.Printf("User Attributes After sorting: %v", attributes)
}

func GenerateProbablityMapForAttributes(attributes []config.AttributeData) (AttributeProbMap){
	var attributeProbMap AttributeProbMap
	attributeProbMap.Attributes_Order1 = make(map[string]RangeMapMultiplierTuple)
	attributeProbMap.Attributes_OrderOver1 = make(map[string]map[string]RangeMapMultiplierTuple)
	for _, element := range attributes {
		if(element.Order_Level == 1){
			attributeProbMap.Attributes_Order1[element.Key] = 
				ComputeRangeMap(Convert1(element.Values))
		}
		if(element.Order_Level > 1){
			attributeProbMap.Attributes_OrderOver1[element.Key] = make(map[string]RangeMapMultiplierTuple)
			for key, value := range element.Values {
				attributeProbMap.Attributes_OrderOver1[element.Key][key] = 
					ComputeRangeMap(Convert2(value.(map[interface{}]interface{})))
			}
		}
	}
	Log.Debug.Printf("UserAttributes Probablity Map: \n %v", attributeProbMap)
	return attributeProbMap
}

func GenerateProbablityMapForCustomUserAttributes(){
// Both orderLevel 1 and over 1
}

func PickAttributes(attributes []config.AttributeData, attributeProbMap AttributeProbMap)(map[string]string){
// Compulsory choose each default attributes
// iterate over each attribute Eg: gender, age
// For orderLevel = 1 choose rand from probablity map directly
// For orderLevel = 2 choose the value of orderLevel = 1 attribute and go one level down 
// to get the actual value of the required attribute from the probablity map
// OrderLevel can be to any level
// but key is currently assumed to be only one and not multiple

// for custom attributes
// check if yes/no
// check how many attributes to pick (random again) - for now we will assume all
//append both and return
// have a global map for user to attributes map
	attr := make(map[string]string)
	for _, element := range attributes {
		if(element.Order_Level == 1) {
			attr[element.Key] = GetRandomValueWithProbablity(
				attributeProbMap.Attributes_Order1[element.Key], 
				fmt.Sprintf("UserAttribute:%s",element.Key))
		}
		if(element.Order_Level > 1) {
			attr[element.Key] = GetRandomValueWithProbablity(
				attributeProbMap.Attributes_OrderOver1[element.Key][attr[element.Dependency]], 
				fmt.Sprintf("UserAttribute:%s",element.Dependency))
		}
	}
	return attr
}

func SelectEventAttributes(){
// Pick the default attributes at individual event level
// Pick the custom attributes depending on yes/no
// for now we will pick all
// nothing to assign to global map
}

func ReassignProbOrNot(){
// Take the event - user attributes
// Check if it matches
// Global level check need to happen to check if the sum of probablities = 1
// Reassign probablity
// Pick next event accordingly - adjusted to new probablity
}

func ReassignProbablity(){
	// Compute new probablities - Should be straight forward
	// Generate new range map
	// check 
}
