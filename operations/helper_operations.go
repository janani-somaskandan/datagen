package operations

import(
	"../utils"
	"../config"
	"bufio"
	"strings"
	"encoding/json"
	"fmt"
)

func ExtractUserData(data string)(string, map[string]string){
	split := strings.Split(data, " ")
	var op UserDataOutput
	if(len(split) == 3){
		json.Unmarshal([]byte(split[2]), &op)
	}else {
		//ignore: Fix this
	}
	return op.UserId, op.UserAttributes
}

func LoadExistingUsers()map[string]map[string]string {
	files := utils.GetAllUnreadFiles(".",config.ConfigV2.User_data_file_name_prefix)
	userData := make(map[string]map[string]string)
	for _, element := range files {
		reader := utils.GetFileHandlegz(element)
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			s := scanner.Text()
			userId, attributes := ExtractUserData(s)
			userData[userId] = attributes
		}
	}
	reader := utils.GetFileHandle(fmt.Sprintf("%s.log",config.ConfigV2.User_data_file_name_prefix))
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		s := scanner.Text()
		userId, attributes := ExtractUserData(s)
		userData[userId] = attributes
	}
	return userData
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

func UserAlreadyExists(userId string,attributes map[string]map[string]string) bool{
	if(attributes[userId] != nil){
		return true
	}
	return false
}