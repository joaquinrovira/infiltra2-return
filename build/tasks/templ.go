package tasks

import (
	"fmt"
	"os"

	"github.com/goyek/goyek/v2"
	"github.com/goyek/x/cmd"
)

var Templ = goyek.Define(goyek.Task{
	Name:  "templ",
	Usage: "runs a-h/templ compiler",
	Deps:  goyek.Deps{installTempl},
	Action: func(a *goyek.A) {
		cmd.Exec(a, fmt.Sprintf("%s/templ generate", *toolDir))
	},
})

var installTempl = goyek.Define(goyek.Task{
	Name:  "templ-install",
	Usage: "installs the a-h/templ compiler",
	Action: func(a *goyek.A) {
		cwd, err := os.Getwd()
		if err != nil {
			a.Fatal(err)
		}
		command := fmt.Sprintf("GOBIN=\"%s/%s\" go install github.com/a-h/templ/cmd/templ@latest", cwd, *toolDir)
		cmd.Exec(a, command)
	},
})
