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

	if data.UserID == getContext(c).User.ID {
		addMessage(c, errorMessage{T(c, "You can't report yourself")})
		getSession(c).Save()
		c.Redirect(302, "/u/"+c.Param("user"))
		return
	}

	if data.UserID == 999 {
		addMessage(c, errorMessage{T(c, "You can't report our bot ;3")})
		getSession(c).Save()
		c.Redirect(302, "/u/"+c.Param("user"))
		return
	}

	if data.UserID == 0 {
		addMessage(c, errorMessage{T(c, "User not found")})
		getSession(c).Save()
		c.Redirect(302, "/")
		return
	}

	defer resp(c, 200, "report.html", data)

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

	if data.UserID == getContext(c).User.ID {
		addMessage(c, errorMessage{T(c, "You can't report yourself")})
		getSession(c).Save()
		c.Redirect(302, "/u/"+c.Param("user"))
		return
	}

	if data.UserID == 999 {
		addMessage(c, errorMessage{T(c, "You can't report our bot ;3")})
		getSession(c).Save()
		c.Redirect(302, "/u/"+c.Param("user"))
		return
	}

	if data.UserID == 0 {
		addMessage(c, errorMessage{T(c, "User not found")})
		getSession(c).Save()
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
		addMessage(c, errorMessage{T(c, "Not correct reason")})
		c.Redirect(302, "/u/"+c.Param("user")+"/report")
		return
	}

	_, err1 := db.Exec(`INSERT INTO reports (id, from_uid, to_uid, chatlog, reason, time, assigned) VALUES (NULL, ?, ?, ?, ?, ?, '');`,
		getContext(c).User.ID, c.Param("user"), c.PostForm("addition"), reasonStr, currentTime)
	if err1 != nil {
		c.Error(err)
		return
	}

	addMessage(c, successMessage{T(c, "User has been reported.")})
	getSession(c).Save()
	c.Redirect(302, "/u/"+c.Param("user"))
}
