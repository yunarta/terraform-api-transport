package transport

import "encoding/json"

// Authentication abstraction
type Authentication interface {
}

// PayloadTransport
// Abstraction of the transport
type PayloadTransport interface {
	Send(request *PayloadRequest) (*PayloadResponse, error)
	SendWithExpectedStatus(request *PayloadRequest, expectedStatus ...int) (*PayloadResponse, error)
}

type PayloadData interface {
	Accept() string
	ContentType() string
	Content() ([]byte, error)
}

type PayloadRequest struct {
	Method string
	Url    string

	Headers map[string]string
	Payload PayloadData
}

type PayloadResponse struct {
	StatusCode int
	Body       string
}

func (p *PayloadResponse) Object(v any) error {
	return json.Unmarshal([]byte(p.Body), &v)
}
