package main

import (
	"github.com/vinarius/vin/cmd"
	_ "github.com/vinarius/vin/cmd/clean/loggroups"
	_ "github.com/vinarius/vin/cmd/clean/stacks"
	_ "github.com/vinarius/vin/cmd/configure"
	_ "github.com/vinarius/vin/cmd/configure/profile"
	_ "github.com/vinarius/vin/cmd/tms/watch"
)

func main() {
	cmd.Execute()
}
