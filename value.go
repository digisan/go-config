package goconfig

import (
	"fmt"
	"regexp"

	lk "github.com/digisan/logkit"
)

var (
	fPathCfg = ""
	mCfg     = make(map[string]any)
)

func Show() {
	fmt.Println()
	r1 := regexp.MustCompile(`^_\w+`)
	r2 := regexp.MustCompile(`\._\w+`)
	for k, v := range mCfg {
		if !(r1.MatchString(k) || r2.MatchString(k)) {
			fmt.Printf("%-16v: %v\n", k, v)
		}
	}
	fmt.Println()
}

func Val[T any](field string) T {

	valAny, ok := mCfg[field]
	lk.FailP1OnErrWhen(!ok, "%v", fmt.Errorf("[%v] is NOT in file [%s]", field, fPathCfg))

	val, ok := valAny.(T)
	lk.FailP1OnErrWhen(!ok, "%v", fmt.Errorf("[%v] cannot convert to needed type, (number must be [float64])", field))

	return val
}

func ValInt(field string) int {
	return int(Val[float64](field))
}

func ValUint(field string) uint {
	return uint(Val[float64](field))
}
