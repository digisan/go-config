package goconfig

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	ff "github.com/digisan/fileflatter"
	. "github.com/digisan/go-generics/v2"
	dt "github.com/digisan/gotk/data-type"
	"github.com/digisan/gotk/strs"
	lk "github.com/digisan/logkit"
	"github.com/tidwall/sjson"
)

type confirmOrder uint8

const (
	first confirmOrder = 1
	final confirmOrder = 2
)

func (ct confirmOrder) String() string {
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

// m: original config map on 'first'; modified config map on 'final'
// return flat-map is without prompt fields
func confirm(cfgName, cfgType string, m map[string]any, co confirmOrder) (map[string]any, bool) {
	fmt.Printf(`
--------------------------------------------
    --- %s [%s] values ---        
--------------------------------------------`, co, cfgName)
	fmt.Println()

	mn := MapFlatToNested(m, func(path string, value any) (p string, v any) {
		p, v = path, value
		if strings.HasPrefix(path, "_") || strings.Contains(path, "._") {
			p = ""
		}
		return
	})

	// preview content ...
	switch cfgType {

	case "json":
		jsBytes, err := json.MarshalIndent(mn, "", "    ")
		lk.FailOnErr("%v", err)
		fmt.Println(string(jsBytes))

	case "toml":
		// "github.com/BurntSushi/toml"
		buf := new(bytes.Buffer)
		lk.FailOnErr("%v", toml.NewEncoder(buf).Encode(mn))
		fmt.Println(buf.String())

	default:
		panic("TODO: add more type to confirm preview")
	}

	//////////////////////////////////////

	if inputJudge("confirm?") {
		return MapNestedToFlat(mn), true
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

func Init(id string, prompt bool, fPaths ...string) (err error) {
	mtx.Lock()
	defer mtx.Unlock()

	MapCfg[id] = &Cfg{path: "", str: "", fm: make(map[string]any)}
	cfg = MapCfg[id]

	var (
		data []byte
	)

	for _, fpath := range fPaths {
		if bytes, err := os.ReadFile(fpath); err == nil {
			data, cfg.path, cfg.str = bytes, fpath, string(bytes)
			cfg.typ = dt.DataType(cfg.str) // "json", "toml", etc.
			// fmt.Printf("config type: %s\n", cfg.typ)
			break
		}
	}
	if err != nil || data == nil {
		return fmt.Errorf("failed to load configure file from [%v]", fPaths)
	}

	//
	cfg.fm, err = ff.FlatContent(data, false)
	lk.FailOnErr("%v", err)

	if !prompt {
		return
	}

	//////////////////////////////////////////////////////////////////////////

	fields, prompts := getPrompts(cfg.fm)

	// if no prompt fields, return
	if len(prompts) == 0 {
		return
	}

	lk.Log("prompts: %v", prompts)

	// check config value & type
	// for k, v := range cfg.fm {
	// 	fmt.Printf("%v(%T) - %v(%T)\n", k, k, v, v)
	// }

	if m, ok := confirm(filepath.Base(cfg.path), cfg.typ, cfg.fm, first); ok {
		cfg.fm = m
		return
	}

RE_INPUT_ALL:
	fmt.Printf(`
----------------------------------------------------------------
input value for [%s]. if <ENTER>, default value applies
----------------------------------------------------------------`, filepath.Base(cfg.path))
	fmt.Println()

	for i, f := range prompts {
		// f is prompt name (e.g. _IP)

		field := fields[i]           // real value field name
		var fVal any = cfg.fm[field] // real field value

		switch fVal.(type) {
		case int, int64, float32, float64, bool:
			fmt.Printf("--> %-20v\t\tvalue: %v\t\tnew value: ", cfg.fm[f], fVal)
		default:
			fmt.Printf("--> %-20v\t\tvalue: '%v'\t\tnew value: ", cfg.fm[f], fVal)
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
			if cfg.fm[field], err = strconv.ParseInt(iVal, 10, 64); err != nil {
				if cfg.fm[field], err = strconv.ParseFloat(iVal, 64); err != nil {
					fmt.Printf("[%v] is invalid, MUST be number, try again\n", iVal)
					goto RE_INPUT
				}
			}
		case bool:
			if cfg.fm[field], err = strconv.ParseBool(iVal); err != nil {
				fmt.Printf("[%v] is invalid, MUST be bool, try again\n", iVal)
				goto RE_INPUT
			}
		default:
			cfg.fm[field] = iVal
		}
		if err != nil {
			panic(err)
		}
	}

	if fm, ok := confirm(filepath.Base(cfg.path), cfg.typ, cfg.fm, final); ok {
		if inputJudge("Overwrite Original File? (format & order may change & comments will disappear!)") {

			// fetch original file content
			ori := string(data)

			// original flat map, which has prompt fields
			fmOri, err := ff.FlatContent(ori, false)
			lk.FailOnErr("%v", err)

			switch cfg.typ {

			case "json":
				// modify original
				for k, v := range fm {
					ori, err = sjson.Set(ori, k, v)
					lk.FailOnErr("%v", err)
				}
				lk.FailOnErr("%v", os.WriteFile(cfg.path, []byte(ori), os.ModePerm))

			case "toml":
				// modify original
				for k, v := range fm {
					fmOri[k] = v
				}

				// "github.com/BurntSushi/toml"
				buf := new(bytes.Buffer)
				lk.FailOnErr("%v", toml.NewEncoder(buf).Encode(MapFlatToNested(fmOri, nil)))
				lk.FailOnErr("%v", os.WriteFile(cfg.path, buf.Bytes(), os.ModePerm))

			default:
				panic("TODO: add more type to confirm preview")
			}
		}
		cfg.fm = fm
		return
	}
	fmt.Println("INPUT AGAIN PLEASE:")
	goto RE_INPUT_ALL
}
