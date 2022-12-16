package goconfig

import (
	"fmt"
	"strings"
	"sync"

	ff "github.com/digisan/fileflatter"
	. "github.com/digisan/go-generics/v2"
	"github.com/digisan/gotk/strs"
	lk "github.com/digisan/logkit"
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

// return first part as number; remainder string; if first part can be number
func startWithNum(path string) (int, string, bool) {
	ss := strings.Split(path, ".")
	switch {
	case len(ss) == 0:
		return -1, "", false
	default:
		n, ok := AnyTryToType[int](ss[0])
		return n, strings.Join(ss[1:], "."), ok
	}
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

	field := path(paths...)

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
}

// for object value
func Object(paths ...any) map[string]any {

	field := path(paths...)

	ks, _ := MapToKVs(cfg.fm, nil, nil)
	ks = FilterMap(strs.SortPaths(ks...),
		func(i int, e string) bool { return strings.HasPrefix(e, field+".") },
		func(i int, e string) string { return e },
	)
	mr := make(map[string]any)
	for _, k := range ks {
		key := strings.TrimPrefix(k, field+".")
		_, _, ok := startWithNum(key)
		lk.FailOnErrWhen(ok, "%v", fmt.Errorf("not an object under '%s'", key))
		mr[key] = cfg.fm[k]
	}
	return MapFlatToNested(mr, nil)
}

func Objects(paths ...any) []map[string]any {

	field := path(paths...)

	ks, _ := MapToKVs(cfg.fm, nil, nil)
	ks = FilterMap(strs.SortPaths(ks...),
		func(i int, e string) bool { return strings.HasPrefix(e, field+".") },
		func(i int, e string) string { return e },
	)

	indices := []int{}
	for _, k := range ks {
		key := strings.TrimPrefix(k, field+".")
		idx, _, ok := startWithNum(key)
		lk.FailOnErrWhen(!ok, "%v", fmt.Errorf("not an array under '%s'", key))
		indices = append(indices, idx)
	}

	// create flat map array
	N := Max(indices...) + 1
	fms := make([]map[string]any, N)
	for i := 0; i < N; i++ {
		fms[i] = make(map[string]any)
	}
	for _, k := range ks {
		key := strings.TrimPrefix(k, field+".")
		idx, key, _ := startWithNum(key)
		fms[idx][key] = cfg.fm[k]
	}

	// prepare return map array
	rt := make([]map[string]any, N)
	for i := 0; i < N; i++ {
		rt[i] = MapFlatToNested(fms[i], nil)
	}
	return rt
}

func CntValArr[T any](paths ...any) int {
	field := path(paths...)
	return len(ValArr[T](field))
}

func CntObjects(paths ...any) int {
	return len(Objects(paths...))
}

//////////////////////////////////////////////////////////

var (
	Bool   = Val[bool]
	Bools  = ValArr[bool]
	Str    = Val[string]
	Strs   = ValArr[string]
	Int    = Val[int]
	Ints   = ValArr[int]
	Uint   = Val[uint]
	Uints  = ValArr[uint]
	Float  = Val[float64]
	Floats = ValArr[float64]
)
