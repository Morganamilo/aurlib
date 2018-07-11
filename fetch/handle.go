package fetch

import (
	"github.com/Morganamilo/aurlib/exec/git"
)

type Handle struct {
	AURURL string
	CacheDir string
	GitHandle git.Handle
}

func MakeHandle(cacheDir string) Handle {
	handle := Handle{
		AURURL: "https://aur.archlinux.org",
		CacheDir: cacheDir,
		GitHandle: git.MakeHandle(),
	}

	return handle
}
