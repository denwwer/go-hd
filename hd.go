// Humanize Duration (hd) – Go package that works like Duration.String() but returns
// a calendar-accurate difference (years, months, weeks, days, hours, minutes, seconds) for human-friendly output.
package hd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type options struct {
	loc       *time.Location
	lang      Lang
	delimiter string
}

// Option for customizing formatting.
type Option func(*options)

func defaultOptions() *options {
	return &options{
		loc:       time.UTC,
		delimiter: ", ",
	}
}

// Location sets the time zone used for the calendar calculation.
func Location(loc *time.Location) Option {
	return func(c *options) {
		c.loc = loc
	}
}

// Language sets the language, e.g. "de", "es".
func Language(l Lang) Option {
	return func(c *options) {
		c.lang = l
	}
}

// Delimiter sets the delimiter used for formatting date.
func Delimiter(d string) Option {
	return func(c *options) {
		c.delimiter = d
	}
}

// Duration represents the difference in calendar units.
type Duration struct {
	Years, Months, Weeks, Days int
	Hours, Minutes, Seconds    int

	*options
}

// String returns a string representing the duration in format "1y 6m 2w 1d 5h 8m 17s".
func (d Duration) String() string {
	if d.options != nil {
		if lang, ok := languages[d.lang]; ok {
			return d.localize(lang)
		}
	}

	var s strings.Builder

	write := func(format string, v int) {
		if v != 0 {
			if s.Len() > 0 {
				s.WriteByte(' ')
			}
			fmt.Fprintf(&s, format, v)
		}
	}

	write("%dy", d.Years)
	write("%dmo", d.Months)
	write("%dw", d.Weeks)
	write("%dd", d.Days)
	write("%dh", d.Hours)
	write("%dmin", d.Minutes)
	write("%ds", d.Seconds)

	if s.Len() == 0 {
		return "0s"
	}
	return s.String()
}

// MarshalJSON implements json.Marshaler for Duration.
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// Since returns a calendar-accurate duration from t to time.Now().
func Since(t time.Time, opts ...Option) Duration {
	return Between(t, time.Now(), opts...)
}

// Between calculates the calendar-accurate duration between any two times.
func Between(start, end time.Time, opts ...Option) Duration {
	opt := defaultOptions()
	for _, o := range opts {
		o(opt)
	}

	start = start.In(opt.loc)
	end = end.In(opt.loc)

	if start.After(end) {
		start, end = end, start
	}

	years, months, days := 0, 0, 0

	for !start.AddDate(years+1, 0, 0).After(end) {
		years++
	}
	for !start.AddDate(years, months+1, 0).After(end) {
		months++
	}
	for !start.AddDate(years, months, days+1).After(end) {
		days++
	}

	partial := end.Sub(start.AddDate(years, months, days))
	hours := int(partial.Hours())
	minutes := int(partial.Minutes()) % 60
	seconds := int(partial.Seconds()) % 60

	weeks := days / 7
	days %= 7

	return Duration{
		Years:   years,
		Months:  months,
		Weeks:   weeks,
		Days:    days,
		Hours:   hours,
		Minutes: minutes,
		Seconds: seconds,

		options: opt,
	}
}
