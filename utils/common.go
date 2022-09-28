package utils

import (
	"strconv"
	"time"
)

const (
	TimeParseLayout      = "2006-01-02 15:04:05"
	FilePrefixTimeLayout = "2006-01-02_150405"
)

var TimeLoc, _ = time.LoadLocation("Asia/Shanghai")

func GenerateId(origin uint64) string {
	return strconv.FormatUint(origin, 10) +
		":" +
		strconv.FormatInt(time.Now().UnixNano(), 10) + ":"
}

// TimeCmp TimeCmp(t1,t2)
//d=abs(t1-t2)
//t1==t2 TimeCmp(t1,t2)==0
//t1>t2  TimeCmp(t1,t2)>0
//t1>t2  TimeCmp(t1,t2)<0
func TimeCmp(t1 time.Time, t2 time.Time) (cmp int, d time.Duration) {
	if t1.Equal(t2) {
		return 0, 0
	}
	if t1.After(t2) {
		return 1, t1.Sub(t2)
	}
	return -1, t2.Sub(t1)
}
