//go:generate statik -src=./static/dist

package main

import (
	"github.com/issummary/issummary/cmd"
	_ "github.com/issummary/issummary/statik"
)

func main() {
	cmd.Execute()
}
