package output

import (
	"bytes"
	"os"
	"os/exec"
)

func Capture(cmd *exec.Cmd) (string, string, error) {
	var outbuf, errbuf bytes.Buffer

	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err := cmd.Run()
	stdout := outbuf.String()
	stderr := errbuf.String()

	return stdout, stderr, err
}

type LineWriter struct {
	prefix        string
	out           *os.File
	seenLine      bool
	printedPrefix bool
}

func MakeLineWriter(prefix string, file *os.File) LineWriter {
	return LineWriter{
		prefix,
		file,
		false,
		false,
	}
}

func (w LineWriter) Write(b []byte) (int, error) {
	var p int

	if w.seenLine {
		return 0, nil
	}

	if !w.printedPrefix {
		var e error
		p, e = w.out.Write([]byte(w.prefix))
		if e != nil {
			return p, e
		}
	}

	for n, r := range b {
		if r == '\n' {
			w.seenLine = true
			n, err := w.out.Write(b[:n+1])
			return n + p, err
		}
	}

	n, err := w.out.Write(b)
	return n + p, err
}

func Show(cmd *exec.Cmd) *exec.Cmd {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func Hide(cmd *exec.Cmd) *exec.Cmd {
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd
}

func ToStderr(cmd *exec.Cmd) *exec.Cmd {
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd
}
