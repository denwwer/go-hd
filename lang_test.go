package hd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormIndex(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		rule   pluralRule
		count  int
		expect int
	}{
		{"none", pluralNone, 5, 0},
		{"one-many singular", pluralOneMany, 1, 0},
		{"one-many plural", pluralOneMany, 2, 1},
		{"slavic many (teens)", pluralSlavic, 12, 0},
		{"slavic many (5-9)", pluralSlavic, 27, 0},
		{"slavic many (0)", pluralSlavic, 30, 0},
		{"slavic one", pluralSlavic, 21, 1},
		{"slavic few", pluralSlavic, 3, 2},
		{"arabic one", pluralArabic, 1, 0},
		{"arabic two", pluralArabic, 2, 1},
		{"arabic few", pluralArabic, 7, 2},
		{"arabic many", pluralArabic, 15, 0},
		{"czech one", pluralCzechOrSlovak, 1, 0},
		{"czech few", pluralCzechOrSlovak, 3, 2},
		{"czech many", pluralCzechOrSlovak, 15, 3},
		{"french one", pluralFrench, 1, 0},
		{"french plural", pluralFrench, 2, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expect, formIndex(tc.rule, tc.count))
		})
	}
}

func TestUnit_WordClampsMissingForms(t *testing.T) {
	t.Parallel()

	// rule asks for index 1, but only one form exists
	u := unit{[]string{"day"}}
	assert.Equal(t, "day", u.word(pluralOneMany, 5))
}

func TestLanguage_WordFirst(t *testing.T) {
	t.Parallel()
	l := language{wordFirst: true, rule: pluralOneMany}
	u := unit{[]string{"mwaka", "miaka"}}

	assert.Equal(t, "miaka 2", l.string(2, u))
}
