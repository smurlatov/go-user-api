package update_test

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
	"user-api-service/internals/http-server/handlers/user/update"
	"user-api-service/internals/http-server/handlers/user/update/mocks"
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
			name:      "Empty user",
			user:      models.User{},
			id:        "8801a593-b781-4388-ae03-232872c8fd90",
			respError: "field FirstName is a required field, field LastName is a required field, field Email is a required field, field Age is a required field",
		},
		{
			name:      "Invalid id",
			user:      testUserData,
			id:        "-12314",
			respError: "invalid uuid",
		},
		{
			name:      "Invalid user data - FirstName",
			user:      models.User{FirstName: "firstName0000", LastName: "LastName", Age: 50, Email: "Email@mail.com"},
			id:        "8801a593-b781-4388-ae03-232872c8fd90",
			respError: "field FirstName is not valid",
		},
		{
			name:      "Invalid user data - Age",
			user:      models.User{FirstName: "firstName", LastName: "LastName", Age: 1000, Email: "Email@mail.com"},
			id:        "8801a593-b781-4388-ae03-232872c8fd90",
			respError: "field Age is not valid",
		}, //TODO add cases for all user fields
		{
			name:      "SaveURL Error",
			user:      testUserData,
			id:        "8801a593-b781-4388-ae03-232872c8fd90",
			respError: "failed to update user",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			userUpdaterMock := mocks.NewUserUpdater(t)

			if tc.respError == "" || tc.mockError != nil {
				userUpdaterMock.On("UpdateUser", tc.user, tc.id).
					Return(tc.mockError).
					Once()
			}

			handler := update.New(slogdiscard.NewDiscardLogger(), userUpdaterMock)

			input := fmt.Sprintf(
				`{"first-name": "%s","last-name": "%s","e-mail": "%s","age": %d }`,
				tc.user.FirstName,
				tc.user.LastName,
				tc.user.Email,
				tc.user.Age,
			)

			url := fmt.Sprintf("/user/%s", tc.id)

			req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			r := chi.NewRouter()
			r.Patch("/user/{id}", handler)

			ts := httptest.NewServer(r)
			defer ts.Close()

			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			body := rr.Body.String()
			var resp update.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
