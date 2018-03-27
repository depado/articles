package main

import "github.com/Depado/articles/code/gochecklist/cmd"

// Build number and versions injected at compile time
var (
	Version = "unknown"
	Build   = "unknown"
)

func main() {
	cmd.Execute(Version, Build)
}
