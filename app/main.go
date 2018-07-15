package app

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var lynlabHeaders map[string]string

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func init() {
	lynlabHeaders = make(map[string]string)
	lynlabHeaders["X-Api-Key"] = os.Getenv("LYNLAB_API_KEY")
}

func getKeyValue(key string) KeyValue {
	req, _ := http.NewRequest("GET", "https://api.lynlab.co.kr/v1/key_values/"+key, nil)
	for key, value := range lynlabHeaders {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	res, _ := client.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var keyValue KeyValue
	_ = json.Unmarshal(body, &keyValue)
	return keyValue
}

func postKeyValue(key string, value string) {
	dataJson, _ := json.Marshal(KeyValue{Key: key, Value: value})
	req, _ := http.NewRequest("POST", "https://api.lynlab.co.kr/v1/key_values", bytes.NewBuffer(dataJson))
	for key, value := range lynlabHeaders {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	res, _ := client.Do(req)
	defer res.Body.Close()
}

func GetFollowedUsers() []string {
	kv := getKeyValue("shuttle.follows")
	return strings.Split(kv.Value, ",")
}

func AddFollowedUser(loginId string) {
	followers := append(strings.Split(getKeyValue("shuttle.follows").Value, ","), loginId)
	postKeyValue("shuttle.follows", strings.Join(followers, ","))
}

func DeleteFollowedUser(loginId string) {
	followers := strings.Split(getKeyValue("shuttle.follows").Value, ",")
	for i, v := range followers {
		if v == loginId {
			followers = append(followers[:i], followers[i+1:]...)
			break
		}
	}

	postKeyValue("shuttle.follows", strings.Join(followers, ","))
}

func GetLastPlayedGameId(loginId string) string {
	kv := getKeyValue("shuttle.last_game_id." + loginId)
	return kv.Value
}

func UpdateLastPlayedGameId(loginId string, gameId string) {
	postKeyValue("shuttle.last_game_id."+loginId, gameId)
}
