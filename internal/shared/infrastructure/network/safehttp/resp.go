package safehttp

import (
	"io/ioutil"
	"log"
	"net/http"
)

type Resp struct {
	req  *Req
	resp *http.Response
}

func NewResp(resp *http.Response, req *Req) (*Resp, error) {
	return &Resp{resp: resp, req: req}, nil
}

func (r *Resp) Req() *Req {
	return r.req
}

func (r *Resp) Origin() *http.Response {
	return r.resp
}

func (r *Resp) Status() string {
	return r.resp.Status
}

func (r *Resp) StatusCode() int {
	return r.resp.StatusCode
}

func (r *Resp) Body() (string, error) {
	defer func() {
		if err := r.Close(); err != nil {
			log.Fatalln("resp: error occurred while closing body. " + err.Error())
		}
	}()
	body, err := ioutil.ReadAll(r.resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (r *Resp) Close() error {
	return r.resp.Body.Close()
}
