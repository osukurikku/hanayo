package main

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type TopgVoter struct {
	UserID int `db:"userid"`
	Count  int `db:"count"`
	Time   int `db:"time"`
}

type topgVotePageStruct struct {
	baseTemplateData

	CanUserSubmit  bool
	NextTimeSubmit int64
	Balance        int
}

func SubmiterVotes(c *gin.Context) {
	addr, err := net.LookupIP("monitor.topg.org")
	if err != nil {
		resp403(c)
		return
	}

	remoteAddr := strings.Split(c.Request.Header.Get("X-FORWARDED-FOR"), ",")
	if remoteAddr[0] != addr[0].String() {
		fmt.Println(remoteAddr)
		fmt.Println(addr[0].String())
		fmt.Println("ip not correct?")
		resp403(c)
		return
	}

	var reResp = regexp.MustCompile(`/[^A-Za-z0-9\_\-]+/`)
	PResp := reResp.ReplaceAllString(c.Query("p_resp"), ``)

	userID, err := strconv.Atoi(PResp)
	if err != nil {
		resp403(c)
		return
	}

	data := new(topgVotePageStruct)
	data.KyutGrill = "donate.jpg"
	data.TitleBar = T(c, "Voting")
	defer resp(c, 200, "topgvoter.html", data)

	currentTime := time.Now().Unix()
	voter := TopgVoter{}
	err = db.Get(&voter, "SELECT userid, count, time FROM topg_votes WHERE userid = ?", userID)
	if err != nil {
		db.Exec("INSERT INTO topg_votes (userid, count, time) VALUES (?, ?, ?)", userID, 0, int(currentTime))
		err = db.Get(&voter, "SELECT userid, count, time FROM topg_votes WHERE userid = ?", userID)
		if err != nil {
			data.Messages = append(data.Messages, errorMessage{T(c, "Sorry, but we can't register you(")})
			return
		}
	}

	if voter.Time > int(currentTime) {
		data.Messages = append(data.Messages, errorMessage{T(c, "Not now, try later!")})
		data.NextTimeSubmit = int64(voter.Time)
	} else {
		voter.Time = int(currentTime + 172800)
		voter.Count++
		db.Exec("UPDATE topg_votes SET time = ?, count = ? WHERE userid = ?", voter.Time, voter.Count, voter.UserID)
	}

	// if getContext(c).User.ID == 0 {

	// }
	// data := new(topgVotePageStruct)
	// defer resp(c, 200, "topgvoter.html", data)

	// currentTime := time.Now().Unix()
	// isNewbie := false
	// voter := TopgVoter{}
	// err := db.Get(&voter, "SELECT userid, count, time FROM topg_votes WHERE id = ?", getContext(c).User.ID)
	// if err != nil {
	// 	db.Exec("INSERT INTO topg_votes (userid, count, time) VALUES (?, ?, ?)", getContext(c).User.ID, 0, int(currentTime))
	// 	err = db.Get(&voter, "SELECT userid, count, time FROM topg_votes WHERE id = ?", getContext(c).User.ID)
	// 	if err != nil {
	// 		data.Messages = append(data.Messages, errorMessage{T(c, "Sorry, but we can't register you(")})
	// 		return
	// 	}
	// 	isNewbie = true
	// }

	// data.Balance = voter.Count
	// data.CanUserSubmit = false

	// // 604800 - 1 week!
	// if isNewbie {
	// 	// submit him vote

	// 	// TODO: Topg Publisher!
	// 	voter.Time += 604800
	// 	voter.Count++
	// 	db.Exec("UPDATE topg_votes SET time = ?, count = ? WHERE userid = ?", voter.Time, voter.Count, voter.UserID)
	// 	data.NextTimeSubmit = int64(voter.Time)
	// 	data.Messages = append(data.Messages, successMessage{T(c, "You successfully voted for us. You have +1 vote")})
	// 	data.Balance = voter.Count
	// } else {
	// 	if voter.Time > int(currentTime) {
	// 		data.Messages = append(data.Messages, errorMessage{T(c, "Not now, try later!")})

	// 		data.NextTimeSubmit = int64(voter.Time)
	// 	} else {
	// 		voter.Time = int(currentTime + 604800)
	// 		voter.Count++
	// 		db.Exec("UPDATE topg_votes SET time = ?, count = ? WHERE userid = ?", voter.Time, voter.Count, voter.UserID)
	// 		data.NextTimeSubmit = int64(voter.Time)
	// 		data.Messages = append(data.Messages, successMessage{T(c, "You successfully voted for us. You have +1 vote")})
	// 		data.Balance = voter.Count
	// 	}
	// }
	return
}

func GetVoter(c *gin.Context) {
	if getContext(c).User.ID == 0 {
		resp403(c)
		return
	}
	data := new(topgVotePageStruct)
	data.KyutGrill = "donate.jpg"
	data.TitleBar = T(c, "Voting")
	defer resp(c, 200, "topgvoter.html", data)

	voter := TopgVoter{}
	err := db.Get(&voter, "SELECT userid, count, time FROM topg_votes WHERE userid = ?", getContext(c).User.ID)
	if err != nil {
		voter.Count = 0
		voter.Time = 0
		data.CanUserSubmit = true
		data.NextTimeSubmit = 0
		return
	}

	data.Balance = voter.Count
	data.CanUserSubmit = false
	currentTime := time.Now().Unix()
	if voter.Time > int(currentTime) {
		data.CanUserSubmit = false
		data.NextTimeSubmit = int64(voter.Time)
	} else {
		data.CanUserSubmit = true
	}
	return
}

func exchangeVotes(c *gin.Context) {
	if getContext(c).User.ID == 0 {
		resp403(c)
		return
	}

	voter := TopgVoter{}
	err := db.Get(&voter, "SELECT userid, count, time FROM topg_votes WHERE userid = ?", getContext(c).User.ID)
	if err != nil {
		respEmpty(c, "Forbidden", warningMessage{T(c, "You haven't enough votes. Min votes to exchange = 15. 15 votes = 150 rub")})
		return
	}

	if voter.Count < 25 {
		respEmpty(c, "Forbidden", warningMessage{T(c, "You haven't enough votes. Min votes to exchange = 15. 15 votes = 150 rub")})
		return
	}
	voter.Count -= 25
	db.Exec("UPDATE topg_votes SET count = ? WHERE userid = ?", voter.Count, voter.UserID)
	db.Exec(`UPDATE users SET balance = balance+150 WHERE id = ?`, voter.UserID)

	c.Redirect(302, "/vote")
	return
}
