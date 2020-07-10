package operations

import(
	"time"
	"sync"
	Log "../utils/Log"
	"../config"
)

var globalTimer bool

func WaitForNSeconds(wg *sync.WaitGroup, duration int){
	defer wg.Done()
	Log.Debug.Printf("Waiting for Total Activity Time")
	WaitIfRealTime(duration)
	globalTimer = true
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