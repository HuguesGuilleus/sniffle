//go:build devmode

package tool

import "fmt"

const DevMode = true

func init() {
	fmt.Println("[DEV MODE ENABLED]")
}
