package get_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-api-service/internals/http-server/handlers/user/get"
	"user-api-service/internals/http-server/handlers/user/get/mocks"
	"user-api-service/internals/lib/logger/slogdiscard"
	"user-api-service/internals/models"
)

func TestUpdateHandler(t *testing.T) {
	testUserData := models.User{FirstName: "firstName", LastName: "LastName", Age: 50, Email: "Email@mail.com"}
	cases := []struct {
		name      string
		user      models.User
		id        string
		respError string
		mockError error
	}{
		{
			name: "Success",
			user: testUserData,
			id:   "8801a593-b781-4388-ae03-232872c8fd90",
		},
		{
			name:      "SaveURL Error",
			user:      testUserData,
			id:        "8801a593-b781-4388-ae03-232872c8fd90",
			respError: "failed to get user",
			mockError: errors.New("unexpected error"),
		},
		//TODO add more checks
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			userGetterMock := mocks.NewUserGetter(t)

			if tc.mockError != nil {
				userGetterMock.On("GetUser", tc.id).
					Return(nil, tc.mockError).
					Once()
			} else {
				userGetterMock.On("GetUser", tc.id).
					Return(tc.user, nil).
					Once()
			}
			handler := get.New(slogdiscard.NewDiscardLogger(), userGetterMock)

			url := fmt.Sprintf("/user/%s", tc.id)

			req, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(nil))
			require.NoError(t, err)

			r := chi.NewRouter()
			r.Get("/user/{id}", handler)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			body := rr.Body.String()
			var resp get.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
