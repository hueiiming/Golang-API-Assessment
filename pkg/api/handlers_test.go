package api

import (
	"Golang-API-Assessment/pkg/repository/mocks"
	"Golang-API-Assessment/pkg/types"
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
		getServer  func(t *testing.T, repository *mocks.Repository, mockMethod string, request interface{}) *httptest.Server
		expStatus  int
		expBody    interface{}
		verifyBody func(t *testing.T, body []byte, expBody interface{})
	}{
		{
			name:       "When HandleRegister is called with correct request via POST method, return status code 204",
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
			getServer: func(t *testing.T, repository *mocks.Repository, mockMethod string, request interface{}) *httptest.Server {
				s := NewServer(":3000", repository)
				repository.On("GetTeacherID", mock.Anything).Once().Return(1, nil)
				repository.On("GetStudentID", mock.Anything).Once().Return(1, nil)
				repository.On("GetStudentID", mock.Anything).Once().Return(2, nil)
				repository.On(mockMethod, mock.Anything, mock.Anything).Once().Return(nil)
				return httptest.NewServer(MakeHTTPHandle(s.HandleRegister))
			},
			expStatus:  http.StatusNoContent,
			expBody:    nil,
			verifyBody: func(t *testing.T, body []byte, expBody interface{}) {},
		},
		{
			name:       "When HandleRegister is called with wrong request via POST method, return status code 400 with error message",
			restMethod: http.MethodPost,
			mockMethod: "Registration",
			path:       "/api/register",
			request: &types.RegisterRequest{
				Students: []string{
					"studentjon@gmail.com",
					"studenthon@gmail.com",
				},
			},
			getServer: func(t *testing.T, repository *mocks.Repository, mockMethod string, request interface{}) *httptest.Server {
				s := NewServer(":3000", repository)
				return httptest.NewServer(MakeHTTPHandle(s.HandleRegister))
			},
			expStatus: http.StatusBadRequest,
			expBody: ApiError{
				Message: "error: invalid teacher email",
			},
			verifyBody: func(t *testing.T, body []byte, expBody interface{}) {
				var respBody ApiError
				err := json.Unmarshal(body, &respBody)
				if err != nil {
					t.Errorf("Error parsing JSON response body: %s", err)
				}
				contains := expBody.(ApiError)
				assert.Contains(t, respBody.Message, contains.Message)
			},
		},
		{
			name:       "When HandleSuspend is called with correct request via POST method, return status code 204",
			restMethod: http.MethodPost,
			mockMethod: "Suspension",
			path:       "/api/suspend",
			request: &types.SuspendRequest{
				Student: "studentmary@gmail.com",
			},
			getServer: func(t *testing.T, repository *mocks.Repository, mockMethod string, request interface{}) *httptest.Server {
				s := NewServer(":3000", repository)
				repository.On("GetStudentID", mock.Anything).Once().Return(1, nil)
				repository.On(mockMethod, mock.Anything).Once().Return(nil)
				return httptest.NewServer(MakeHTTPHandle(s.HandleSuspend))
			},
			expStatus:  http.StatusNoContent,
			expBody:    nil,
			verifyBody: func(t *testing.T, body []byte, expBody interface{}) {},
		},
		{
			name:       "When HandleSuspend is called with wrong request via POST method, return status code 400 with error message",
			restMethod: http.MethodPost,
			mockMethod: "Suspension",
			path:       "/api/suspend",
			request:    &types.SuspendRequest{},
			getServer: func(t *testing.T, repository *mocks.Repository, mockMethod string, request interface{}) *httptest.Server {
				s := NewServer(":3000", repository)
				return httptest.NewServer(MakeHTTPHandle(s.HandleSuspend))
			},
			expStatus: http.StatusBadRequest,
			expBody: ApiError{
				Message: "error: invalid request",
			},
			verifyBody: func(t *testing.T, body []byte, expBody interface{}) {
				var respBody ApiError
				err := json.Unmarshal(body, &respBody)
				if err != nil {
					t.Errorf("Error parsing JSON response body: %s", err)
				}
				contains := expBody.(ApiError)
				assert.Contains(t, respBody.Message, contains.Message)
			},
		},
		{
			name:       "When HandleNotification is called with correct request via POST method, return status code 204 and recipients body",
			restMethod: http.MethodPost,
			mockMethod: "GetNotification",
			path:       "/api/retrievefornotifications",
			request: &types.NotificationRequest{
				Teacher: "teacherken@gmail.com",
				Message: "Hello students! @studentagnes@gmail.com @studentmiche@gmail.com",
			},
			getServer: func(t *testing.T, repository *mocks.Repository, mockMethod string, request interface{}) *httptest.Server {
				s := NewServer(":3000", repository)
				repository.On(mockMethod, request).Once().Return([]string{
					"studenthon@gmail.com", "studentjon@gmail.com"}, nil)
				return httptest.NewServer(MakeHTTPHandle(s.HandleRetrieveNotifications))
			},
			expStatus: http.StatusOK,
			expBody: types.NotificationResponse{
				Recipients: []string{
					"studenthon@gmail.com",
					"studentjon@gmail.com",
				},
			},
			verifyBody: func(t *testing.T, body []byte, expBody interface{}) {
				var respBody types.NotificationResponse
				err := json.Unmarshal(body, &respBody)
				if err != nil {
					t.Errorf("Error parsing JSON response body: %s", err)
				}
				expected := expBody.(types.NotificationResponse)
				assert.Equal(t, expected, respBody)
			},
		},
		{
			name:       "When HandleNotification is called with wrong request via POST method, return status code 400 and error message",
			restMethod: http.MethodPost,
			mockMethod: "GetNotification",
			path:       "/api/retrievefornotifications",
			request: &types.NotificationRequest{
				Teacher: "teacherken@gmail.com",
			},
			getServer: func(t *testing.T, repository *mocks.Repository, mockMethod string, request interface{}) *httptest.Server {
				s := NewServer(":3000", repository)
				return httptest.NewServer(MakeHTTPHandle(s.HandleRetrieveNotifications))
			},
			expStatus: http.StatusBadRequest,
			expBody: ApiError{
				Message: "error: invalid request",
			},
			verifyBody: func(t *testing.T, body []byte, expBody interface{}) {
				var respBody ApiError
				err := json.Unmarshal(body, &respBody)
				if err != nil {
					t.Errorf("Error parsing JSON response body: %s", err)
				}
				contains := expBody.(ApiError)
				assert.Contains(t, respBody.Message, contains.Message)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repository := new(mocks.Repository)

			server := tc.getServer(t, repository, tc.mockMethod, tc.request)
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

			assert.Equal(t, tc.expStatus, resp.StatusCode)
			tc.verifyBody(t, body, tc.expBody)
		})
	}
}

