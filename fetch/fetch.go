package fetch

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

import (
	"github.com/Morganamilo/aurlib/utils/errors"
	"github.com/Morganamilo/aurlib/utils/slice"
)

func (h Handle) URL(pkgbase string) string {
	return h.AURURL + "/" + pkgbase + ".git"
}

func (h Handle) Download(pkgs []string) ([]string, error) {
	var wg sync.WaitGroup
	err := &errors.MultiError{}
	var mux sync.Mutex
	fetched := make([]string, 0)

	f := func(pkg string) {
		defer wg.Done()
		out, clone, e := h.GitHandle.Download(h.URL(pkg), h.CacheDir, pkg)
		err.Add(e)

		if !clone {
			mux.Lock()
			fetched = append(fetched, pkg)
			mux.Unlock()
		}

		fmt.Print(out)
	}

	for _, pkg := range pkgs {
		wg.Add(1)
		go f(pkg)
	}

	wg.Wait()

	return fetched, err.Return()
}

func (h Handle) NeedMerge(pkgs []string) ([]string, error) {
	toMerge := make([]string, 0)

	for _, pkg := range pkgs {
		needMerge, err := h.GitHandle.NeedMerge(h.CacheDir, pkg)
		if err != nil {
			return toMerge, err
		}

		if needMerge {
			toMerge = append(toMerge, pkg)
		}
	}

	return toMerge, nil
}

func (h Handle) Merge(pkgs []string) error {
	var wg sync.WaitGroup
	err := &errors.MultiError{}
	format := fmt.Sprintf("%%-%ds :  %%s", slice.Longest(pkgs))

	f := func(pkg string) {
		defer wg.Done()
		out, e := h.GitHandle.Merge(h.CacheDir, pkg)
		err.Add(e)

		index := strings.Index(out, "\n")
		if index != -1 {
			fmt.Printf(format, pkg, out[:index+1])
		} else {
			fmt.Printf(format+"\n", pkg, out)
		}
	}

	for _, pkg := range pkgs {
		wg.Add(1)
		go f(pkg)
	}

	wg.Wait()

	return err.Return()
}

func (h Handle) Diffs(pkgs []string) error {
	for _, pkg := range pkgs {
		out, err := h.GitHandle.Diff(h.CacheDir, true, pkg)
		if err != nil {
			return err
		}

		fmt.Print(out)
	}

	return nil
}

func (h Handle) DiffsToFile(pkgs []string) error {
	for _, pkg := range pkgs {
		out, err := h.GitHandle.Diff(h.CacheDir, false, pkg)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(filepath.Join(h.PatchDir, pkg)+".diff", []byte(out), 0644)

		if err != nil {
			return err
		}
	}

	return nil
}

func (h Handle) Link(pkgs []string) (string, error) {
	tmp, err := ioutil.TempDir("", "aur.")
	linked := false

	for _, pkg := range pkgs {
		path := filepath.Join(h.CacheDir, pkg)
		tmpPath := filepath.Join(tmp, pkg)
		_, err = os.Stat(path)

		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			goto err
		}

		err = os.Symlink(path, tmpPath)
		if err != nil && !os.IsNotExist(err) {
			goto err
		}

		linked = true

		path = filepath.Join(h.CacheDir, pkg, "PKGBUILD")
		tmpPath = filepath.Join(tmp, pkg+"-"+"PKGBUILD")
		_, err = os.Stat(path)

		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			goto err
		}

		err = os.Symlink(path, tmpPath)
		if err != nil && !os.IsNotExist(err) {
			goto err
		}

		path = filepath.Join(h.CacheDir, pkg, ".SRCINFO")
		tmpPath = filepath.Join(tmp, pkg+"-"+"SRCINFO")
		_, err = os.Stat(path)

		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			goto err
		}

		err = os.Symlink(path, tmpPath)
		if err != nil && !os.IsNotExist(err) {
			goto err
		}

		path = filepath.Join(h.PatchDir, pkg) + ".diff"
		tmpPath = filepath.Join(tmp, pkg) + ".diff"
		_, err = os.Stat(path)

		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			goto err
		}

		err = os.Symlink(path, tmpPath)
		if err != nil && !os.IsNotExist(err) {
			goto err
		}
	}

	if !linked {
		os.RemoveAll(tmp)
		tmp = ""
	}

	return tmp, nil

err:
	os.RemoveAll(tmp)
	return "", fmt.Errorf("Failed to link build files: %s", err.Error())
}

func (h Handle) CleanDiffs(pkgs []string) error {
	errs := &errors.MultiError{}

	for _, pkg := range pkgs {
		path := filepath.Join(h.PatchDir, pkg+".diff")
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			errs.Add(err)
			continue
		}

		err = os.Remove(path)
		if err != nil {
			errs.Add(err)
		}
	}

	return errs.Return()
}
