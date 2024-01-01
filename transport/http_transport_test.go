package transport

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTTPTransport(t *testing.T) {
	// Create a local http server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte("OK"))
	}))

	// Close the server when test finishes
	defer server.Close()

	// Instantiate your HTTP transport instance
	transport := NewHttpPayloadTransport(server.URL, BasicAuthentication{
		Username: "",
		Password: "",
	})

	// Call your transport function
	send, err := transport.Send(&PayloadRequest{
		Method: "GET",
		Url:    "/test",
	})

	assert.Nil(t, err, "got error: %v", err)
	assert.Equal(t, http.StatusOK, send.StatusCode, "got status code: %d, want: %d")
	assert.Equal(t, "OK", send.Body, "got body: %s, want: OK")
}

func TestHTTPTransportSendWithExpectedStatus(t *testing.T) {
	// Create a local http server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if strings.Contains(req.URL.Path, "ok") {
			rw.WriteHeader(http.StatusOK)
			_, _ = rw.Write([]byte("OK"))
		} else if strings.Contains(req.URL.Path, "created") {
			rw.WriteHeader(http.StatusCreated)
			_, _ = rw.Write([]byte("OK"))
		} else if strings.Contains(req.URL.Path, "moved") {
			rw.WriteHeader(http.StatusMovedPermanently)
			_, _ = rw.Write([]byte("OK"))
		} else {
			rw.WriteHeader(http.StatusBadRequest)
		}
	}))
	// Close the server when test finishes
	defer server.Close()

	// Instantiate your HTTP transport instance
	transport := NewHttpPayloadTransport(server.URL, BasicAuthentication{
		Username: "",
		Password: "",
	})

	// Call your transport function
	reply, err := transport.SendWithExpectedStatus(&PayloadRequest{
		Method: "GET",
		Url:    "/ok",
	}, http.StatusOK, http.StatusCreated)

	assert.Nil(t, err, "got error: %v", err)
	assert.Equal(t, "OK", reply.Body, "got body: %s, want: OK", reply.Body)

	reply, err = transport.SendWithExpectedStatus(&PayloadRequest{
		Method: "GET",
		Url:    "/created",
	}, http.StatusOK, http.StatusCreated)

	assert.Nil(t, err, "got error: %v", err)
	assert.Equal(t, http.StatusCreated, reply.StatusCode)

	reply, err = transport.SendWithExpectedStatus(&PayloadRequest{
		Method: "GET",
		Url:    "/moved",
	}, http.StatusOK)

	assert.NotNil(t, err, "got error: %v", err)
	assert.Equal(t, http.StatusMovedPermanently, reply.StatusCode)

	reply, err = transport.SendWithExpectedStatus(&PayloadRequest{
		Method: "GET",
		Url:    "/error",
	}, http.StatusOK, http.StatusCreated)

	assert.NotNil(t, err, "expected error, got nil")
	assert.Equal(t, http.StatusBadRequest, reply.StatusCode, "got status code: %d, want: %d", reply.StatusCode, http.StatusBadRequest)
}

func TestHTTPTransportSendWithBody(t *testing.T) {
	// Create a local http server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		body, _ := io.ReadAll(req.Body)
		bodyString := string(body)

		assert.Equal(t, "\"OK\"", bodyString)

		rw.WriteHeader(http.StatusOK)
	}))
	// Close the server when test finishes
	defer server.Close()

	// Instantiate your HTTP transport instance
	transport := NewHttpPayloadTransport(server.URL, BasicAuthentication{
		Username: "",
		Password: "",
	})

	// Call your transport function
	reply, _ := transport.SendWithExpectedStatus(&PayloadRequest{
		Method:  "GET",
		Url:     "/ok",
		Payload: JsonPayloadData{Payload: "OK"},
	}, http.StatusOK)
	assert.Equal(t, http.StatusOK, reply.StatusCode, "got status: %s, want: OK", reply.StatusCode)
}

func TestHTTPTransportBasicAuthentication(t *testing.T) {
	// Create a local http server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")
		assert.Equal(t, "Basic dXNlcm5hbWU6cGFzc3dvcmQ=", auth)

		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte("OK"))
	}))

	defer server.Close()

	transport := NewHttpPayloadTransport(server.URL, BasicAuthentication{
		Username: "username",
		Password: "password",
	})

	_, _ = transport.Send(&PayloadRequest{
		Method: "GET",
		Url:    "/test",
		Headers: map[string]string{
			"Additional": "Header",
		},
	})
}

func TestHTTPTransportBearerAuthentication(t *testing.T) {
	// Create a local http server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")
		assert.Equal(t, "Bearer mytoken123", auth)

		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte("OK"))
	}))

	defer server.Close()

	transport := NewHttpPayloadTransport(server.URL, BearerAuthentication{
		Token: "mytoken123",
	})

	_, _ = transport.Send(&PayloadRequest{
		Method: "GET",
		Url:    "/test",
		Headers: map[string]string{
			"Additional": "Header",
		},
	})
}

func TestHTTPTransportAdditionaHeaders(t *testing.T) {
	// Create a local http server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Additional")
		assert.Equal(t, "Header", auth)

		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte("OK"))
	}))

	defer server.Close()

	transport := NewHttpPayloadTransport(server.URL, BearerAuthentication{
		Token: "mytoken123",
	})

	_, _ = transport.Send(&PayloadRequest{
		Method: "GET",
		Url:    "/test",
		Headers: map[string]string{
			"Additional": "Header",
		},
	})
}
