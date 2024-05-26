package main

import (
	"github.com/joaquinrovira/infiltra2-returns/build/tasks"

	"github.com/goyek/goyek/v2"
	"github.com/goyek/x/boot"
)

// More info: https://github.com/goyek/goyek

func main() {
	goyek.SetDefault(tasks.GoBuild)
	boot.Main()
}
