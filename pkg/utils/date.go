package utils

import (
	"fmt"
	"time"
)

func GetPastMonthToday(t time.Time, month int) string {
	// today
	fmt.Printf("today: [%s]\n", t)
	// 判断天数范围 小于等于28天的计算,覆盖大多数情况
	if t.Day() <= 28 {
		return t.AddDate(0, -month, 0).Format("2006-01-02 15:04:05")
	}
	// 月份的天数数组
	monthDay := [13]int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	// 计算目标所在日期
	target := t.AddDate(0, 0, 1-t.Day()).AddDate(0, -month, 0)
	// 计算当月最大天数
	targetDay := monthDay[target.Month()]
	// 计算闰年
	if target.Month() == time.February && (target.Year()%400 == 0 || (target.Year()%100 != 0 && target.Year()%4 == 0)) {
		targetDay++
	}
	if t.Day() > targetDay {
		return target.AddDate(0, 0, targetDay-1).Format("2006-01-02 15:04:05")
	}
	return target.AddDate(0, 0, t.Day()-1).Format("2006-01-02 15:04:05")
}
