package transport

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// BasicAuthentication represents username/password credentials
type BasicAuthentication struct {
	Username string
	Password string
}

// Ensure BasicAuthentication implements Authentication interface
var _ Authentication = &BasicAuthentication{}

// BearerAuthentication represents bearer token credentials
type BearerAuthentication struct {
	Token string
}

// Ensure BearerAuthentication implements Authentication interface
var _ Authentication = &BearerAuthentication{}

// HttpPayloadTransport is a transport layer using HTTP protocol
type HttpPayloadTransport struct {
	baseUrl        string
	authentication Authentication
}

// NewHttpPayloadTransport creates a new HTTP transport instance
func NewHttpPayloadTransport(baseUrl string, authentication Authentication) *HttpPayloadTransport {
	return &HttpPayloadTransport{
		baseUrl:        baseUrl,
		authentication: authentication,
	}
}

// Ensure HttpPayloadTransport implements PayloadTransport interface
var _ PayloadTransport = &HttpPayloadTransport{}

// Send is used send http request and returns http response
func (h *HttpPayloadTransport) Send(request *PayloadRequest) (*PayloadResponse, error) {
	// Start tracking request time
	startTime := time.Now()

	// Calculating the time taken by request at end of function execution
	defer func() {
		duration := time.Since(startTime)
		milliseconds := duration.Milliseconds()
		fmt.Printf("HTTPX %s %s, time = %d\n", request.Method, request.Url, milliseconds)
	}()

	var body io.Reader

	// If Payload is not nil, read content to Reader
	if request.Payload != nil {
		content, err := request.Payload.Content()
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(content)
	}

	// Creating http request
	// #nosec G107 - low level api transport
	httpRequest, err := http.NewRequest(request.Method, h.baseUrl+request.Url, body)
	if err != nil {
		return nil, err
	}

	// If authentication exists, set it in http Request
	if h.authentication != nil {
		switch h.authentication.(type) {
		case BasicAuthentication:
			authentication := h.authentication.(BasicAuthentication)
			httpRequest.SetBasicAuth(authentication.Username, authentication.Password)
		case BearerAuthentication:
			authentication := h.authentication.(BearerAuthentication)
			httpRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authentication.Token))
		}
	}

	// Setting headers in http request
	for key, value := range request.Headers {
		httpRequest.Header.Set(key, value)
	}

	// If Payload exists, set content type and Accept Header
	// else set default Accept Header
	if request.Payload != nil {
		httpRequest.Header.Set("Content-Type", request.Payload.ContentType())
		httpRequest.Header.Set("Accept", request.Payload.Accept())
	} else {
		httpRequest.Header.Set("Accept", "application/json")
	}

	// Send http request using default http client
	client := http.DefaultClient
	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	// Read content from http response
	content, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	// Create PayloadResponse using status and body from http response
	return &PayloadResponse{
		StatusCode: httpResponse.StatusCode,
		Body:       string(content),
	}, err
}

// SendWithExpectedStatus sends http request and verifies the response's status
func (h *HttpPayloadTransport) SendWithExpectedStatus(request *PayloadRequest, expectedStatus ...int) (*PayloadResponse, error) {

	// Send http request and get response
	reply, err := h.Send(request)
	if err != nil {
		return nil, err
	}

	// Verify if the response status is as expected
	for _, v := range expectedStatus {
		if v == reply.StatusCode {
			return reply, nil
		}
	}

	// If status is different from expected return error
	return reply, fmt.Errorf(reply.Body)
}
