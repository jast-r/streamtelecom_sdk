package streamtelecom

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

func (c *Client) SendSingleSMS(ctx context.Context, sms SingleSMSRequest) (string, error) {
	form := url.Values{}
	form.Add("login", c.login)
	form.Add("pass", c.password)

	if sms.SourceAddress == "" {
		return "", errors.New("source address can`t be empty")
	}
	if sms.DestinationAddress == "" {
		return "", errors.New("destination address can`t be empty")
	}
	if sms.Text == "" {
		return "", errors.New("message text can`t be empty")
	}

	form.Add("sourceAddress", sms.SourceAddress)
	form.Add("destinationAddress", sms.DestinationAddress)
	form.Add("data", sms.Text)

	if sms.SendDate.Unix() > 0 {
		form.Add("sendDate", sms.SendDate.Format("2006-01-02T15:04:05"))
	}
	if sms.TTL != "" {
		form.Add("validity", sms.TTL)
	}
	if sms.CallbackUrl != "" {
		form.Add("callback_url", sms.CallbackUrl)
	}
	if sms.UserID != "" {
		form.Add("user_id", sms.UserID)
	}
	if sms.NameDeliver != "" {
		form.Add("name_deliver", sms.NameDeliver)
	}

	val, err := c.doHTTP(ctx, endpointSendSingleSMS, http.MethodPost, form)
	if err != nil {
		return "", err
	}

	var streamtelecomID []string
	err = json.Unmarshal(val, &streamtelecomID)
	if err != nil {
		return "", err
	}

	if len(streamtelecomID) == 0 {
		return "", errors.New("no id`s into response")
	}

	return streamtelecomID[0], nil
}
