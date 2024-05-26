package tasks

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/goyek/goyek/v2"
	"github.com/goyek/x/cmd"
)

func tailwindFile() string {
	return fmt.Sprintf("%s/tailwind", *toolDir)
}

var Tailwind = goyek.Define(goyek.Task{
	Name:  "tailwind",
	Usage: "runs tailwind compiler",
	Deps:  goyek.Deps{getTailwind},
	Action: func(a *goyek.A) {
		cmd.Exec(a, fmt.Sprintf("%s -i static/input.css -o static/public/tailwind.css", tailwindFile()))
	},
})

var getTailwind = goyek.Define(goyek.Task{
	Name:  "tailwind-install",
	Usage: "installs the tailwind compiler",
	Action: func(a *goyek.A) {
		if err := createBinDir(a); err != nil {
			a.Fatal(err)
		}
		if err := downloadTailwind(a); err != nil {
			a.Fatal(err)
		}
	},
})

func createBinDir(_ *goyek.A) error {
	return os.MkdirAll(*toolDir, os.ModePerm)
}

func downloadTailwind(a *goyek.A) error {
	file := tailwindFile()

	if _, err := os.Stat(file); err == nil {
		a.Log("binary already exists")
		return nil
	} else {
		a.Log("downloading binary")
	}

	out, err := initTailwindFile()
	if err != nil {
		return err
	}
	defer out.Close()

	url := tailwindDownloadURL()
	a.Log(url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func tailwindDownloadURL() string {
	system := runtime.GOOS
	system = strings.ReplaceAll(system, "darwin", "macos")
	arch := runtime.GOARCH
	arch = strings.ReplaceAll(arch, "amd64", "x64")

	return fmt.Sprintf("https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-%s-%s", system, arch)
}

func initTailwindFile() (*os.File, error) {
	file := tailwindFile()

	out, err := os.Create(file)
	if err != nil {
		return nil, err
	}

	err = os.Chmod(file, 0777)
	if err != nil {
		return nil, err
	}

	return out, nil
}
