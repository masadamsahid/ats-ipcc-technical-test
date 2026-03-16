package service

import (
	"tsilodot/dto"
	"tsilodot/helpers"
	"tsilodot/model"
	"tsilodot/repository"
)

type AuthService struct {
	userRepository repository.IUserRepository
}

type IAuthService interface {
	Register(params *model.User) (*dto.RegisterResponseData, error)
	Login(params *dto.LoginRequest) (*dto.LoginResponseData, error)
}

func NewAuthService(userRepository repository.IUserRepository) IAuthService {
	return &AuthService{userRepository: userRepository}
}

func (a AuthService) Register(newUser *model.User) (*dto.RegisterResponseData, error) {
	hashPwd, err := helpers.HashPassword(newUser.Password)
	newUser.Password = hashPwd

	user, err := a.userRepository.CreateUser(nil, newUser)
	if err != nil {
		return nil, err
	}

	accessToken, err := helpers.CreateAuthToken(helpers.AuthTokenClaims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
	if err != nil {
		return nil, err
	}

	return &dto.RegisterResponseData{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt.String(),
		UpdatedAt:   user.UpdatedAt.String(),
		AccessToken: accessToken,
	}, nil
}

// Login implements IAuthService.
func (a *AuthService) Login(params *dto.LoginRequest) (*dto.LoginResponseData, error) {
	user, err := a.userRepository.FindUserByEmail(nil, params.Email)
	if err != nil {
		return nil, err
	}

	err = helpers.CompareHashPassword(user.Password, params.Password)
	if err != nil {
		return nil, err
	}

	accessToken, err := helpers.CreateAuthToken(helpers.AuthTokenClaims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponseData{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt.String(),
		UpdatedAt:   user.UpdatedAt.String(),
		AccessToken: accessToken,
	}, nil
}
