package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/russross/blackfriday"
	"github.com/thehowl/qsql"
	"golang.org/x/oauth2"
	discordoauth "zxq.co/ripple/go-discord-oauth"
	"zxq.co/ripple/hanayo/modules/bbcode"
	"zxq.co/ripple/hanayo/modules/btcaddress"
	"zxq.co/ripple/hanayo/modules/btcconversions"
	"zxq.co/ripple/hanayo/modules/doc"
	fasuimappings "zxq.co/ripple/hanayo/modules/fa-semantic-mappings"
	"zxq.co/ripple/playstyle"
	"zxq.co/ripple/rippleapi/common"

	// vkapi "github.com/SevereCloud/vksdk/api"
)

type RankRequest struct {
	Bid             int
	MapType         string
	ModeratorID     int
	ModeratorReason string
	ModeratorName   string
	Status          int
}

type FancyGraphSimple struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type BetterNews struct {
	Name    string
	News    string
	Preview string
	Date    string
	Link    string
}

// funcMap contains useful functions for the various templates.
var funcMap = template.FuncMap{
	// html disables HTML escaping on the values it is given.
	"html": func(value interface{}) template.HTML {
		return template.HTML(fmt.Sprint(value))
	},
	// safe html attr CSS
	"attrCSS": func(s string) template.CSS {
		return template.CSS(s)
	},
	// avatars is a function returning the configuration constant AvatarURL
	"config": func(key string) interface{} {
		return configMap[key]
	},
	// navbarItem is a function to generate an item in the navbar.
	// The reason why this exists is that I wanted to have the currently
	// selected element in the navbar having the "active" class.
	"navbarItem": func(currentPath, name, path string) template.HTML {
		var act string
		if path == currentPath {
			act = "active "
		}
		return template.HTML(fmt.Sprintf(`<a class="%sitem" href="%s">%s</a>`, act, path, name))
	},
	// curryear returns the current year.
	"curryear": func() int {
		return time.Now().Year()
	},
	// hasAdmin returns, based on the user's privileges, whether they should be
	// able to see the RAP button (aka AdminPrivilegeAccessRAP).
	"hasAdmin": func(privs common.UserPrivileges) bool {
		return privs&common.AdminPrivilegeAccessRAP > 0
	},
	// isRAP returns whether the current page is in RAP.
	"isRAP": func(p string) bool {
		parts := strings.Split(p, "/")
		return len(parts) > 1 && parts[1] == "admin"
	},
	// favMode is just a helper function for user profiles. Basically checks
	// whether a float and an int are ==, and if they are it will return "active ",
	// so that the element in the mode menu of a user profile can be marked as
	// the current active element.
	"favMode": func(favMode float64, current int) string {
		if int(favMode) == current {
			return "active "
		}
		return ""
	},
	// slice generates a []interface{} with the elements it is given.
	// useful to iterate over some elements, like this:
	//  {{ range slice 1 2 3 }}{{ . }}{{ end }}
	"slice": func(els ...interface{}) []interface{} {
		return els
	},
	// int converts a float/int to an int.
	"int": func(f interface{}) int {
		if f == nil {
			return 0
		}

		switch f := f.(type) {
		case int:
			return f
		case float64:
			return int(f)
		case float32:
			return int(f)
		case string:
			a, err := strconv.Atoi(f)
			if err != nil {
				return 0
			}
			return a
		case qsql.String:
			a, err := strconv.Atoi(f.String())
			if err != nil {
				return 0
			}
			return a
		default:
			fmt.Println(reflect.TypeOf(f))
			fmt.Println("unknown type")
		}
		return 0
	},
	// float converts an int to a float.
	"float": func(i int) float64 {
		return float64(i)
	},
	// atoi converts a string to an int and then a float64.
	// If s is not an actual int, it returns nil.
	"atoi": func(s string) interface{} {
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil
		}
		return float64(i)
	},
	// atoint is like atoi but returns always an int.
	"atoint": func(s string) int {
		i, _ := strconv.Atoi(s)
		return i
	},

	"isBday": func(i int) bool {
		_, month, day := time.Now().Date()
		if month == 10 && day == 19 {
			return true
		} else {
			return false
		}
	},

	"isSocialExists": func(social string, uID int) bool {
		var sId string
		db.Get(&sId, "SELECT account_id FROM social_networks WHERE user_id = ? AND account_type = ?", uID, social)
		if len(sId) < 1 {
			return false
		}
		return true
	},

	"socialConfig": func(socialConfig string) (tSC typicalSocialConfig) {
		r := reflect.ValueOf(&config)
		f := reflect.Indirect(r).FieldByName(socialConfig)
		fString := string(f.String())

		if len(fString) < 1 {
			return tSC
		}
		splitted := strings.Split(fString, "|")
		tSC.ClientID = splitted[0]
		tSC.ClientSecret = splitted[1]
		tSC.RedirectURI = splitted[2]
		return tSC
	},

	"socialAccounts": func(userID int) (sDI []socialDBInfo) {
		rows, err := db.Query(`
		SELECT
			user_id, account_type, account_id
		FROM
			social_networks
		WHERE
			user_id = ?
		`, userID)
		if err != nil {
			return sDI
		}

		for rows.Next() {
			var sAcc socialDBInfo
			err := rows.Scan(
				&sAcc.UserID,
				&sAcc.AccountType,
				&sAcc.AccountID,
			)
			if err != nil {
				break
			}
			sDI = append(sDI, sAcc)
		}

		sort.Slice(sDI, func(i, j int) bool {
			return sDI[i].AccountType < sDI[j].AccountType
		})

		return sDI
	},

	// parseUserpage compiles BBCode to HTML.
	"parseUserpage": func(s string) template.HTML {
		return template.HTML(bbcode.Compile(s))
	},
	// time converts a RFC3339 timestamp to the HTML element <time>.
	"time": func(s string) template.HTML {
		t, _ := time.Parse(time.RFC3339, s)
		return _time(s, t)
	},
	// time generates a time from a native Go time.Time
	"timeFromTime": func(t time.Time) template.HTML {
		return _time(t.Format(time.RFC3339), t)
	},
	// timeAddDay is basically time but adds a day.
	"timeAddDay": func(s string) template.HTML {
		t, _ := time.Parse(time.RFC3339, s)
		t = t.Add(time.Hour * 24)
		return _time(t.Format(time.RFC3339), t)
	},
	// nativeTime creates a native Go time.Time from a RFC3339 timestamp.
	"nativeTime": func(s string) time.Time {
		t, _ := time.Parse(time.RFC3339, s)
		return t
	},
	// band is a bitwise AND.
	"band": func(i1 int, i ...int) int {
		for _, el := range i {
			i1 &= el
		}
		return i1
	},
	// countryReadable converts a country's ISO name to its full name.
	"countryReadable": countryReadable,
	"country": func(s string, name bool) template.HTML {
		var c string
		if name {
			c = countryReadable(s)
			if c == "" {
				return ""
			}
		}
		return template.HTML(fmt.Sprintf(`<img src="https://s.kurikku.pw/flags/%s.png" class="ui kr-flag" alt="img flag"></img>%s`, strings.ToUpper(s), c))
	},
	"ToUpper": strings.ToUpper,
	// humanize pretty-prints a float, e.g.
	//     humanize(1000) == "1,000"
	"humanize": func(f float64) string {
		return humanize.Commaf(f)
	},
	// levelPercent basically does this:
	//     levelPercent(56.23215) == "23"
	"levelPercent": func(l float64) string {
		_, f := math.Modf(l)
		f *= 100
		return fmt.Sprintf("%.0f", f)
	},
	// level removes the decimal part from a float.
	"level": func(l float64) string {
		i, _ := math.Modf(l)
		return fmt.Sprintf("%.0f", i)
	},
	// faIcon converts a fontawesome icon to a semantic ui icon.
	"faIcon": func(i string) string {
		classes := strings.Split(i, " ")
		for i, class := range classes {
			if v, ok := fasuimappings.Mappings[class]; ok {
				classes[i] = v
			}
		}
		return strings.Join(classes, " ")
	},
	// log fmt.Printf's something
	"log": fmt.Printf,
	// has returns whether priv1 has all 1 bits of priv2, aka priv1 & priv2 == priv2
	"has": func(priv1 interface{}, priv2 float64) bool {
		var p1 uint64
		switch priv1 := priv1.(type) {
		case common.UserPrivileges:
			p1 = uint64(priv1)
		case float64:
			p1 = uint64(priv1)
		case int:
			p1 = uint64(priv1)
		}
		return p1&uint64(priv2) == uint64(priv2)
	},
	// _range is like python range's.
	// If it is given 1 argument, it returns a []int containing numbers from 0
	// to x.
	// If it is given 2 arguments, it returns a []int containing numers from x
	// to y if x < y, from y to x if y < x.
	"_range": func(x int, y ...int) ([]int, error) {
		switch len(y) {
		case 0:
			r := make([]int, x)
			for i := range r {
				r[i] = i
			}
			return r, nil
		case 1:
			nums, up := pos(y[0] - x)
			r := make([]int, nums)
			for i := range r {
				if up {
					r[i] = i + x + 1
				} else {
					r[i] = i + y[0]
				}
			}
			if !up {
				// reverse r
				sort.Sort(sort.Reverse(sort.IntSlice(r)))
			}
			return r, nil
		}
		return nil, errors.New("y must be at maximum 1 parameter")
	},
	// blackfriday passes some markdown through blackfriday.
	"blackfriday": func(m string) template.HTML {
		// The reason of m[strings.Index...] is to remove the "header", where
		// there is the information about the file (namely, title, old_id and
		// reference_version)
		return template.HTML(
			blackfriday.Run(
				[]byte(
					m[strings.Index(m, "\n---\n")+5:],
				),
			),
		)
	},
	// i is an inline if.
	// i (cond) (true) (false)
	"i": func(a bool, x, y interface{}) interface{} {
		if a {
			return x
		}
		return y
	},
	// modes returns an array containing all the modes (in their string representation).
	"modes": func() []string {
		return []string{
			"osu! standard",
			"Taiko",
			"Catch the Beat",
			"osu!mania",
		}
	},
	// _or is like or, but has only false and nil as its "falsey" values
	"_or": func(args ...interface{}) interface{} {
		for _, a := range args {
			if a != nil && a != false {
				return a
			}
		}
		return nil
	},
	// unixNano returns the UNIX timestamp of when hanayo was started in nanoseconds.
	"unixNano": func() string {
		return strconv.FormatInt(hanayoStarted, 10)
	},
	"unixNanoInt": func() int {
		return int(time.Now().Unix())
	},
	// playstyle returns the string representation of a playstyle.
	"playstyle": func(i float64, f *profileData) string {
		var parts []string

		p := int(i)
		for k, v := range playstyle.Styles {
			if p&(1<<uint(k)) > 0 {
				parts = append(parts, f.T(v))
			}
		}

		return strings.Join(parts, ", ")
	},
	// arithmetic plus/minus
	"plus": func(i ...float64) float64 {
		var sum float64
		for _, i := range i {
			sum += i
		}
		return sum
	},
	"minus": func(i1 float64, i ...float64) float64 {
		for _, i := range i {
			i1 -= i
		}
		return i1
	},
	// rsin - Return Slice If Nil
	"rsin": func(i interface{}) interface{} {
		if i == nil {
			return []struct{}{}
		}
		return i
	},
	// loadjson loads a json file.
	"loadjson": func(jsonfile string) interface{} {
		f, err := ioutil.ReadFile(jsonfile)
		if err != nil {
			return nil
		}
		var x interface{}
		err = json.Unmarshal(f, &x)
		if err != nil {
			return nil
		}
		return x
	},
	// loadChangelog loads the changelog.
	"loadChangelog": loadChangelog,
	// teamJSON returns the data of team.json
	"teamJSON": func() map[string]interface{} {
		f, err := ioutil.ReadFile("team.json")
		if err != nil {
			return nil
		}
		var m map[string]interface{}
		json.Unmarshal(f, &m)
		return m
	},
	// in returns whether the first argument is in one of the following
	"in": func(a1 interface{}, as ...interface{}) bool {
		for _, a := range as {
			if a == a1 {
				return true
			}
		}
		return false
	},
	"capitalise": strings.Title,
	// servicePrefix gets the prefix of a service, like github.
	"servicePrefix": func(s string) string { return servicePrefixes[s] },
	// randomLogoColour picks a "random" colour for ripple's logo.
	// Not even used anymore
	"randomLogoColour": func() string {
		if rand.Int()%4 == 0 {
			return logoColours[rand.Int()%len(logoColours)]
		}
		return "pink"
	},
	"randInt": func(min int, max int) int {
		return rand.Intn(max-min) + min
	},
	// after checks whether a certain time is after time.Now()
	"after": func(s string) bool {
		t, _ := time.Parse(time.RFC3339, s)
		return t.After(time.Now())
	},

	// qsql functions
	"qb": func(q string, p ...interface{}) map[string]qsql.String {
		r, err := qb.QueryRow(q, p...)
		if err != nil {
			fmt.Println(err)
		}
		if r == nil {
			return make(map[string]qsql.String, 0)
		}
		return r
	},
	"qba": func(q string, p ...interface{}) []map[string]qsql.String {
		r, err := qb.Query(q, p...)
		if err != nil {
			fmt.Println(err)
		}
		return r
	},
	"qbe": func(q string, p ...interface{}) int {
		i, _, err := qb.Exec(q, p...)
		if err != nil {
			fmt.Println(err)
		}
		return i
	},

	"getRankRequests": func(userid int) (requests []RankRequest) {
		rows, err := db.Query("SELECT rank_requests.bid, rank_requests.type, rank_requests.moderator_id, rank_requests.moderator_reason, rank_requests.status FROM rank_requests WHERE rank_requests.userid = ? ORDER BY rank_requests.id DESC", userid)
		if err != nil {
			return requests
		}
		for rows.Next() {
			var r RankRequest
			err = rows.Scan(
				&r.Bid, &r.MapType, &r.ModeratorID, &r.ModeratorReason, &r.Status,
			)

			if err != nil {
				fmt.Println(err)
				return requests
			}
			if r.ModeratorID == 0 {
				r.ModeratorID = 0
				r.ModeratorName = "???"
				r.ModeratorReason = "???"
			} else {
				db.QueryRow("SELECT username FROM users WHERE id = ?", r.ModeratorID).Scan(&r.ModeratorName)
			}
			requests = append(requests, r)
		}

		return requests
	},

	"getImageBadges": func(userid int) (images []string) {
		rows, err := db.Query("SELECT image_url FROM user_image_badges WHERE user_id = ? ORDER BY id DESC", userid)
		if err != nil {
			return images
		}
		for rows.Next() {
			var s string
			err = rows.Scan(&s)

			if err != nil {
				fmt.Println(err)
				return images
			}
			images = append(images, s)
		}
		return images
	},

	"getFancyUsersStats": func() (graph string) {
		graphs := []FancyGraphSimple{}

		rows, err := db.Query("SELECT users_osu FROM bancho_stats WHERE banchostats_id mod 10 = 0 ORDER BY banchostats_id DESC LIMIT 144")
		if err != nil {
			fmt.Println(err)
			return graph
		}
		iter := 0
		for rows.Next() {
			var graphItem FancyGraphSimple
			err = rows.Scan(&graphItem.Y)
			if err != nil {
				fmt.Println(err)
				return graph
			}

			graphItem.X = iter
			iter += 1

			graphs = append(graphs, graphItem)
		}

		for left, right := 0, len(graphs)-1; left < right; left, right = left+1, right-1 {
			posLeft := graphs[left].X
			graphs[left].X = graphs[right].X
			graphs[right].X = posLeft

			graphs[left], graphs[right] = graphs[right], graphs[left]
		}
		jsonStr, _ := json.Marshal(graphs)
		return string(jsonStr)
	},

	// "getBetterNews": func() (news []BetterNews) {
	// 	newsLetter, err := vkWrapper.WallGet(vkapi.Params{
	// 		"owner_id": -175152353,
	// 		"count":    10,
	// 		"filter":   " owner",
	// 	})

	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return news
	// 	}

	// 	for _, el := range newsLetter.Items {
	// 		betterNew := BetterNews{}

	// 		var (
	// 			finalStr string
	// 			count    int
	// 		)
	// 		count = 50
	// 		var lines []string = strings.Split(el.Text, "\n")
	// 		if len(lines) > 1 {
	// 			n := ReplaceUnCorrectSymbols(lines[0])
	// 			betterNew.Name = n
	// 			betterNew.News = ReplaceUnCorrectSymbols(GetFirstN(strings.Join(lines[1:], "\n"), 256))
	// 		} else {
	// 			s := ReplaceUnCorrectSymbols(el.Text)
	// 			if len(s) > 0 {
	// 				if count >= len(s) {
	// 					count = len(s) - 1
	// 				}
	// 				finalStr = s[:count]
	// 			}
	// 			betterNew.Name = finalStr
	// 			betterNew.News = ReplaceUnCorrectSymbols(GetFirstN(s, 256))
	// 		}

	// 		if len(el.Attachments) > 0 && el.Attachments[0].Type == "photo" {
	// 			betterNew.Preview = el.Attachments[0].Photo.MaxSize().URL
	// 		} else {
	// 			betterNew.Preview = "/static/headers/unnamed.png"
	// 		}

	// 		betterNew.Date = "02.02"
	// 		betterNew.Link = "https://vk.com/osukurikku?w=wall-175152353_" + strconv.Itoa(el.ID)

	// 		news = append(news, betterNew)
	// 	}
	// 	return news
	// },

	"getBetterNewsV2": func(language string) GhostPostsResponse {
		country := strings.ToLower(language)
		if country != "ru" {
			country = "en"
		}

		x := GhostPostsResponse{}
		d, err := http.Get(fmt.Sprintf("https://blog.kurikku.pw/ghost/api/v4/content/posts/?key=%s&filter=tag:%s", config.GhostApiToken, country))
		if err != nil {
			return x
		}
		data, _ := ioutil.ReadAll(d.Body)
		json.Unmarshal(data, &x)
		return x
	},

	"incrValue": func(value int, incr int) (newValue int) {
		return value + incr
	},

	// bget makes a request to the bancho api
	// https://docs.ripple.moe/docs/banchoapi/v1
	"bget": func(ept string, qs ...interface{}) map[string]interface{} {
		d, err := http.Get(fmt.Sprintf(config.BanchoAPI+"/api/v1/"+ept, qs...))
		if err != nil {
			return nil
		}
		x := make(map[string]interface{})
		data, _ := ioutil.ReadAll(d.Body)
		json.Unmarshal(data, &x)
		return x
	},
	"apiv2": func(ept string, qs ...interface{}) map[string]interface{} {
		d, err := http.Get(fmt.Sprintf(config.BaseURL+"/api/v2/"+ept, qs...))
		if err != nil {
			return nil
		}
		x := make(map[string]interface{})
		data, _ := ioutil.ReadAll(d.Body)
		json.Unmarshal(data, &x)
		return x
	},
	// styles returns playstyle.Styles
	"styles": func() []string {
		return playstyle.Styles[:]
	},
	// shift shifts n1 by n2
	"shift": func(n1, n2 int) int {
		return n1 << uint(n2)
	},
	// calculateDonorPrice calculates the price of x donor months in euros.
	"calculateDonorPrice": func(a float64) string {
		currencies := btcconversions.CurrencyRates
		return fmt.Sprintf("%.2f", a*(150/currencies.EUR))
	},
	// is2faEnabled checks 2fa is enabled for an user
	"is2faEnabled": is2faEnabled,
	// get2faConfirmationToken retrieves the current confirmation token for a certain user.
	"get2faConfirmationToken": get2faConfirmationToken,
	// csrfGenerate creates a csrf token input
	"csrfGenerate": func(u int) template.HTML {
		return template.HTML(`<input type="hidden" name="csrf" value="` + mustCSRFGenerate(u) + `">`)
	},
	// csrfURL creates a CSRF token for GET requests.
	"csrfURL": func(u int) template.URL {
		return template.URL("csrf=" + mustCSRFGenerate(u))
	},
	// systemSetting retrieves some information from the table system_settings
	"systemSettings": systemSettings,
	// authCodeURL gets the auth code for discord
	"authCodeURL": func(u int) string {
		return getDiscord().AuthCodeURL(mustCSRFGenerate(u))
	},
	// perc returns a percentage
	"perc": func(i, total float64) string {
		return fmt.Sprintf("%.0f", i/total*100)
	},
	// atLeastOne returns 1 if i < 1, or i otherwise.
	"atLeastOne": func(i int) int {
		if i < 1 {
			i = 1
		}
		return i
	},
	// ieForm fixes forms in IE/Trident being immensely fucked up. I hate microsoft.
	"ieForm": func(c *gin.Context) template.HTML {
		if !isIE(c.Request.UserAgent()) {
			return ""
		}
		return ieUnfucker
	},
	// version gets what's the current Hanayo version.
	"version": func() string {
		return version
	},
	"generateKey": generateKey,
	// getKeys gets the recovery 2fa keys for an user
	"getKeys": func(id int) []string {
		var keyRaw string
		db.Get(&keyRaw, "SELECT recovery FROM 2fa_totp WHERE userid = ?", id)
		s := make([]string, 0, 8)
		json.Unmarshal([]byte(keyRaw), &s)
		return s
	},
	// rediget retrieves a value from redis.
	"rediget": func(k string) string {
		x := rd.Get(k)
		if x == nil {
			return ""
		}
		if err := x.Err(); err != nil {
			fmt.Println(err)
		}
		return x.Val()
	},
	"getBitcoinAddress": btcaddress.Get,
	"languageInformation": func() []langInfo {
		return languageInformation
	},
	"languageInformationByNameShort": func(s string) langInfo {
		for _, lang := range languageInformation {
			if lang.NameShort == s {
				return lang
			}
		}
		return langInfo{}
	},
	"countryList": func(n int64) []string {
		return rd.ZRevRange("hanayo:country_list", 0, n-1).Val()
	},
	"documentationFiles": doc.GetDocs,
	"documentationData": func(slug string, language string) doc.File {
		if i, err := strconv.Atoi(slug); err == nil {
			slug = doc.SlugFromOldID(i)
		}
		return doc.GetFile(slug, language)
	},
	"privilegesToString": func(privs float64) string {
		return common.Privileges(privs).String()
	},
	"shouldShowJP": func(c *gin.Context, cont context) bool {
		for _, x := range getLang(c) {
			if x == "ja" {
				return cont.User.ID%5 == 4
			}
		}
		return false
	},
	"htmlescaper": template.HTMLEscaper,
}

