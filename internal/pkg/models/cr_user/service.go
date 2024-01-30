package cr_user

import (
	"errors"
	"example/internal/pkg"
	"example/internal/pkg/entities"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"os"
	"regexp"
	"time"
)

type Service interface {
	AuthenticateUser(identity, password string) (string, error)
	FetchProfile(user interface{}) (*jwt.MapClaims, error)
	UpdateProfile(user interface{}, password string, payload *entities.CrUser) (*entities.CrUser, error)
	InsertUser(payload *entities.CrUser) (*entities.CrUser, error)
	FetchAllUser(page, limit int) (*[]entities.CrUser, int64, error)
	FetchDetailUser(ID uint) (*entities.CrUser, error)
	UpdateUser(ID uint, payload *entities.CrUser) (*entities.CrUser, error)
	UpdateUserPassword(ID uint, oldPassword, newPassword string) (*entities.CrUser, error)
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

func (s service) AuthenticateUser(identity, password string) (string, error) {
	var user *entities.CrUser

	isValidEmail := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(identity)

	if isValidEmail {
		user, _ = s.repository.ReadUserByEmail(identity)
	} else {
		user, _ = s.repository.ReadUserByUsername(identity)
	}

	if user == nil {
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
	claims["user_id"] = user.ID
	claims["username"] = user.Username
	claims["email"] = user.Email
	claims["sex"] = user.Sex
	claims["status"] = user.Status
	claims["phone"] = user.Phone
	claims["avatar"] = user.Avatar
	claims["role_id"] = user.Role.ID
	claims["role_name"] = user.Role.Name
	claims["team_id"] = user.Team.ID
	claims["team_name"] = user.Team.Name
	claims["permissions"] = permissions
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	t, err := token.SignedString([]byte(pkg.GetEnv("SECRET")))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (s service) FetchProfile(user interface{}) (*jwt.MapClaims, error) {
	token, ok := user.(*jwt.Token)
	if !ok {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return &claims, nil
}

func (s service) UpdateProfile(user interface{}, password string, payload *entities.CrUser) (*entities.CrUser, error) {
	token, ok := user.(*jwt.Token)
	if !ok {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	userID, ok := claims["user_id"].(uint)
	if !ok {
		return nil, errors.New("invalid user ID")
	}

	item, err := s.repository.ReadUserByID(userID)
	if err != nil {
		return nil, err
	}

	if item.Avatar != "" {
		if err = os.Remove(item.Avatar); err != nil {
			return nil, err
		}
	}

	if !s.checkPasswordHash(password, item.Password) {
		return nil, errors.New("invalid password")
	}

	updateUser, err := s.repository.UpdateUser(item, payload)
	if err != nil {
		return nil, err
	}

	response := updateUser.ToResponse()
	return &response, nil
}

func (s service) InsertUser(payload *entities.CrUser) (*entities.CrUser, error) {
	hashPass, err := s.hashPassword(payload.Password)
	if err != nil {
		return nil, err
	}
	payload.Password = hashPass

	createUser, err := s.repository.CreateUser(payload)
	if err != nil {
		return nil, err
	}

	response := createUser.ToResponse()
	return &response, nil
}

func (s service) FetchAllUser(page, limit int) (*[]entities.CrUser, int64, error) {
	users, count, err := s.repository.ReadUser(page, limit)
	if err != nil {
		return nil, 0, err
	}

	var results []entities.CrUser
	for _, item := range *users {
		results = append(results, item.ToResponse())
	}

	return &results, count, nil
}

func (s service) FetchDetailUser(ID uint) (*entities.CrUser, error) {
	user, err := s.repository.ReadUserByID(ID)
	if err != nil {
		return nil, err
	}
	response := user.ToResponse()
	return &response, nil
}

func (s service) UpdateUser(ID uint, payload *entities.CrUser) (*entities.CrUser, error) {
	user, err := s.repository.ReadUserByID(ID)
	if err != nil {
		return nil, err
	}

	if user.Avatar != "" {
		if err = os.Remove(user.Avatar); err != nil {
			return nil, err
		}
	}

	updateUser, err := s.repository.UpdateUser(user, payload)
	if err != nil {
		return nil, err
	}

	response := updateUser.ToResponse()
	return &response, nil
}

func (s service) UpdateUserPassword(ID uint, oldPassword, newPassword string) (*entities.CrUser, error) {
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
