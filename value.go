package goconfig

import (
	"encoding/json"
	"fmt"
	"regexp"

	. "github.com/digisan/go-generics/v2"
	lk "github.com/digisan/logkit"
	"github.com/tidwall/gjson"
)

var (
	fPathCfg = ""
	jsCfg    = ""
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

// only for primitive value
func Val[T any](field string) T {

	valAny, ok := mCfg[field]
	lk.FailP1OnErrWhen(!ok, "%v", fmt.Errorf("[%v] is NOT in file [%s]", field, fPathCfg))

	t := fmt.Sprintf("%T", new(T))
	var ret any
	switch t {

	case "*float32":
		ret = float32(valAny.(float64))

	case "*int":
		ret = int(valAny.(float64))
	case "*int8":
		ret = int8(valAny.(float64))
	case "*int16":
		ret = int16(valAny.(float64))
	case "*int32":
		ret = rune(valAny.(float64))
	case "*int64":
		ret = int64(valAny.(float64))

	case "*uint":
		ret = uint(valAny.(float64))
	case "*uint8":
		ret = byte(valAny.(float64))
	case "*uint16":
		ret = uint16(valAny.(float64))
	case "*uint32":
		ret = uint32(valAny.(float64))
	case "*uint64":
		ret = uint64(valAny.(float64))

		// ... more number types

	default:
		val, ok := valAny.(T)
		lk.FailP1OnErrWhen(!ok, "%v", fmt.Errorf("[%v] cannot convert to type [%s]", field, t[1:]))
		return val
	}
	return ret.(T)
}

func ValArr[T any](field string) []T {

	lk.FailP1OnErrWhen(len(jsCfg) == 0, "%v", fmt.Errorf("config data is empty, Must Init"))

	if r := gjson.Get(jsCfg, field); r.IsArray() {
		ret := FilterMap(r.Array(), nil, func(i int, e gjson.Result) any {
			switch fmt.Sprintf("%T", new(T)) {

			case "*float64":
				return e.Num
			case "*float32":
				return float32(e.Num)

			case "*int":
				return int(e.Num)
			case "*int8":
				return int8(e.Num)
			case "*int16":
				return int16(e.Num)
			case "*int32":
				return rune(e.Num)
			case "*int64":
				return int64(e.Num)

			case "*uint":
				return uint(e.Num)
			case "*uint8":
				return byte(e.Num)
			case "*uint16":
				return uint16(e.Num)
			case "*uint32":
				return uint32(e.Num)
			case "*uint64":
				return uint64(e.Num)

			case "*bool":
				return e.Bool()

			case "*string":
				return e.Str

			case "*interface {}", "*map[string]interface {}":
				rt := make(map[string]any)
				lk.FailP1OnErr("%v", json.Unmarshal([]byte(e.Raw), &rt))
				return rt

				// ... more

			default:
				lk.FailP1OnErr("%v", fmt.Errorf("more type case is needed in [ValArr]"))
				return nil
			}
		})
		return SlcCvt[T](ret)
	}

	lk.FailP1OnErr("%v", fmt.Errorf("[%s] is not array type", field))
	return nil
}

func ValObj(field string) map[string]any {

	lk.FailP1OnErrWhen(len(jsCfg) == 0, "%v", fmt.Errorf("config data is empty, Must Init"))

	if r := gjson.Get(jsCfg, field); r.IsObject() {
		rt := make(map[string]any)
		lk.FailP1OnErr("%v", json.Unmarshal([]byte(r.Raw), &rt))
		return rt
	}

	lk.FailP1OnErr("%v", fmt.Errorf("[%s] is not object type", field))
	return nil
}

func CntArr[T any](field string) int {
	return len(ValArr[T](field))
}
