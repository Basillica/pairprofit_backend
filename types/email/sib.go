package email

import (
	"fmt"
	"time"

	"github.com/antihax/optional"
	"github.com/gin-gonic/gin"
	sib_api_v3_sdk "github.com/sendinblue/APIv3-go-library/lib"
)

type SibObject struct {
	Email       string
	Firstname   *string
	Lastname    *string
	RegisterAt  *time.Time
	LastLogin   *time.Time
	PhoneNumber *int16
	StatusNo    int
}

type EmailSender struct {
	Recipients  *[]sib_api_v3_sdk.SendSmtpEmailTo
	Attachments *[]sib_api_v3_sdk.SendSmtpEmailAttachment
	Params      interface{}
	Urls        *[]sib_api_v3_sdk.SendSmtpEmailAttachment
	Subject     string
	HtmlContent string
}

type SendSmtpEmailAttachment struct {
	// Absolute url of the attachment (no local file).
	Url string `json:"url,omitempty"`
	// Base64 encoded chunk data of the attachment generated on the fly
	Content string `json:"content,omitempty"`
	// Required if content is passed. Name of the attachment
	Name string `json:"name,omitempty"`
}

func (s *SibObject) GetAllContacts(c *gin.Context) {
	sibClient := c.MustGet("sibClient").(*sib_api_v3_sdk.APIClient)
	var params = &sib_api_v3_sdk.ContactsApiGetContactsOpts{
		Limit:         optional.NewInt64(50),
		Offset:        optional.NewInt64(0),
		ModifiedSince: optional.NewString("2020-09-20T19:20:30+01:00"),
	}

	obj, resp, err := sibClient.ContactsApi.GetContacts(c, params)
	if err != nil {
		fmt.Println("Error when calling ContactsApi->GetContacts: ", err.Error())
		return
	}
	fmt.Println("GetAllContact Object:", obj, " GetAllContact Response: ", resp)
	return
}

func (s *SibObject) CreateContact(c *gin.Context) {
	sibClient := c.MustGet("sibClient").(*sib_api_v3_sdk.APIClient)
	var params = sib_api_v3_sdk.CreateContact{
		Email:   s.Email,
		ListIds: []int64{1},
	}

	obj, resp, err := sibClient.ContactsApi.CreateContact(c, params)
	if err != nil {
		fmt.Println("Error in ContactsApi->CreateContact", err.Error())
		return
	}
	fmt.Println("CreateContact Object:", obj, " CreateContact Response: ", resp)
	return
}

func (s *SibObject) GetContact(c *gin.Context) *sib_api_v3_sdk.GetAccount {
	sibClient := c.MustGet("sibClient").(*sib_api_v3_sdk.APIClient)

	result, resp, err := sibClient.AccountApi.GetAccount(c)
	if err != nil {
		fmt.Println("Error when calling AccountApi->get_account: ", err.Error())
		return &sib_api_v3_sdk.GetAccount{}
	}
	fmt.Println("GetAccount Object:", result, " GetAccount Response: ", resp)
	return &result
}

func (s *SibObject) GetContactDetails(c *gin.Context) *sib_api_v3_sdk.GetExtendedContactDetails {
	sibClient := c.MustGet("sibClient").(*sib_api_v3_sdk.APIClient)
	exRes, resp, err := sibClient.ContactsApi.GetContactInfo(c, s.Email)
	if err != nil {
		fmt.Println("Error in ContactsApi->GetContactInfo ", err.Error())
		return &sib_api_v3_sdk.GetExtendedContactDetails{}
	}
	fmt.Println("GetContactInfo extended data: ", exRes, "response:", resp)
	return &exRes
}

func (s *SibObject) DeleteContact(c *gin.Context) {
	sibClient := c.MustGet("sibClient").(*sib_api_v3_sdk.APIClient)
	resp, err := sibClient.ContactsApi.DeleteContact(c, s.Email)
	if err != nil {
		fmt.Println("Error in ContactsApi->DeleteContact", err.Error())
		return
	}
	fmt.Println("Delete contact response:", resp)
	return
}

