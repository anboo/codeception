package codeception

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

type Actor struct {
	baseURL      string
	headers      map[string]string
	lastResponse *http.Response
	client       *http.Client
	t            *testing.T
}

func NewActor(t *testing.T, baseURL string, headers map[string]string) *Actor {
	return &Actor{
		baseURL: baseURL,
		headers: headers,
		client:  &http.Client{},
		t:       t,
	}
}

func (a *Actor) SendGet(endpoint string, body map[string]interface{}) *Actor {
	queryParameters := url.Values{}
	for field, value := range body {
		queryParameters.Add(field, fmt.Sprintf("%v", value))
	}

	req, err := http.NewRequest("GET", a.formUrl(endpoint)+"?"+queryParameters.Encode(), nil)
	if err != nil {
		a.t.Fatalf("create request: %s", err)
	}

	err = a.makeRequest(req)
	if err != nil {
		a.t.Fatalf("make request: %s", err)
	}
	return a
}

func (a *Actor) SendPost(endpoint string, body interface{}) *Actor {
	reqBody, err := json.Marshal(body)
	if err != nil {
		a.t.Fatalf("json marshal: %s", err)
	}

	req, err := http.NewRequest("POST", a.formUrl(endpoint), bytes.NewBuffer(reqBody))
	if err != nil {
		a.t.Fatalf("new request: %s", err)
	}

	err = a.makeRequest(req)
	if err != nil {
		a.t.Fatalf("make request %s", err)
	}
	return a
}

func (a *Actor) SendPatch(endpoint string, body interface{}) *Actor {
	reqBody, err := json.Marshal(body)
	if err != nil {
		a.t.Fatalf("json marshal: %s", err)
	}

	req, err := http.NewRequest("PATCH", a.formUrl(endpoint), bytes.NewBuffer(reqBody))
	if err != nil {
		a.t.Fatalf("new request: %s", err)
	}

	err = a.makeRequest(req)
	if err != nil {
		a.t.Fatalf("make request: %s", err)
	}
	return a
}

func (a *Actor) SeeResponseCodeIs(status int) *Actor {
	if a.lastResponse.StatusCode != status {
		a.t.Fatalf("expected response code is %d got %d", status, a.lastResponse.StatusCode)
	}
	return a
}

func (a *Actor) DontSeeResponseCodeIs(status int) *Actor {
	if a.lastResponse.StatusCode == status {
		a.t.Fatalf("expected response code is not %d got %d", status, a.lastResponse.StatusCode)
	}
	return a
}

func (a *Actor) SeeJSON(expected map[string]interface{}) *Actor {
	resp := map[string]interface{}{}

	responseBody, err := io.ReadAll(a.lastResponse.Body)
	if err != nil {
		a.t.Fatalf("cannot read response")
	}

	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		a.t.Fatalf("unmarhsal json response")
	}

	if !reflect.DeepEqual(resp, expected) {
		a.t.Fatalf("expected response contain %v got %v", expected, resp)
	}

	return a
}

func (a *Actor) makeRequest(req *http.Request) error {
	response, err := a.client.Do(req)
	if err != nil {
		a.t.Fatalf("http client do: %s", err)
	}
	a.lastResponse = response
	return nil
}

func (a *Actor) formUrl(endpoint string) string {
	return a.baseURL + endpoint
}
