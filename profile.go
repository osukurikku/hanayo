package main

import (
	"database/sql"
	"strconv"

	"fmt"

	"github.com/gin-gonic/gin"
	"zxq.co/ripple/rippleapi/common"
)

// TODO: replace with simple ResponseInfo containing userid
type profileData struct {
	baseTemplateData
	UserID    int
	UserStyle string

	ProfileGrill         string
	ProfileGrillAbsolute bool
	ProfileColour        string
}

func userProfile(c *gin.Context) {
	var (
		userID     int
		username   string
		userStyle  string
		privileges uint64
		country    string
	)

	var (
		favMode  int
		pp_std   int
		pp_taiko int
		pp_ctb   int
		pp_mania int
	)

	ctx := getContext(c)

	u := c.Param("user")
	if _, err := strconv.Atoi(u); err != nil {
		err := db.QueryRow("SELECT  users.id, users.username, users_stats.user_style, users.privileges, users_stats.country, users_stats.favourite_mode, users_stats.pp_std, users_stats.pp_taiko, users_stats.pp_ctb, users_stats.pp_mania FROM users INNER JOIN users_stats ON users_stats.id = users.id WHERE username = ? AND "+ctx.OnlyUserPublic()+" LIMIT 1", u).Scan(&userID, &username, &userStyle, &privileges, &country, &favMode, &pp_std, &pp_taiko, &pp_ctb, &pp_mania)
		if err != nil && err != sql.ErrNoRows {
			c.Error(err)
		}
	} else {
		err := db.QueryRow(`SELECT users.id, users.username, users_stats.user_style, users.privileges, users_stats.country, users_stats.favourite_mode, users_stats.pp_std, users_stats.pp_taiko, users_stats.pp_ctb, users_stats.pp_mania FROM users INNER JOIN users_stats ON users_stats.id = users.id WHERE users.id = ? AND `+ctx.OnlyUserPublic()+` LIMIT 1`, u).Scan(&userID, &username, &userStyle, &privileges, &country, &favMode, &pp_std, &pp_taiko, &pp_ctb, &pp_mania)
		switch {
		case err == nil:
		case err == sql.ErrNoRows:
			err := db.QueryRow(`SELECT users.id, users.username, users_stats.user_style, users.privileges, users_stats.country, users_stats.favourite_mode, users_stats.pp_std, users_stats.pp_taiko, users_stats.pp_ctb, users_stats.pp_mania FROM users INNER JOIN users_stats ON users_stats.id = users.id WHERE username = ? AND `+ctx.OnlyUserPublic()+` LIMIT 1`, u).Scan(&userID, &username, &userStyle, &privileges, &country, &favMode, &pp_std, &pp_taiko, &pp_ctb, &pp_mania)
			if err != nil && err != sql.ErrNoRows {
				c.Error(err)
			}
		default:
			c.Error(err)
		}
	}

	data := new(profileData)
	if len(userStyle) > 0 {
		data.UserStyle = userStyle
	}

	if ctx.User.Privileges&common.AdminPrivilegeManageUsers == common.AdminPrivilegeManageUsers || ctx.User.ID == userID {
		data.UserID = userID
	} else {
		if privileges&2 == 0 {
			data.UserID = 0
		} else {
			data.UserID = userID
		}
	}

	defer resp(c, 200, "profile.html", data)

	if data.UserID == 0 {
		data.TitleBar = "User not found"
		data.Messages = append(data.Messages, warningMessage{T(c, "That user could not be found.")})
		return
	}

	if common.UserPrivileges(privileges)&common.UserPrivilegeDonor > 0 {
		var profileBackground struct {
			Type  int
			Value string
		}
		db.Get(&profileBackground, "SELECT type, value FROM profile_backgrounds WHERE uid = ?", data.UserID)
		switch profileBackground.Type {
		case 1:
			data.ProfileGrill = "/static/profbackgrounds/" + profileBackground.Value
			data.ProfileGrillAbsolute = true
		case 2:
			data.ProfileColour = profileBackground.Value
		}
	}

	modesPp := [4]int{pp_std, pp_taiko, pp_ctb, pp_mania}
	modesNames := [4]string{
		"osu!standart",
		"osu!taiko",
		"osu!ctb",
		"osu!mania",
	}

	data.MetaDescription = fmt.Sprintf("%s is %s player from %s with %dpp", username, modesNames[favMode], countryReadable(country), modesPp[favMode])
	data.MetaImage = fmt.Sprintf("https://a.kurikku.pw/%d", userID)

	data.TitleBar = T(c, "%s's profile", username)
	data.DisableHH = true
	data.Scripts = append(data.Scripts, "/static/profile.js")
}
