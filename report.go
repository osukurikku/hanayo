package main

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func reportForm(c *gin.Context) {
	var (
		userID     int
		username   string
		privileges uint64
	)

	ctx := getContext(c)

	u := c.Param("user")
	if _, err := strconv.Atoi(u); err != nil {
		err := db.QueryRow("SELECT id, username, privileges FROM users WHERE username = ? AND "+ctx.OnlyUserPublic()+" LIMIT 1", u).Scan(&userID, &username, &privileges)
		if err != nil && err != sql.ErrNoRows {
			c.Error(err)
		}
	} else {
		err := db.QueryRow(`SELECT id, username, privileges FROM users WHERE id = ? AND `+ctx.OnlyUserPublic()+` LIMIT 1`, u).Scan(&userID, &username, &privileges)
		switch {
		case err == nil:
		case err == sql.ErrNoRows:
			err := db.QueryRow(`SELECT id, username, privileges FROM users WHERE username = ? AND `+ctx.OnlyUserPublic()+` LIMIT 1`, u).Scan(&userID, &username, &privileges)
			if err != nil && err != sql.ErrNoRows {
				c.Error(err)
			}
		default:
			c.Error(err)
		}
	}

	data := new(profileData)
	data.UserID = userID

	defer resp(c, 200, "report.html", data)

	if data.UserID == 0 {
		data.TitleBar = "User not found"
		data.Messages = append(data.Messages, warningMessage{T(c, "That user could not be found.")})
		return
	}

	data.TitleBar = T(c, "Report %s", username)
	data.DisableHH = true
}

func reportSubmit(c *gin.Context) {
	var (
		userID     int
		username   string
		privileges uint64
	)
	if getContext(c).User.ID == 0 {
		resp403(c)
		return
	}

	ctx := getContext(c)

	u := c.Param("user")
	if _, err := strconv.Atoi(u); err != nil {
		err := db.QueryRow("SELECT id, username, privileges FROM users WHERE username = ? AND "+ctx.OnlyUserPublic()+" LIMIT 1", u).Scan(&userID, &username, &privileges)
		if err != nil && err != sql.ErrNoRows {
			c.Error(err)
		}
	} else {
		err := db.QueryRow(`SELECT id, username, privileges FROM users WHERE id = ? AND `+ctx.OnlyUserPublic()+` LIMIT 1`, u).Scan(&userID, &username, &privileges)
		switch {
		case err == nil:
		case err == sql.ErrNoRows:
			err := db.QueryRow(`SELECT id, username, privileges FROM users WHERE username = ? AND `+ctx.OnlyUserPublic()+` LIMIT 1`, u).Scan(&userID, &username, &privileges)
			if err != nil && err != sql.ErrNoRows {
				c.Error(err)
			}
		default:
			c.Error(err)
		}
	}

	data := new(profileData)
	data.UserID = userID

	if data.UserID == 0 {
		addMessage(c, errorMessage("User not found"))
		getSession(c).save()
		c.Redirect(302, "/u/"+c.Param("user")+"/report")
		return
	}

	reason, err := strconv.Atoi(c.PostForm("reason"))
	if err != nil {
		c.Error(err)
	}
	currentTime := int(time.Now().Unix())
	reasonStr := ""
	switch reason {
	case 1:
		reasonStr = "Foul play / Cheating"
		break
	case 2:
		reasonStr = "Spamming"
		break
	case 3:
		reasonStr = "Other"
		break
	default:
		addMessage(c, errorMessage("Not correct reason"))
		c.Redirect(302, "/u/"+c.Param("user")+"/report")
		return
	}

	res, err := db.Exec(`INSERT INTO reports (id, from_uid, to_uid, chatlog, reason, time, assigned) VALUES (NULL, ?, ?, ?, ?, ?, '');`,
		getContext(c).User.ID, c.Param("user"), c.PostForm("addition"), reasonStr, currentTime)
	if err != nil {
		c.Error(err)
		return
	}

	addMessage(c, successMessage{T(c, "User has been reported.")})
	getSession(c).Save()
	c.Redirect(302, "/u/"+c.Param("user"))
}
