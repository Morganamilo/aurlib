package main

import "fmt"
import "os/exec"
import "os"
import "github.com/Morganamilo/aurlib/fetch"

func main() {
	pkgs := []string{"yay", "aurman", "aurutils", "yaourt", "pikaur", "libc++"}

	pwd, err := os.Getwd()
	checkerr(err)

	handle := fetch.MakeHandle(pwd+"/.clone", pwd+"/.patch")
	err = os.MkdirAll(pwd+"/.clone", 0755)
	checkerr(err)

	err = os.MkdirAll(".patch", 0755)
	checkerr(err)

	fetched, err := handle.Download(pkgs)
	checkerr(err)

	toMerge, err := handle.NeedMerge(fetched)
	checkerr(err)

	//err = handle.Diffs(toMerge)
	err = handle.DiffsToFile(toMerge)
	checkerr(err)

	tmp, err := handle.Link(pkgs)
	checkerr(err)

	if tmp != "" {
		cmd := exec.Command("vifm", tmp)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		checkerr(err)
		os.RemoveAll(tmp)
	}

	err = handle.Merge(toMerge)
	checkerr(err)

	handle.CleanDiffs(pkgs)
}

func checkerr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
