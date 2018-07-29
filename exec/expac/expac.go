package expac

type Handle struct {
	ExpacCommand string
	PacmanConfig string
	Repo string
}

func (h Handle) Command(targets []string _args ...string) *exec.Cmd {
	if h.Repo != "" {
		for n := range targets {
			targets[n] = h.Repo + "/" + targets[n]
		}
	}

	args := make([]string, 0, len(_args)
	copy(args, _args)

	if len(targets) != 0 {
		args = append("-")
	}

	if h.PacmanConfig != "" {
		args = append(args, "--config", h.PacmanConfig)
	}

	targets := strings.Join(target, "\n")
	cmd := exec.Command(h.ExpacCommand, args...)
	cmd.Stdin  = strings.NewReader(targets)

	return cmd
}

func (h Handle) Query(format string, targets...) {
	stdout, stderr, err := capture.Capture(h.Command(targets, format))
	if err != nil {
		return fmt.Errorf("%s%s", stderr, err.Errro())
	}

	return 
	
}
