package api

import (
	"Golang-API-Assessment/repository/mocks"
	"Golang-API-Assessment/types"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers(t *testing.T) {
	tests := []struct {
		name      string
		request   *types.RegisterRequest
		expStatus int
	}{
		{
			name: "When handleRegister is called with correct request, return status code 204",
			request: &types.RegisterRequest{
				Teacher: "teacherken@gmail.com",
				Students: []string{
					"studentjon@gmail.com",
					"studenthon@gmail.com",
				},
			},
			expStatus: 204,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repository := new(mocks.Repository)
			repository.On("Registration", tc.request).Once().Return(nil)
			defer repository.AssertExpectations(t)

			s := NewServer(":3000", repository)
			server := httptest.NewServer(makeHTTPHandle(s.handleRegister))
			requestByte, err := json.Marshal(tc.request)
			if err != nil {
				t.Errorf("Error marshalling request: %s", err)
			}
			resp, err := http.Post(server.URL+"/api/register", "application/json", bytes.NewBuffer(requestByte))
			if err != nil {
				t.Errorf("Error posting request: %s", err)
			}

			assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		})
	}
}
