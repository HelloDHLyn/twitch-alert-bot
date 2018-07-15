package twitch

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// https://dev.twitch.tv/docs/api/reference/#get-users
type Users struct {
	Data []struct {
		ID              string `json:"id"`
		Login           string `json:"login"`
		DisplayName     string `json:"display_name"`
		Type            string `json:"type"`
		BroadcasterType string `json:"broadcaster_type"`
		Description     string `json:"description"`
		ProfileImageURL string `json:"profile_image_url"`
		OfflineImageURL string `json:"offline_image_url"`
		ViewCount       int    `json:"view_count"`
		Email           string `json:"email"`
	} `json:"data"`
}

// https://dev.twitch.tv/docs/api/reference/#get-streams
type Streams struct {
	Data []struct {
		ID           string    `json:"id"`
		UserID       string    `json:"user_id"`
		GameID       string    `json:"game_id"`
		CommunityIds []string  `json:"community_ids"`
		Type         string    `json:"type"`
		Title        string    `json:"title"`
		ViewerCount  int       `json:"viewer_count"`
		StartedAt    time.Time `json:"started_at"`
		Language     string    `json:"language"`
		ThumbnailURL string    `json:"thumbnail_url"`
	} `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}

// https://dev.twitch.tv/docs/api/reference/#get-games
type Games struct {
	Data []struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		BoxArtURL string `json:"box_art_url"`
	} `json:"data"`
}

func getAndMarshal(url string, v interface{}) {
	headers := make(map[string]string)
	headers["Client-ID"] = os.Getenv("TWITCH_CLIENT_ID")

	req, _ := http.NewRequest("GET", url, nil)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	res, _ := client.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(body, &v)
}

func GetStreams(userLoginIds []string) Streams {
	url := "https://api.twitch.tv/helix/streams?"
	for idx, id := range userLoginIds {
		if idx > 0 {
			url = url + "&"
		}
		url = url + "user_login=" + id
	}

	var streams Streams
	getAndMarshal(url, &streams)
	return streams
}

func GetUsers(loginIds []string) Users {
	url := "https://api.twitch.tv/helix/users?"
	for idx, id := range loginIds {
		if idx > 0 {
			url = url + "&"
		}
		url = url + "login=" + id
	}

	var users Users
	getAndMarshal(url, &users)
	return users
}

func GetGames(ids []string) Games {
	url := "https://api.twitch.tv/helix/games?"
	for idx, id := range ids {
		if idx > 0 {
			url = url + "&"
		}
		url = url + "id=" + id + "&"
	}

	var games Games
	getAndMarshal(url, &games)
	return games
}
