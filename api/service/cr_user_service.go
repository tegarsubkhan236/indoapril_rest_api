package service

import (
	"errors"
	"example/pkg"
	"example/pkg/model"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strconv"
)

func GetAllUser(page, limit int, user model.CrUserRequest) ([]model.CrUser, int64, error) {
	var db = pkg.DB
	var data []model.CrUser
	var count int64
	var offset = (page - 1) * limit

	db = db.Model(&data)
	db = db.Preload("Role")
	db = db.Preload("Role.Permissions")

	//if ids != nil {
	//	db = db.Where("id IN (?)", ids)
	//}
	if user.Username != "" {
		db = db.Where("username LIKE ?", "%"+user.Username+"%")
	}
	if user.Email != "" {
		db = db.Where("email LIKE ?", "%"+user.Email+"%")
	}

	if err := db.Count(&count).Offset(offset).Limit(limit).Order("id desc").Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, count, nil
}

func GetUserById(id uint) (*model.CrUser, error) {
	var db = pkg.DB
	var user model.CrUser

	db = db.Preload("Role")
	db = db.Preload("Role.Permissions")

	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item not found")
		}
		return nil, err
	}

	return &user, nil
}

func GetUserByEmail(e string) (*model.CrUser, error) {
	var db = pkg.DB
	var user model.CrUser

	if err := db.Where(&model.CrUser{Email: e}).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func GetUserByUsername(u string) (*model.CrUser, error) {
	var db = pkg.DB
	var user model.CrUser

	if err := db.Where(&model.CrUser{Username: u}).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func CreateUser(data model.CrUser) (*model.CrUser, error) {
	var db = pkg.DB

	if err := db.Create(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func UpdateUser(user model.CrUser, payload model.CrUser) (*model.CrUser, error) {
	var db = pkg.DB

	if err := db.Model(&user).Updates(payload).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func DestroyUser(ids []int) error {
	var db = pkg.DB
	var users []model.CrUser
	var tx = db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := db.Unscoped().Delete(&users, ids).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func ValidToken(t *jwt.Token, id string) bool {
	n, err := strconv.Atoi(id)
	if err != nil {
		return false
	}
	claims := t.Claims.(jwt.MapClaims)
	uid := int(claims["user_id"].(float64))
	if uid != n {
		return false
	}
	return true
}
