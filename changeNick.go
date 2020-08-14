package main

import (
	"regexp"
	"strings"
	"time"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"zxq.co/ripple/rippleapi/common"
)

func submitNick(c *gin.Context) {
	data := new(baseTemplateData)
	defer resp(c, 200, "settings/change_username.html", data)
	ctx := getContext(c)
	data.TitleBar = "Change Username"

	if ctx.User.ID == 0 {
		resp403(c)
		return
	}

	lastNickChange := 0
	db.Get(&lastNickChange, "SELECT latest_nickchange WHERE id = ? LIMIT 1", ctx.User.ID)
	if lastNickChange > 0 && ((ctx.User.Privileges & common.UserPrivilegeDonor) == 0) {
		// User used one free change by himself
		data.Messages = append(data.Messages, errorMessage{T(c, "You used your free change! Now you can change nickname only with Donor privileges!")})
		return
	}
	// 259200 secs == 3 days
	if lastNickChange > 0 && ((ctx.User.Privileges & common.UserPrivilegeDonor) > 0) && lastNickChange > int(time.Now().Unix()) {
		data.Messages = append(data.Messages, errorMessage{T(c, "You can change your nickname only after 3 days. Plz wait!")})
		return
	}

	newu := c.PostForm("newu")
	if len(newu) < 2 || len(newu) > 15 {
		data.Messages = append(data.Messages, errorMessage{T(c, "Must be at least 2 characters long and mustn't be longer than 15 characters")})
		return
	}
	// No username with mixed spaces
	if strings.Index(newu, " ") != -1 && strings.Index(newu, "_") != -1 {
		data.Messages = append(data.Messages, errorMessage{T(c, "Mustn't contain mixed underscores and spaces")})
		return
	}
	match, _ := regexp.MatchString("^[A-Za-z0-9 _\\[\\]-]+$", newu)
	if match != true {
		data.Messages = append(data.Messages, errorMessage{T(c, "Must contain only letter, symbols, spaces, underscores, square brackets or dashes")})
		return
	}

	// Check if username is already in db
	safe := strings.ReplaceAll(strings.ToLower(newu), " ", "_")
	userDB := ""
	db.Get(&userDB, "SELECT username_safe FROM users WHERE LOWER(username_safe) = LOWER(?) AND id != ? LIMIT 1", safe, ctx.User.ID)
	if len(userDB) > 0 {
		data.Messages = append(data.Messages, errorMessage{T(c, "Username already used by another user. No changes have been made.")})
		return
	}

	jsonStr, _ := json.Marshal(struct {
		UserID      int    `json:"userID"`
		NewUsername string `json:"newUsername"`
	}{
		ctx.User.ID, newu,
	})

	rd.Publish("peppy:change_username", string(jsonStr))
	db.Exec("UPDATE users SET latest_nickchange = ? WHERE id = ?", int(time.Now().Unix())+259200, ctx.User.ID)

	data.Messages = append(data.Messages, successMessage{T(c, "Your nickname successfully changed!")})
	return
}
