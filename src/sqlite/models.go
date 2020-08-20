package sqlite

// ServerModel Structure of the server obj
type ServerModel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ServerType  string `json:"type"`
	IPAddress   string `json:"ip_address"`
	City        string `json:"city"`
	ControlPort string `json:"control_port"`
	UserID      string `json:"user_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Os          string `json:"os"`
	Enabled     string `json:"enabled"`
	KeyPort     int    `json:"key_port"`
}

// ExtensionModel Structure of the extension obj
type ExtensionModel struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	Icon      string `json:"icon"`
	Service   string `json:"service"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Order     int    `json:"order"`
	SslPorts  string `json:"sslPorts"`
	Issuer    string `json:"issuer"`
	Language  string `json:"language"`
	Support   string `json:"support"`
	Displays  string `json:"displays"`
	Status    string `json:"status"`
}

// SettingsModel Structure of the user settings obj
type SettingsModel struct {
	ID        string `json:"id"`
	ServerID  string `json:"server_id"`
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	Value     string `json:"value"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// TokenModel Structure of the tokens table
type TokenModel struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Token     string `json:"token"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

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

// AccessToken Structure of the personal access token
type AccessToken struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	UserID     string `json:"user_id"`
	LastUsedAt string `json:"last_used_at"`
	LastUsedIP string `json:"last_used_ip"`
	Token      string `json:"token"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// License Structure of the personal extension license
type License struct {
	ID          string `json:"id"`
	Data        string `json:"data"`
	ExtensionID string `json:"extension_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// Widget Structure of the widget
type Widget struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	UserID      string `json:"user_id"`
	Type        string `json:"type"`
	ExtensionID string `json:"extension_id"`
	ServerID    string `json:"serrver_id"`
	Function    string `json:"function"`
	Text        string `json:"text"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Order       int    `json:"order"`
}

// Permission Structure of the permissions
type Permission struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Type      string `json:"type"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Extra     string `json:"extra"`
	Blame     string `json:"blame"`
	MorphID   string `json:"morph_id"`
	MorphType string `json:"morph_type"`
}

// RoleUsers Structure of the role users
type RoleUsers struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	RoleID    string `json:"role_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Type      string `json:"type"`
}
