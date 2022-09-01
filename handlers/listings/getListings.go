package listings

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	// "pairprofit.com/x/helpers"
	"pairprofit.com/x/types/listing"
	listing_type "pairprofit.com/x/types/listing"

	"github.com/golang/protobuf/proto"
	uuid "github.com/google/uuid"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func GetListings(c *gin.Context) {
	// access, err := c.Cookie("access_token")
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"success": false})
	// 	return
	// }

	// _, ok := helpers.GetUserOutput(c, access)
	// if !ok {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
	// 	return
	// }

	// var req listing_type.BaseListing
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	helpers.ValidatePayload(err, c)
	// 	return
	// }

	// cognito := c.MustGet("cognito").(*cognitoidentityprovider.Client)
	// appenv := c.MustGet("appenv").(*appenv.AppConfig)

	// category := listing_type.Listing_ListingCategory_name[0]

	//*req.Description,
	l := getListing()
	ls := &listing_type.Listings{
		Listing: getListings(l),
	}

	data, err := proto.Marshal(l)
	if err != nil {
		fmt.Println("marshaling error: ", err)
	}

	newL := &listing.Listing{}
	err = proto.Unmarshal(data, newL)
	if err != nil {
		fmt.Println("unmarshaling error: ", err)
	}

	fmt.Println(ls)
	fmt.Println(newL)
	fmt.Println(data)

	c.JSON(http.StatusAccepted, gin.H{"status": "sent"})
}

func getListing() *listing_type.Listing {
	l := &listing_type.Listing{
		Description:     "some boring description",
		Id:              uuid.New().String(),
		ListingCategory: listing_type.Listing_AutoRepairs,
		Owner: &listing_type.Listing_OwnerDetail{
			Email:       "test@gmail.com",
			ImageUri:    "someboringimageuri",
			ProfileLink: "someboringprofilelink",
		},
		Location: &listing_type.Listing_Location{
			Longitude: "longitude",
			Latitude:  "latitude",
		},
		CreationTime: timestamppb.Now(),
	}
	return l
}

func getListings(listing *listing_type.Listing) []*listing_type.Listing {
	var listings []*listing_type.Listing
	return append(listings, listing)
}
