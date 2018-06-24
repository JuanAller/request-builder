package caller

import (
	"github.com/JuanAller/request-builder/src/api/builder"
	"golang.org/x/sync/errgroup"
)

type Executable interface {
	Execute(entityResponse interface{}) *builder.Response
}

type Caller interface {
	ExecuteCall() error
}

/**
 Execute N calls, if any fail, return error, and abort
 */
func InParallelCalls(callers ...Caller) error {
	var g errgroup.Group
	for _, caller := range callers {
		caller := caller
		g.Go(func() error {
			return caller.ExecuteCall()
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
