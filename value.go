package goconfig

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	ff "github.com/digisan/fileflatter"
	. "github.com/digisan/go-generics/v2"
	"github.com/digisan/gotk/strs"
	lk "github.com/digisan/logkit"
	"github.com/tidwall/gjson"
)

type Cfg struct {
	path string
	typ  string
	str  string
	fm   map[string]any
}

var (
	mtx    = sync.Mutex{}
	MapCfg = make(map[string]*Cfg)
	cfg    *Cfg
)

func Use(id string) error {
	mtx.Lock()
	defer mtx.Unlock()
	p, ok := MapCfg[id]
	if !ok {
		return fmt.Errorf("[%v] is uninitialized, do 'Init' before using it", id)
	}
	cfg = p
	return nil
}

func Show() {
	ff.PrintFlat(cfg.fm)
}

func path(paths ...any) string {
	sp := FilterMap(paths, nil, func(i int, e any) string { return fmt.Sprint(e) })
	return strings.Join(sp, ".")
}

// for primitive value
func Val[T any](paths ...any) T {

	field := path(paths...)

	valAny, ok := cfg.fm[field]
	lk.FailP1OnErrWhen(!ok, "%v", fmt.Errorf("[%v] is NOT in file [%s]", field, cfg.path))

	ret, ok := AnyTryToType[T](valAny)
	lk.FailOnErrWhen(!ok, "%v", fmt.Errorf("value of path [%v] cannot be type [%T]", valAny, *new(T)))
	return ret
}

// for array value, array must be primitive array
func ValArr[T any](paths ...any) []T {

	lk.FailP1OnErrWhen(len(cfg.str) == 0, "%v", fmt.Errorf("config data is empty, Must Init"))

	field := path(paths...)

	switch cfg.typ {

	case "json":
		if r := gjson.Get(cfg.str, field); r.IsArray() {
			ret := FilterMap(r.Array(), nil, func(i int, e gjson.Result) any {
				if rt, ok := AnyTryToType[T](e.Value()); ok {
					return rt
				}
				return e.Value()
			})
			return AnysToTypes[T](ret)
		}

	case "toml":
		ks, _ := MapToKVs(cfg.fm, nil, nil)
		ks = FilterMap(strs.SortPaths(ks...),
			func(i int, e string) bool { return strings.HasPrefix(e, field+".") },
			func(i int, e string) string { return e },
		)
		r := []any{}
		for _, k := range ks {
			r = append(r, cfg.fm[k])
		}
		ret, ok := AnysTryToTypes[T](r)
		lk.FailOnErrWhen(!ok, "%v", fmt.Errorf("cannot get value as array of [%T]", *new(T)))
		return ret

	default:
		panic("TODO: add more type for text content to get value array")
	}

	lk.FailP1OnErr("%v", fmt.Errorf("[%s] is not array type", field))
	return nil
}

// for object value
func ValObj(paths ...any) map[string]any {

	lk.FailP1OnErrWhen(len(cfg.str) == 0, "%v", fmt.Errorf("config data is empty, Must do 'Init'"))

	field := path(paths...)

	switch cfg.typ {

	case "json":
		if r := gjson.Get(cfg.str, field); r.IsObject() {
			rt := make(map[string]any)
			lk.FailP1OnErr("%v", json.Unmarshal([]byte(r.Raw), &rt))
			return rt
		}

	case "toml":
		ks, _ := MapToKVs(cfg.fm, nil, nil)
		ks = FilterMap(strs.SortPaths(ks...),
			func(i int, e string) bool { return strings.HasPrefix(e, field+".") },
			func(i int, e string) string { return e },
		)
		mr := make(map[string]any)
		for _, k := range ks {
			mr[k] = cfg.fm[k]
		}
		return MapFlatToNested(mr, nil)

	default:

	}

	lk.FailP1OnErr("%v", fmt.Errorf("[%s] is not object type", field))
	return nil
}

func CntArr[T any](paths ...any) int {
	field := path(paths...)
	return len(ValArr[T](field))
}
