package hd

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helper to create times in UTC quickly
func utc(y int, m time.Month, d, h, min, s int) time.Time {
	return time.Date(y, m, d, h, min, s, 0, time.UTC)
}

func TestBetween(t *testing.T) {
	t.Parallel()

	start := utc(2020, 1, 1, 0, 0, 0)
	end := utc(2023, 4, 15, 6, 30, 15)

	d := Between(start, end)
	assert.Equal(t, 3, d.Years)
	assert.Equal(t, 3, d.Months)
	assert.Equal(t, 2, d.Weeks)
	assert.Equal(t, 0, d.Days)
	assert.Equal(t, 6, d.Hours)
	assert.Equal(t, 30, d.Minutes)
	assert.Equal(t, 15, d.Seconds)
}

func TestBetween_WeeksAndDays(t *testing.T) {
	t.Parallel()

	// 9 days apart => 1 week 2 days
	start := utc(2024, 5, 1, 0, 0, 0)
	end := utc(2024, 5, 10, 0, 0, 0)

	d := Between(start, end)
	assert.Equal(t, 1, d.Weeks)
	assert.Equal(t, 2, d.Days)
	assert.Equal(t, "1w 2d", d.String())
}

func TestBetween_DaysUnderWeek(t *testing.T) {
	t.Parallel()

	// 6 days apart => no weeks, 6 days
	start := utc(2024, 5, 1, 0, 0, 0)
	end := utc(2024, 5, 7, 0, 0, 0)

	d := Between(start, end)
	assert.Equal(t, 0, d.Weeks)
	assert.Equal(t, 6, d.Days)
	assert.Equal(t, "6d", d.String())
}

func TestBetween_WithLocation(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	require.Nil(t, err)

	start := utc(2025, 9, 22, 0, 0, 0)
	end := utc(2025, 9, 23, 12, 30, 15)

	d := Between(start, end, Location(loc))
	assert.Equal(t, 0, d.Years)
	assert.Equal(t, 0, d.Months)
	assert.Equal(t, 1, d.Days)
	assert.Equal(t, 12, d.Hours)
	assert.Equal(t, 30, d.Minutes)
	assert.Equal(t, 15, d.Seconds)
}

func TestBetween_ReverseOrder(t *testing.T) {
	t.Parallel()
	// Start is after end; should still produce positive difference
	start := utc(2025, 1, 1, 0, 0, 0)
	end := utc(2023, 1, 1, 0, 0, 0)

	d := Between(start, end)
	assert.Equal(t, 2, d.Years)
}

func TestBetween_LangZeroDuration(t *testing.T) {
	t.Parallel()

	// start == end
	start := utc(2024, 5, 1, 0, 0, 0)

	d := Between(start, start, Language(LangEN))
	assert.Equal(t, "0 seconds", d.String())
}

func TestBetween_LeapYear(t *testing.T) {
	t.Parallel()
	// Feb 29, 2020 → Feb 28, 2021 should be 1 year minus 1 day
	start := utc(2020, 2, 29, 0, 0, 0)
	end := utc(2021, 2, 28, 0, 0, 0)

	d := Between(start, end, Location(time.UTC))
	assert.Equal(t, 0, d.Years)
	assert.Equal(t, 11, d.Months)
}

func TestBetween_Lang(t *testing.T) {
	t.Parallel()

	start := utc(2020, 1, 1, 0, 0, 0)
	end := utc(2023, 4, 15, 6, 30, 15) // 3y 3m 2w 6h 30m 15s

	testCases := []struct {
		lang   Lang
		expect string
	}{
		{LangEN, "3 years, 3 months, 2 weeks, 6 hours, 30 minutes, 15 seconds"},
		{LangDE, "3 Jahre, 3 Monate, 2 Wochen, 6 Stunden, 30 Minuten, 15 Sekunden"},
		{LangJA, "3 年, 3 ヶ月, 2 週間, 6 時間, 30 分, 15 秒"},
	}

	for _, tc := range testCases {
		t.Run(strconv.Itoa(int(tc.lang)), func(t *testing.T) {
			t.Parallel()
			d := Between(start, end, Language(tc.lang))
			assert.Equal(t, tc.expect, d.String())
		})
	}
}

func TestBetween_LangEnDelimiter(t *testing.T) {
	t.Parallel()

	// en has no delimiter -> space-joined; singular forms for count 1
	start := utc(2024, 5, 1, 0, 0, 0)
	end := utc(2024, 5, 9, 1, 0, 0)

	d := Between(start, end, Language(LangEN), Delimiter(" / "))
	assert.Equal(t, "1 week / 1 day / 1 hour", d.String())
}

func TestBetween_LangArabic(t *testing.T) {
	t.Parallel()

	// count 2 hides the number (dual form), digits are replaced
	start := utc(2024, 5, 1, 0, 0, 0)
	end := utc(2024, 5, 15, 3, 0, 0)

	d := Between(start, end, Language(LangAR))
	assert.Equal(t, "أسبوعين ٣ ساعات", d.String())
}

func TestBetween_LangUnknownFallsBack(t *testing.T) {
	t.Parallel()

	start := utc(2024, 5, 1, 0, 0, 0)
	end := utc(2024, 5, 10, 0, 0, 0)

	d := Between(start, end, Language(Lang(33)))
	assert.Equal(t, "1w 2d", d.String())
}

func TestBetween_LangCzech(t *testing.T) {
	t.Parallel()

	// 3 days => few form "dny"
	start := utc(2024, 5, 1, 0, 0, 0)
	end := utc(2024, 5, 4, 0, 0, 0)

	d := Between(start, end, Language(LangCS))
	assert.Equal(t, "3 dny", d.String())
}

func TestBetween_LangFrench(t *testing.T) {
	t.Parallel()

	start := utc(2024, 5, 1, 0, 0, 0)
	end := utc(2024, 5, 10, 2, 0, 0)

	d := Between(start, end, Language(LangFR))
	assert.Equal(t, "1 semaine, 2 jours, 2 heures", d.String())
}

func TestSince(t *testing.T) {
	t.Parallel()

	d := Since(time.Now().AddDate(0, 0, -2).Add(-time.Hour))
	assert.Equal(t, "2d 1h", d.String())
}

func TestDuration_String(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		dur    Duration
		expect string
	}{
		{
			name:   "All non-zero",
			dur:    Duration{Years: 1, Months: 2, Weeks: 3, Days: 4, Hours: 5, Minutes: 6, Seconds: 7},
			expect: "1y 2mo 3w 4d 5h 6min 7s",
		},
		{
			name:   "Some zeros",
			dur:    Duration{Years: 0, Months: 0, Days: 3, Hours: 0, Minutes: 5, Seconds: 0},
			expect: "3d 5min",
		},
		{
			name:   "All zeros",
			dur:    Duration{},
			expect: "0s",
		},
		{
			name:   "Single unit",
			dur:    Duration{Hours: 7},
			expect: "7h",
		},
		{
			name:   "weeks only",
			dur:    Duration{Weeks: 2},
			expect: "2w",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expect, tc.dur.String())
		})
	}
}

func TestDuration_MarshalJSON(t *testing.T) {
	t.Parallel()

	start := utc(2020, 1, 1, 0, 0, 0)
	end := utc(2023, 4, 15, 6, 30, 15)

	data := struct {
		Duration Duration `json:"duration"`
	}{
		Duration: Between(start, end),
	}

	b, err := json.Marshal(data)
	require.Nil(t, err)
	assert.Equal(t, `{"duration":"3y 3mo 2w 6h 30min 15s"}`, string(b))
}
