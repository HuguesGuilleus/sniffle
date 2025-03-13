package sch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssertKey(t *testing.T) {
	name, f := AssertKey("q", func(a any) int { return a.(int) })
	assert.Equal(t, "[*].q is unique", name)
	assert.NoError(t, f(nil, []any{1, 2}))
	assert.Error(t, f(nil, []any{1, 1}))
}

func TestAssertOnlyOneTrue(t *testing.T) {
	name, f := AssertOnlyOneTrue(`q`, func(a any) bool { return a.(bool) })
	assert.Equal(t, "exact one [$].q is true", name)
	assert.NoError(t, f(nil, []any{true, false, false}))
	assert.Error(t, f(nil, []any{false, false, false}))
	assert.Error(t, f(nil, []any{true, true, false}))
}
