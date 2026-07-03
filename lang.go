package hd

import (
	"strconv"
	"strings"
)

type pluralRule int

const (
	pluralNone pluralRule = iota
	pluralOneMany
	pluralSlavic
	pluralArabic
	pluralCzechOrSlovak
	pluralFrench
)

type unit struct {
	forms []string
}

func (u unit) word(rule pluralRule, c int) string {
	i := formIndex(rule, c)
	if i >= len(u.forms) {
		i = len(u.forms) - 1
	}

	return u.forms[i]
}

// plural form index:
//
//	pluralNone          [word]
//	pluralOneMany       [one, many]
//	pluralSlavic        [many, one, few]
//	pluralArabic        [one, two, few]
//	pluralCzechOrSlovak [one, fractional, few, many]
//	pluralFrench        [one, many]
func formIndex(rule pluralRule, c int) int {
	switch rule {
	case pluralOneMany:
		if c == 1 {
			return 0
		}
		return 1
	case pluralSlavic:
		mod10, mod100 := c%10, c%100
		switch {
		case (mod100 >= 5 && mod100 <= 20) || (mod10 >= 5 && mod10 <= 9) || mod10 == 0:
			return 0
		case mod10 == 1:
			return 1
		default:
			return 2
		}
	case pluralArabic:
		switch {
		case c == 2:
			return 1
		case c > 2 && c < 11:
			return 2
		default:
			return 0
		}
	case pluralCzechOrSlovak:
		mod10, mod100 := c%10, c%100
		switch {
		case c == 1:
			return 0
		case mod10 >= 2 && mod10 <= 4 && mod100 < 10:
			return 2
		default:
			return 3
		}
	case pluralFrench:
		if c >= 2 {
			return 1
		}
		return 0
	default: // pluralNone
		return 0
	}
}

func (d Duration) localize(l language) string {
	pieces := []struct {
		count int
		unit  unit
	}{
		{d.Years, l.years},
		{d.Months, l.months},
		{d.Weeks, l.weeks},
		{d.Days, l.days},
		{d.Hours, l.hours},
		{d.Minutes, l.minutes},
		{d.Seconds, l.seconds},
	}

	var parts []string
	for _, p := range pieces {
		if p.count != 0 {
			parts = append(parts, l.string(p.count, p.unit))
		}
	}

	if len(parts) == 0 {
		return l.string(0, l.seconds)
	}

	if len(parts) <= 2 {
		return strings.Join(parts, " ")
	}

	return strings.Join(parts, d.delimiter)
}

type language struct {
	years   unit
	months  unit
	weeks   unit
	days    unit
	hours   unit
	minutes unit
	seconds unit

	hideCountIfTwo   bool
	wordFirst        bool
	digitReplacement []string
	rule             pluralRule
}

func (l language) string(c int, u unit) string {
	word := u.word(l.rule, c)

	if l.hideCountIfTwo && c == 2 {
		return word
	}

	count := strconv.Itoa(c)
	if len(l.digitReplacement) == 10 {
		var s strings.Builder
		for _, r := range count {
			s.WriteString(l.digitReplacement[r-'0'])
		}
		count = s.String()
	}

	if l.wordFirst {
		return word + " " + count
	}
	return count + " " + word
}

// Lang represents a language code.
type Lang int

// List of support languages.
const (
	LangEN Lang = iota
	LangDE
	LangAR
	LangCS
	LangES
	LangFR
	LangIT
	LangJA
	LangKO
	LangUK
	LangSK
	LangTR
	LangZH
	LangZHHant
	LangNO
	LangNL
	LangHI
)