var localeLanguages = []string{"en", "by", "ua", "de", "pl", "he", "it", "es", "ru", "th", "fr", "nl", "ro", "fi", "sv", "vi", "ko"}

var hanayoStarted = time.Now().UnixNano()

var servicePrefixes = map[string]string{
	"github":  "https://github.com/",
	"twitter": "https://twitter.com/",
	"mail":    "mailto:",
	"vk":      "https://vk.com/",
}

var logoColours = [...]string{
	"blue",
	"green",
	"orange",
	"red",
}

// we still haven't got jquery when the script is here, so well shit.
const ieUnfucker = `<input type="submit" class="ie" name="submit" value="submit">
<script>
var deferredToPageLoad = function() {
	$("button[form]").click(function() {
		$("form#" + $(this).attr("form") + " input.ie").click();
	});
};
</script>`

func pos(x int) (int, bool) {
	if x > 0 {
		return x, true
	}
	return x * -1, false
}
func _time(s string, t time.Time) template.HTML {
	return template.HTML(fmt.Sprintf(`<time class="timeago" datetime="%s">%v</time>`, s, t))
}

// Fantastic IEs And Where To Find Them
var ieUserAgentsContain = []string{
	"MSIE ",
	"Trident/",
	"Edge/",
}

func isIE(s string) bool {
	for _, v := range ieUserAgentsContain {
		if strings.Contains(s, v) {
			return true
		}
	}
	return false
}

