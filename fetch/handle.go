package fetch

import (
	"github.com/Morganamilo/aurlib/exec/git"
)

type Handle struct {
	AURURL    string
	CacheDir  string
	PatchDir  string
	GitHandle git.Handle
}

func MakeHandle(cacheDir string, patchDir string) Handle {
	handle := Handle{
		AURURL:    "https://aur.archlinux.org",
		CacheDir:  cacheDir,
		PatchDir:  patchDir,
		GitHandle: git.MakeHandle(),
	}

	return handle
}
