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

func TestJsonStrValid(t *testing.T) {
	jstr1 := `"test1"`
	jstr2 := `test2`
	jint := `1`
	jarray := `[1]`
	jarray2 := `[1"]`
	fmt.Println(JsonStrValid(jstr1))
	fmt.Println(JsonStrValid(jstr2))
	fmt.Println(JsonStrValid(jint))
	fmt.Println(JsonStrValid(jarray))
	fmt.Println(JsonStrValid(jarray2))
}
