package deps

import (
	"github.com/Morganamilo/aurlib/rpc"
	"github.com/mikkeloscar/aur"
)

type DepLevel int

const (
	DepLevelAll DepLevel = iota
	DepLevelNeeded
	DepLevelNone
)

type Handle struct {
	DepLevel DepLevel
	Versioned bool
	InfoLimit int
}

func MakeHandle() Handle {
	return Handle{
		DepLevel: DepLevelAll,
		Versioned: true,
		InfoLimit: 100,
	}
}


func (h Handle) Depends(pkgs []string) ([]*aur.Pkg, error) {
	info, err := rpc.Info(pkgs, h.InfoLimit)

	if err != nil {
		return nil, err
	}

	if h.DepLevel == DepLevelNone {
		return info, nil
	}

	return info, nil
}
