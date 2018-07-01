package caller

import (
	"github.com/JuanAller/request-builder/src/api/builder"
)

type RestCaller struct {
	RequestBuilder  Executable
	Entity          interface{}
	ResponseHandler func(resp *builder.Response) (err error, retry bool)
	Retries         int
}

func (c *RestCaller) ExecuteCall() error {
	var err error
	var retry bool
	for i := 0; i <= c.Retries; i++ {
		err, retry = c.ResponseHandler(c.RequestBuilder.Execute(c.Entity))
		if err == nil {
			return nil
		}
		if !retry {
			return err
		}
	}
	return err
}
