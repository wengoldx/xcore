// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2025/05/28   youhei         New version
// -------------------------------------------------------------------

package utils

import (
	"time"
)

// Time convert utils for create time.Time object from string
// with local timezone, and format to string with indicated layout.
//
//	ut := utils.Now()
//	ut := utils.FromTime(time.Now())
//	ut := utils.FromString('2026-01-02 03:04:05')
//	ut := utils.FromRFC3339('2026-01-02T03:04:05Z08:00')
//	ut := utils.AddYear(1)  // after 1 year.
//	ut := utils.AddMonth(1) // after 1 month.
//	ut := utils.AddDate(1)  // after 1 day.
//
// Format now time to string as:
//
//	utils.Now().ToString()  // 2026-01-02 03:04:05
//	utils.Now().ToRFC3339() // 2026-01-02T03:04:05Z08:00
//	utils.Now().ToDate()    // 2026-01-02
//	utils.Now().ToTime()    // 03:04:05
type UTime struct {
	ins time.Time
}

const (
	MsLayout       = "2006-01-02 15:04:05.000" // Time layout with million seconds.
	ZoneLayout     = "2006-01-02T15:04:05Z"    // Time layout with local timezone, see RFC3339.
	DateNoneHyphen = "20060102"                // Time layout date only without hyphen char.
	TimeNoneHyphen = "20060102150405"          // Time layout datetime without hyphen char.
	HourNoneHyphen = "150405"                  // Time layout time only without hyphen char.
	MSNoneHyphen   = "20060102150405000"       // Time layout with million seconds without hyphen char.
)

/* ------------------------------------------------------------------- */
/* Create UTime Object Utils                                           */
/* ------------------------------------------------------------------- */

// Now time with local timezone.
func Now() *UTime { return &UTime{ins: time.Now()} }

// Today date at 00:00:00 clock as '2006-01-02 00:00:00'.
func Today() *UTime { return &UTime{ins: Now().Date()} }

// Yesterday date at 00:00:00 clock as '2006-01-01 00:00:00'.
func Yesterday() *UTime { return &UTime{ins: FromTime(time.Now().AddDate(0, 0, -1)).Date()} }

// Tomorrow date at 00:00:00 clock as '2006-01-03 00:00:00'.
func Tomorrow() *UTime { return &UTime{ins: FromTime(time.Now().AddDate(0, 0, 1)).Date()} }

// New UTime objec from time.Time object.
func FromTime(t time.Time) *UTime { return &UTime{ins: t} }

// New UTime objec from datetime string format as '2006-01-02 15:04:05'.
func FromString(s string) *UTime { return FromTime(parseLocalTime(s, time.DateTime)) }

// New UTime objec from datetime string format as '2006-01-02T15:04:05Z07:00'.
func FromRFC3339(s string) *UTime { return FromTime(parseLocalTime(s, time.RFC3339)) }

// New UTime object from now and offset input years.
func AddYear(y int) *UTime { return Now().AddYear(y) }

// New UTime object from now and offset input months.
func AddMonth(m int) *UTime { return Now().AddMonth(m) }

// New UTime object from now and offset input days.
func AddDay(d int) *UTime { return Now().AddDay(d) }

/* ------------------------------------------------------------------- */
/* UTime Methods                                                       */
/* ------------------------------------------------------------------- */

func (t *UTime) AddYear(y int) *UTime      { t.ins = t.ins.AddDate(y, 0, 0); return t }
func (t *UTime) AddMonth(m int) *UTime     { t.ins = t.ins.AddDate(0, m, 0); return t }
func (t *UTime) AddDay(d int) *UTime       { t.ins = t.ins.AddDate(0, 0, d); return t }
func (t *UTime) AddHour(h int) *UTime      { t.ins = t.ins.Add(time.Duration(h) * time.Hour); return t }
func (t *UTime) AddMinute(m int) *UTime    { t.ins = t.ins.Add(time.Duration(m) * time.Minute); return t }
func (t *UTime) AddSecond(s int) *UTime    { t.ins = t.ins.Add(time.Duration(s) * time.Second); return t }
func (t *UTime) Offset(y, m, d int) *UTime { t.ins = t.ins.AddDate(y, m, d); return t }

// Return time.Time object.
func (t *UTime) Time() time.Time { return t.ins }

