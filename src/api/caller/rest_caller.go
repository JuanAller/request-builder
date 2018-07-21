package caller

import (
	"time"
	"github.com/JuanAller/request-builder/src/api/builder"
)

type ResponseHandler func(resp *builder.Response) (err error, retry bool)
type BackOffStrategy func(retryNumber int) time.Duration

type RestCaller struct {
	RequestBuilder  ExecutableRequest
	Entity          interface{}
	ResponseHandler ResponseHandler
	Retries         int
	BackOff         BackOffStrategy
}

func (c *RestCaller) ExecuteCall() error {
	err, retry := c.ResponseHandler(c.RequestBuilder.Execute(c.Entity))
	for i := 1; i <= c.Retries; i++ {
		if err == nil {
			return nil
		}
		if !retry {
			return err
		}
		time.Sleep(c.BackOff(i))
		err, retry = c.ResponseHandler(c.RequestBuilder.Execute(c.Entity))
	}
	return err
}
