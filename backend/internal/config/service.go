package config

import (
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// API Keys
func (s *Service) CreateAPIKey(req CreateAPIKeyRequest) (*APIKey, error) {
	// Validate request
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Value == "" {
		return nil, fmt.Errorf("value is required")
	}
	if req.Type == "" {
		return nil, fmt.Errorf("type is required")
	}
	if req.Provider == "" {
		return nil, fmt.Errorf("provider is required")
	}

	return s.repo.CreateAPIKey(req)
}

func (s *Service) GetAPIKeys() ([]APIKey, error) {
	return s.repo.GetAPIKeys()
}

func (s *Service) GetAPIKeyByID(id string) (*APIKey, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}
	return s.repo.GetAPIKeyByID(id)
}

func (s *Service) UpdateAPIKey(id string, req UpdateAPIKeyRequest) (*APIKey, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}
	return s.repo.UpdateAPIKey(id, req)
}

func (s *Service) DeleteAPIKey(id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}
	return s.repo.DeleteAPIKey(id)
}

// App Types
func (s *Service) CreateAppType(req CreateAppTypeRequest) (*AppType, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	return s.repo.CreateAppType(req)
}

func (s *Service) GetAppTypes() ([]AppType, error) {
	return s.repo.GetAppTypes()
}

// Chat Models
func (s *Service) CreateChatModel(req CreateChatModelRequest) (*ChatModel, error) {
	// Validate request
	if req.Type == "" {
		return nil, fmt.Errorf("type is required")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Value == "" {
		return nil, fmt.Errorf("value is required")
	}
	if req.Provider == "" {
		return nil, fmt.Errorf("provider is required")
	}

	return s.repo.CreateChatModel(req)
}

func (s *Service) GetChatModels() ([]ChatModel, error) {
	return s.repo.GetChatModels()
}

func (s *Service) GetChatModelsByType(modelType string) ([]ChatModel, error) {
	if modelType == "" {
		return nil, fmt.Errorf("model type is required")
	}
	return s.repo.GetChatModelsByType(modelType)
}

func (s *Service) GetChatModelByID(id string) (*ChatModel, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}
	return s.repo.GetChatModelByID(id)
}

func (s *Service) UpdateChatModel(id string, req UpdateChatModelRequest) (*ChatModel, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}
	return s.repo.UpdateChatModel(id, req)
}

func (s *Service) DeleteChatModel(id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}
	return s.repo.DeleteChatModel(id)
}

func (s *Service) GetAvailableAPIKeyForProvider(provider string) (*APIKey, error) {
	apiKeys, err := s.repo.GetAPIKeys()
	if err != nil {
		return nil, err
	}

	for _, key := range apiKeys {
		if key.Provider == provider && key.Enabled {
			return &key, nil
		}
	}

	return nil, fmt.Errorf("no available API key found for provider: %s", provider)
}