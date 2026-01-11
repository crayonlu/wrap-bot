package plugins

import (
	"strconv"
	"strings"
	"time"
)

type TimeParser struct{}

func NewTimeParser() *TimeParser {
	return &TimeParser{}
}

func (p *TimeParser) Parse(expression string) (time.Time, error) {
	now := time.Now()

	expression = strings.ToLower(strings.TrimSpace(expression))

	num, unit := p.parseNumberAndUnit(expression)

	switch unit {
	case "天":
		if strings.Contains(expression, "后") {
			return now.AddDate(num, 0, 0), nil
		}
		return now.AddDate(-num, 0, 0), nil
	case "周", "星期":
		weekday := p.parseWeekday(expression)
		return p.getNextWeekday(now, weekday), nil
	case "月":
		if strings.Contains(expression, "下") {
			return now.AddDate(0, num+1, 0), nil
		}
		return now.AddDate(0, num, 0), nil
	case "小时":
		if strings.Contains(expression, "后") {
			return now.Add(time.Duration(num) * time.Hour), nil
		}
		return now.Add(-time.Duration(num) * time.Hour), nil
	case "分钟":
		if strings.Contains(expression, "后") {
			return now.Add(time.Duration(num) * time.Minute), nil
		}
		return now.Add(-time.Duration(num) * time.Minute), nil
	}

	return now, nil
}

func (p *TimeParser) parseNumberAndUnit(expression string) (int, string) {
	if strings.Contains(expression, "明天") {
		return 1, "天"
	}
	if strings.Contains(expression, "后天") {
		return 2, "天"
	}

	words := []string{"天", "周", "星期", "月", "小时", "分钟"}
	for _, word := range words {
		if strings.Contains(expression, word) {
			idx := strings.Index(expression, word)
			numStr := strings.TrimSpace(expression[:idx])
			num, _ := strconv.Atoi(numStr)
			if num == 0 {
				num = 1
			}
			return num, word
		}
	}

	return 0, ""
}

func (p *TimeParser) parseWeekday(expression string) time.Weekday {
	weekdayMap := map[string]time.Weekday{
		"周日":    time.Sunday,
		"星期日": time.Sunday,
		"周一":    time.Monday,
		"星期一": time.Monday,
		"周二":    time.Tuesday,
		"星期二": time.Tuesday,
		"周三":    time.Wednesday,
		"星期三": time.Wednesday,
		"周四":    time.Thursday,
		"星期四": time.Thursday,
		"周五":    time.Friday,
		"星期五": time.Friday,
		"周六":    time.Saturday,
		"星期六": time.Saturday,
	}

	for key, wd := range weekdayMap {
		if strings.Contains(expression, key) {
			return wd
		}
	}

	return time.Sunday
}

func (p *TimeParser) getNextWeekday(now time.Time, weekday time.Weekday) time.Time {
	daysUntil := int((weekday - now.Weekday() + 7) % 7)
	if daysUntil == 0 {
		daysUntil = 7
	}
	return now.AddDate(0, 0, daysUntil)
}
