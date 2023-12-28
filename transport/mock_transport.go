package transport

import (
	"errors"
	"fmt"
)

// MockPayloadTransport is a mock of the PayloadTransport, which stores payloads in memory instead of sending them over network
type MockPayloadTransport struct {
	// Stores the responses to the specific payloads
	Payloads map[string]PayloadResponse
}

// Send method tries to fetch the payload from the memory and returns it
func (m *MockPayloadTransport) Send(request *PayloadRequest) (*PayloadResponse, error) {
	// First, we try to get the response based on the URL provided in the request.
	response, ok := m.Payloads[request.Url]
	if ok {
		// If the URL matches, we return the response
		return &response, nil
	}

	// If the URL doesn't match, we try to get the response based on the request Method and URL.
	response, ok = m.Payloads[fmt.Sprintf("%s:%s", request.Method, request.Url)]
	if ok {
		// If the Method and URL matches, we return the response
		return &response, nil
	} else {
		// If none of the URL/Method+URL matches, we return an error
		return nil, errors.New("no payload for specified endpoint")
	}
}

// SendWithExpectedStatus method sends the payload and checks response status code
func (m *MockPayloadTransport) SendWithExpectedStatus(request *PayloadRequest, expectedStatus ...int) (*PayloadResponse, error) {
	// Sending the payload
	reply, err := m.Send(request)
	if err != nil {
		// If there is an error, return it immediately
		return nil, err
	}

	// Checking if the response status code in within the expected ones
	for _, v := range expectedStatus {
		if v == reply.StatusCode {
			return reply, nil
		}
	}

	// If response status code is not one of expected, return an error
	return reply, fmt.Errorf(reply.Body)
}

// Making sure MockPayloadTransport fully implements the PayloadTransport interface
var _ PayloadTransport = &MockPayloadTransport{}
