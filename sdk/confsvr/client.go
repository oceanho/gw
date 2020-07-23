package confsvr

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/oceanho/gw/sdk/confsvr/req"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	options *Option
}

type Option struct {
	Addr         string
	Service      string
	Version      string
	Proto        string
	MaxRetries   int
	PullInterval int
}

func (c *Client) Do(reqobj req.IRequest, output interface{}) (int, error) {
	var reader io.Reader
	var buffer bufio.ReadWriter
	body := reqobj.Body()
	if len(body) > 0 {
		reader = &buffer
		buffer.Write(reqobj.Body())
	}

	uri := reqobj.Url()
	uri = strings.TrimLeft(uri, "/")
	uri = strings.TrimRight(uri, "/")

	url := fmt.Sprintf("%s/%s/%s/%s", c.options.Addr, c.options.Version, c.options.Service, uri)
	oriReq, _ := http.NewRequest(reqobj.Method(), url, reader)
	for k, header := range reqobj.Headers() {
		oriReq.Header.Set(k, header)
	}

	client := http.DefaultClient
	resp, err := client.Do(oriReq)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("do http request: %v", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("read resp body: %v", err)
	}
	if !req.IsHttpSucc(resp.StatusCode) {
		return resp.StatusCode, fmt.Errorf("status code [%d] are not excepted, resp text: \"%s\"", resp.StatusCode, string(b))
	}
	err = json.Unmarshal(b, output)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("decode str to json: \"%s\", %v", string(b), err)
	}
	return resp.StatusCode, nil
}

var (
	defaultAddr         = "https://127.0.0.1:8080"
	defaultVersion      = "/api/v1"
	defaultService      = "confsvr"
	defaultProto        = "http"
	defaultMaxRetries   = 10
	defaultPullInterval = 15
)

func DefaultOptions() *Option {
	return &Option{
		Addr:         defaultAddr,
		Proto:        defaultProto,
		Service:      defaultService,
		Version:      defaultVersion,
		MaxRetries:   defaultMaxRetries,
		PullInterval: defaultPullInterval,
	}
}

func NewClient(opts *Option) *Client {
	opts.Addr = strings.TrimRight(opts.Addr, "/")
	opts.Version = strings.TrimLeft(opts.Version, "/")
	opts.Version = strings.TrimRight(opts.Version, "/")
	client := &Client{
		options: opts,
	}
	return client
}