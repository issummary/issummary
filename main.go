//go:generate statik -src=./static/dist

package main

import (
	"github.com/mpppk/issummary/cmd"
	_ "github.com/mpppk/issummary/statik"
)

func main() {
	cmd.Execute()
}
