package caller

import (
	"time"
	"github.com/JuanAller/request-builder/src/api/builder"
)

type ResponseHandler func(resp *builder.Response) (err error, retry bool)
type BackOffStrategy func(retryNumber int) time.Duration

type restCaller struct {
	requestBuilder  ExecutableRequest
	Entity          interface{}
	responseHandler ResponseHandler
	retries         int
	backOff         BackOffStrategy
}

func NewRestCaller(r ExecutableRequest, entity interface{}, rh ResponseHandler, retries int, bos BackOffStrategy) *restCaller {
	return &restCaller{
		requestBuilder:  r,
		Entity:          entity,
		responseHandler: rh,
		retries:         retries,
		backOff:         bos,
	}
}

func (c *restCaller) ExecuteCall() error {
	err, retry := c.responseHandler(c.requestBuilder.Execute(c.Entity))
	for i := 1; i <= c.retries; i++ {
		if err == nil {
			return nil
		}
		if !retry {
			return err
		}
		time.Sleep(c.backOff(i))
		err, retry = c.responseHandler(c.requestBuilder.Execute(c.Entity))
	}
	return err
}
