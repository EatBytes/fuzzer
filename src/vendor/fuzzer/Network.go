package fuzzer

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

type callback func(string)

var NET = NETWORK{
	host:      "http://localhost",
	method:    0,
	parameter: "fuzzer",
	crypt:     false,
	status:    false,
}

type config struct {
	url    string
	method string
	form   *bytes.Buffer
	jar    []string
	proxy  string
}

type NETWORK struct {
	host      string
	method    int
	parameter string
	crypt     bool
	status    bool
	cmd       string

	_body         url.Values
	_lastResponse *http.Response
}

func (n *NETWORK) GetMethod() int {
	return n.method
}

func (n *NETWORK) GetParameter() string {
	return n.parameter
}

func (n *NETWORK) IsSetup() bool {
	return n.status
}

func (n *NETWORK) SetConfig(url string, method int, parameter string, crypt bool) {
	n.host = url
	n.method = method
	n.parameter = parameter
	n.crypt = crypt

	n.status = true
}

func (n *NETWORK) Send(r string, f callback) {
	var httpResponse *http.Response
	var response string

	if n.method == 0 {
		httpResponse = n.get(r)
		buffer := GetBody(httpResponse)
		response = string(buffer)
	}

	if n.method == 1 {
		httpResponse = n.post(r)
		buffer := GetBody(httpResponse)
		response = string(buffer)
	}

	if n.method == 3 {
		n.getWithHeader(r)
	}

	if n.method == 4 {
		n.getWithCookie(r)
	}

	if httpResponse != nil && httpResponse.StatusCode < 400 {
		f(response)
	} else {
		fmt.Println("Error with the response: " + httpResponse.Status)
	}
}

func (n *NETWORK) post(r string) *http.Response {
	n.status = true

	request := Encode(r)
	n.cmd = request

	form := url.Values{}
	form.Set(n.parameter, request)
	n._body = form

	data := bytes.NewBufferString(form.Encode())

	c := config{
		url:    n.host,
		method: "POST",
		form:   data,
	}

	return n._send(&c)
}

func (n *NETWORK) get(r string) *http.Response {
	n.status = true

	request := Encode(r)
	n.cmd = request

	url := n.host + "?" + n.parameter + "=" + request

	c := config{
		url:    url,
		method: "GET",
		form:   nil,
	}

	return n._send(&c)
}

func (n *NETWORK) getWithHeader(r string) {

}

func (n *NETWORK) getWithCookie(r string) {

}

func (n *NETWORK) _send(c *config) *http.Response {
	client := &http.Client{}
	data := c.form

	var req *http.Request
	var err error

	if c.form != nil {
		req, err = http.NewRequest(c.method, c.url, data)
	} else {
		req, err = http.NewRequest(c.method, c.url, nil)
	}

	if err != nil {
		panic(err)
	}

	n._headerConfig(req)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	n._lastResponse = resp

	return resp
}

func (n *NETWORK) _headerConfig(req *http.Request) {
	if n.method == 1 {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	} else if n.method == 2 {
		req.Header.Add(n.parameter, n.cmd)
	}
}

func (n *NETWORK) GetResponse() *http.Response {
	return n._lastResponse
}

func (n *NETWORK) GetRequest() *http.Request {
	n._lastResponse.Request.PostForm = n._body
	return n._lastResponse.Request
}