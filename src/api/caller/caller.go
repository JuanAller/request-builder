package caller

import (
	"github.com/JuanAller/request-builder/src/api/builder"
	"golang.org/x/sync/errgroup"
)

type ExecutableRequest interface {
	Execute(entityResponse interface{}) *builder.Response
}

type Caller interface {
	ExecuteCall() error
}

func InParallelCalls(callers ...Caller) error {
	var g errgroup.Group
	for _, caller := range callers {
		caller := caller
		g.Go(caller.ExecuteCall())
	}
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
