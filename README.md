# hd
![Coverage](https://img.shields.io/badge/Coverage-99.0%25-brightgreen)
Humanize Duration (hd) – Go package that works like `Duration.String()` but returns a calendar-accurate difference (years, months, weeks, days, hours, minutes, seconds) for human-friendly output.

## Usage Example

### Since

Like `time.Since`, defaults to UTC.

```go
date := time.Date(2024, 3, 14, 10, 30, 0, 0, time.UTC) 

d := hd.Since(date)
d.Years // => 1
d.Months // => 6
d.Weeks // => 1
d.Days // => 1
d.String() // => 1y 6m 1w 1d 5h 8m 17s
```

### Between

Between two specific times.

```go
start := time.Date(2022, 3, 14, 10, 0, 0, 0, time.UTC)
end := time.Date(2023, 3, 14, 12, 0, 0, 0, time.UTC)

d := hd.Between(start, end)
d.String() // => 1y 2h
```

### Options

Configure via options; pass zero or more, later values win: `hd.Between(start, end, hd.Location(loc), hd.Language(hd.LangDE))`.

With location:

```go
loc, _ := time.LoadLocation("America/New_York")
start := time.Date(2025, 9, 22, 0, 0, 0, 0, time.UTC)
end := time.Date(2025, 9, 23, 12, 30, 15, 0, time.UTC)

d := hd.Between(start, end, hd.Location(loc))
d.String() // => 1d 12h 30m 15s
```
With custom delimiter:

```go
d := hd.Between(start, end, hd.Language(hd.LangEN), hd.Delimiter(" / "))
d.String() // => 1 week / 1 day / 1 hour
```

With language:

```go
d := hd.Between(start, end, hd.Language(hd.LangDE))
d.String() // => 3 Jahre, 3 Monate, 2 Wochen, 6 Stunden, 30 Minuten, 15 Sekunden
```

#### Supported languages

| Code | Language | Code | Language |
|------| --- | --- | --- |
| `ar` | Arabic | `ko` | Korean |
| `cs` | Czech | `nl` | Dutch |
| `de` | German | `no` | Norwegian |
| `en` | English | `sk` | Slovak |
| `es` | Spanish | `tr` | Turkish |
| `fr` | French | `uk` | Ukrainian |
| `hi` | Hindi | `zh` | Chinese (Simplified) |
| `it` | Italian | `zh-Hant` | Chinese (Traditional) |
| `ja` | Japanese | | |

Don't see your language? PRs adding a new language to [languages](lang.go) are very welcome.
