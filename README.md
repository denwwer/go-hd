# hd
![Coverage](https://img.shields.io/badge/Coverage-100.0%25-brightgreen)
[![Go Reference](https://pkg.go.dev/badge/github.com/denwwer/go-hd.svg)](https://pkg.go.dev/github.com/denwwer/go-hd)

Humanize Duration (hd) – Go package that works like `Duration.String()` but returns a calendar-accurate difference (years, months, days, hours, minutes, seconds) for human-friendly output.

```
go get github.com/denwwer/go-hd
```

## Usage Example

### Since

Like `time.Since`, with nil location defaults to UTC.

```go
date := time.Date(2024, 3, 14, 10, 30, 0, 0, time.UTC) 

d := hd.Since(date, nil)
d.Years // => 1
d.Months // => 6
d.Days // => 8
d.String() // => 1y 6m 8d 5h 8m 17s
```

### Between

Between two specific times with nil location.

```go
start := time.Date(2022, 3, 14, 10, 0, 0, 0, time.UTC)
end := time.Date(2023, 3, 14, 12, 0, 0, 0, time.UTC)

d := Between(start, end, nil)
d.String() // => 1y 2h
```

With location.

```go
loc, _ := time.LoadLocation("America/New_York")
start := time.Date(2025, 9, 22, 0, 0, 0, 0, time.UTC)
end := time.Date(2025, 9, 23, 12, 30, 15, 0, time.UTC)

d := Between(start, end, loc)
d.String() // => 1d 12h 30m 15s
```
