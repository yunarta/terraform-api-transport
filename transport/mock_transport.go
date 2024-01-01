package transport

import (
	"fmt"
)

// MockPayloadTransport is a mock of the PayloadTransport, which stores payloads in memory instead of sending them over network
type MockPayloadTransport struct {
	// Stores the responses to the specific payloads
	Payloads     map[string]PayloadResponse
	SentRequests map[string][]PayloadRequest
}

// Send method tries to fetch the payload from the memory and returns it
func (m *MockPayloadTransport) Send(request *PayloadRequest) (*PayloadResponse, error) {
	response, err := m.findResponse(request)
	if err == nil {
		if m.SentRequests == nil {
			m.SentRequests = make(map[string][]PayloadRequest)
		}

		key := fmt.Sprintf("%s:%s", request.Method, request.Url)
		if sentRequest, ok := m.SentRequests[key]; ok {
			sentRequest = append(sentRequest, *request)
			m.SentRequests[key] = sentRequest
		} else {
			m.SentRequests[key] = []PayloadRequest{*request}
		}
	}

	return response, err
}

func (m *MockPayloadTransport) findResponse(request *PayloadRequest) (*PayloadResponse, error) {
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
		return nil, MockRequestError{
			error: fmt.Sprintf("no payload for specified endpoint %s:%s", request.Method, request.Url),
			path:  fmt.Sprintf("%s:%s", request.Method, request.Url),
		}
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
	return handleResponseException(reply)
}

// Making sure MockPayloadTransport fully implements the PayloadTransport interface
var _ PayloadTransport = &MockPayloadTransport{}
