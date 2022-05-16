package streamtelecom

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	apiURL = "http://gateway.api.sc/rest"

	endpointBalance       = "/Balance/balance.php"
	endpointSenderList    = "/Statistic/originator.php"
	endpointTariffList    = "/Balance/price_list.php"
	endpointSendSingleSMS = "/Send/SendSms/"

	defaultTimeout = 5 * time.Second
)

type (
	tariffListReponse struct {
		SMS map[string]struct {
			Price string `json:"R"`
		} `json:"sms"`
		Hlr map[string]struct {
			Price string `json:"R"`
		} `json:"hlr"`
		Email map[string]struct {
			Price string `json:"R"`
		} `json:"email"`
		Tg map[string]struct {
			Price string `json:"R"`
		} `json:"tg"`
		Vk map[string]struct {
			Price string `json:"R"`
		} `json:"vk"`
		Bot map[string]struct {
			Price string `json:"R"`
		} `json:"bot"`
	}

	SingleSMSRequest struct {
		login              string
		password           string
		DestinationAddress string
		SendDate           time.Time
		Text               string
		SourceAddress      string
		TTL                string
		CallbackUrl        string
		UserID             string
		NameDeliver        string
	}

	errorResponse struct {
		Code int    `json:"Code"`
		Desc string `json:"Desc"`
	}
)

type Client struct {
	client   *http.Client
	login    string
	password string
}

func NewClient(login, password string) (*Client, error) {
	if login == "" || password == "" {
		return nil, errors.New("login or passwords canno`t be empty")
	}

	return &Client{
		client: &http.Client{
			Timeout: defaultTimeout, // 5 sec
		},
		login:    login,
		password: password,
	}, nil
}

func (c *Client) GetBalance(ctx context.Context) (float64, error) {
	form := url.Values{}
	form.Add("login", c.login)
	form.Add("pass", c.password)

	val, err := c.doHTTP(ctx, endpointBalance, http.MethodPost, form)
	if err != nil {
		return -1, err
	}

	var balance float64
	err = json.Unmarshal(val, &balance)
	if err != nil {
		return -1, err
	}

	return balance, nil
}

func (c *Client) GetSenderList(ctx context.Context) ([]string, error) {
	form := url.Values{}
	form.Add("login", c.login)
	form.Add("pass", c.password)

	val, err := c.doHTTP(ctx, endpointSenderList, http.MethodPost, form)
	if err != nil {
		return nil, err
	}

	var senders []string

	err = json.Unmarshal(val, &senders)
	if err != nil {
		return nil, err
	}
	return senders, nil
}

func (c *Client) GetTariffList(ctx context.Context) (*tariffListReponse, error) {
	form := url.Values{}
	form.Add("login", c.login)
	form.Add("pass", c.password)

	val, err := c.doHTTP(ctx, endpointTariffList, http.MethodPost, form)
	if err != nil {
		return nil, err
	}

	var tarrifs tariffListReponse

	err = json.Unmarshal(val, &tarrifs)
	if err != nil {
		return nil, err
	}

	return &tarrifs, nil
}

func (c *Client) doHTTP(ctx context.Context, endpoint, method string, form url.Values) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, apiURL+endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create new request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to send http request")
	}

	respB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read request body")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var failedResponse errorResponse
		err = json.Unmarshal(respB, &failedResponse)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to unmarshal input body")
		}
		return nil, errors.New(failedResponse.Desc)
	}

	return respB, nil
}
