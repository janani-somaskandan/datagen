package main

import(	
	"os"
	"fmt"
	"strings"
	"net/http"
	"bytes"
	"encoding/json"
	"bufio"
	"io/ioutil"
	"../operations"
	"../utils"
	Log "../utils/Log"
)

var endpoint string
var authToken string
var bulkLoadUrl string
var getUserIdUrl string
var clientUserIdToUserIdMap map[string]string = make(map[string]string)

func ExtractEventData(data string) interface{} {
	split := strings.Split(data, " ")
	var op operations.EventOutput
	if(len(split) >= 3){
		json.Unmarshal([]byte(split[2]), &op)
	}else {
		//ignore: Fix this
	}
	return op
}

func IngestData(obj interface{}){
	reqBody, _ := json.Marshal(obj)
	url := fmt.Sprintf("%s%s", endpoint, bulkLoadUrl)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		Log.Error.Fatal(err)
	}
	req.Header.Add("Authorization", authToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		//TODO: Janani Handle Retry
		Log.Error.Fatal(err)
	}
	fmt.Println(resp)
}

func getUserId(clientUserId string, eventTimestamp int64) (string, error) {
	userId, found := clientUserIdToUserIdMap[clientUserId]
	if !found {
		// Create a user.
		userRequestMap := make(map[string]interface{})
		userRequestMap["c_uid"] = clientUserId
		userRequestMap["join_timestamp"] = eventTimestamp

		reqBody, _ := json.Marshal(userRequestMap)
		url := fmt.Sprintf("%s%s", endpoint, getUserIdUrl)
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		if err != nil {
			Log.Error.Fatal(err)
		}
		req.Header.Add("Authorization", authToken)
		resp, err := client.Do(req)
		if err != nil {
			Log.Error.Fatal(fmt.Sprintf(
				"Http Post user creation failed. Url: %s, reqBody: %s, response: %+v, error: %+v", url, reqBody, resp, err))
			return "", err
		}
		// always close the response-body, even if content is not required
		defer resp.Body.Close()
		jsonResponse, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			Log.Error.Fatal("Unable to parse http user create response.")
			return "", err
		}
		var jsonResponseMap map[string]interface{}
		json.Unmarshal(jsonResponse, &jsonResponseMap)
		userId = jsonResponseMap["user_id"].(string)
		clientUserIdToUserIdMap[clientUserId] = userId
	}
	return userId, nil
}

func main(){
	Log.RegisterLogFiles()
	// TODO : Check how log can be initialized to a file
	//GET these variables from Config
	root := "."
	maxBatchSize := 1000
	filePrefix := "livspace"
	processedFiles := "ProcessedFiles"
	counter := 0
	var events []operations.EventOutput
	endpoint = os.Args[1]
	authToken = os.Args[2]
	bulkLoadUrl = "/sdk/event/track/bulk"
	getUserIdUrl = "/sdk/user/identify"

	utils.CreateDirectoryIfNotExists(processedFiles)
	files := utils.GetAllUnreadFiles(root, filePrefix)
	Log.Debug.Printf("Files to be processed %v", files)

	for _, element := range files {
		Log.Debug.Printf("Processing contents of File: %s", element)
		zipReader := utils.GetFileHandlegz(fmt.Sprintf("%s/%s",root,element))
		scanner := bufio.NewScanner(zipReader)
		for scanner.Scan() {
			s := scanner.Text()
			op := ExtractEventData(s).(operations.EventOutput)
			if(op.UserId != ""){
				op.UserId, _ = getUserId(op.UserId, (int64)(op.Timestamp))
				events = append(events, op)				
			}
			counter++ 

			if(counter == maxBatchSize){
				Log.Debug.Printf("Processing %v records", len(events))
				//IngestData(events)
				counter = 0
				events = nil
			}
		}
		if(counter != 0){
			Log.Debug.Printf("Processing %v records", len(events))
			//IngestData(events)
			counter = 0
			events = nil
		}
		Log.Debug.Printf("Done !!! Processing contents of File: %s", element)
		// Set clean up for files in Processed folder
		utils.MoveFiles(
			fmt.Sprintf("%s/%s",root,element),
			fmt.Sprintf("%s/%s/%s",root, processedFiles, element))
	} 
}
