package model

import (
	"time"
)

// NowUnix 返回当前时间的Unix时间戳（秒）
func NowUnix() int64 {
	return time.Now().Unix()
}

// UnixToTime 将Unix时间戳转换为time.Time
func UnixToTime(unix int64) time.Time {
	return time.Unix(unix, 0)
}

// SoftDelete 软删除标记
func SoftDelete() int64 {
	return NowUnix()
}
