package main

import (
	"fmt"

	cfg "github.com/digisan/go-config"
)

func main() {
	cfg.Init(true, "../config1.json", "../config.json")

	cfg.Show()

	fmt.Println("Bool", cfg.Val[bool]("Bool"))
	fmt.Println("IP", cfg.Val[string]("web", "IP"))
	fmt.Println("Port", cfg.Val[int16]("web.Port"))
	fmt.Println("element", cfg.Val[int8]("web.Array", 1))
	fmt.Println("array", cfg.ValArr[int8]("web.Array"))
	fmt.Println("object(expert-field2)", cfg.ValObj("expert")["field2"])
	fmt.Println("array count (expert)", cfg.CntArr[any]("expert-array"))
	fmt.Println("array count (simple)", cfg.CntArr[int]("simple-array"))

	// fmt.Println("Not Existing", cfg.Val[string]("Myname"))
	// fmt.Println(cfg.Val[int]("Port"))
}