var languages = map[Lang]language{
	LangEN: {
		years:   unit{[]string{"year", "years"}},
		months:  unit{[]string{"month", "months"}},
		weeks:   unit{[]string{"week", "weeks"}},
		days:    unit{[]string{"day", "days"}},
		hours:   unit{[]string{"hour", "hours"}},
		minutes: unit{[]string{"minute", "minutes"}},
		seconds: unit{[]string{"second", "seconds"}},
		rule:    pluralOneMany,
	},
	LangDE: {
		years:   unit{[]string{"Jahr", "Jahre"}},
		months:  unit{[]string{"Monat", "Monate"}},
		weeks:   unit{[]string{"Woche", "Wochen"}},
		days:    unit{[]string{"Tag", "Tage"}},
		hours:   unit{[]string{"Stunde", "Stunden"}},
		minutes: unit{[]string{"Minute", "Minuten"}},
		seconds: unit{[]string{"Sekunde", "Sekunden"}},
		rule:    pluralOneMany,
	},
	LangAR: {
		years:            unit{[]string{"سنة", "سنتان", "سنوات"}},
		months:           unit{[]string{"شهر", "شهران", "أشهر"}},
		weeks:            unit{[]string{"أسبوع", "أسبوعين", "أسابيع"}},
		days:             unit{[]string{"يوم", "يومين", "أيام"}},
		hours:            unit{[]string{"ساعة", "ساعتين", "ساعات"}},
		minutes:          unit{[]string{"دقيقة", "دقيقتان", "دقائق"}},
		seconds:          unit{[]string{"ثانية", "ثانيتان", "ثواني"}},
		hideCountIfTwo:   true,
		digitReplacement: []string{"۰", "١", "٢", "٣", "٤", "٥", "٦", "٧", "٨", "٩"},
		rule:             pluralArabic,
	},
	LangCS: {
		years:   unit{[]string{"rok", "roku", "roky", "let"}},
		months:  unit{[]string{"měsíc", "měsíce", "měsíce", "měsíců"}},
		weeks:   unit{[]string{"týden", "týdne", "týdny", "týdnů"}},
		days:    unit{[]string{"den", "dne", "dny", "dní"}},
		hours:   unit{[]string{"hodina", "hodiny", "hodiny", "hodin"}},
		minutes: unit{[]string{"minuta", "minuty", "minuty", "minut"}},
		seconds: unit{[]string{"sekunda", "sekundy", "sekundy", "sekund"}},
		rule:    pluralCzechOrSlovak,
	},
	LangES: {
		years:   unit{[]string{"año", "años"}},
		months:  unit{[]string{"mes", "meses"}},
		weeks:   unit{[]string{"semana", "semanas"}},
		days:    unit{[]string{"día", "días"}},
		hours:   unit{[]string{"hora", "horas"}},
		minutes: unit{[]string{"minuto", "minutos"}},
		seconds: unit{[]string{"segundo", "segundos"}},
		rule:    pluralOneMany,
	},
	LangFR: {
		years:   unit{[]string{"an", "ans"}},
		months:  unit{[]string{"mois", "mois"}},
		weeks:   unit{[]string{"semaine", "semaines"}},
		days:    unit{[]string{"jour", "jours"}},
		hours:   unit{[]string{"heure", "heures"}},
		minutes: unit{[]string{"minute", "minutes"}},
		seconds: unit{[]string{"seconde", "secondes"}},
		rule:    pluralFrench,
	},
	LangIT: {
		years:   unit{[]string{"anno", "anni"}},
		months:  unit{[]string{"mese", "mesi"}},
		weeks:   unit{[]string{"settimana", "settimane"}},
		days:    unit{[]string{"giorno", "giorni"}},
		hours:   unit{[]string{"ora", "ore"}},
		minutes: unit{[]string{"minuto", "minuti"}},
		seconds: unit{[]string{"secondo", "secondi"}},
		rule:    pluralOneMany,
	},
	LangJA: {
		years:   unit{[]string{"年"}},
		months:  unit{[]string{"ヶ月"}},
		weeks:   unit{[]string{"週間"}},
		days:    unit{[]string{"日"}},
		hours:   unit{[]string{"時間"}},
		minutes: unit{[]string{"分"}},
		seconds: unit{[]string{"秒"}},
	},
	LangKO: {
		years:   unit{[]string{"년"}},
		months:  unit{[]string{"개월"}},
		weeks:   unit{[]string{"주일"}},
		days:    unit{[]string{"일"}},
		hours:   unit{[]string{"시간"}},
		minutes: unit{[]string{"분"}},
		seconds: unit{[]string{"초"}},
	},
	LangUK: {
		years:   unit{[]string{"років", "рік", "роки"}},
		months:  unit{[]string{"місяців", "місяць", "місяці"}},
		weeks:   unit{[]string{"тижнів", "тиждень", "тижні"}},
		days:    unit{[]string{"днів", "день", "дні"}},
		hours:   unit{[]string{"годин", "година", "години"}},
		minutes: unit{[]string{"хвилин", "хвилина", "хвилини"}},
		seconds: unit{[]string{"секунд", "секунда", "секунди"}},
		rule:    pluralSlavic,
	},
	LangSK: {
		years:   unit{[]string{"rok", "roky", "roky", "rokov"}},
		months:  unit{[]string{"mesiac", "mesiace", "mesiace", "mesiacov"}},
		weeks:   unit{[]string{"týždeň", "týždne", "týždne", "týždňov"}},
		days:    unit{[]string{"deň", "dni", "dni", "dní"}},
		hours:   unit{[]string{"hodina", "hodiny", "hodiny", "hodín"}},
		minutes: unit{[]string{"minúta", "minúty", "minúty", "minút"}},
		seconds: unit{[]string{"sekunda", "sekundy", "sekundy", "sekúnd"}},
		rule:    pluralCzechOrSlovak,
	},
	LangTR: {
		years:   unit{[]string{"yıl"}},
		months:  unit{[]string{"ay"}},
		weeks:   unit{[]string{"hafta"}},
		days:    unit{[]string{"gün"}},
		hours:   unit{[]string{"saat"}},
		minutes: unit{[]string{"dakika"}},
		seconds: unit{[]string{"saniye"}},
	},
	LangZH: {
		years:   unit{[]string{"年"}},
		months:  unit{[]string{"个月"}},
		weeks:   unit{[]string{"周"}},
		days:    unit{[]string{"天"}},
		hours:   unit{[]string{"小时"}},
		minutes: unit{[]string{"分钟"}},
		seconds: unit{[]string{"秒"}},
	},
	LangZHHant: {
		years:   unit{[]string{"年"}},
		months:  unit{[]string{"個月"}},
		weeks:   unit{[]string{"周"}},
		days:    unit{[]string{"天"}},
		hours:   unit{[]string{"小時"}},
		minutes: unit{[]string{"分鐘"}},
		seconds: unit{[]string{"秒"}},
	},
	LangNO: {
		years:   unit{[]string{"år", "år"}},
		months:  unit{[]string{"måned", "måneder"}},
		weeks:   unit{[]string{"uke", "uker"}},
		days:    unit{[]string{"dag", "dager"}},
		hours:   unit{[]string{"time", "timer"}},
		minutes: unit{[]string{"minutt", "minutter"}},
		seconds: unit{[]string{"sekund", "sekunder"}},
		rule:    pluralOneMany,
	},
	LangNL: {
		years:   unit{[]string{"jaar", "jaar"}},
		months:  unit{[]string{"maand", "maanden"}},
		weeks:   unit{[]string{"week", "weken"}},
		days:    unit{[]string{"dag", "dagen"}},
		hours:   unit{[]string{"uur", "uur"}},
		minutes: unit{[]string{"minuut", "minuten"}},
		seconds: unit{[]string{"seconde", "seconden"}},
		rule:    pluralOneMany,
	},
	LangHI: {
		years:   unit{[]string{"साल", "साल"}},
		months:  unit{[]string{"महीना", "महीने"}},
		weeks:   unit{[]string{"हफ़्ता", "हफ्ते"}},
		days:    unit{[]string{"दिन", "दिन"}},
		hours:   unit{[]string{"घंटा", "घंटे"}},
		minutes: unit{[]string{"मिनट", "मिनट"}},
		seconds: unit{[]string{"सेकंड", "सेकंड"}},
		rule:    pluralOneMany,
	},
}
