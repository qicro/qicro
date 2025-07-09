package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/github"
)

type OAuthService struct {
	googleConfig *oauth2.Config
	githubConfig *oauth2.Config
}

func NewOAuthService(googleClientID, googleClientSecret, githubClientID, githubClientSecret, redirectURL string) *OAuthService {
	return &OAuthService{
		googleConfig: &oauth2.Config{
			ClientID:     googleClientID,
			ClientSecret: googleClientSecret,
			RedirectURL:  redirectURL + "/google",
			Scopes:       []string{"email", "profile"},
			Endpoint:     google.Endpoint,
		},
		githubConfig: &oauth2.Config{
			ClientID:     githubClientID,
			ClientSecret: githubClientSecret,
			RedirectURL:  redirectURL + "/github",
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		},
	}
}

func (o *OAuthService) GetAuthURL(provider, state string) (string, error) {
	switch provider {
	case "google":
		return o.googleConfig.AuthCodeURL(state), nil
	case "github":
		return o.githubConfig.AuthCodeURL(state), nil
	default:
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}
}

func (o *OAuthService) ExchangeToken(provider, code string) (*OAuthUserInfo, error) {
	switch provider {
	case "google":
		return o.exchangeGoogleToken(code)
	case "github":
		return o.exchangeGitHubToken(code)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

func (o *OAuthService) exchangeGoogleToken(code string) (*OAuthUserInfo, error) {
	ctx := context.Background()
	token, err := o.googleConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	client := o.googleConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &OAuthUserInfo{
		ID:    userInfo.ID,
		Email: userInfo.Email,
		Name:  userInfo.Name,
	}, nil
}

func (o *OAuthService) exchangeGitHubToken(code string) (*OAuthUserInfo, error) {
	ctx := context.Background()
	token, err := o.githubConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	client := o.githubConfig.Client(ctx, token)
	
	// 获取用户基本信息
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var userInfo struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	// 如果用户没有公开邮箱，获取邮箱列表
	email := userInfo.Email
	if email == "" {
		email, err = o.getGitHubUserEmail(client)
		if err != nil {
			return nil, fmt.Errorf("failed to get user email: %w", err)
		}
	}

	return &OAuthUserInfo{
		ID:    fmt.Sprintf("%d", userInfo.ID),
		Email: email,
		Name:  userInfo.Name,
	}, nil
}

func (o *OAuthService) getGitHubUserEmail(client *http.Client) (string, error) {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var emails []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}

	if err := json.Unmarshal(body, &emails); err != nil {
		return "", err
	}

	for _, email := range emails {
		if email.Primary {
			return email.Email, nil
		}
	}

	if len(emails) > 0 {
		return emails[0].Email, nil
	}

	return "", fmt.Errorf("no email found")
}