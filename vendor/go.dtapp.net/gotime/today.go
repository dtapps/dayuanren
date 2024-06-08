package gotime

import (
	"log"
	"strconv"
	"strings"
	"time"
)

// Current 获取当前的时间
func Current() Pro {

	var err error
	p := NewPro()

	p.loc, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Printf("【gotime】时区错误：%v\n", err)
		p.Time = time.Now().Add(time.Hour * 8)
	} else {
		p.Time = time.Now().In(p.loc)
	}

	return p
}

// SetCurrent 设置当前的时间
func SetCurrent(sTime time.Time) Pro {
	p := NewPro()
	p.Time = sTime
	return p
}

// SetCurrentParse 设置当前的时间
func SetCurrentParse(str string) Pro {

	var err error
	p := NewPro()

	p.loc, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Printf("【gotime】时区错误：%v\n", err)
		p.Time = time.Now().Add(time.Hour * 8)
	}

	layout := DateTimeFormat
	if str == "" || str == "0" || str == "0000-00-00 00:00:00" || str == "0000-00-00" || str == "00:00:00" {
		return p
	}
	if len(str) == 10 && strings.Count(str, "-") == 2 {
		layout = DateFormat
	}
	if strings.Index(str, "T") == 10 {
		layout = RFC3339Format
	}
	if _, err := strconv.ParseInt(str, 10, 64); err == nil {
		switch len(str) {
		case 8:
			layout = ShortDateFormat
		case 14:
			layout = ShortDateTimeFormat
		}
	}
	location, _ := time.ParseInLocation(layout, str, p.loc)

	p.Time = location
	return p
}

// SetCurrentUnix 设置当前的时间 Unix时间戳
func SetCurrentUnix(ts int64) Pro {
	p := NewPro()
	p.Time = time.Unix(ts, 0)
	return p
}

// SetCurrentMillisecondUnix 设置当前的时间 毫秒Unix时间戳
func SetCurrentMillisecondUnix(ts int64) Pro {
	p := NewPro()
	p.Time = time.Unix(ts/1000, 0)
	return p
}

// Now 今天此刻
func (p Pro) Now() time.Time {
	return p.Time
}

// Format 今天此刻格式化
func (p Pro) Format() string {
	return p.Time.Format(DateTimeFormat)
}

// FormatFilter 今天此刻格式化 带 过滤无效时间
func (p Pro) FormatFilter() string {
	if p.Time.Format(DateTimeFormat) == "0001-01-01 00:00:00" || p.Time.Format(DateTimeFormat) == "0001-01-01 08:05:43" {
		return ""
	} else {
		return p.Time.Format(DateTimeFormat)
	}
}

// ToDateFormat 今天此刻日期
func (p Pro) ToDateFormat() string {
	return p.Time.Format(DateFormat)
}

// ToDateFormatTime 今天此刻日期
func (p Pro) ToDateFormatTime() time.Time {
	return SetCurrentParse(p.Time.Format(DateFormat)).Time
}

// ToTimeFormat 今天此刻时间
func (p Pro) ToTimeFormat() string {
	return p.Time.Format(TimeFormat)
}

// Timestamp 今天此刻时间戳
func (p Pro) Timestamp() int64 {
	return p.Time.Unix()
}

// TimestampWithSecond 今天此刻时间戳
func (p Pro) TimestampWithSecond() int64 {
	return p.Time.Unix()
}

// TimestampWithMillisecond 今天毫秒级时间戳
func (p Pro) TimestampWithMillisecond() int64 {
	return p.Time.UnixNano() / int64(time.Millisecond)
}

// TimestampWithMicrosecond 今天微秒级时间戳
func (p Pro) TimestampWithMicrosecond() int64 {
	return p.Time.UnixNano() / int64(time.Microsecond)
}

// TimestampWithNanosecond 今天纳秒级时间戳
func (p Pro) TimestampWithNanosecond() int64 {
	return p.Time.UnixNano()
}
