package translate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLangEnglish(t *testing.T) {
	l := LangEnglish{}
	date := time.Date(2026, time.September, 25, 23, 31, 19, 0, time.FixedZone("UTC+2", 2*60*60))
	assert.EqualValues(t,
		`<time datetime=2026-09-25T23:31:19Z>Friday, 25 September 2026 23:31:19 UTC+2</time>`,
		l.DateHourLong(date),
	)
	assert.EqualValues(t, `<time datetime=2026-09-25>2026-09-25</time>`, l.DateShort(date))
	assert.EqualValues(t, `<time datetime=2026-09-25>September 25, 2026</time>`, l.DateLong(date))
}
