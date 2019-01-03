// Copyright 2018 yejiantai Authors
//
// package time 公共方法
package time

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// 获得当前格式化的数据,形如2017-03-15 16:07:32.236
func GetFormatTime() string {
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.%03d", time.Now().Year(),
		time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(),
		time.Now().Second(), time.Now().Nanosecond()/1e6)
}

// 获得当前纯数字类型的时间数据,形如20170315160732236
func GetOnlyNumTime() string {
	return fmt.Sprintf("%04d%02d%02d%02d%02d%02d%03d", time.Now().Year(),
		time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(),
		time.Now().Second(), time.Now().Nanosecond()/1e6)
}

// 获得当前时间,格式为16:07:32.236
func GetCurTime() string {
	return fmt.Sprintf("%02d:%02d:%02d.%03d", time.Now().Hour(), time.Now().Minute(),
		time.Now().Second(), time.Now().Nanosecond()/1e6)
}

// 获取昨天的日期
func GetYesterday() string {
	return time.Now().AddDate(0, 0, -1).Format("20060102")
}

// 获取今天的日期
func GetNowDay() string {
	return time.Now().Format("20060102")
}

// 获取明天的日期
func GetTomorrow() string {
	return time.Now().AddDate(0, 0, 1).Format("20060102")
}

// 获取n_add天后的日期
func GetAddDay(str_day string, n_add int) (string, error) {
	var (
		day time.Time
		err error
		i   int
	)

	if day, err = ConvertDayToTime(str_day); err != nil {
		return "", err
	}
	for i = 0; i < n_add; i++ {
		day = day.AddDate(0, 0, 1)
	}
	return day.Format("20060102"), nil
}

// 获取nSecond秒后的时间差
func GetAfterDuration(nSecond int64) time.Duration {
	d, _ := time.ParseDuration(fmt.Sprintf("+%ds", nSecond))

	return d
}

// 获取日期范围内的全部日期列表,使用类似GetDateStringList("20180102", "20180205")
func GetDateStringList(str_start_date, str_end_date string) ([]string, error) {
	var (
		t_start time.Time
		err     error
		t_end   time.Time
	)

	if t_start, err = ConvertDayToTime(str_start_date); err != nil {
		return nil, err
	}

	if t_end, err = ConvertDayToTime(str_end_date); err != nil {
		return nil, err
	}

	return GetDateList(t_start, t_end)
}

// 获得开始时间和结束时间范围内的全部日期列表
func GetDateList(t_start, t_end time.Time) ([]string, error) {
	str_end_date, str_start_date := t_end.Format("20060102"), t_start.Format("20060102")
	if str_start_date > str_end_date {
		return nil, errors.New("开始时间大于结束时间")
	}
	var arr_date = []string{str_start_date}
	if str_start_date == str_end_date {
		return arr_date, nil
	}
	t := t_start
	for {
		//加一天
		t = t.AddDate(0, 0, 1)

		if t.Format("20060102") <= str_end_date {
			arr_date = append(arr_date, t.Format("20060102"))
		} else {
			break
		}
	}

	return arr_date, nil
}

// 将20180815格式的日期转换为time类型
func ConvertDayToTime(str_date string) (time.Time, error) {
	var t time.Time
	if len(str_date) != 8 {
		return t, errors.New("参数" + str_date + "格式错误，应该类似20180101")
	}
	year, err := strconv.Atoi(str_date[:4])
	if err != nil {
		return t, errors.New("参数" + str_date + "格式错误，应该类似20180101")
	}
	month, err := strconv.Atoi(str_date[4:6])
	if err != nil || month < 1 || month > 12 {
		return t, errors.New("参数" + str_date + "格式错误，应该类似20180101")
	}
	day, err := strconv.Atoi(str_date[6:])
	if err != nil || day < 1 || day > 31 {
		return t, errors.New("参数" + str_date + "格式错误，应该类似20180101")
	}
	t = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return t, nil
}
