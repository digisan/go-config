package main

import (
	"fmt"

	cfg "github.com/digisan/go-config"
)

func main() {

	cfg.Init("toml1", false, "../data/config1.toml")
	cfg.Use("toml1")
	cfg.Show()

	fmt.Println("InFolder:", cfg.Path("InFolder"))
	fmt.Println("Trim.OutFolder:", cfg.Path("Trim.OutFolder"))
	fmt.Println("Split.Enabled:", cfg.Bool("Split.Enabled"))
	fmt.Println("Split.Schema.0:", cfg.Str("Split.Schema.0"))
	fmt.Println("Split.Schema:", cfg.Strs("Split.Schema"))
	fmt.Println("Split.Schema:", cfg.ValArr[string]("Split.Schema"))

	fmt.Println("Trim:", cfg.Object("Trim"))

	fmt.Println("Merge:", cfg.Objects("Merge"))
	fmt.Println("Merge Count:", cfg.CntObjects("Merge"))

	fmt.Println("Trim.OutFolder:", cfg.Path("Trim.OutFolder"))
	fmt.Println("Trim.OutFolder:", cfg.PathAbs("Trim.OutFolder"))

	fmt.Println("--------------------------------------------------------------------------")

	cfg.Init("toml", false, "../data/config.toml")
	cfg.Use("toml")
	cfg.Show()

	fmt.Println("enabled:", cfg.Val[bool]("database.enabled"))
	fmt.Println("ports:", cfg.Val[int]("database.ports", 0))
	fmt.Println("ports:", cfg.Val[int64]("database.ports.0"))
	fmt.Println("array:", cfg.ValArr[int]("database.ports"))
	fmt.Println("array:", cfg.ValArr[int8]("clients.data.1"))
	fmt.Println("servers.alpha:", cfg.Object("servers", "alpha"))
	fmt.Println("servers:", cfg.Object("servers"))
	fmt.Println("time:", cfg.DateTime("owner.dob"))

	// // fmt.Println("array", cfg.ValArr[any]("clients.data")) // error@ [[], []] !!!

	fmt.Println("--------------------------------------------------------------------------")

	cfg.Init("json", false, "../data/config.json")
	cfg.Use("json")
	cfg.Show()

	fmt.Println("Bool", cfg.Val[bool]("Bool"))
	fmt.Println("IP", cfg.Val[string]("web", "IP"))
	fmt.Println("Port", cfg.Val[int16]("web.Port"))
	fmt.Println("element", cfg.Val[int8]("web.Array", 1))
	fmt.Println("array", cfg.ValArr[int8]("web.Array"))
	fmt.Println("object(expert-field2)", cfg.Object("expert")["field2"])
	fmt.Println("array (expert)", cfg.Objects("expert-array"))
	fmt.Println("array count (simple)", cfg.CntValArr[int]("simple-array"))
	fmt.Println("array", cfg.ValArr[string]("string-array"))

	if missing := "Myname"; cfg.HasVal(missing) {
		fmt.Println(cfg.Val[string](missing))
	} else {
		fmt.Printf("Not Existing: %s\n", missing)
	}

	fmt.Println("--------------------------------------------------------------------------")
}
