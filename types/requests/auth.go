package requests

type LoginRequest struct {
	Username  string `json:"username" binding:"required,email,gte=6,lte=100"`
	Password  string `json:"password" binding:"required,gte=6,lte=100"`
	GrantType string `json:"grant_type" binding:"required,oneof=password refresh_token otp_request otp_token"`
}

type LoginViaLinkRequest struct {
	Username  string `json:"username" binding:"required,email,gte=6,lte=100"`
	GrantType string `json:"grant_type" binding:"required,oneof=password refresh_token otp_request otp_token"`
	Code      string `json:"code" binding:"required,gte=6,lte=100"`
	Token     string `json:"token" binding:"required,gte=6,lte=100"`
}

type RefreshRequest struct {
	GrantType string `json:"grant_type" binding:"required,oneof=password refresh_token otp_request otp_token"`
}

type RegisterRequest struct {
	Username      *string `json:"username" binding:"required,email,gte=6,lte=100"`
	Password      string  `json:"password" binding:"required,gte=6,lte=100"`
	Firstname     string  `json:"firstname" binding:"required,alpha,gte=1,lte=50"`
	Lastname      string  `json:"lastname" binding:"required,alpha,gte=1,lte=50"`
	InvitedUserID *int    `json:"invited_user_id"`
	ImageUri      *string `json:"image_uri"`
	Token         *string `json:"token"`
	AccountType   string  `json:"account_type" binding:"required,alpha,gte=1,lte=50,oneof=regularAccount serviceProvider"`
}

type ResetPasswordRequest struct {
	Username string `json:"username" binding:"required,email,gte=6,lte=100"`
	Hash     string `json:"hash" binding:"required,gte=6,lte=100"`
	Token    string `json:"token" binding:"required,gte=6,lte=100"`
}

type ForgotPasswordRequest struct {
	Username string `json:"username" binding:"required,email,gte=6,lte=100"`
}

type UpdatePasswordRequest struct {
	Username string `json:"username" binding:"required,email,gte=6,lte=100"`
	Password string `json:"password" binding:"required,gte=6,lte=100"`
}
