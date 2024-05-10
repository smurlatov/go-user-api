package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-api-service/internals/http-server/handlers/user/save"
	"user-api-service/internals/http-server/handlers/user/save/mocks"
	"user-api-service/internals/lib/logger/slogdiscard"
	"user-api-service/internals/models"
)

func TestSaveHandler(t *testing.T) {
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
			respError: "field FirstName is a required field, field LastName is a required field, field Email is a required field, field Age is a required field",
		},
		{
			name:      "Invalid user data - FirstName",
			user:      models.User{FirstName: "firstName0000", LastName: "LastName", Age: 50, Email: "Email@mail.com"},
			respError: "field FirstName is not valid",
		},
		{
			name:      "SaveURL Error",
			user:      testUserData,
			respError: "failed to save user",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			userSaverMock := mocks.NewUserSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				userSaverMock.On("SaveUser", tc.user).
					Return("", tc.mockError).
					Once()
			} else if tc.id != "" {
				userSaverMock.On("SaveUser", tc.user).
					Return(tc.id, nil).
					Once()
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), userSaverMock)

			input := fmt.Sprintf(
				`{"first-name": "%s","last-name": "%s","e-mail": "%s","age": %d }`,
				tc.user.FirstName,
				tc.user.LastName,
				tc.user.Email,
				tc.user.Age,
			)

			req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
