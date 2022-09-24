package main

import (
	"fmt"

	cfg "github.com/digisan/go-config"
)

func main() {
	cfg.Init(true, "../config1.json", "../config.json")

	cfg.Show()

	fmt.Println("Bool", cfg.Val[bool]("Bool"))
	fmt.Println("IP", cfg.Val[string]("web.IP"))
	fmt.Println("Port", cfg.Val[int16]("web.Port"))
	fmt.Println("array", cfg.ValArr[int8]("web.Array"))

	fmt.Println("Not Existing", cfg.Val[string]("Myname"))
	fmt.Println(cfg.Val[int]("Port"))
}
