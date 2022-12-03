package goconfig

import "testing"

func TestAnalyzeConfig(t *testing.T) {
	Init("test", true, "./config.json")
}
