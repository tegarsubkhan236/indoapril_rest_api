package cr_user

import (
	"errors"
	"example/internal/pkg"
	"example/internal/pkg/entities"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type Service interface {
	AuthenticateUser(identity, password string) (string, error)
	FetchProfile() (*entities.CrUserResp, error)
	UpdateProfile(user *entities.CrUser) (*entities.CrUserResp, error)
	UpdateProfilePassword(oldPassword, newPassword string) (*entities.CrUserResp, error)

	InsertUser(user *entities.CrUser) (*entities.CrUserResp, error)
	FetchAllUser(page, limit int) (*[]entities.CrUserResp, int64, error)
	FetchDetailUser(ID uint) (*entities.CrUserResp, error)
	UpdateUser(ID uint, user *entities.CrUser) (*entities.CrUserResp, error)
	UpdateUserPassword(ID uint, oldPassword, newPassword string) (*entities.CrUserResp, error)
	DeleteUser(ID []uint) error

	checkPasswordHash(password, hash string) bool
	hashPassword(password string) (string, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s service) FetchProfile() (*entities.CrUserResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) UpdateProfile(user *entities.CrUser) (*entities.CrUserResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) UpdateProfilePassword(oldPassword, newPassword string) (*entities.CrUserResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) AuthenticateUser(identity, password string) (string, error) {
	var user *entities.CrUser

	userWithEmail, err := s.repository.ReadUserByEmail(identity)
	if err == nil {
		user = userWithEmail
	}

	userWithUsername, err := s.repository.ReadUserByUsername(identity)
	if err == nil {
		user = userWithUsername
	}

	if userWithEmail == nil && userWithUsername == nil {
		return "", errors.New("user not found")
	}

	if !s.checkPasswordHash(password, user.Password) {
		return "", errors.New("invalid password")
	}

	token := jwt.New(jwt.SigningMethodHS256)

	var permissions []string
	for _, permission := range user.Role.Permissions {
		permissions = append(permissions, permission.Name)
	}

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["user_id"] = user.ID
	claims["role_id"] = user.RoleID
	claims["permissions"] = permissions
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	t, err := token.SignedString([]byte(pkg.GetEnv("SECRET")))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (s service) InsertUser(user *entities.CrUser) (*entities.CrUserResp, error) {
	hashPass, err := s.hashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashPass

	createUser, err := s.repository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	response := createUser.ToResponse()
	return &response, nil
}

func (s service) FetchAllUser(page, limit int) (*[]entities.CrUserResp, int64, error) {
	users, count, err := s.repository.ReadUser(page, limit)
	if err != nil {
		return nil, 0, err
	}

	var results []entities.CrUserResp
	for _, item := range *users {
		results = append(results, item.ToResponse())
	}

	return &results, count, nil
}

func (s service) FetchDetailUser(ID uint) (*entities.CrUserResp, error) {
	user, err := s.repository.ReadUserByID(ID)
	if err != nil {
		return nil, err
	}
	response := user.ToResponse()
	return &response, nil
}

func (s service) UpdateUser(ID uint, user *entities.CrUser) (*entities.CrUserResp, error) {
	oldUser, err := s.repository.ReadUserByID(ID)
	if err != nil {
		return nil, err
	}

	if oldUser.Avatar != "" {
		if err = os.Remove(oldUser.Avatar); err != nil {
			return nil, err
		}
	}

	updateUser, err := s.repository.UpdateUser(oldUser, user)
	if err != nil {
		return nil, err
	}

	response := updateUser.ToResponse()
	return &response, nil
}

func (s service) UpdateUserPassword(ID uint, oldPassword, newPassword string) (*entities.CrUserResp, error) {
	user, err := s.repository.ReadUserByID(ID)
	if err != nil {
		return nil, err
	}

	if !s.checkPasswordHash(oldPassword, user.Password) {
		return nil, errors.New("old password is incorrect")
	}

	hashedNewPass, err := s.hashPassword(newPassword)
	if err != nil {
		return nil, err
	}

	var payload entities.CrUser
	payload.Password = hashedNewPass
	updateUser, err := s.repository.UpdateUser(user, &payload)
	if err != nil {
		return nil, err
	}

	response := updateUser.ToResponse()
	return &response, nil
}

func (s service) DeleteUser(ID []uint) error {
	return s.repository.DeleteUser(ID)
}

func (s service) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s service) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}
