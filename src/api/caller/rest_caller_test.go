package caller

import (
	"testing"
	"github.com/JuanAller/request-builder/src/api/builder"
	"net/http"
	"errors"
)

type ExecutableMock struct {
	mockExecute func(entityResponse interface{}) *builder.Response
	totalCalls  int
}

func (rbm *ExecutableMock) restartCalls() {
	rbm.totalCalls = 0
}

func (rbm *ExecutableMock) Execute(entityResponse interface{}) *builder.Response {
	rbm.totalCalls++
	return rbm.mockExecute(entityResponse)
}

func TestRestCaller_ExecuteCall(t *testing.T) {
	cases := []struct {
		name            string
		executable      *ExecutableMock
		entity          interface{}
		retries         int
		responseHandler func(resp *builder.Response) error
		expectedError   string
		withRetries     bool
	}{
		{
			name: "not error",
			executable: &ExecutableMock{
				mockExecute: func(entityResponse interface{}) *builder.Response {
					return &builder.Response{
						StatusCode: http.StatusOK,
					}
				},
			},
			entity:  nil,
			retries: 0,
			responseHandler: func(resp *builder.Response) error {
				if resp.StatusCode != http.StatusOK {
					return errors.New("an error")
				}
				return nil
			},
			expectedError: "",
		},
		{
			name: "error_with_no_retries",
			executable: &ExecutableMock{
				mockExecute: func(entityResponse interface{}) *builder.Response {
					return &builder.Response{
						StatusCode: http.StatusServiceUnavailable,
					}
				},
			},
			entity:  nil,
			retries: 0,
			responseHandler: func(resp *builder.Response) error {
				if resp.StatusCode == http.StatusServiceUnavailable {
					return errors.New("an error")
				}
				return nil
			},
			expectedError: "an error",
		},
		{
			name: "error_with_retries",
			executable: &ExecutableMock{
				mockExecute: func(entityResponse interface{}) *builder.Response {
					return &builder.Response{
						StatusCode: http.StatusServiceUnavailable,
					}
				},
			},
			entity:  nil,
			retries: 2,
			responseHandler: func(resp *builder.Response) error {
				if resp.StatusCode == http.StatusServiceUnavailable {
					return errors.New("an error")
				}
				return nil
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
			}

			err := restCaller.ExecuteCall()

			if c.withRetries {
				if restCaller.Retries != c.retries {
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
