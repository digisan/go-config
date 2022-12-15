package main

import (
	"fmt"

	cfg "github.com/digisan/go-config"
)

func main() {
	cfg.Init("toml", true, "../data/config.toml")
	cfg.Use("toml")
	cfg.Show()

	fmt.Println("enabled:", cfg.Val[bool]("database.enabled"))
	fmt.Println("ports:", cfg.Val[int]("database.ports", 0))
	fmt.Println("ports:", cfg.Val[int64]("database.ports.0"))
	fmt.Println("array:", cfg.ValArr[int]("database.ports"))
	fmt.Println("array:", cfg.ValArr[int8]("clients.data.1"))
	fmt.Println("object:", cfg.ValObj("servers.alpha"))
	fmt.Println("object:", cfg.ValObj("servers"))

	// fmt.Println("array", cfg.ValArr[any]("clients.data")) // error !!!

	fmt.Println("------------------------------------------------------")

	// cfg.Init("json", true, "../data/config.json")
	// cfg.Use("json")
	// cfg.Show()

	// fmt.Println("Bool", cfg.Val[bool]("Bool"))
	// fmt.Println("IP", cfg.Val[string]("web", "IP"))
	// fmt.Println("Port", cfg.Val[int16]("web.Port"))
	// fmt.Println("element", cfg.Val[int8]("web.Array", 1))
	// fmt.Println("array", cfg.ValArr[int8]("web.Array"))
	// fmt.Println("object(expert-field2)", cfg.ValObj("expert")["field2"])
	// fmt.Println("array count (expert)", cfg.CntArr[any]("expert-array"))
	// fmt.Println("array count (simple)", cfg.CntArr[int]("simple-array"))
	// fmt.Println("array", cfg.ValArr[string]("string-array"))
	// fmt.Println("Not Existing", cfg.Val[string]("Myname"))
}
