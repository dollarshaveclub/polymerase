package httpreq

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type HTTPService interface {
	HTTPRequest(string, string, io.Reader, map[string]string, bool) (*HTTPResponse, error)
}

type HTTPResponse struct {
	Body      string
	BodyBytes []byte
	Resp      *http.Response
}

type HTTPRequestConfig struct {
	URL                   string
	Method                string
	Body                  io.Reader
	Headers               map[string]string
	FailOnError           bool
	InsecureSkipTLSVerify bool
	TimeoutSeconds        uint
}

func getRespBody(resp *http.Response) (string, []byte, error) {
	defer resp.Body.Close()
	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", []byte{}, err
	}
	return string(bb), bb, nil
}

// HTTPRequest executes a given HTTP API request, returning response body
func HTTPRequest(url string, method string, body io.Reader, headers map[string]string, failOnError bool) (*HTTPResponse, error) {
	c := &HTTPRequestConfig{
		URL:            url,
		Method:         method,
		Body:           body,
		Headers:        headers,
		FailOnError:    failOnError,
		TimeoutSeconds: 30,
	}
	return HTTPComplexRequest(c)
}

//HTTPComplexRequest allows more control over request options
func HTTPComplexRequest(c *HTTPRequestConfig) (*HTTPResponse, error) {
	hresp := &HTTPResponse{}
	req, err := http.NewRequest(c.Method, c.URL, c.Body)
	if err != nil {
		return hresp, err
	}
	for k, v := range c.Headers {
		req.Header.Add(k, v)
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: c.InsecureSkipTLSVerify},
	}
	hc := http.Client{
		Transport: tr,
		Timeout:   time.Duration(c.TimeoutSeconds) * time.Second,
	}
	resp, err := hc.Do(req)
	if err != nil {
		return hresp, err
	}
	bs, bb, err := getRespBody(resp)
	if err != nil {
		return hresp, err
	}
	hresp.Body = bs
	hresp.BodyBytes = bb
	hresp.Resp = resp
	if resp.StatusCode > 399 && c.FailOnError {
		return hresp, fmt.Errorf("Server response indicates failure: %v %v", resp.StatusCode, bs)
	}
	return hresp, nil
}
