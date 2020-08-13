package main

import (
	"bytes"
	"encoding/json"
	_ "image/gif"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

func avatarSubmit(c *gin.Context) {
	ctx := getContext(c)
	if ctx.User.ID == 0 {
		resp403(c)
		return
	}
	var m message
	defer func() {
		simpleReply(c, m)
	}()
	file, multipartFileHeader, err := c.Request.FormFile("avatar")
	if err != nil {
		m = errorMessage{T(c, "An error occurred.")}
		return
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("avatar", multipartFileHeader.Filename)
	avatarContent, err := ioutil.ReadAll(file)
	part.Write(avatarContent)

	rt, err := writer.CreateFormField("rt")
	rt.Write([]byte(ctx.Token))

	if err != nil {
		log.Fatal(err)
	}

	//io.Copy(rt, strings.NewReader(ctx.Token))
	writer.Close()
	request, err := http.NewRequest("POST", "http://127.0.0.1:5010/avatars/uploadAvatar", body)

	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	errorRed := struct {
		Wrong string `json:"error"`
	}{}
	_ = json.Unmarshal(content, &errorRed)

	if len(errorRed.Wrong) > 1 {
		m = errorMessage{T(c, errorRed.Wrong)}
		return
	}

	// okay thats' not error
	// now try to unmarshal in repsonse

	goodResponse := struct {
		Response string `json:"response"`
	}{}
	err = json.Unmarshal(content, &goodResponse)

	if err != nil {
		m = errorMessage{T(c, "Something happend, try again later")}
		return
	}

	m = successMessage{T(c, goodResponse.Response)}
}
