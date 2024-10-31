package sch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArray(t *testing.T) {
	assert.NoError(t, Array(True()).Match([]any{true}))
	assert.Error(t, Array(True()).Match([]any{1}))
	assert.Error(t, Array(True()).Match(1))

	assert.Equal(t, `[]<span class=sch-base>true</span>`, genHTML(Array(True())))
}
