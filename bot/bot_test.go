package bot

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	cmn "github.com/jaimeteb/chatto/common"
)

func TestBot1(t *testing.T) {
	path := "../examples/00_test/"

	bot := LoadBot(&path)
	if bot.Name != "test_bot" {
		t.Errorf("bot name is incorrect, got: %v, want: %v.", bot.Name, "test_bot")
	}

	ans := bot.Answer(cmn.Message{
		Sender: "bar",
		Text:   "on",
	})

	if ans.(string) != "Turning on." {
		t.Errorf("answer is incorrect, got: %v, want: %v.", ans.(string), "Turning on.")
	}
}

func TestBot2(t *testing.T) {
	path := "../examples/00_test/"
	bot := LoadBot(&path)

	bot.Answer(cmn.Message{
		Sender: "baz",
		Text:   "on",
	})

	jsonStr := []byte(`{"sender": "42", "text": "on"}`)
	req, _ := http.NewRequest("POST", "", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	bot.restEndpointHandler(w, req)

	jsonStr2 := []byte(`{"update_id": 1, "message": {"message_id": 0, "from": {"id": 42, "first_name": "", "username": ""}, "date": 0, "text": "off"}}`)
	req2, _ := http.NewRequest("POST", "", bytes.NewBuffer(jsonStr2))
	w2 := httptest.NewRecorder()
	bot.telegramEndpointHandler(w2, req2)

	formData := url.Values{
		"From":             {"42"},
		"Body":             {"?"},
		"To":               {""},
		"MediaUrl":         {""},
		"MediaContentType": {""},
		"MessageSid":       {""},
		"SmsStatus":        {""},
		"AccountSid":       {""},
		"Sid":              {""},
		"SmsSid":           {""},
		"SmsMessageSid":    {""},
		"NumMedia":         {"0"},
		"NumSegments":      {"0"},
		"ApiVersion":       {""},
	}
	req3, _ := http.NewRequest("POST", "", strings.NewReader(formData.Encode()))
	w3 := httptest.NewRecorder()
	bot.twilioEndpointHandler(w3, req3)

	req4, _ := http.NewRequest("GET", "/senders/42", nil)
	w4 := httptest.NewRecorder()
	bot.detailsHandler(w4, req4)

	jsonStr5 := []byte(`{"text": "."}`)
	req5, _ := http.NewRequest("POST", "", bytes.NewBuffer(jsonStr5))
	w5 := httptest.NewRecorder()
	bot.predictHandler(w5, req5)

	jsonStr6 := []byte(`{"event": {"channel": "43", "text": "on"}}`)
	req6, _ := http.NewRequest("POST", "", bytes.NewBuffer(jsonStr6))
	w6 := httptest.NewRecorder()
	bot.slackEndpointHandler(w6, req6)

	jsonStr7 := []byte(`{"challenge": "challenge"}`)
	req7, _ := http.NewRequest("POST", "", bytes.NewBuffer(jsonStr7))
	w7 := httptest.NewRecorder()
	bot.slackEndpointHandler(w7, req7)
}

func TestBotNoClientsAndImages(t *testing.T) {
	path := "../examples/01_moodbot/"

	bot := LoadBot(&path)
	if bot.Clients.Telegram.Client != nil || bot.Clients.Twilio.Client != nil {
		t.Errorf("bot.Clients is incorrect, got: %v, want: %v.", bot.Clients, "{}")
	}

	wREST := httptest.NewRecorder()
	messages := []interface{}{
		cmn.Message{
			Text: "only text",
		},
		cmn.Message{
			Text:  "text and image",
			Image: "https://i.imgur.com/8MU0IUT.jpeg",
		},
		"string in the wild",
		map[string]interface{}{
			"text": "text in map",
		},
		map[interface{}]interface{}{
			"text": "text in interface map",
		},
	}
	SendMessages(messages, &bot.Clients.REST, "8809", wREST)

	SendMessages(new(interface{}), &bot.Clients.REST, "8809", wREST)
}

func TestServeBot(t *testing.T) {
	path := "../examples/00_test/"
	port := 9999

	go ServeBot(&path, &port)
}

func TestExtFromBot(t *testing.T) {
	path := "../examples/00_test/"
	bot := LoadBot(&path)
	bot.Clients = Clients{}
	bot.Answer(cmn.Message{
		Sender: "ext_tester",
		Text:   "hello",
		Image:  "",
	})
}
