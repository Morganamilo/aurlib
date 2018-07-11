package capture

import (
	"bytes"
	"os/exec"
)

func Capture(cmd *exec.Cmd) (string, string, error) {
	var outbuf, errbuf bytes.Buffer

	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	cmd.Stdin = nil
	err := cmd.Run()
	stdout := outbuf.String()
	stderr := errbuf.String()

	return stdout, stderr, err
}

