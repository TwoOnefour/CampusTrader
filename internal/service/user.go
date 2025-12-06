package service

import (
	"CampusTrader/internal/model"
	"CampusTrader/internal/util/jwtUtils"
	"CampusTrader/pkg/str"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) Register(ctx context.Context, req *model.User) error {
	db := s.db.WithContext(ctx)
	var count int64
	db.Model(&model.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		return errors.New("用户名已存在")
	}

	db.Model(&model.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		return errors.New("邮箱已被注册")
	}

	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	req.Password = string(hashedPwd)

	if err := db.Create(req).Error; err != nil {
		return errors.New("注册失败，请稍后重试")
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, account string, password string) (string, error) {
	var user model.User
	var err error
	db := s.db.WithContext(ctx)
	if strings.Contains(account, "@") {
		err = db.Where("email = ?", account).First(&user).Error
	} else if len(account) == 11 && str.IsAllDigits(account) {
		err = db.Where("phone = ?", account).First(&user).Error
	} else {
		err = db.Where("username = ?", account).First(&user).Error
	}

	if err != nil {
		return "", errors.New("账号不存在")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("密码不正确")
	}
	return jwtUtils.GenerateToken(user.ID, user.Username)
}

func (s *UserService) GetUserInfo(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	db := s.db.WithContext(ctx)
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