func TestHandlersGet(t *testing.T) {
	tests := []struct {
		name       string
		restMethod string
		mockMethod string
		path       string
		getServer  func(t *testing.T, repository *mocks.Repository, repoResp []string, mockMethod string) *httptest.Server
		repoResp   []string
		expStatus  int
		expBody    interface{}
		verifyBody func(t *testing.T, body []byte, expBody interface{})
	}{
		{
			name:       "When getCommonStudents is called with 1 teacher request via GET method, return status code 200 and response body",
			restMethod: http.MethodGet,
			mockMethod: "GetCommonStudents",
			path:       "/api/commonstudents?teacher=teacherken%40gmail.com",
			repoResp: []string{
				"studentjon@gmail.com",
				"studenthon@gmail.com",
			},
			getServer: func(t *testing.T, repository *mocks.Repository, repoResp []string, mockMethod string) *httptest.Server {
				s := NewServer(":3000", repository)
				repository.On(mockMethod, mock.Anything).Return(repoResp, nil)
				return httptest.NewServer(MakeHTTPHandle(s.HandleCommonStudents))
			},
			expStatus: http.StatusOK,
			expBody: types.CommonStudentsResponse{
				Students: []string{
					"studentjon@gmail.com",
					"studenthon@gmail.com",
					"student_only_under_teacherken@gmail.com",
				},
			},
			verifyBody: func(t *testing.T, body []byte, expBody interface{}) {
				var respBody types.CommonStudentsResponse
				err := json.Unmarshal(body, &respBody)
				if err != nil {
					t.Errorf("Error parsing JSON response body: %s", err)
				}
				expected := expBody.(types.CommonStudentsResponse)
				assert.Equal(t, expected, respBody)
			},
		},
		{
			name:       "When getCommonStudents is called with 1 invalid teacher request via GET method, return status code 400 and error message",
			restMethod: http.MethodGet,
			mockMethod: "GetCommonStudents",
			path:       "/api/commonstudents?feature=teacherken%40gmail.com",
			repoResp: []string{
				"studentjon@gmail.com",
				"studenthon@gmail.com",
			},
			getServer: func(t *testing.T, repository *mocks.Repository, repoResp []string, mockMethod string) *httptest.Server {
				s := NewServer(":3000", repository)
				return httptest.NewServer(MakeHTTPHandle(s.HandleCommonStudents))
			},
			expStatus: http.StatusBadRequest,
			expBody: ApiError{
				Message: "invalid query param",
			},
			verifyBody: func(t *testing.T, body []byte, expBody interface{}) {
				var respBody ApiError
				err := json.Unmarshal(body, &respBody)
				if err != nil {
					t.Errorf("Error parsing JSON response body: %s", err)
				}
				contains := expBody.(ApiError)
				assert.Contains(t, respBody.Message, contains.Message)
			},
		},
		{
			name:       "When getCommonStudents is called with 1 empty params teacher request via GET method, return status code 400 and error message",
			restMethod: http.MethodGet,
			mockMethod: "GetCommonStudents",
			path:       "/api/commonstudents?teacher=",
			repoResp: []string{
				"studentjon@gmail.com",
				"studenthon@gmail.com",
			},
			getServer: func(t *testing.T, repository *mocks.Repository, repoResp []string, mockMethod string) *httptest.Server {
				s := NewServer(":3000", repository)
				return httptest.NewServer(MakeHTTPHandle(s.HandleCommonStudents))
			},
			expStatus: http.StatusBadRequest,
			expBody: ApiError{
				Message: "empty query param",
			},
			verifyBody: func(t *testing.T, body []byte, expBody interface{}) {
				var respBody ApiError
				err := json.Unmarshal(body, &respBody)
				if err != nil {
					t.Errorf("Error parsing JSON response body: %s", err)
				}
				contains := expBody.(ApiError)
				assert.Contains(t, respBody.Message, contains.Message)
			},
		},
		{
			name:       "When getCommonStudents is called with 2 teacher request via GET method, return status code 200 and response body",
			restMethod: http.MethodGet,
			mockMethod: "GetCommonStudents",
			path:       "/api/commonstudents?teacher=teacherken%40gmail.com&teacher=teacherjoe%40gmail.com",
			repoResp: []string{
				"studentjon@gmail.com",
				"studenthon@gmail.com",
			},
			getServer: func(t *testing.T, repository *mocks.Repository, repoResp []string, mockMethod string) *httptest.Server {
				s := NewServer(":3000", repository)
				repository.On(mockMethod, mock.Anything).Return(repoResp, nil)
				return httptest.NewServer(MakeHTTPHandle(s.HandleCommonStudents))
			},
			expStatus: http.StatusOK,
			expBody: types.CommonStudentsResponse{
				Students: []string{
					"studentjon@gmail.com",
					"studenthon@gmail.com",
				},
			},
			verifyBody: func(t *testing.T, body []byte, expBody interface{}) {
				var respBody types.CommonStudentsResponse
				err := json.Unmarshal(body, &respBody)
				if err != nil {
					t.Errorf("Error parsing JSON response body: %s", err)
				}
				expected := expBody.(types.CommonStudentsResponse)
				assert.Equal(t, expected, respBody)
			},
		},
		{
			name:       "When getCommonStudents is called with 2 invalid teacher request via GET method, return status code 200 and response body",
			restMethod: http.MethodGet,
			mockMethod: "GetCommonStudents",
			path:       "/api/commonstudents?feature=teacherken%40gmail.com&teacher=teacherjoe%40gmail.com",
			repoResp: []string{
				"studentjon@gmail.com",
				"studenthon@gmail.com",
			},
			getServer: func(t *testing.T, repository *mocks.Repository, repoResp []string, mockMethod string) *httptest.Server {
				s := NewServer(":3000", repository)
				return httptest.NewServer(MakeHTTPHandle(s.HandleCommonStudents))
			},
			expStatus: http.StatusBadRequest,
			expBody: ApiError{
				Message: "invalid query param",
			},
			verifyBody: func(t *testing.T, body []byte, expBody interface{}) {
				var respBody ApiError
				err := json.Unmarshal(body, &respBody)
				if err != nil {
					t.Errorf("Error parsing JSON response body: %s", err)
				}
				contains := expBody.(ApiError)
				assert.Contains(t, respBody.Message, contains.Message)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repository := new(mocks.Repository)

			server := tc.getServer(t, repository, tc.repoResp, tc.mockMethod)
			defer repository.AssertExpectations(t)
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

			assert.Equal(t, tc.expStatus, resp.StatusCode)
			tc.verifyBody(t, body, tc.expBody)
		})
	}
}
