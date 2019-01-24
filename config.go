package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type config struct {
	requests    []*http.Request
	req         []Request     `json:"requests"`
	connections int           `json:"connections`
	host        string        `json:"host"`
	duration    time.Duration `json:"duration"`
}

type Request struct {
	Method      string `json:"method"`
	Endpoint    string `json:"endpoint"`
	Data        string `json:"data"`
	Header      string `json:"header"`
	Connections string `json:"connections"`
}

// FromJSON returns a config from a json file
func fromJSON(path string) (*config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var conf config

	err = json.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}

	for _, req := range conf.req {
		var buf io.Reader
		addr := conf.host + req.Endpoint

		if req.Data != "" {
			buf = bytes.NewBufferString(req.Data)
		}

		r, err := http.NewRequest(req.Method, addr, buf)
		if err != nil {
			return nil, err
		}

		_ = append(conf.requests, r)
	}

	return &conf, nil

}
