package safehttp

import (
	"io/ioutil"
	"log"
	"net/http"
)

type Resp struct {
	origin *http.Response
}

func NewResp(resp *http.Response) (*Resp, error) {
	return &Resp{origin: resp}, nil
}

func (r *Resp) Origin() *http.Response {
	return r.origin
}

func (r *Resp) Status() string {
	return r.origin.Status
}

func (r *Resp) StatusCode() int {
	return r.origin.StatusCode
}

func (r *Resp) Body() (string, error) {
	defer func() {
		if err := r.Close(); err != nil {
			log.Fatalln("resp: error occurred while closing body. " + err.Error())
		}
	}()
	body, err := ioutil.ReadAll(r.origin.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (r *Resp) Close() error {
	return r.origin.Body.Close()
}
