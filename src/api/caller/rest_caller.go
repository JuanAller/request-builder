package caller

import (
	"github.com/JuanAller/request-builder/src/api/builder"
	"errors"
)

var ErrNotRetry = errors.New("Not retry")

/**
 Entity : struct where marshal response body
 ResponseHandler : handle status code, or errors
 Retries : number of retries, if request fail
 */
type RestCaller struct {
	RequestBuilder  Executable
	Entity          interface{}
	ResponseHandler func(resp *builder.Response) error
	Retries         int
}

func (c *RestCaller) ExecuteCall() error {
	var err error
	for i := 0; i <= c.Retries; i++ {
		err = c.ResponseHandler(c.RequestBuilder.Execute(c.Entity))
		if err == nil {
			return nil
		}
		if err == ErrNotRetry {
			return ErrNotRetry
		}
	}
	return err
}
