package sch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexp(t *testing.T) {
	assert.NoError(t, Regexp(`^ECI\(\d{4}\)\d{6}$`).Match("ECI(2024)000008"))
	assert.Error(t, Regexp(`^ECI\(\d{4}\)\d{6}$`).Match("ECI(2024)0000008"))
	assert.Error(t, Regexp(`^ECI\(\d{4}\)\d{6}$`).Match(1))
	assert.Equal(t, `<span class=sch-xstr>regexp/<u>^ECI\(\d{4}\)\d{6}$</u>/</span>`, genHTML(Regexp(`^ECI\(\d{4}\)\d{6}$`)))
}

func TestTime(t *testing.T) {
	assert.NoError(t, Time("2006-01-02 15:04:05").Match("2024-10-31 00:45:39"))
	assert.Error(t, Time("2006-01-02 15:04:05").Match("2024-10-31T00:45:39"))
	assert.Error(t, Time("2006-01-02 15:04:05").Match(1))

	assert.Equal(t, `<span class=sch-xstr title="A time value encoded into a string">string(<u>2006-01-02 15:04:05</u>)</span>`,
		genHTML(Time("2006-01-02 15:04:05")))
}
