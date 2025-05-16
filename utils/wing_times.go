// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package utils

import (
	"fmt"
	"time"
)

// Day, week, duration on nanosecond
const (
	Day  = time.Hour * 24
	Week = Day * 7
)

// Day, week, duration on millisecond
const (
	DayMsDur  = Day / time.Millisecond
	WeekMsDur = Week / time.Millisecond
)

// Day, week duration on second
const (
	DaySecondsDur  = Day / time.Second
	WeekSecondsDur = Week / time.Second
)

const (
	// DateLayout standery date layout format at day minimum
	DateLayout = "2006-01-02"

	// TimeLayout standery time layout format at second minimum
	TimeLayout = "2006-01-02 15:04:05"

	// HourLayout standery time layout format as hour
	HourLayout = "15:04:05"

	// MSLayout standery time layout format at million second minimum
	MSLayout = "2006-01-02 15:04:05.000"

	// DateNoneHyphen standery time layout format at second minimum without hyphen char
	DateNoneHyphen = "20060102"

	// TimeNoneHyphen standery time layout format at second minimum without hyphen char
	TimeNoneHyphen = "20060102150405"

	// HourNoneHyphen standery time layout format as hour without hyphen char
	HourNoneHyphen = "150405"

	// MSNoneHyphen standery time layout format at million second minimum without hyphen char
	MSNoneHyphen = "20060102150405000"
)

// ParseTime parse time with zoom, the src time string maybe
// formated from time.Format() return value.
//
// `WARNING` :
//
// time.Now(), time.Format() are local time with timezoom offset,
// but time.Parse() parse time string without timezoom just UTC time,
// you can parse local time by use
//
//	time.ParseInLocation(layout, timestring, time.Local).
func ParseTime(layout, src string) (time.Time, error) {
	return time.ParseInLocation(layout, src, time.Local)
}

// IsToday check the given day string if today, the des time
// string must format by time.Format() or offseted timezoom
func IsToday(des string) bool {
	now := time.Now().Format(DateLayout)
	st, _ := ParseTime(DateLayout, now)
	dt, _ := ParseTime(DateLayout, des)
	return st.Unix() == dt.Unix()
}

// IsTodayUnix check the given time string if today, the des time
// unix seconds from time.Now() or contian timezoom
func IsTodayUnix(des int64) bool {
	deslayout := time.Unix(des, 0).Format(DateLayout)
	return IsToday(deslayout)
}

// IsSameDay equal given days string based on TimeLayout, the src and
// des time string must format by time.Format() or offseted timezoom
func IsSameDay(src string, des string) bool {
	st, _ := ParseTime(DateLayout, src)
	dt, _ := ParseTime(DateLayout, des)
	return st.Unix() == dt.Unix()
}

// IsSameTime equal given time string based on TimeLayout, the src and
// des time string must format by time.Format() or offseted timezoom
func IsSameTime(src string, des string) bool {
	st, _ := ParseTime(TimeLayout, src)
	dt, _ := ParseTime(TimeLayout, des)
	return st.Unix() == dt.Unix()
}

// Today return today unix time with offseted location timezoom
func Today() int64 {
	return time.Now().Unix()
}

// Yesterday return yesterday unix time base on current location time
func Yesterday() int64 {
	return time.Now().AddDate(0, 0, -1).Unix()
}

// Tommorrow return tommorrow unix time srart from current location time
func Tommorrow() int64 {
	return time.Now().AddDate(0, 0, 1).Unix()
}

// NextWeek return next week unix time start from current,
// or from the given unix seconds and nano seconds
func NextWeek(start ...int64) int64 {
	if len(start) > 1 && start[0] > 0 {
		return time.Unix(start[0], start[1]).AddDate(0, 0, 7).Unix()
	}
	return time.Now().AddDate(0, 0, 7).Unix()
}

// NextMonth return next month unix time start from current,
// or from the given unix seconds and nano seconds
func NextMonth(start ...int64) int64 {
	if len(start) > 1 && start[0] > 0 {
		return time.Unix(start[0], start[1]).AddDate(0, 1, 0).Unix()
	}
	return time.Now().AddDate(0, 1, 0).Unix()
}

