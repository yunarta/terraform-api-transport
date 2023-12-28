package transport

import (
	"net/http"
	"testing"
)

func TestMockTransport(t *testing.T) {
	mockTransport := MockPayloadTransport{
		Payloads: map[string]PayloadResponse{
			"/rest/api/latest/deploy/project": {
				StatusCode: 200,
				Body:       "OK",
			},
			"POST:/rest/api/latest/deploy": {
				StatusCode: 200,
				Body:       "OK",
			},
			"PUT:/rest/api/latest/error": {
				StatusCode: 400,
				Body:       "OK",
			},
		},
	}

	send, _ := mockTransport.Send(&PayloadRequest{
		Url: "/rest/api/latest/deploy/project",
	})

	if send.Body != "OK" {
		t.Error("Body was not 'OK'")
	}

	send, _ = mockTransport.Send(&PayloadRequest{
		Method: http.MethodPost,
		Url:    "/rest/api/latest/deploy",
	})

	if send.Body != "OK" {
		t.Error("Body was not 'OK'")
	}

	_, err := mockTransport.SendWithExpectedStatus(&PayloadRequest{
		Method: http.MethodPut,
		Url:    "/rest/api/latest/error",
	}, 200)

	if err == nil {
		t.Error("Error should be trigger")
	}

	_, err = mockTransport.SendWithExpectedStatus(&PayloadRequest{
		Method: http.MethodPut,
		Url:    "/none",
	}, 200)

	if err == nil {
		t.Error("Error should be trigger")
	}
}
