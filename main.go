package main

import "fmt"

var (
	// Version will be populated with binary semver by the linker
	// during the build process.
	// See https://blog.cloudflare.com/setting-go-variables-at-compile-time/
	// and the README file.
	Version string

	// Commit will be populated with correct git commit id by the linker
	// during the build process.
	// See https://blog.cloudflare.com/setting-go-variables-at-compile-time/
	// and the README file.
	Commit string
)

func main() {
	fmt.Println(fmt.Sprintf("Hello! from version %s and commit %s", Version, Commit))
}
