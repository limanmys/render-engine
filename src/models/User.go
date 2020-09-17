package models

// UserModel Structure of the users table
type UserModel struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Status        int    `json:"status"`
	LastLoginAt   string `json:"last_login_at"`
	RememberToken string `json:"remember_token"`
	LastLoginIP   string `json:"last_login_ip"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	ForceChange   int    `json:"forceChange"`
	ObjectGUID    string `json:"objectguid"`
	AuthType      string `json:"auth_type"`
}
