// Package gotwilio is a library for interacting with http://www.twilio.com/ API.
package gotwilio

import (
	"net/http"
	"net/url"
	"strings"
	"appengine"
	"appengine/urlfetch"
	"io/ioutil"
)

// Twilio stores basic information important for connecting to the
// twilio.com REST api such as AccountSid and AuthToken.
type Twilio struct {
	AccountSid string
	AuthToken  string
	BaseUrl    string
	Client     *http.Client
	Context    appengine.Context
}

// Exception is a representation of a twilio exception.
type Exception struct {
	Status   int    `json:"status"`    // HTTP specific error code
	Message  string `json:"message"`   // HTTP error message
	Code     int    `json:"code"`      // Twilio specific error code
	MoreInfo string `json:"more_info"` // Additional info from Twilio
}

// Create a new Twilio struct.
func NewTwilioClient(accountSid, authToken string, c appengine.Context) *Twilio {
	twilioUrl := "https://api.twilio.com/2010-04-01" // Should this be moved into a constant?

	return &Twilio{accountSid, authToken, twilioUrl, urlfetch.Client(c), c}
}

func (twilio *Twilio) post(formValues url.Values, twilioUrl string) (*http.Response, error) {
	req, err := http.NewRequest("POST", twilioUrl, strings.NewReader(formValues.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(twilio.AccountSid, twilio.AuthToken)

	twilio.Context.Infof("request: %v",req)
	
	resp , err := twilio.Client.Do(req)
	if err != nil {
		twilio.Context.Infof("error: %v",err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	twilio.Context.Infof("body: %s",string(body))

	return resp,err
}

func (twilio *Twilio) myPost(formValues url.Values, twilioUrl string) (*http.Response, error) {
	req, err := http.NewRequest("POST", twilioUrl, strings.NewReader(formValues.Encode()))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(twilio.AccountSid, twilio.AuthToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return twilio.Client.Do(req)
}