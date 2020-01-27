package assert

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

//Nil ...
func Nil(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

//Equals equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}, msgs ...interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		if len(msgs) > 0 {
			fmt.Printf("\033[31m%s\033[39m\n\n", fmt.Sprintln(msgs...))
		}
		tb.FailNow()
	}
}