func (s *SibObject) UpdateContact(c *gin.Context, lastname, firstname *string, status *int) {
	sibClient := c.MustGet("sibClient").(*sib_api_v3_sdk.APIClient)
	var a interface{}

	attr := map[string]interface{}{}
	if lastname != nil {
		attr["LASTNAME"] = &lastname
	}
	if firstname != nil {
		attr["FIRSTNAME"] = &firstname
	}
	if firstname != nil {
		attr["ACCOUNTSTATUS"] = 2
	}

	a = attr
	var params = &sib_api_v3_sdk.UpdateContact{
		ListIds:    []int64{10},
		Attributes: &a,
	}
	resp, err := sibClient.ContactsApi.UpdateContact(c, *params, s.Email)
	if err != nil {
		fmt.Println("Error in ContactsApi->UpdateContact", err.Error())
		return
	}
	fmt.Println("UpdateContact Response: ", resp)
	return
}

func (s *SibObject) ListAllAttributes(c *gin.Context) *sib_api_v3_sdk.GetAttributes {
	sibClient := c.MustGet("sibClient").(*sib_api_v3_sdk.APIClient)
	obj, resp, err := sibClient.ContactsApi.GetAttributes(c)
	if err != nil {
		fmt.Println("Error in ContactsApi->GetAttributes ", err.Error())
		return &sib_api_v3_sdk.GetAttributes{}
	}
	fmt.Println("GetAttributes response:", resp, "GetAttributes object:", obj)
	return &obj
}

func (s *SibObject) CreateContactAttributes(c *gin.Context) {
	sibClient := c.MustGet("sibClient").(*sib_api_v3_sdk.APIClient)
	body := sib_api_v3_sdk.CreateAttribute{
		Enumeration: []sib_api_v3_sdk.CreateAttributeEnumeration{
			{Value: 1, Label: "Invited"},
			{Value: 2, Label: "Registered"},
			{Value: 3, Label: "Active"},
			{Value: 4, Label: "Suspended"},
			{Value: 5, Label: "Inactive"},
		},
		Type_: "category",
	}
	attributeCategory := "category"
	attributeName := "accountstatus"

	resp, err := sibClient.ContactsApi.CreateAttribute(c, body, attributeCategory, attributeName)
	if err != nil {
		fmt.Println("Error in ContactsApi.CreateAttribute", err.Error())
		return
	}
	fmt.Println("CreateAttribute response:", resp)
	return
}

func (s *SibObject) DeleteAttribute(c *gin.Context, category, attribute string) {
	sibClient := c.MustGet("sibClient").(*sib_api_v3_sdk.APIClient)
	attributeCategory := category //"category"
	attributeName := attribute    //"levelOfExpertise"

	resp, err := sibClient.ContactsApi.DeleteAttribute(c, attributeCategory, attributeName)
	if err != nil {
		fmt.Println("Error in ContactsApi.DeleteAttribute", err.Error())
		return
	}
	fmt.Println("DeleteAttribute response:", resp)
	return
}

func (s *SibObject) CreateList(c *gin.Context, listName string, folderId int64) {
	sibClient := c.MustGet("sibClient").(*sib_api_v3_sdk.APIClient)
	body := sib_api_v3_sdk.CreateList{
		Name:     listName, //"MyListName",
		FolderId: folderId, //1,
	}
	obj, resp, err := sibClient.ContactsApi.CreateList(c, body)
	if err != nil {
		fmt.Println("Error in ContactsApi->CreateList ", err.Error())
		return
	}
	fmt.Println("CreateList response:", resp, "CreateList object:", obj)
	return
}

func (s *SibObject) GetListDetail(c *gin.Context) {
	sibClient := c.MustGet("sibClient").(*sib_api_v3_sdk.APIClient)
	listId := int64(2)

	obj, resp, err := sibClient.ContactsApi.GetList(c, listId)
	if err != nil {
		fmt.Println("Error in ContactsApi->GetList ", err.Error())
		return
	}
	fmt.Println("GetList response:", resp, "GetList object:", obj)
	return
}

