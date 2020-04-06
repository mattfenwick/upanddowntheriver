package game

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"time"
)

type Client struct {
	Host  string
	Port  int
	Resty *resty.Client
}

func NewClient(host string, port int) *Client {
	restyClient := resty.New()
	restyClient.SetRetryCount(3)
	restyClient.SetRetryWaitTime(500 * time.Millisecond)
	restyClient.SetTimeout(time.Duration(5 * time.Second))
	return &Client{Host: host, Port: port, Resty: restyClient}
}

func (client *Client) url(path string) string {
	return fmt.Sprintf("http://%s:%d/%s", client.Host, client.Port, path)
}

func (client *Client) GetModel() (string, error) {
	url := client.url("model")
	resp, err := client.Resty.R(). /*.SetBody()*/ Get(url)
	if err != nil {
		return resp.String(), err
	}
	if resp.StatusCode() < 200 || resp.StatusCode() > 299 {
		return resp.String(), errors.New(fmt.Sprintf("bad status code from %s: %d", url, resp.StatusCode()))
	}
	return resp.String(), nil
}

func (client *Client) postJson(path string, body interface{}, result interface{}) (string, error) {
	url := client.url(path)
	req := client.Resty.R().SetHeader("Content-Type", "application/json")
	if result != nil {
		req = req.SetResult(result)
	}
	if body != nil {
		req = req.SetBody(body)
	}

	resp, err := req.Post(url)
	if err != nil {
		return resp.String(), err
	}
	if resp.StatusCode() < 200 || resp.StatusCode() > 299 {
		return resp.String(), errors.New(fmt.Sprintf("bad status code from %s: %d", url, resp.StatusCode()))
	}

	return resp.String(), nil
}

func (client *Client) postAction(action *PlayerAction) (*PlayerModel, error) {
	result := &PlayerModel{}
	_, err := client.postJson("action", action, result)
	return result, err
}

func (client *Client) GetMyModel(me string) (*PlayerModel, error) {
	body := &PlayerAction{Me: me, GetModel: &GetPlayerModelAction{}}
	return client.postAction(body)
}
