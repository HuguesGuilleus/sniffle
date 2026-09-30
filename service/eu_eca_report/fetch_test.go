package eu_eca_report

import (
	"testing"

	"github.com/HuguesGuilleus/sniffle/common/language"
	"github.com/stretchr/testify/assert"
)

func TestMergeReports(t *testing.T) {
	merge := mergeReports([]report{
		{ID: "1", L: [26]*reportL{language.English: {Title: "1.en"}}},
		{ID: "2", L: [26]*reportL{language.English: {Title: "2.en"}}},
		{ID: "4", L: [26]*reportL{language.English: {Title: "4.en"}}},
	}, []report{
		{ID: "1", L: [26]*reportL{language.French: {Title: "1.fr"}}},
		{ID: "3", L: [26]*reportL{language.French: {Title: "3.fr"}}},
		{ID: "4", L: [26]*reportL{language.French: {Title: "4.fr"}}},
		{ID: "5", L: [26]*reportL{language.French: {Title: "5.fr"}}},
	})
	assert.Equal(t, []report{
		{ID: "1", L: [26]*reportL{
			language.English: {Title: "1.en"},
			language.French:  {Title: "1.fr"},
		}},
		{ID: "2", L: [26]*reportL{language.English: {Title: "2.en"}}},
		{ID: "3", L: [26]*reportL{language.French: {Title: "3.fr"}}},
		{ID: "4", L: [26]*reportL{
			language.English: {Title: "4.en"},
			language.French:  {Title: "4.fr"},
		}},
		{ID: "5", L: [26]*reportL{language.French: {Title: "5.fr"}}},
	}, merge)
}
