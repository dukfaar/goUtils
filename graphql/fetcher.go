package graphql

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

type Request struct {
	Query         string                 `json:"query,omitempty" url:"query" schema:"query"`
	Variables     map[string]interface{} `json:"variables,omitempty" url:"variables" schema:"variables"`
	OperationName string                 `json:"operationName,omitempty" url:"operationName" schema:"operationName"`
}

type Fetcher interface {
	Fetch(request Request) (interface{}, error)
}

type WebSocketFetcher struct {
	Connection *websocket.Conn
}

func NewWebSocketFetcher(host string, path string) (*WebSocketFetcher, error) {
	u := url.URL{Scheme: "ws", Host: host, Path: path}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		fmt.Printf("Error connecting: %v\n", err)
		return nil, err
	}

	return &WebSocketFetcher{
		Connection: conn,
	}, nil
}

func (f *WebSocketFetcher) Fetch(request Request) (interface{}, error) {
	err := f.Connection.WriteJSON(request)

	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}

	_, msg, err := f.Connection.ReadMessage()

	var response interface{}
	err = json.Unmarshal(msg, &response)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}

	return response, nil
}

type HttpFetcher struct {
	url     string
	client  http.Client
	headers map[string]string
}

func NewHttpFetcher(host string, path string) (*HttpFetcher, error) {
	return &HttpFetcher{
		url:     "http://" + host + path,
		client:  http.Client{},
		headers: make(map[string]string),
	}, nil
}

func (f *HttpFetcher) SetHeader(name string, value string) {
	f.headers[name] = value
}

func (f *HttpFetcher) Fetch(request Request) (interface{}, error) {
	jsonValue, _ := json.Marshal(request)

	httpRequest, err := http.NewRequest("POST", f.url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}

	for key, value := range f.headers {
		httpRequest.Header.Set(key, value)
	}

	response, err := f.client.Do(httpRequest)

	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if result["errors"] != nil {
		errorString, _ := json.Marshal(result["errors"])
		return nil, errors.New(string(errorString))
	}

	return result["data"], nil
}
