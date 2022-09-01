package requests

import (
	"mime/multipart"
	"net/http"

	"github.com/mholt/binding"
)

type InviteUserConfirmRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Token     string `json:"token"`
	Workspace string `json:"workspace"`
}

type InviteUserRequest struct {
	Emails []string `json:"emails"`
}

type UpdateProfileDataRequest struct {
	Email             *string               `form:"email"`
	Firstname         *string               `form:"firstname"`
	Lastname          *string               `form:"lastname"`
	Password          *string               `form:"password"`
	Role              *string               `form:"role"`
	Company           *string               `form:"company"`
	PreferredLanguage *string               `form:"preferred_language"`
	About             *string               `form:"about"`
	Url               *string               `form:"url"`
	Picture           *multipart.FileHeader `form:"picture"`
}

func (f *UpdateProfileDataRequest) FieldMap(c *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&f.Email:             "email",
		&f.Firstname:         "firstname",
		&f.Lastname:          "lastname",
		&f.Password:          "password",
		&f.Role:              "role",
		&f.Company:           "company",
		&f.PreferredLanguage: "preferred_language",
		&f.About:             "about",
		&f.Url:               "url",
		&f.Picture:           "picture",
	}
}
