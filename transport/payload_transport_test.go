package transport

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPayloadResponseJson(t *testing.T) {
	type Test struct {
		Key string `json:"key"`
	}

	response := PayloadResponse{
		StatusCode: 0,
		Body:       `{"key":"value"}`,
	}

	var test Test
	err := response.Object(&test)

	assert.Nil(t, err, "There should not have been an error")
	assert.Equal(t, "value", test.Key, "The key should have been value")
}
