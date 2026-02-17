package main

import (
	"github.com/vinarius/vin/cmd"
	_ "github.com/vinarius/vin/cmd/tms/watch"
)

func main() {
	cmd.Execute()
}
