package sch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArray(t *testing.T) {
	assert.NoError(t, Array(True()).Match([]any{true}))
	assert.NoError(t, Array(True()).Match([]any{true, true}))
	assert.Error(t, Array(True()).Match([]any{1}))
	assert.Error(t, Array(True()).Match(1))
	assert.Equal(t, `[]<span class=sch-base>true</span>`, genHTML(Array(True())))

	assert.NoError(t, ArraySize(1, True()).Match([]any{true}))
	assert.Error(t, ArraySize(1, True()).Match([]any{true, true}))
	assert.Error(t, ArraySize(1, True()).Match([]any{}))
	assert.Equal(t, `[1]<span class=sch-base>true</span>`, genHTML(ArraySize(1, True())))

	assert.NoError(t, ArrayMin(1, True()).Match([]any{true}))
	assert.NoError(t, ArrayMin(1, True()).Match([]any{true, true}))
	assert.Error(t, ArrayMin(1, True()).Match([]any{}))
	assert.Equal(t, `[1..]<span class=sch-base>true</span>`, genHTML(ArrayMin(1, True())))

	assert.NoError(t, ArrayRange(1, 2, True()).Match([]any{true}))
	assert.NoError(t, ArrayRange(1, 2, True()).Match([]any{true, true}))
	assert.Error(t, ArrayRange(1, 2, True()).Match([]any{true, true, true}))
	assert.Error(t, ArrayRange(1, 2, True()).Match([]any{true, false}))
	assert.Error(t, ArrayRange(1, 2, True()).Match([]any{}))
	assert.Equal(t, `[1..2]<span class=sch-base>true</span>`, genHTML(ArrayRange(1, 2, True())))
}
