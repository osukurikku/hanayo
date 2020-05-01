package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

//baseOauthStruct - vk.go
//socialDBInfo - vk.go

func twitchAuthorizeUser(code string) (nickname string, err error) {
	twitchSettings := config.TwitchSettings
	twitchArray := strings.Split(twitchSettings, "|")
	if len(twitchArray) < 3 {
		return nickname, errors.New("twitch settings length < 3 should be client_id|client_secret|redirect_uri")
	}

	client_id, client_secret, redirect_uri := twitchArray[0], twitchArray[1], twitchArray[2]

	resp, err := http.Post("https://id.twitch.tv/oauth2/token?client_id="+client_id+"&grant_type=authorization_code&client_secret="+client_secret+"&redirect_uri="+redirect_uri+"&code="+code, "application/json", nil)
	if err != nil {
		return nickname, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nickname, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nickname, err
	}

	if result["access_token"] == nil {
		return nickname, errors.New("twitch oauth failed, bad code!")
	}
	token := result["access_token"].(string)

	client := http.Client{}
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		return nickname, err
	}
	req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Set("Client-ID", client_id)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	if err != nil {
		return nickname, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nickname, err
	}

	var result2 map[string]interface{}
	err = json.Unmarshal(body, &result2)
	if err != nil {
		return nickname, err
	}

	if result2["data"] == nil {
		return nickname, errors.New("twitch oauth failed, data array not found!")
	}
	data := result2["data"].([]interface{})

	nickname = data[0].(map[string]interface{})["login"].(string)

	return nickname, nil
}

func twitchAuthorizeHandler(c *gin.Context) {
	if getContext(c).User.ID == 0 {
		simpleReply(c, errorMessage{T(c, "Please log in before")})
		return
	}
	data := new(baseOAuthInfo)
	data.TitleBar = T(c, "OAuth Gate")
	data.KyutGrill = "2fa.jpg"
	defer resp(c, 200, "socialOAuth.html", data)

	twitchCode := c.Request.URL.Query().Get("code")
	if len(twitchCode) < 1 {
		data.Authorized = false
		data.RedirectBackPage = "/settings/socialAuth"
		data.Messages = append(data.Messages, errorMessage{T(c, "One of arguments is empty. Don't hack!")})
		return
	}

	twitchNickname, err := twitchAuthorizeUser(twitchCode)
	if err != nil {
		c.Error(err)

		data.Authorized = false
		data.RedirectBackPage = "/settings/socialAuth"
		data.Messages = append(data.Messages, errorMessage{T(c, "Something goes wrong, try again later!")})
		return
	}

	socialDBInfo := socialDBInfo{}
	if err := db.QueryRow(`
	SELECT
		user_id, account_type, account_id
	FROM
		social_networks
	WHERE
		account_id = '`+twitchNickname+`'
	AND
		account_type = 'twitch'
	`).Scan(&socialDBInfo.UserID, &socialDBInfo.AccountType, &socialDBInfo.AccountID); err != nil {
		db.Exec("INSERT INTO social_networks (user_id, account_type, account_id) VALUES (?, ?, ?)", getContext(c).User.ID, "twitch", twitchNickname)
	} else {
		if socialDBInfo.UserID != getContext(c).User.ID {
			data.Authorized = false
			data.RedirectBackPage = "/settings/socialAuth"
			data.Messages = append(data.Messages, errorMessage{T(c, "This account already connected to another player! ")})
			return
		}

		data.Authorized = true
		data.RedirectBackPage = "/settings/socialAuth"
		data.Messages = append(data.Messages, successMessage{T(c, "You connected this account earlier!")})
		return
	}

	data.Authorized = true
	data.RedirectBackPage = "/settings"
	data.Messages = append(data.Messages, successMessage{T(c, "You successfully connected twitch account to your profile")})
}

func twitchUnAuthorizeHandler(c *gin.Context) {
	if getContext(c).User.ID == 0 {
		simpleReply(c, errorMessage{T(c, "Please log in before")})
		return
	}
	data := new(baseOAuthInfo)
	data.TitleBar = T(c, "OAuth Gate")
	data.KyutGrill = "2fa.jpg"
	defer resp(c, 200, "socialOAuth.html", data)

	socialDBInfo := socialDBInfo{}
	if err := db.QueryRow(`
	SELECT
		user_id, account_type, account_id
	FROM
		social_networks
	WHERE
		user_id = `+strconv.Itoa(getContext(c).User.ID)+`
	AND
		account_type = 'twitch'
	`).Scan(&socialDBInfo.UserID, &socialDBInfo.AccountType, &socialDBInfo.AccountID); err != nil {
		data.Authorized = true
		data.RedirectBackPage = "/settings"
		data.Messages = append(data.Messages, successMessage{T(c, "Your profile hasn't linked account ;d Ok, BRUH!~~~~")})
		return
	}

	db.Exec("DELETE FROM social_networks WHERE user_id = ? AND account_type  = 'twitch'", getContext(c).User.ID)

	data.Authorized = true
	data.RedirectBackPage = "/settings/socialAuth"
	data.Messages = append(data.Messages, successMessage{T(c, "Your profile has successfully unlinked!")})
}
