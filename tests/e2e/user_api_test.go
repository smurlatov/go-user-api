package e2e

import (
	"fmt"
	"github.com/bxcodec/faker/v4"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
	"user-api-service/internals/http-server/handlers/user/save"
	"user-api-service/internals/models"
)

const (
	host = "localhost:8090"
)

func TestPositiveCase(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	testUser := models.User{
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		Age:       10,
		Email:     faker.Email(),
	}
	//Create user
	response := e.POST("/users").
		WithJSON(save.Request{
			FirstName: testUser.FirstName,
			LastName:  testUser.LastName,
			Age:       testUser.Age,
			Email:     testUser.Email,
		}).
		Expect().
		Status(200).
		JSON().Object().
		ContainsKey("id")

	//Get user by id
	id := response.Value("id").String().Raw()
	urlWithId := fmt.Sprintf("/user/%s", id)
	response = e.GET(urlWithId).
		Expect().
		Status(200).
		JSON().Object().
		ContainsKey("user")

	userObj := response.Value("user").Object()
	var actualUser models.User
	userObj.Decode(&actualUser)
	require.Equal(t, testUser, actualUser)

	//Update user firstName
	testUser.FirstName = faker.FirstName() //modify testUserData
	e.PATCH(urlWithId).
		WithJSON(save.Request{
			FirstName: testUser.FirstName,
			LastName:  testUser.LastName,
			Age:       testUser.Age,
			Email:     testUser.Email,
		}).
		Expect().
		Status(200).
		JSON().Object()
	//Get user one more time and check new data
	response = e.GET(urlWithId).
		Expect().
		Status(200).
		JSON().Object().
		ContainsKey("user")

	userObj = response.Value("user").Object()
	var updatedUser models.User
	userObj.Decode(&updatedUser)
	require.Equal(t, testUser, updatedUser)
	require.NotEqual(t, actualUser.FirstName, updatedUser.FirstName) //check that firstname changed
}
