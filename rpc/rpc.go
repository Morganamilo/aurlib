package rpc

import (
	"sync"
)

import (
	"github.com/mikkeloscar/aur"
	"github.com/Morganamilo/aurlib/utils/errors"
	"github.com/Morganamilo/aurlib/utils/math"
)

func Info(names []string, limit int) ([]*aur.Pkg, error) {
	info := make([]*aur.Pkg, 0, len(names))
	err := &errors.MultiError{}
	var mux sync.Mutex
	var wg sync.WaitGroup

	makeRequest := func(n, max int) {
		defer wg.Done()
		tempInfo, requestErr := aur.Info(names[n:max])

		if requestErr != nil {
			err.Add(requestErr)
			return
		}

		mux.Lock()
		for _, _i := range tempInfo {
			i := _i
			info = append(info, &i)
		}
		mux.Unlock()
	}

	for n := 0; n < len(names); n += limit {
		max := math.Min(len(names), n+limit)
		wg.Add(1)
		go makeRequest(n, max)
	}

	wg.Wait()

	return info, err.Return()
}

