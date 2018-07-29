package main

import (
	"fmt"
)

import (
	"github.com/Morganamilo/aurlib/deps"
)

func main() {
	pkgs := []string{"yay", "discord"}

	handle := deps.MakeHandle()
	deps, err := handle.Depends(pkgs)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(deps)
}
