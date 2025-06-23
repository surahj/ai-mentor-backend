package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleoauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

// GoogleUserInfo represents the user information from Google
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

// GoogleOAuthService handles Google OAuth2 operations
type GoogleOAuthService struct {
	config *oauth2.Config
}

// NewGoogleOAuthService creates a new Google OAuth service
func NewGoogleOAuthService() (*GoogleOAuthService, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET must be set")
	}

	log.Printf("INFO: Using Google OAuth Redirect URL: %s", redirectURL)

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleOAuthService{
		config: config,
	}, nil
}

// GetAuthURL returns the Google OAuth authorization URL
func (s *GoogleOAuthService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state)
}

// ExchangeCodeForToken exchanges authorization code for access token
func (s *GoogleOAuthService) ExchangeCodeForToken(code string) (*oauth2.Token, error) {
	ctx := context.Background()
	token, err := s.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %v", err)
	}
	return token, nil
}

// GetUserInfo retrieves user information from Google using the access token
func (s *GoogleOAuthService) GetUserInfo(token *oauth2.Token) (*GoogleUserInfo, error) {
	ctx := context.Background()

	// Create OAuth2 service
	oauth2Service, err := googleoauth2.NewService(ctx, option.WithTokenSource(s.config.TokenSource(ctx, token)))
	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth2 service: %v", err)
	}

	// Get user info
	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	// Convert to our GoogleUserInfo struct
	googleUser := &GoogleUserInfo{
		ID:            userInfo.Id,
		Email:         userInfo.Email,
		VerifiedEmail: userInfo.VerifiedEmail != nil && *userInfo.VerifiedEmail,
		Name:          userInfo.Name,
		GivenName:     userInfo.GivenName,
		FamilyName:    userInfo.FamilyName,
		Picture:       userInfo.Picture,
	}

	return googleUser, nil
}

// GetUserInfoFromTokenString retrieves user info using a token string
func (s *GoogleOAuthService) GetUserInfoFromTokenString(tokenString string) (*GoogleUserInfo, error) {
	// Make a request to Google's userinfo endpoint
	url := "https://www.googleapis.com/oauth2/v2/userinfo"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+tokenString)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %v", err)
	}

	return &userInfo, nil
}
