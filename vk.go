package main

import (
	"strings"
	"errors"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type baseOAuthInfo struct {
	baseTemplateData

	Authorized      	bool
	RedirectBackPage	string
}

type vkAuthStruct struct {
	AccessToken		string `json:"access_token"`
	ExpiresIn		int `json:"expires_in"`
	UserId			int `json:"user_id"`
}

type typicalSocialConfig struct {
	ClientID		string
	ClientSecret	string
	RedirectURI		string
}

type socialDBInfo struct {
	UserID			int
	AccountType		string
	AccountID		string
}

func vKAuthorizeUser(code string) (vkData vkAuthStruct, err error) {
	vkSettings := config.VKSettings
	vkArray := strings.Split(vkSettings, "|")
	if len(vkArray) < 3 {
		return vkData, errors.New("vk settings length < 3 should be client_id|client_secret|redirect_uri")
	}

	client_id, client_secret, redirect_uri := vkArray[0], vkArray[1], vkArray[2]

	resp, err := http.Get("https://oauth.vk.com/access_token?client_id="+client_id+"&client_secret="+client_secret+"&redirect_uri="+redirect_uri+"&code="+code)
	if err != nil {
		return vkData, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return vkData, err
	}

	err = json.Unmarshal(body, &vkData)
	if err != nil {
		return vkData, errors.New("vk error ;d")
	}

	if vkData.UserId < 1 {
		return vkData, errors.New("vk error ;d")
	}

	return vkData, nil
}

func vKAuthorizeHandler(c *gin.Context) {
	if getContext(c).User.ID == 0 {
		simpleReply(c, errorMessage{T(c, "Please log in before")})
		return
	}
	data := new(baseOAuthInfo)
	data.TitleBar = T(c, "OAuth Gate")
	data.KyutGrill = "2fa.jpg"
	defer resp(c, 200, "socialOAuth.html", data)

	vkCode := c.Request.URL.Query().Get("code")
	if len(vkCode) < 1 {
		data.Authorized = false
		data.RedirectBackPage = "/settings/socialAuth"
		data.Messages = append(data.Messages, errorMessage{T(c, "One of arguments is empty. Don't hack!")})
		return
	}

	vkData, err := vKAuthorizeUser(vkCode)
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
		account_id = '`+strconv.Itoa(vkData.UserId)+`'
	AND
		account_type = 'vk'
	`).Scan(&socialDBInfo.UserID, &socialDBInfo.AccountType, &socialDBInfo.AccountID); err != nil {
		db.Exec("INSERT INTO social_networks (user_id, account_type, account_id) VALUES (?, ?, ?)", getContext(c).User.ID, "vk", vkData.UserId)
	} else {
		if socialDBInfo.UserID != getContext(c).User.ID {
			data.Authorized = false
			data.RedirectBackPage = "/settings/socialAuth"
			data.Messages = append(data.Messages,  errorMessage{T(c, "This account already connected to another player! ")})
			return
		}

		data.Authorized = true
		data.RedirectBackPage = "/settings/socialAuth"
		data.Messages = append(data.Messages,  successMessage{T(c, "You connected this account earlier!")})
		return
	}

	data.Authorized = true
	data.RedirectBackPage = "/settings"
	data.Messages = append(data.Messages,  successMessage{T(c, "You successfully connected VK account to your profile")})
}

func vKUnAuthorizeHandler(c *gin.Context) {
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
		account_type = 'vk'
	`).Scan(&socialDBInfo.UserID, &socialDBInfo.AccountType, &socialDBInfo.AccountID); err != nil {
		data.Authorized = true
		data.RedirectBackPage = "/settings"
		data.Messages = append(data.Messages,  successMessage{T(c, "Your profile hasn't linked account ;d Ok, BRUH!~~~~")})
		return
	}
	
	db.Exec("DELETE FROM social_networks WHERE user_id = ? AND account_type  = 'vk'", getContext(c).User.ID)

	data.Authorized = true
	data.RedirectBackPage = "/settings/socialAuth"
	data.Messages = append(data.Messages,  successMessage{T(c, "Your profile has successfully unlinked!")})
}