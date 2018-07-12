package caller

import (
	"testing"
	"github.com/JuanAller/request-builder/src/api/builder"
	"net/http"
	"errors"
	"time"
	"fmt"
)

type executableMock struct {
	mockExecute func(entityResponse interface{}) *builder.Response
	totalCalls  int
}

func (rbm *executableMock) restartCalls() {
	rbm.totalCalls = 0
}

func (rbm *executableMock) Execute(entityResponse interface{}) *builder.Response {
	rbm.totalCalls++
	return rbm.mockExecute(entityResponse)
}

func TestRestCaller_ExecuteCall(t *testing.T) {
	cases := []struct {
		name            string
		executable      *executableMock
		entity          interface{}
		retries         int
		responseHandler func(resp *builder.Response) (error, bool)
		backOff         func(retry int) time.Duration
		expectedError   string
		withRetries     bool
	}{
		{
			name: "not error",
			executable: &executableMock{
				mockExecute: func(entityResponse interface{}) *builder.Response {
					return &builder.Response{
						StatusCode: http.StatusOK,
					}
				},
			},
			entity:  nil,
			retries: 0,
			responseHandler: func(resp *builder.Response) (error, bool) {
				if resp.StatusCode != http.StatusOK {
					return errors.New("an error"), true
				}
				return nil, false
			},
			backOff: func(retry int) time.Duration {
				return time.Millisecond * 0
			},
			expectedError: "",
		},
		{
			name: "error_with_no_retries",
			executable: &executableMock{
				mockExecute: func(entityResponse interface{}) *builder.Response {
					return &builder.Response{
						StatusCode: http.StatusServiceUnavailable,
					}
				},
			},
			entity:  nil,
			retries: 0,
			responseHandler: func(resp *builder.Response) (error, bool) {
				if resp.StatusCode == http.StatusServiceUnavailable {
					return errors.New("an error"), false
				}
				return nil, false
			},
			backOff: func(retry int) time.Duration {
				return time.Millisecond * 0
			},
			expectedError: "an error",
		},
		{
			name: "error_with_retries",
			executable: &executableMock{
				mockExecute: func(entityResponse interface{}) *builder.Response {
					return &builder.Response{
						StatusCode: http.StatusServiceUnavailable,
					}
				},
			},
			entity:  nil,
			retries: 2,
			responseHandler: func(resp *builder.Response) (error, bool) {
				if resp.StatusCode == http.StatusServiceUnavailable {
					return errors.New("an error"), true
				}
				return nil, false
			},
			backOff: func(retry int) time.Duration {
				return time.Millisecond * 0
			},
			expectedError: "an error",
			withRetries:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			restCaller := &RestCaller{
				RequestBuilder:  c.executable,
				Entity:          c.entity,
				ResponseHandler: c.responseHandler,
				Retries:         c.retries,
				BackOff:         c.backOff,
			}

			err := restCaller.ExecuteCall()

			if c.withRetries {
				if c.executable.totalCalls-1 != c.retries {
					fmt.Println(c.executable.totalCalls)
					c.executable.restartCalls()
					t.Errorf("Retries fail")
				}
			}
			c.executable.restartCalls()

			if err != nil {
				if err.Error() != c.expectedError {
					t.Errorf("Error expected %s , but got %s", c.expectedError, err.Error())
				}
			} else {
				if c.expectedError != "" {
					t.Errorf("Expected error")
				}
			}
		})
	}
}