// NextQuarter return next quarter unix time start from current,
// or from the given unix seconds and nano seconds
func NextQuarter(start ...int64) int64 {
	if len(start) > 1 && start[0] > 0 {
		return time.Unix(start[0], start[1]).AddDate(0, 3, 0).Unix()
	}
	return time.Now().AddDate(0, 3, 0).Unix()
}

// NextHalfYear return next half a year unix time start from current,
// or from the given unix seconds and nano seconds
func NextHalfYear(start ...int64) int64 {
	if len(start) > 1 && start[0] > 0 {
		return time.Unix(start[0], start[1]).AddDate(0, 6, 0).Unix()
	}
	return time.Now().AddDate(0, 6, 0).Unix()
}

// NextYear return next year unix time start from current,
// or from the given unix seconds and nano seconds
func NextYear(start ...int64) int64 {
	if len(start) > 1 && start[0] > 0 {
		return time.Unix(start[0], start[1]).AddDate(1, 0, 0).Unix()
	}
	return time.Now().AddDate(1, 0, 0).Unix()
}

// NextTime return next unix time start from current,
// or from the given unix seconds and nano seconds
func NextTime(duration time.Duration, start ...int64) int64 {
	if len(start) > 1 && start[0] > 0 {
		return time.Unix(start[0], start[1]).Add(duration).Unix()
	}
	return time.Now().Add(duration).Unix()
}

// DayUnix return the given day unix time at 0:00:00
func DayUnix(src int64) int64 {
	dt := time.Unix(src, 0).Format(DateLayout)
	st, _ := ParseTime(DateLayout, dt)
	return st.Unix()
}

// TodayUnix return today unix time at 0:00:00
func TodayUnix() int64 {
	now := time.Now().Format(DateLayout)
	st, _ := ParseTime(DateLayout, now)
	return st.Unix()
}

// YesterdayUnix return yesterday unix time at 0:00:00
func YesterdayUnix() int64 {
	return NextUnix(-Day)
}

// TommorrowUnix return tommorrow unix time at 0:00:00
func TommorrowUnix() int64 {
	return NextUnix(Day)
}

// WeekUnix return next week day (same as current weekday) unix time at 0:00:00
func WeekUnix() int64 {
	return NextUnix(Week)
}

// MonthUnix return next month day (same as current day of month) unix time at 0:00:00
func MonthUnix() int64 {
	return NextUnix2(0, 1, 0)
}

// QuarterUnix return next quarter unix time at 0:00:00
func QuarterUnix() int64 {
	return NextUnix2(0, 3, 0)
}

// YearUnix return next year unix time at 0:00:00
func YearUnix() int64 {
	return NextUnix2(1, 0, 0)
}

// NextUnix return next 0:00:00 unix time at day after given duration
func NextUnix(duration time.Duration) int64 {
	nt := time.Now().Add(duration).Format(DateLayout)
	st, _ := ParseTime(DateLayout, nt)
	return st.Unix()
}

// DaysUnix return the unix time at 0:00:00, by offset days
func DaysUnix(days int) int64 {
	return NextUnix2(0, 0, days)
}

// WeekUnix return unix time at 0:00:00, by offset weeks
func WeeksUnix(weeks int) int64 {
	return NextUnix2(0, 0, weeks*7)
}

// MonthsUnix return unix time at 0:00:00, by offset months
func MonthsUnix(months int) int64 {
	return NextUnix2(0, months, 0)
}

// QuartersUnix return unix time at 0:00:00, by offset quarters
func QuartersUnix(quarters int) int64 {
	return NextUnix2(0, 3*quarters, 0)
}

// YearsUnix return unix time at 0:00:00, by offset years
func YearsUnix(years int) int64 {
	return NextUnix2(years, 0, 0)
}

// NextUnix2 return next 0:00:00 unix time at day after given years, months and days
func NextUnix2(years, months, days int) int64 {
	nt := time.Now().AddDate(years, months, days).Format(DateLayout)
	st, _ := ParseTime(DateLayout, nt)
	return st.Unix()
}

