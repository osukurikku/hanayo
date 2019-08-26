package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/osuripple/cheesegull/models"
)

type beatmapPageData struct {
	baseTemplateData

	Found      bool
	Beatmap    models.Beatmap
	Beatmapset models.Set
	SetJSON    string
}

func beatmapSetInfo(c *gin.Context) {
	data := new(beatmapPageData)

	sid := c.Param("sid")
	if _, err := strconv.Atoi(sid); err != nil {
		c.Error(err)
	} else {
		result := make(map[string]interface{})
		resp, err := http.Get(config.CheesegullAPI + "/s/" + sid)
		if err != nil {
			c.Error(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.Error(err)
		}

		err = json.Unmarshal(body, &result)
		if err != nil {
			c.Error(err)
		}

		if len(result["ChildrenBeatmaps"]) < 1 {
			if data.Beatmapset.ID == 0 {
				data.TitleBar = T(c, "Beatmap not found.")
				data.Messages = append(data.Messages, errorMessage{T(c, "Beatmap could not be found.")})
				resp(c, 200, "beatmap.html", data)
				return
			}
		}

		beatmapID := result["ChildrenBeatmaps"][0]["BeatmapID"]
		c.Redirect(302, "/b/"+strconv.Itoa(beatmapID))
	}
}

func beatmapInfo(c *gin.Context) {
	data := new(beatmapPageData)
	defer resp(c, 200, "beatmap.html", data)

	b := c.Param("bid")
	if _, err := strconv.Atoi(b); err != nil {
		c.Error(err)
	} else {
		data.Beatmap, err = getBeatmapData(b)
		if err != nil {
			c.Error(err)
			return
		}
		data.Beatmapset, err = getBeatmapSetData(data.Beatmap)
		if err != nil {
			c.Error(err)
			return
		}
		sort.Slice(data.Beatmapset.ChildrenBeatmaps, func(i, j int) bool {
			if data.Beatmapset.ChildrenBeatmaps[i].Mode != data.Beatmapset.ChildrenBeatmaps[j].Mode {
				return data.Beatmapset.ChildrenBeatmaps[i].Mode < data.Beatmapset.ChildrenBeatmaps[j].Mode
			}
			return data.Beatmapset.ChildrenBeatmaps[i].DifficultyRating < data.Beatmapset.ChildrenBeatmaps[j].DifficultyRating
		})
	}

	if data.Beatmapset.ID == 0 {
		data.TitleBar = T(c, "Beatmap not found.")
		data.Messages = append(data.Messages, errorMessage{T(c, "Beatmap could not be found.")})
		return
	}

	data.KyutGrill = fmt.Sprintf("https://assets.ppy.sh/beatmaps/%d/covers/cover.jpg?%d", data.Beatmapset.ID, data.Beatmapset.LastUpdate.Unix())
	data.KyutGrillAbsolute = true

	setJSON, err := json.Marshal(data.Beatmapset)
	if err == nil {
		data.SetJSON = string(setJSON)
	} else {
		data.SetJSON = "[]"
	}

	data.TitleBar = T(c, "%s - %s", data.Beatmapset.Artist, data.Beatmapset.Title)
	data.Scripts = append(data.Scripts, "/static/tablesort.js", "/static/beatmap.js")
}

func getBeatmapData(b string) (beatmap models.Beatmap, err error) {
	resp, err := http.Get(config.CheesegullAPI + "/b/" + b)
	if err != nil {
		return beatmap, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return beatmap, err
	}

	err = json.Unmarshal(body, &beatmap)
	if err != nil {
		return beatmap, err
	}

	return beatmap, nil
}

func getBeatmapSetData(beatmap models.Beatmap) (bset models.Set, err error) {
	resp, err := http.Get(config.CheesegullAPI + "/s/" + strconv.Itoa(beatmap.ParentSetID))
	if err != nil {
		return bset, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return bset, err
	}

	err = json.Unmarshal(body, &bset)
	if err != nil {
		return bset, err
	}

	return bset, nil
}
