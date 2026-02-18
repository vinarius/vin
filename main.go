package main

import (
	"github.com/vinarius/vin/v2/cmd"
	_ "github.com/vinarius/vin/v2/cmd/clean/loggroups"
	_ "github.com/vinarius/vin/v2/cmd/clean/stacks"
	_ "github.com/vinarius/vin/v2/cmd/configure"
	_ "github.com/vinarius/vin/v2/cmd/configure/profile"
	_ "github.com/vinarius/vin/v2/cmd/tms/watch"
)

func main() {
	cmd.Execute()
}
