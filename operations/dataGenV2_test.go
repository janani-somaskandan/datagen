package operations

import(
	"testing"
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