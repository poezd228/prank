package domain

type (
	Role string
	// Account - структура содержащая иммутабельную информацию о клиенте
	Account struct {
		Id    int64
		Login string
		Name  string
		Role  string
	}
	TokenInfo struct {
		Token     string `json:"token"`      // токен
		ExpiresAt int64  `json:"expires_at"` // время конца жизни (timestamp в миллисекундах)
	}
	TokenResponse struct {
		AccessToken  TokenInfo
		RefreshToken TokenInfo
	}

	AuthRequestData struct {
		Login    string `json:"login"`
		Password string `json:"password"`
		Version  string `json:"version"`
		Token    string `json:"token"`
	}
	AuthRequest struct {
		Id   uint64          `json:"id"`
		Data AuthRequestData `json:"auth"`
	}
	AuthAdminRequest struct {
		Id   uint64
		Data AuthAdminRequestData
	}
	AuthAdminRequestData struct {
		Id   int64
		Role string
	}
)
type AuthPurpose int32 // enum для типа jwt-токена
const (
	PurposeAccess  = AuthPurpose(iota) // access токен
	PurposeRefresh                     // refresh токен
)
