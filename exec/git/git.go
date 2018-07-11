package git

import (
	"fmt"
	"path/filepath"
	"os"
	"os/exec"
	"strings"
)

import (
	"github.com/Morganamilo/aurlib/exec/capture"
)

func (h Handle) Command(dir string, command string, _args ...string) *exec.Cmd {
	args := make([]string, 0, len(h.GitArgs) + 2 + len(_args) + len(h.GitCommandArgs))
	args = append(args, h.GitArgs...)
	args = append(args, "-C", dir)
	args = append(args, command)
	args = append(args, h.GitCommandArgs...)
	args = append(args, _args...)

	cmd := exec.Command(h.GitCommand, args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.Env = append(cmd.Env, h.GitEnvironment...)

	return cmd
}

func (h Handle) Download(url string, path string, name string) (bool, error) {
	_, err := os.Stat(filepath.Join(path, name, ".git"))
	if os.IsNotExist(err) {
		err = h.Command(path, "clone", url, name).Run()
		if err != nil {
			return false, fmt.Errorf("error cloning %s", name)
		}

		return true, nil
	} else if err != nil {
		return false, fmt.Errorf("error reading %s", filepath.Join(path, name, ".git"))
	}

	err = h.Command(filepath.Join(path, name), "fetch", "-v").Run()
	if err != nil {
		return false, fmt.Errorf("error fetching %s", name)
	}

	return false, nil
}

func (h Handle) Merge(path string, name string) error {
	err := h.Command(filepath.Join(path, name), "reset", "--hard", "HEAD").Run()
	if err != nil {
		return fmt.Errorf("error resetting %s", name)
	}

	err = h.Command(filepath.Join(path, name), "merge", "--no-edit", "--ff").Run()
	if err != nil {
		return fmt.Errorf("error merging %s", name)
	}

	return nil
}

func (h Handle) NeedMerge(path string, name string) (bool, error) {
	lines, err := h.RevParse(filepath.Join(path, name), "HEAD", "HEAD@{upstream}")
	if err != nil {
		return false, err
	}

	head := lines[0]
	upstream := lines[1]

	return head != upstream, nil
}

func (h Handle) RevParse(path string, args ...string) ([]string, error) {
	stdout, stderr, err := capture.Capture(h.Command(path, "rev-parse", "HEAD", "HEAD@{upstream}"))
	if err != nil {
		return nil, fmt.Errorf("%s%s", stderr, err)
	}

	return strings.Split(stdout, "\n"), nil
}
