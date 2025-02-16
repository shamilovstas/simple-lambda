package transport

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessage_String(t *testing.T) {
	message := NewMessage("Hello")
	expected := "{\"message\":\"Hello\"}"
	actual := message.String()
	assert.Equal(t, expected, actual)
}
