package service

type AuthService struct {
	validUsername string
	validPassword string
	validToken    string
}

type VerifyResult struct {
	Valid   bool   `json:"valid"`
	Subject string `json:"subject,omitempty"`
	Error   string `json:"error,omitempty"`
}

func New() *AuthService {
	return &AuthService{
		validUsername: "student",
		validPassword: "student",
		validToken:    "demo-token",
	}
}

func (s *AuthService) Login(username, password string) (string, bool) {
	if username == s.validUsername && password == s.validPassword {
		return s.validToken, true
	}
	return "", false
}

func (s *AuthService) Verify(token string) VerifyResult {
	if token == s.validToken {
		return VerifyResult{
			Valid:   true,
			Subject: "student",
		}
	}

	return VerifyResult{
		Valid: false,
		Error: "unauthorized",
	}
}
