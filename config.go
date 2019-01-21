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
	connections int
	host        string
	duration    time.Duration
}

type File struct {
	Requests    []Request     `json:"requests"`
	Connections int           `json:"connections"`
	Host        string        `json:"host"`
	Duration    time.Duration `json:"Duration"`
}

type Request struct {
	Method   string `json:"method"`
	Endpoint string `json:"endpoint"`
	Data     string `json:"data"`
	Header   string `json:"header"`
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

	var file File

	err = json.Unmarshal(data, &file)
	if err != nil {
		return nil, err
	}

	conf := &config{
		connections: file.Connections,
		host:        file.Host,
		duration:    file.Duration,
	}

	for _, req := range file.Requests {
		var buf io.Reader
		addr := file.Host + req.Endpoint

		if req.Data != "" {
			buf = bytes.NewBufferString(req.Data)
		}

		r, err := http.NewRequest(req.Method, addr, buf)
		if err != nil {
			return nil, err
		}

		_ = append(conf.requests, r)
	}

	return conf, nil

}
