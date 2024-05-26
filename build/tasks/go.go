package tasks

import (
	"fmt"

	"github.com/goyek/goyek/v2"
	"github.com/goyek/x/cmd"
)

var GoBuild = goyek.Define(goyek.Task{
	Name:  "go-build",
	Usage: "runs go build",
	Action: func(a *goyek.A) {
		ExtraArgs := ""
		if *debug {
			ExtraArgs += " -gcflags=all=\"-N -l\""
		}
		cmd.Exec(a, fmt.Sprintf("go build -o %s/%s %s ./%s", *outdir, *outname, ExtraArgs, *target))
	},
	Deps: goyek.Deps{Templ, Tailwind},
})
