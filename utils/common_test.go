package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestTimecmp(t *testing.T) {
	t1, _ := time.Parse(TimeParseLayout, "2006-01-02 15:06:05")
	t2, _ := time.Parse(TimeParseLayout, "2006-01-02 15:05:05")
	fmt.Println(TimeCmp(t1, t2))
	fmt.Println(t2.Sub(t1))
}
