package webhook

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/hazcod/go-intigriti"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"
)

type Endpoint struct {
	webhookURL 	string
	headers 	map[string]string
	format		string
	method 		string
	clientVersion string
}

func NewEndpoint(url *url.URL, method string, extraHeaders map[string]string, format string, clientVersion string) Endpoint {
	return Endpoint{
		webhookURL: url.String(),
		method:  strings.ToLower(method),
		headers:    extraHeaders,
		clientVersion: clientVersion,
		format: strings.ToLower(format),
	}
}

func formatPayload(format string, submission intigriti.Submission) ([]byte, error) {
	switch strings.ToLower(format) {

	case "xml":
		b, err := xml.Marshal(&submission)
		if err != nil {
			return nil, errors.Wrap(err, "could not convert to xml")
		}
		return b, nil

	case "json":
		b, err := json.Marshal(&submission)
		if err != nil {
			return nil, errors.Wrap(err, "could not convert to json")
		}
		return b, nil
	}

	return nil, errors.New("unknown output format: " + format)
}

func (e *Endpoint) Send(submission intigriti.Submission) error {
	client := &http.Client{}

	payload, err := formatPayload(e.format, submission)
	if err != nil {
		return errors.Wrap(err, "unknown output format")
	}

	req, err := http.NewRequest(strings.ToUpper(e.method), e.webhookURL, bytes.NewReader(payload))
	if err != nil {
		return errors.Wrap(err, "could not create http request")
	}

	req.Header.Set("User-Agent", e.clientVersion)
	for key, value := range e.headers {
		req.Header.Set(key, value)
	}

	if _, err = client.Do(req); err != nil {
		return errors.Wrap(err, "could not send webhook")
	}

	return nil
}