type systemSetting struct {
	Name   string
	Int    int
	String string
}

func systemSettings(names ...string) map[string]systemSetting {
	var settingsRaw []systemSetting
	q, p, _ := sqlx.In("SELECT name, value_int as `int`, value_string as `string` FROM system_settings WHERE name IN (?)", names)
	err := db.Select(&settingsRaw, q, p...)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	settings := make(map[string]systemSetting, len(names))
	for _, s := range settingsRaw {
		settings[s.Name] = s
	}
	return settings
}

func getDiscord() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.DiscordOAuthID,
		ClientSecret: config.DiscordOAuthSecret,
		RedirectURL:  config.BaseURL + "/settings/discord/finish",
		Endpoint:     discordoauth.Endpoint,
		Scopes:       []string{"identify"},
	}
}

func getLanguageFromGin(c *gin.Context) string {
	for _, l := range getLang(c) {
		if in(l, localeLanguages) {
			return l
		}
	}
	return ""
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

type langInfo struct {
	Name, CountryShort, NameShort string
}

var languageInformation = []langInfo{
	{"Беларускі", "by", "by"},
	{"Deutsch", "de", "de"},
	{"English (UK)", "gb", "en"},
	{"Español", "es", "es"},
	{"Suomi", "fi", "fi"},
	{"Français", "fr", "fr"},
	{"עִבְרִית", "il", "he"},
	{"Italiano", "it", "it"},
	{"한국어", "kr", "ko"},
	{"Nederlands", "nl", "nl"},
	{"Polski", "pl", "pl"},
	{"Русский", "ru", "ru"},
	{"Română", "ro", "ro"},
	{"Svenska", "se", "sv"},
	{"ภาษาไทย", "th", "th"},
	{"Українська", "ua", "ua"},
	{"Tiếng Việt Nam", "vn", "vi"},
}