// HourDiff return diff hours, minutes, seconds
func HourDiff(start, end time.Time) (int, int, int) {
	v, h, m := int(end.Unix()-start.Unix()), 3600, 60
	return (v / h), (v % h / m), (v % m)
}

// DayDiff return diff days, hours, minutes, seconds
func DayDiff(start, end time.Time) (int, int, int, int) {
	v := int(end.Unix() - start.Unix())
	var d, h, m int = int(DaySecondsDur), 3600, 60
	return (v / d), (v % d / h), (v % h / m), (v % m)
}

// DurHours return readable time during start to end like 06:25:48,
// you can see the format string, but it must contain 3 %0xd to parse numbers
func DurHours(start, end time.Time, format ...string) string {
	h, m, s := HourDiff(start, end)
	return fmt.Sprintf(VarString(format, "%02d:%02d:%02d"), h, m, s)
}

// DurDays return readable time during start to end like 2d 6h 25m 48s,
// you can set the format string, but it must contain 4 %0xd to parse numbers
func DurDays(start, end time.Time, format ...string) string {
	d, h, m, s := DayDiff(start, end)
	return fmt.Sprintf(VarString(format, "%dd %dh %dm %ds"), d, h, m, s)
}

// DurNowNs return formated second string from given start in unix nanoseconds
func DurNowNs(start int64) string {
	dur := time.Now().UnixNano() - start
	s := (int64)(time.Duration(dur) / time.Second)
	ms := (int64)((time.Duration(dur) % time.Second) / time.Millisecond)
	ws := (int64)((time.Duration(dur) % time.Millisecond) / time.Microsecond)
	ns := (int64)(time.Duration(dur) % time.Microsecond)
	return fmt.Sprintf("%ds %d.%d.%d", s, ms, ws, ns)
}

// FormatTime format unix time to TimeLayout or MSLayout layout
func FormatTime(sec int64, nsec ...int64) string {
	return time.Unix(sec, VarInt64(nsec, 0)).Format(TimeLayout)
}

// FormatUnix format unix time to given time layout with location timezoom
func FormatUnix(layout string, sec int64, nsec ...int64) string {
	switch layout {
	case DateLayout, TimeLayout, HourLayout, DateNoneHyphen, TimeNoneHyphen, HourNoneHyphen:
		return time.Unix(sec, 0).Format(layout)

	case MSLayout, MSNoneHyphen:
		return time.Unix(sec, VarInt64(nsec, 0)).Format(layout)
	}

	// TimeLayout as the default time layout
	return time.Unix(sec, 0).Format(TimeLayout)
}

// FormatNow format now to given time layout, it may format as TimeLayout
// when input param not set, and the formated time contain location timezoom.
func FormatNow(layout ...string) string {
	nowns := time.Now().UnixNano()
	if l := VarString(layout, ""); l != "" {
		return FormatUnix(l, nowns/1e9, (nowns%1e9)/1e6)
	}
	return FormatUnix(TimeLayout, nowns/1e9)
}

// FormatToday format now to today string as '2006-01-02', or custom
// format set by layout like '2006/01/02', '2006.01.02' and so on.
func FormatToday(layout ...string) string {
	return time.Now().Format(VarString(layout, DateLayout))
}

// FormatDay format day string from given unix seconds as '2006-01-02',
// or custom format like '2006/01/02', '2006.01.02' and so on.
func FormatDay(sec int64, layout ...string) string {
	return time.Unix(sec, 0).Format(VarString(layout, DateLayout))
}

// FormatDur format the time before or after now by add given duration, it
// may format as TimeLayout when input layout not set, and the formated
// time contain location timezoom.
func FormatDur(d time.Duration, layout ...string) string {
	nowns := time.Now().Add(d).UnixNano()
	if l := VarString(layout, ""); l != "" {
		return FormatUnix(l, nowns/1e9, (nowns%1e9)/1e6)
	}
	return FormatUnix(TimeLayout, nowns/1e9)
}
