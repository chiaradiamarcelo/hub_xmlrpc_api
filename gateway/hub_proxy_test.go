package gateway

import (
	"errors"
	"reflect"
	"testing"
)

func Test_ProxyCallToHub(t *testing.T) {
	tt := []struct {
		name             string
		args             []interface{}
		mockExecuteCall  func(hubAPIEndpoint string, call string, args []interface{}) (response interface{}, err error)
		expectedResponse interface{}
		expectedErr      string
	}{
		{
			name: "ProxyCallToHub call_successful",
			args: []interface{}{"arg1", "arg2"},
			mockExecuteCall: func(hubAPIEndpoint string, call string, args []interface{}) (response interface{}, err error) {
				return "success_response", nil
			},
			expectedResponse: "success_response",
		},
		{
			name: "ProxyCallToHub call_error",
			args: []interface{}{"arg1", "arg2"},
			mockExecuteCall: func(hubAPIEndpoint string, call string, args []interface{}) (response interface{}, err error) {
				return nil, errors.New("call_error")
			},
			expectedErr: "call_error",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mockUyuniCallExecutor := new(mockUyuniCallExecutor)
			mockUyuniCallExecutor.mockExecuteCall = tc.mockExecuteCall

			hubProxy := NewHubProxy("hub_API_endpoint", mockUyuniCallExecutor)

			response, err := hubProxy.ProxyCallToHub("call", tc.args)

			if err != nil && tc.expectedErr != err.Error() {
				t.Fatalf("Error during executing request: %v", err)
			}
			if err == nil && !reflect.DeepEqual(response, tc.expectedResponse) {
				t.Fatalf("Expected and actual values don't match, Expected value is: %v", tc.expectedResponse)
			}
		})
	}
}
