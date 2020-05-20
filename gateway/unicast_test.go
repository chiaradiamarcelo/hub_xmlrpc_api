package gateway

import (
	"errors"
	"reflect"
	"strconv"
	"testing"
)

func Test_Unicast(t *testing.T) {
	mockRetrieveServerSessionByServerIDFound := func(hubSessionKey string, serverID int64) *ServerSession {
		strServerID := strconv.FormatInt(serverID, 10)
		return &ServerSession{serverID, strServerID + "serverAPIEndpoint", strServerID + "serverSessionkey", hubSessionKey}
	}

	mockRetrieveServerSessionByServerIDNotFound := func(hubSessionKey string, serverID int64) *ServerSession {
		return nil
	}

	tt := []struct {
		name                                string
		serverID                            int64
		serverArgs                          []interface{}
		mockRetrieveServerSessionByServerID func(hubSessionKey string, serverID int64) *ServerSession
		mockExecuteCall                     func(serverEndpoint string, call string, args []interface{}) (response interface{}, err error)
		expectedResponse                    interface{}
		expectedErr                         string
	}{
		{
			name:                                "Unicast call_successful",
			serverID:                            1,
			serverArgs:                          []interface{}{"arg1", "arg2"},
			mockRetrieveServerSessionByServerID: mockRetrieveServerSessionByServerIDFound,
			mockExecuteCall: func(serverEndpoint string, call string, args []interface{}) (response interface{}, err error) {
				return "success_response", nil
			},
			expectedResponse: "success_response",
		},
		{
			name:                                "Unicast call_error",
			serverID:                            1,
			serverArgs:                          []interface{}{"arg1", "arg2"},
			mockRetrieveServerSessionByServerID: mockRetrieveServerSessionByServerIDFound,
			mockExecuteCall: func(serverEndpoint string, call string, args []interface{}) (response interface{}, err error) {
				return nil, errors.New("call_error")
			},
			expectedErr: "call_error",
		},
		{
			name:                                "Unicast serverSession_not_found",
			serverID:                            1,
			serverArgs:                          []interface{}{"arg1", "arg2"},
			mockRetrieveServerSessionByServerID: mockRetrieveServerSessionByServerIDNotFound,
			expectedErr:                         "Authentication error: provided session key is invalid",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mockSession := new(mockSession)
			mockSession.mockRetrieveServerSessionByServerID = tc.mockRetrieveServerSessionByServerID

			mockUyuniCallExecutor := new(mockUyuniCallExecutor)
			mockUyuniCallExecutor.mockExecuteCall = tc.mockExecuteCall

			unicaster := NewUnicaster(mockUyuniCallExecutor, mockSession)

			response, err := unicaster.Unicast("hubSessionKey", "call", tc.serverID, tc.serverArgs)

			if err != nil && tc.expectedErr != err.Error() {
				t.Fatalf("Error during executing request: %v", err)
			}
			if err == nil && !reflect.DeepEqual(response, tc.expectedResponse) {
				t.Fatalf("Expected and actual values don't match, Expected is: %v", tc.expectedResponse)
			}
		})
	}
}
