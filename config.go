package goconfig

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	// . "github.com/digisan/go-generics/v2"
	"github.com/digisan/gotk/strs"
	jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
	"github.com/tidwall/sjson"
)

type confirmType uint8

const (
	first confirmType = 1
	final confirmType = 2
)

func (ct confirmType) String() string {
	switch ct {
	case first:
		return "default"
	case final:
		return "review"
	default:
		return "unknown"
	}
}

func inputJudge(prompt string) bool {
	fmt.Println(prompt + " [y/N]")
	input := ""
	_, err := fmt.Scanf("%s", &input)
	switch {
	case err == nil && strs.IsIn(true, true, input, "YES", "Y", "OK", "TRUE"):
		return true
	case err != nil && err.Error() == "unexpected newline" && len(input) == 0:
		return false
	default:
		return false
	}
}

// m: original config map on 'first'
//
//	modified config map on 'final'
func confirm(cfgName string, m map[string]any, ct confirmType) (map[string]any, bool) {
	fmt.Printf(`
--------------------------------------------
    --- %s [%s] values ---        
--------------------------------------------`, ct, cfgName)
	fmt.Println()

	cfg := jt.Composite(m, func(path string, value any) (p string, v any, raw bool) {
		p, v, raw = path, value, false                                    // if return 'raw' is true, it must be <string> type
		if strings.HasPrefix(path, "_") || strings.Contains(path, "._") { // && unicode.IsUpper(rune(path[0])) {
			p = ""
		}
		return
	})
	fmt.Println(jt.FmtStr(cfg, "    "))

	// trimmed config map
	m, err := jt.Flatten([]byte(cfg))
	lk.FailOnErr("%v", err)

	if inputJudge("confirm?") {
		return m, true
	}
	return nil, false
}

func getPrompts(m map[string]any) (fields, prompts []string) {
	r1 := regexp.MustCompile(`^_\w+`)
	r2 := regexp.MustCompile(`\._\w+`)
	for k := range m {
		if r1.MatchString(k) || r2.MatchString(k) {
			fk := strings.TrimPrefix(k, "_")
			fk = strings.ReplaceAll(fk, "._", ".")
			if _, ok := m[fk]; ok {
				fields = append(fields, fk)
				prompts = append(prompts, k)
			}
		}
	}
	return
}

func Init(prompt bool, fPaths ...string) (err error) {

	var (
		data []byte
	)

	for _, fpath := range fPaths {
		if bytes, err := os.ReadFile(fpath); err == nil {
			data, fPathCfg, jsCfg = bytes, fpath, string(bytes)
			break
		}
	}
	lk.FailP1OnErrWhen(err != nil || data == nil, "%v from %v", fmt.Errorf("failed to load configure file"), fPaths)

	mCfg, err = jt.Flatten(data)
	lk.FailOnErr("%v", err)

	if !prompt {
		return
	}

	//////////////////////////////////////////////////////////////////////////

	fields, prompts := getPrompts(mCfg)

	// if no prompt fields, return config json map
	if len(prompts) == 0 {
		return
	}

	lk.Log("prompts: %v", prompts)

	//
	// check config value & type
	//
	// for k, v := range mCfg {
	// 	fmt.Printf("%v(%T) - %v(%T)\n", k, k, v, v)
	// }

	if m, ok := confirm(filepath.Base(fPathCfg), mCfg, first); ok {
		mCfg = m
		return
	}

RE_INPUT_ALL:
	fmt.Printf(`
----------------------------------------------------------------
input value for [%s]. if <ENTER>, default value applies
----------------------------------------------------------------`, filepath.Base(fPathCfg))
	fmt.Println()

	for i, f := range prompts {
		// f is prompt name (e.g. _IP)

		field := fields[i]         // real value field name
		var fVal any = mCfg[field] // real field value

		switch fVal.(type) {
		case int, int64, float32, float64, bool:
			fmt.Printf("--> %-20v\t\tvalue: %v\t\tnew value: ", mCfg[f], fVal)
		default:
			fmt.Printf("--> %-20v\t\tvalue: '%v'\t\tnew value: ", mCfg[f], fVal)
		}

	RE_INPUT:
		var iVal string
		if scanner := bufio.NewScanner(os.Stdin); scanner.Scan() {
			iVal = scanner.Text()
		}

		if len(iVal) == 0 {
			continue
		}

		switch fVal.(type) {
		case int, int64, float32, float64:
			if mCfg[field], err = strconv.ParseInt(iVal, 10, 64); err != nil {
				if mCfg[field], err = strconv.ParseFloat(iVal, 64); err != nil {
					fmt.Printf("[%v] is invalid, MUST be number, try again\n", iVal)
					goto RE_INPUT
				}
			}
		case bool:
			if mCfg[field], err = strconv.ParseBool(iVal); err != nil {
				fmt.Printf("[%v] is invalid, MUST be bool, try again\n", iVal)
				goto RE_INPUT
			}
		default:
			mCfg[field] = iVal
		}
		if err != nil {
			panic(err)
		}
	}

	if m, ok := confirm(filepath.Base(fPathCfg), mCfg, final); ok {
		if inputJudge("Overwrite Original File?") {
			ori := string(data)
			for k, v := range m {
				ori, err = sjson.Set(ori, k, v)
				lk.FailOnErr("%v", err)
			}
			lk.FailOnErr("%v", os.WriteFile(fPathCfg, []byte(ori), os.ModePerm))
		}
		mCfg = m
		return
	}
	fmt.Println("INPUT AGAIN PLEASE:")
	goto RE_INPUT_ALL
}
