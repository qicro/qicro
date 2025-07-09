package auth

import (
	"errors"
	"fmt"
	"time"
)

type Service struct {
	repo         *Repository
	jwtService   *JWTService
	oauthService *OAuthService
}

func NewService(repo *Repository, jwtService *JWTService, oauthService *OAuthService) *Service {
	return &Service{
		repo:         repo,
		jwtService:   jwtService,
		oauthService: oauthService,
	}
}

func (s *Service) Register(req RegisterRequest) (*LoginResponse, error) {
	// 检查用户是否已存在
	if _, err := s.repo.GetUserByEmail(req.Email); err == nil {
		return nil, errors.New("user already exists")
	}

	// 创建新用户
	user := NewUser(req.Email)
	
	// 加密密码
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.PasswordHash = &hashedPassword

	// 保存到数据库
	if err := s.repo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 生成JWT token
	token, err := s.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *Service) Login(req LoginRequest) (*LoginResponse, error) {
	// 获取用户
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 检查密码
	if user.PasswordHash == nil {
		return nil, errors.New("user registered with OAuth, please use OAuth login")
	}

	if err := CheckPassword(*user.PasswordHash, req.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 生成JWT token
	token, err := s.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *Service) GetAuthURL(provider, state string) (string, error) {
	return s.oauthService.GetAuthURL(provider, state)
}

func (s *Service) OAuthLogin(provider, code string) (*LoginResponse, error) {
	// 获取OAuth用户信息
	userInfo, err := s.oauthService.ExchangeToken(provider, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	// 检查用户是否已存在
	user, err := s.repo.GetUserByOAuth(provider, userInfo.ID)
	if err != nil {
		// 用户不存在，检查是否已有相同邮箱的用户
		existingUser, emailErr := s.repo.GetUserByEmail(userInfo.Email)
		if emailErr == nil {
			// 邮箱已存在，更新OAuth信息
			existingUser.OAuthProvider = &provider
			existingUser.OAuthID = &userInfo.ID
			existingUser.UpdatedAt = time.Now()
			if err := s.repo.UpdateUser(existingUser); err != nil {
				return nil, fmt.Errorf("failed to update user: %w", err)
			}
			user = existingUser
		} else {
			// 创建新用户
			user = NewOAuthUser(userInfo.Email, provider, userInfo.ID)
			if err := s.repo.CreateUser(user); err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}
		}
	}

	// 生成JWT token
	token, err := s.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *Service) ValidateToken(tokenString string) (*User, error) {
	claims, err := s.jwtService.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetUserByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *Service) RefreshToken(tokenString string) (string, error) {
	return s.jwtService.RefreshToken(tokenString)
}