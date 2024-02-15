package api

import (
	"Golang-API-Assessment/repository/mocks"
	"Golang-API-Assessment/types"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlersPost(t *testing.T) {
	tests := []struct {
		name       string
		restMethod string
		mockMethod string
		path       string
		request    interface{}
		expStatus  int
		expBody    interface{}
	}{
		{
			name:       "When handleRegister is called with correct request via POST method, return status code 204",
			restMethod: http.MethodPost,
			mockMethod: "Registration",
			path:       "/api/register",
			request: &types.RegisterRequest{
				Teacher: "teacherken@gmail.com",
				Students: []string{
					"studentjon@gmail.com",
					"studenthon@gmail.com",
				},
			},
			expStatus: http.StatusNoContent,
			expBody:   nil,
		},
		{
			name:       "When handleSuspend is called with correct request via POST method, return status code 204",
			restMethod: http.MethodPost,
			mockMethod: "Suspension",
			path:       "/api/suspend",
			request: &types.SuspendRequest{
				Student: "studentmary@gmail.com",
			},
			expStatus: http.StatusNoContent,
			expBody:   nil,
		},
		{
			name:       "When handleNotification is called with correct request via POST method, return status code 204 and recipients body",
			restMethod: http.MethodPost,
			mockMethod: "GetNotification",
			path:       "/api/retrievefornotifications",
			request: &types.NotificationRequest{
				Teacher: "teacherken@gmail.com",
				Message: "Hello students! @studentagnes@gmail.com @studentmiche@gmail.com",
			},
			expStatus: http.StatusOK,
			expBody: types.Notification{
				Recipients: []string{
					"studenthon@gmail.com",
					"studentjon@gmail.com",
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repository := new(mocks.Repository)

			s := NewServer(":3000", repository)
			var server *httptest.Server

			switch tc.mockMethod {
			case "Registration":
				repository.On(tc.mockMethod, tc.request).Once().Return(nil)
				server = httptest.NewServer(makeHTTPHandle(s.handleRegister))
			case "Suspension":
				repository.On(tc.mockMethod, tc.request).Once().Return(nil)
				server = httptest.NewServer(makeHTTPHandle(s.handleSuspend))
			case "GetNotification":
				repository.On(tc.mockMethod, tc.request).Once().Return([]string{
					"studenthon@gmail.com", "studentjon@gmail.com"}, nil)
				server = httptest.NewServer(makeHTTPHandle(s.handleRetrieveNotifications))
			}
			defer repository.AssertExpectations(t)
			defer server.Close()

			requestByte, err := json.Marshal(tc.request)
			if err != nil {
				t.Errorf("Error marshalling request: %s", err)
			}
			req, err := http.NewRequest(tc.restMethod, server.URL+tc.path, bytes.NewBuffer(requestByte))
			if err != nil {
				t.Errorf("Error creating POST request: %s", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("Error getting POST response: %s", err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Error reading response body: %s", err)
			}

			switch tc.mockMethod {
			case "GetNotification":
				var respBody types.Notification
				err = json.Unmarshal(body, &respBody)
				if err != nil {
					t.Errorf("Error parsing JSON response body: %s", err)
				}
				assert.Equal(t, tc.expBody, respBody)
			}
			assert.Equal(t, tc.expStatus, resp.StatusCode)
		})
	}
}

func TestHandlersGet(t *testing.T) {
	tests := []struct {
		name       string
		restMethod string
		mockMethod string
		path       string
		repo       []string
		expStatus  int
		expBody    interface{}
	}{
		{
			name:       "When getCommonStudents is called with 1 teacher request via GET method, return status cose 200 and response body",
			restMethod: http.MethodGet,
			mockMethod: "GetCommonStudents",
			path:       "/api/commonstudents?teacher=teacherken%40gmail.com",
			repo: []string{
				"studentjon@gmail.com",
				"studenthon@gmail.com",
			},
			expStatus: 200,
			expBody: types.CommonStudents{
				Students: []string{
					"studentjon@gmail.com",
					"studenthon@gmail.com",
					"student_only_under_teacherken@gmail.com",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repository := new(mocks.Repository)

			repository.On(tc.mockMethod, mock.Anything).Return(tc.repo, nil)
			defer repository.AssertExpectations(t)

			s := NewServer(":3000", repository)
			server := httptest.NewServer(makeHTTPHandle(s.handleCommonStudents))
			defer server.Close()

			req, err := http.NewRequest(tc.restMethod, server.URL+tc.path, nil)
			if err != nil {
				t.Errorf("Error creating GET request: %s", err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("Error getting GET response: %s", err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Error reading response body: %s", err)
			}

			var respBody types.CommonStudents
			err = json.Unmarshal(body, &respBody)
			if err != nil {
				t.Errorf("Error parsing JSON response body: %s", err)
			}

			assert.Equal(t, tc.expStatus, resp.StatusCode)
			assert.Equal(t, tc.expBody, respBody)
		})
	}
}