func (s *SibObject) UpdateList(c *gin.Context) {
	sibClient := c.MustGet("sibClient").(*sib_api_v3_sdk.APIClient)
	listId := int64(2)
	body := sib_api_v3_sdk.UpdateList{
		Name:     "myListName",
		FolderId: int64(1),
	}
	resp, err := sibClient.ContactsApi.UpdateList(c, body, listId)
	if err != nil {
		fmt.Println("Error in ContactsApi->UpdateList ", err.Error())
		return
	}
	fmt.Println("UpdateList response:", resp)
	return
}

func (s *SibObject) DeleteList(c *gin.Context) {
	sibClient := c.MustGet("sibClient").(*sib_api_v3_sdk.APIClient)
	listId := int64(2)

	resp, err := sibClient.ContactsApi.DeleteList(c, listId)
	if err != nil {
		fmt.Println("Error in ContactsApi->DeleteList ", err.Error())
		return
	}
	fmt.Println("DeleteList response:", resp)
	return
}

func (s *SibObject) GetContactsFromList(c *gin.Context) {
	sibClient := c.MustGet("sibClient").(*sib_api_v3_sdk.APIClient)

	listId := int64(2)
	body := &sib_api_v3_sdk.ContactsApiGetContactsFromListOpts{
		Limit:         optional.NewInt64(10),
		Offset:        optional.NewInt64(0),
		ModifiedSince: optional.NewString("2020-01-20T19:20:30+01:00"),
		Sort:          optional.NewString("asc"),
	}

	obj, resp, err := sibClient.ContactsApi.GetContactsFromList(c, listId, body)
	if err != nil {
		fmt.Println("Error in ContactsApi->GetContactsFromList ", err.Error())
		return
	}
	fmt.Println("GetContactsFromList response:", resp, "GetContactsFromList object:", obj)
	return
}

func (s *SibObject) AddExistingContactToList(c *gin.Context) {
	sibClient := c.MustGet("sibClient").(*sib_api_v3_sdk.APIClient)
	listId := int64(2)
	body := sib_api_v3_sdk.AddContactToList{
		Emails: []string{
			"example@example.com",
		},
	}
	obj, resp, err := sibClient.ContactsApi.AddContactToList(c, body, listId)
	if err != nil {
		fmt.Println("Error in ContactsApi->GetContactsFromList ", err.Error())
		return
	}
	fmt.Println("GetContactsFromList response:", resp, "GetContactsFromList object:", obj)
	return
}

func (s *SibObject) DeleteContactFromList(c *gin.Context) {
	sibClient := c.MustGet("sibClient").(*sib_api_v3_sdk.APIClient)
	listId := int64(2)
	body := sib_api_v3_sdk.RemoveContactFromList{
		Emails: []string{
			"example@example.com",
		},
	}
	obj, resp, err := sibClient.ContactsApi.RemoveContactFromList(c, body, listId)
	if err != nil {
		fmt.Println("Error in ContactsApi->RemoveContactFromList ", err.Error())
		return
	}
	fmt.Println("RemoveContactFromList response:", resp, "RemoveContactFromList object:", obj)
	return
}

func (s *SibObject) CreateFolder(c *gin.Context) {
	sibClient := c.MustGet("sibClient").(*sib_api_v3_sdk.APIClient)
	body := sib_api_v3_sdk.CreateUpdateFolder{
		Name: "myFolderName",
	}
	obj, resp, err := sibClient.ContactsApi.CreateFolder(c, body)
	if err != nil {
		fmt.Println("Error in ContactsApi->CreateFolder ", err.Error())
		return
	}
	fmt.Println("CreateFolder response:", resp, "CreateFolder object:", obj)
	return
}

func (s *SibObject) GetFolderDetails(c *gin.Context) *sib_api_v3_sdk.GetFolder {
	sibClient := c.MustGet("sibClient").(*sib_api_v3_sdk.APIClient)
	folderId := int64(1)

	obj, resp, err := sibClient.ContactsApi.GetFolder(c, folderId)
	if err != nil {
		fmt.Println("Error in ContactsApi->GetFolder ", err.Error())
		return &sib_api_v3_sdk.GetFolder{}
	}
	fmt.Println("GetFolder response:", resp, "GetFolder object:", obj)
	return &obj
}

