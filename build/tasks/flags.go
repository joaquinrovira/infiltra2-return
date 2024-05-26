package tasks

import "flag"

var outdir = flag.String("outdir", "bin", `binary output directory`)
var outname = flag.String("outname", "", `binary output name`)
var target = flag.String("t", "./app", `target package (relative path)`)
var debug = flag.Bool("debug", false, "build for debugging")
var toolDir = flag.String("tools", ".tools", `directory for storing utilities`)