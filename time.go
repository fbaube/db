package db

import (
	"fmt"
	"strconv"
	"time"
)

// Ito09az converts its int arg (0..35) to
// a string of length one, in the range
// (for int  0..9)  "0".."9",
// (for int 10..35) "a".."z"
func Ito09az(i int) string {
	if i <= 9 {
		return strconv.Itoa(i)
	}
	var bb = make([]byte, 1, 1)
	bb[0] = byte(i - 10 + 'a')
	return string(bb)
}

// NowAsYMDHM maps (reversibly) the current
// time to "YMDhm" (where m is minutes / 2).
func NowAsYMDHM() string {
	var now = time.Now()
	// year = last digit
	var y string = fmt.Sprintf("%d", now.Year())[3:]
	var m string = Ito09az(int(now.Month()))
	var d string = Ito09az(now.Day())
	var h string = Ito09az(now.Hour())
	var n string = Ito09az(now.Minute() / 2)
	// fmt.Printf("%s-%s-%s-%s-%s", y, m, d, h, n)
	return fmt.Sprintf("%s%s%s%s%s", y, m, d, h, n)
}

// NowAsYM maps (reversibly) the current
// year+month to "YM".
func NowAsYM() string {
	var now = time.Now()
	// year = last digit
	var y string = fmt.Sprintf("%d", now.Year())[3:]
	var m string = Ito09az(int(now.Month()))
	// fmt.Printf("%s-%s-%s-%s-%s", y, m, d, h, n)
	return fmt.Sprintf("%s%s", y, m)
}