func (s *SibObject) UpdateFolder(c *gin.Context) {
	sibClient := c.MustGet("sibClient").(*sib_api_v3_sdk.APIClient)
	folderId := int64(1)
	body := sib_api_v3_sdk.CreateUpdateFolder{
		Name: "mySecondFolder",
	}

	resp, err := sibClient.ContactsApi.UpdateFolder(c, body, folderId)
	if err != nil {
		fmt.Println("Error in ContactsApi->UpdateFolder ", err.Error())
		return
	}
	fmt.Println("UpdateFolder response:", resp)
	return
}

func (s *SibObject) DeleteFolderAndItsLists(c *gin.Context) {
	sibClient := c.MustGet("sibClient").(*sib_api_v3_sdk.APIClient)
	folderId := int64(1)

	resp, err := sibClient.ContactsApi.DeleteFolder(c, folderId)
	if err != nil {
		fmt.Println("Error in ContactsApi->DeleteFolder ", err.Error())
		return
	}
	fmt.Println("DeleteFolder response:", resp)
	return
}

func (s *SibObject) GetAllLists(c *gin.Context) {
	sibClient := c.MustGet("sibClient").(*sib_api_v3_sdk.APIClient)
	params := sib_api_v3_sdk.ContactsApiGetListsOpts{
		Limit:  optional.NewInt64(10),
		Offset: optional.NewInt64(0),
	}

	obj, resp, err := sibClient.ContactsApi.GetLists(c, &params)
	if err != nil {
		fmt.Println("Error in ContactsApi->GetLists ", err.Error())
		return
	}
	fmt.Println("GetLists response:", resp, "GetLists object:", obj)
	return
}

func (e *EmailSender) SendEmail(c *gin.Context) {
	sibClient := c.MustGet("sibClient").(*sib_api_v3_sdk.APIClient)
	e.PopulateRecipients(c)
	if e.Params != nil {
		e.PopulateParams(c)
	}
	if e.Attachments != nil {
		e.PopulateAttachments(c)
	}
	res := e.Send(c, sibClient)
	fmt.Println(res)
}

func (e *EmailSender) Send(c *gin.Context, sibClient *sib_api_v3_sdk.APIClient) *sib_api_v3_sdk.CreateSmtpEmail {
	body := sib_api_v3_sdk.SendSmtpEmail{
		HtmlContent: e.HtmlContent,
		Subject:     e.Subject,
		Sender: &sib_api_v3_sdk.SendSmtpEmailSender{
			Name:  "PairProfit",
			Email: "noreply@pairprofit.com",
		},
		To: *e.Recipients,
		ReplyTo: &sib_api_v3_sdk.SendSmtpEmailReplyTo{
			Name:  "PairProfit Service",
			Email: "service@pairprofit.com",
		},
	}
	if e.Params != nil {
		body.Params = &e.Params
	}
	if e.Attachments != nil {
		body.Attachment = *e.Attachments
	}

	obj, resp, err := sibClient.TransactionalEmailsApi.SendTransacEmail(c, body)
	if err != nil {
		fmt.Println("Error in TransactionalEmailsApi->SendTransacEmail ", err.Error())
		return &sib_api_v3_sdk.CreateSmtpEmail{}
	}
	fmt.Println("SendTransacEmail, response:", resp, "SendTransacEmail object", obj)
	return &obj
}

func (e *EmailSender) PopulateParams(c *gin.Context) {
	params := "{\"FirstName\":\"Joe\", \"Transcription\":\"Doe\"}"
	e.Params = params
}

func (e *EmailSender) PopulateAttachments(c *gin.Context) {
	var urls []sib_api_v3_sdk.SendSmtpEmailAttachment
	for _, url := range *e.Urls {
		urls = append(urls, sib_api_v3_sdk.SendSmtpEmailAttachment{
			Url:     url.Url,
			Content: url.Content,
			Name:    url.Name,
		})
	}
	e.Attachments = &urls
}

func (e *EmailSender) PopulateRecipients(c *gin.Context) {
	var reciepient []sib_api_v3_sdk.SendSmtpEmailTo
	for _, contact := range *e.Recipients {
		reciepient = append(reciepient, sib_api_v3_sdk.SendSmtpEmailTo{
			Email: contact.Email,
			Name:  contact.Name,
		})
	}
	e.Recipients = &reciepient
}
