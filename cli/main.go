package main

import (
	"fmt"

	cfg "github.com/digisan/go-config"
)

func main() {
	cfg.Init(true, "../config.json")

	cfg.Show()

	fmt.Println("Bool", cfg.Val[bool]("Bool"))
	fmt.Println("IP", cfg.Val[string]("web.IP"))
	fmt.Println("Port", cfg.ValInt("web.Port"))

	// fmt.Println(cfg.Val[string]("Myname"))
	// fmt.Println(cfg.Val[int]("Port"))
}
