package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	ext "github.com/jaimeteb/chatto/extension"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/query"
)

var weatherKey = os.Getenv("WEATHER_API_KEY")
var weatherURL = "http://api.weatherapi.com/v1/current.json?key=%s&q=%s"

type weatherResponse struct {
	Location weatherResponseLocation `json:"location"`
	Current  weatherResponseCurrent  `json:"current"`
}

type weatherResponseLocation struct {
	Name    string `json:"name"`
	Region  string `json:"region"`
	Country string `json:"country"`
}

type weatherResponseCurrent struct {
	Condition  weatherResponseCurrentCondition `json:"condition"`
	TempC      float32                         `json:"temp_c"`
	TempF      float32                         `json:"temp_f"`
	FeelsLikeC float32                         `json:"feelslike_c"`
	FeelsLikeF float32                         `json:"feelslike_f"`
	Humidity   int                             `json:"humidity"`
}

type weatherResponseCurrentCondition struct {
	Text string `json:"text"`
}

var jokeURL = "https://v2.jokeapi.dev/joke/Any?blacklistFlags=nsfw,religious,political,racist,sexist,explicit&type=single"

type jokeResponse struct {
	Joke string `json:"joke"`
}

var quoteURL = "http://api.quotable.io/random"

type quoteResponse struct {
	Content string `json:"content"`
	Author  string `json:"author"`
}

var serpKey = os.Getenv("SCALE_SERP_API_KEY")
var serpURL = "https://api.scaleserp.com/search?api_key=%s&q=%s"

type serpResponse struct {
	ResponseInfo serpResponseInfo      `json:"request_info"`
	AnswerBox    serpResponseAnswerBox `json:"answer_box"`
}

type serpResponseInfo struct {
	Success          bool `json:"success"`
	CreditsRemaining int  `json:"credits_remaining"`
}

type serpResponseAnswerBox struct {
	AnswerBoxType int                  `json:"answer_box_type"`
	Answers       []serpResponseAnswer `json:"answers"`
}

type serpResponseAnswer struct {
	Answer string                   `json:"answer"`
	Source serpResponseAnswerSource `json:"source"`
}

type serpResponseAnswerSource struct {
	Link string `json:"link"`
}

func errFunc(req *ext.Request, err error) *ext.Response {
	log.Errorf("%#v", err)
	return &ext.Response{
		FSM:     req.FSM,
		Answers: query.Answers(req.DB.DefaultMessages.Error),
	}
}

func weatherFunc(req *ext.Request) (res *ext.Response) {
	location := url.QueryEscape(req.Question.Text)

	resp, err := http.Get(fmt.Sprintf(weatherURL, weatherKey, location))
	if err != nil {
		return errFunc(req, err)
	}

	defer resp.Body.Close()
	var weatherResp weatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		return errFunc(req, err)
	}

	var message string
	switch resp.StatusCode {
	case 200:
		weatherMessage := "In %s, %s, it is %s. The temperature is %2.1f °C (%2.1f °F) and feels like %2.1f °C (%2.1f °F)."
		message = fmt.Sprintf(
			weatherMessage,
			weatherResp.Location.Name,
			weatherResp.Location.Region,
			strings.ToLower(weatherResp.Current.Condition.Text),
			weatherResp.Current.TempC,
			weatherResp.Current.TempF,
			weatherResp.Current.FeelsLikeC,
			weatherResp.Current.FeelsLikeC,
		)
	case 400:
		message = "Sorry, I couldn't find your location, try with another one please."
		return &ext.Response{
			FSM: &fsm.FSM{
				State: req.DB.StateTable["ask_location"],
				Slots: req.FSM.Slots,
			},
			Answers: query.Answers(message),
		}
	default:
		return errFunc(req, errors.New(resp.Status))
	}

	return &ext.Response{
		FSM:     req.FSM,
		Answers: query.Answers(message),
	}
}

func jokeFunc(req *ext.Request) (res *ext.Response) {
	resp, err := http.Get(jokeURL)
	if err != nil {
		return errFunc(req, err)
	}

	defer resp.Body.Close()
	var jokeResp jokeResponse
	if err := json.NewDecoder(resp.Body).Decode(&jokeResp); err != nil {
		return errFunc(req, err)
	}

	return &ext.Response{
		FSM:     req.FSM,
		Answers: query.Answers(jokeResp.Joke),
	}
}

func quoteFunc(req *ext.Request) (res *ext.Response) {
	resp, err := http.Get(quoteURL)
	if err != nil {
		return errFunc(req, err)
	}

	defer resp.Body.Close()
	var quoteResp quoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&quoteResp); err != nil {
		return errFunc(req, err)
	}

	return &ext.Response{
		FSM:     req.FSM,
		Answers: query.Answers(fmt.Sprintf("%s\n    - %s", quoteResp.Content, quoteResp.Author)),
	}
}

func miscFunc(req *ext.Request) (res *ext.Response) {
	queryURL := url.QueryEscape(strings.ReplaceAll(req.Question.Text, " ", "+"))

	resp, err := http.Get(fmt.Sprintf(serpURL, serpKey, queryURL))
	if err != nil {
		return errFunc(req, err)
	}

	defer resp.Body.Close()
	var serpResp serpResponse
	if err := json.NewDecoder(resp.Body).Decode(&serpResp); err != nil {
		return errFunc(req, err)
	}

	if serpResp.AnswerBox.AnswerBoxType == 0 || len(serpResp.AnswerBox.Answers) == 0 {
		return &ext.Response{
			FSM:     req.FSM,
			Answers: query.Answers("I'm sorry, I couldn't find an answer to that question."),
		}
	}

	answer := serpResp.AnswerBox.Answers[0]

	if answer.Answer == "" {
		return &ext.Response{
			FSM:     req.FSM,
			Answers: query.Answers("I'm sorry, I couldn't find an answer to that question."),
		}
	}

	message := answer.Answer
	if answer.Source.Link != "" {
		message += " \nSource: " + answer.Source.Link
	}

	return &ext.Response{
		FSM:     req.FSM,
		Answers: query.Answers(message),
	}
}

var myExtMap = ext.RegisteredFuncs{
	"ext_weather": weatherFunc,
	"ext_joke":    jokeFunc,
	"ext_quote":   quoteFunc,
	"ext_misc":    miscFunc,
}

func main() {
	log.Fatalln(ext.ServeExtensionREST(myExtMap))
}
