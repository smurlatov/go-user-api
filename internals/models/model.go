package models

type User struct {
	FirstName string `json:"first-name" validate:"required,alpha"`
	LastName  string `json:"last-name" validate:"required,alpha"`
	Email     string `json:"e-mail" validate:"required,email"`
	Age       uint   `json:"age" validate:"gte=0,lte=130"`
}
