package transport

import "encoding/json"

type JsonPayloadData struct {
	Payload any
}

var _ PayloadData = &JsonPayloadData{}

func (j JsonPayloadData) Accept() string {
	return "application/json"
}

func (j JsonPayloadData) ContentType() string {
	return "application/json"
}

func (j JsonPayloadData) Content() ([]byte, error) {
	payload, err := json.Marshal(j.Payload)
	return payload, err
}

func (j JsonPayloadData) ContentMust() []byte {
	payload, err := j.Content()
	if err != nil {
		panic(err)
	}

	return payload
}
