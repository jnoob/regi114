package main

import (
	"fmt"
	"strconv"
	"time"
)

func getTime() string {
	return strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
}

func getDutyDate() string {
	return time.Now().Add(time.Hour * 24 * 7).Format("2006-01-02")
}

const TIME_LAYOUT = "2006-01-02 15:04:05 -0700 MST"

var _7Pre3s int64

func init() {
	t := time.Now()
	str := fmt.Sprintf("%d-%02d-%02d 06:59:57 +0800 CST", t.Year(), t.Month(), t.Day())
	t, _ = time.Parse(TIME_LAYOUT, str)
	_7Pre3s = t.Unix()
}

func get7Pre3Second() int64 {
	return _7Pre3s
}
