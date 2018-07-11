package git

type Handle struct {
	GitCommand string
	GitArgs []string
	GitCommandArgs []string
	GitEnvironment []string
}

func MakeHandle() Handle {
	handle := Handle{
		GitCommand: "git",
		GitEnvironment: []string{"GIT_TERMINAL_PROMPT=0"},
	}

	return handle
}
