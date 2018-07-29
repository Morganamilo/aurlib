package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

import (
	"github.com/Morganamilo/aurlib/exec/output"
)

func (h Handle) Command(dir string, command string, _args ...string) *exec.Cmd {
	args := make([]string, 0, len(h.GitArgs)+2+len(_args)+len(h.GitCommandArgs))
	args = append(args, h.GitArgs...)
	args = append(args, "-C", dir)
	args = append(args, command)
	args = append(args, h.GitCommandArgs...)
	args = append(args, _args...)

	cmd := exec.Command(h.GitCommand, args...)
	if len(h.GitEnvironment) > 0 {
		cmd.Env = append(os.Environ(), h.GitEnvironment...)
	}

	return cmd
}

func (h Handle) Download(url string, path string, name string) (string, bool, error) {
	_, err := os.Stat(filepath.Join(path, name, ".git"))
	if os.IsNotExist(err) {
		cmd := h.Command(path, "clone", "--no-progress", url, name)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return string(out), true, fmt.Errorf("error cloning %s: %s", name, out)
		}

		return string(out), true, nil
	} else if err != nil {
		return "", false, fmt.Errorf("error reading %s: %s", filepath.Join(path, name, ".git"), err.Error())
	}

	out, err := h.Command(filepath.Join(path, name), "fetch", "-v").CombinedOutput()
	if err != nil {
		return string(out), false, fmt.Errorf("error fetching %s: %s", name, out)
	}

	return string(out), false, nil
}

func (h Handle) Merge(path string, name string) (string, error) {
	out, err := h.Command(filepath.Join(path, name), "reset", "--hard", "-q", "HEAD").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error resetting %s: %s", name, out)
	}

	out, err = h.Command(filepath.Join(path, name), "merge", "--no-edit").CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("error merging %s: %s", name, out)
	}

	return string(out), nil
}

func (h Handle) Diff(path string, color bool, name string) (string, error) {
	out, err := h.Command(filepath.Join(path, name), "reset", "--hard", "HEAD").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error resetting %s: %s", name, out)
	}

	// --no-commit doesnt make a commit but still requires a email and name
	// to be set. Placeholder values are used just to that git doesn't
	// error.
	h.GitArgs = append(h.GitArgs, "-c", "user.email=aur", "-c", "user.name=aur")
	out, err = h.Command(filepath.Join(path, name), "merge", "--no-edit", "--no-ff", "--no-commit").CombinedOutput()
	h.GitArgs = h.GitArgs[:len(h.GitArgs)-4]
	if err != nil {
		return "", fmt.Errorf("error merging %s: %s", name, out)
	}

	colorWhen := "--color=always"
	if !color {
		colorWhen = "--color=never"
	}

	out1, err := h.Command(filepath.Join(path, name), "log", "HEAD..HEAD@{upstream}", colorWhen).CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("error diffing %s: %s", name, out)
	}

	out2, err := h.Command(filepath.Join(path, name), "diff", "--stat", "--patch", "--cached", colorWhen).CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("error diffing %s: %s", name, out)
	}

	return string(append(out1, out2...)), nil
}

func (h Handle) NeedMerge(path string, name string) (bool, error) {
	lines, err := h.RevParse(filepath.Join(path, name), "HEAD", "HEAD@{upstream}")
	fmt.Println(name, lines)
	if err != nil {
		return false, err
	}

	head := lines[0]
	upstream := lines[1]

	return head != upstream, nil
}

func (h Handle) RevParse(path string, args ...string) ([]string, error) {
	stdout, stderr, err := output.Capture(h.Command(path, "rev-parse", "HEAD", "HEAD@{upstream}"))
	if err != nil {
		return nil, fmt.Errorf("%s%s", stderr, err)
	}

	return strings.Split(stdout, "\n"), nil
}
