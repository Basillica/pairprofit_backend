package listing

import "time"

type Category int

type BaseListing struct {
	Description  *string   `json:"description" binding:"required,gte=6,lte=100"`
	Id           *string   `json:"id" binding:"required"`
	Category     *int      `json:"category" binding:"required,gte=0,lte=50"`
	Owner        *string   `json:"ower" binding:"required,email,alpha,gte=1,lte=50"`
	Longitude    *string   `json:"longitude" binding:"required,alpha,gte=1,lte=50"`
	Latitude     *string   `json:"latitude" binding:"required,alpha,gte=1,lte=50"`
	CreationTime time.Time `json:"creation_time" binding:"required,alpha,gte=1,lte=50"`
}
