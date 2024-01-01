package transport

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonPayloadData_ContentMust(t *testing.T) {
	test := JsonPayloadData{
		Payload: "test",
	}
	assert.Equal(t, []byte("\"test\""), test.ContentMust())
}
