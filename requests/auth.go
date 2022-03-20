package requests

import (
	"mime/multipart"
	"net/http"

	"github.com/mholt/binding"
	"github.com/shopspring/decimal"
)

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

type UpdateTopicChildRequest struct {
	S3ID          string          `json:"s3_id"`
	Page          int             `json:"page"`
	ParentID      int             `json:"entry_id"`
	Duration      float64         `json:"duration"`
	XCoordinate   decimal.Decimal `json:"x_position"`
	YCoordinate   decimal.Decimal `json:"y_position"`
	Transcription string          `json:"transcription"`
	ChildID       int             `json:"child_id"`
}
