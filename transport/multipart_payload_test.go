package transport

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestMultipartPayload_Accept(t *testing.T) {
	payload := &MultipartPayload{}
	assert.Equal(t, "application/json", payload.Accept())
}

func TestMultipartPayload_ContentType(t *testing.T) {
	payload := &MultipartPayload{}

	expected := fmt.Sprintf("multipart/form-data; boundary=%s", payload.Boundary)
	assert.Equal(t, expected, payload.ContentType())
}

func TestMultipartPayload_FormKeyValue(t *testing.T) {
	payload := &MultipartPayload{
		Form: map[string]string{
			"testKey": "testValue",
		},
	}

	buf, err := payload.Content()
	assert.Nil(t, err)

	bodyStr := string(buf)

	expectedBodyStr := fmt.Sprintf(`--%s
Content-Disposition: form-data; name="testKey"

%s
--%s--
`,
		payload.Boundary, "testValue", payload.Boundary)
	assert.Equal(t, strings.Replace(expectedBodyStr, "\r", "", -1), strings.Replace(bodyStr, "\r", "", -1))
}

func TestMultipartPayload_MultipartFile(t *testing.T) {
	payload := &MultipartPayload{
		File: &MultipartFile{
			Key:     "Key",
			Name:    "Name",
			Content: "test file content",
		},
	}

	buf, err := payload.Content()
	assert.Nil(t, err)

	bodyStr := string(buf)

	expectedBodyStr := fmt.Sprintf(`--%s
Content-Disposition: form-data; name="Key"; filename="Name"
Content-Type: application/octet-stream

test file content
--%s--
`,
		payload.Boundary, payload.Boundary)
	assert.Equal(t, strings.Replace(expectedBodyStr, "\r", "", -1), strings.Replace(bodyStr, "\r", "", -1))
}