// Return time.Time object but only year, month, day values as '2006-01-02 00:00:00'.
func (t *UTime) Date() time.Time {
	y, m, d := t.ins.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

// Return unix seconds.
func (t *UTime) Unix() int64 { return t.ins.Unix() }

// Return unix nano seconds.
func (t *UTime) UnixNano() int64 { return t.ins.UnixNano() }

// Check whether expired.
func (t *UTime) IsExpired() bool { return t.ins.Before(time.Now()) }

// Check whether today.
func (t *UTime) IsToday() bool { return IsToday(t.ins) }

// Check whether the same datetime.
func (t *UTime) IsSame(d time.Time) bool { return t.ins.Equal(d) }

// Check whether the same date, it only check year, month, day.
func (t *UTime) IsSameDay(d time.Time) bool { return IsSameDay(t.ins, d) }

// Check t whether before the input d datetime.
func (t *UTime) Before(d time.Time) bool { return t.ins.Before(d) }

// Check whether same datetime from otner UTime object.
func (t *UTime) Equal(o *UTime) bool { return t.ins.Equal(o.ins) }

// To UTC time string as '2006-01-02 15:04:05'.
func (t *UTime) ToString() string { return t.ins.Format(time.DateTime) }

// To RFC3339 time string as '2006-01-02T15:04:05Z07:00'.
func (t *UTime) ToRFC3339() string { return t.ins.Format(time.RFC3339) }

// To date only string as '2006-01-02'.
func (t *UTime) ToDate() string { return t.ins.Format(time.DateOnly) }

// To time only string as '15:04:05'.
func (t *UTime) ToTime() string { return t.ins.Format(time.TimeOnly) }

// To time string according layout format.
//
//	- utils.ZoneLayout     : '2006-01-02T15:04:05Z'
//	- utils.MSLayout       : '2006-01-02 15:04:05.000'
//	- utils.DateNoneHyphen : '20060102'
//	- utils.TimeNoneHyphen : '20060102150405'
//	- utils.HourNoneHyphen : '150405'
//	- utils.MSNoneHyphen   : '20060102150405000'
//
// See ToString(), ToRFC3339(), ToDate(), ToTime() for more layout formats.
func (t *UTime) ToLayout(layout string) string { return t.ins.Format(layout) }

// Diff with input d datetime, and return hours, minutes, seconds.
//
// # WARNING:
//	- More than 24 hours diff duratons will be trimmed!
func (t *UTime) TimeDiff(oth time.Time) (int, int, int) {
	v, h, m := int(t.Unix()-oth.Unix()), 3600, 60
	return (v / h) % 24, (v % h / m), (v % m)
}

// Diff with input d datetime, and return days, hours, minutes, seconds.
func (t *UTime) DayDiff(oth time.Time) (int, int, int, int) {
	v, d, h, m := int(t.Unix()-oth.Unix()), 3600*24, 3600, 60
	return (v / d), (v % d / h), (v % h / m), (v % m)
}

/* ------------------------------------------------------------------- */
/* Global Export Utils                                                 */
/* ------------------------------------------------------------------- */

// Parse time from string with timezone.
func parseLocalTime(s, layout string) time.Time {
	t, _ := time.ParseInLocation(layout, s, time.Local)
	return t
}

// Parse time from string with timezone, the src time string maybe formated
// from time.Format() returned value.
//
// # WARNING:
//
// time.Now(), time.Format() are local time with timezone offset,
// but time.Parse() parse time string without timezone just UTC time,
// you can parse local time by use.
//
//	time.ParseInLocation(layout, timestring, time.Local).
//
// The layout enable using:
//	- time.DateTime : '2006-01-02 15:04:05'
//	- time.RFC3339  : '2006-01-02T15:04:05Z07:00'
//	- time.DateOnly : '2006-01-02'
//	- time.TimeOnly : '15:04:05'
//	- utils.MSLayout: '2006-01-02 15:04:05.000'
func ParseTime(layout, src string) (time.Time, error) {
	return time.ParseInLocation(layout, src, time.Local)
}

// Check times whether on same date, it just check year, month and day.
//
//	return d1.year == d2.year && d1.month == d2.month && d1.day == d2.day
func IsSameDay(d1, d2 time.Time) bool {
	y, m, d := d1.Year(), d1.Month(), d1.Day()
	return y == d2.Year() && m == d2.Month() && d == d2.Day()
}

// Check the given datatime string, unix sencods, time object
// whether today, the datetime string must format by time.Format()
// or offseted timezone
func IsToday[T any](des T) bool {
	now := time.Now()
	y, m, d := now.Year(), now.Month(), now.Day() 
	switch vt := any(des).(type) {
	case string:
		dt, _ := ParseTime(time.DateOnly, vt)
		return y == dt.Year() && m == dt.Month() && d == dt.Day()
	case int64:
		dt := time.Unix(vt, 0)
		return y == dt.Year() && m == dt.Month() && d == dt.Day()
	case time.Time:
		return y == vt.Year() && m == vt.Month() && d == vt.Day()
	}
	return false
}