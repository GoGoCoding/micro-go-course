package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gogocoding/micro-go-course/dao"
	"github.com/gogocoding/micro-go-course/redis"
	"github.com/jinzhu/gorm"
)

type UserInfoDTO struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RegisterUserVO struct {
	Username string
	Password string
	Email    string
}

var (
	ErrUserExisted = errors.New("user is existed")
	ErrPassword    = errors.New("email and password are not match")
	ErrRegistering = errors.New("email is registering")
)

type UserService interface {
	Login(ctx context.Context, email, password string) (*UserInfoDTO, error)

	Register(ctx context.Context, vo *RegisterUserVO) (*UserInfoDTO, error)
}

type UserServiceImpl struct {
	userDAO dao.UserDAO
}

func MakeUserServiceImpl(userDAO dao.UserDAO) UserService {
	return &UserServiceImpl{
		userDAO: userDAO,
	}
}

func (userservice *UserServiceImpl) Login(ctx context.Context, email, password string) (*UserInfoDTO, error) {
	user, err := userservice.userDAO.SelectByEmail(email)
	if err == nil {
		if user.Password == password {
			return &UserInfoDTO{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
			}, nil
		} else {
			return nil, ErrPassword
		}
	} else {
		log.Printf("err: %s", err)
	}
	return nil, err
}

func (userService UserServiceImpl) Register(ctx context.Context, vo *RegisterUserVO) (*UserInfoDTO, error) {
	lock := redis.GetRedisLock(vo.Email, time.Duration(5)*time.Second)
	err := lock.Lock()
	if err != nil {
		log.Printf("err: %s", err)
		return nil, ErrRegistering
	}

	defer lock.Unlock()

	existUser, err := userService.userDAO.SelectByEmail(vo.Email)

	if (err == nil && existUser == nil) || err == gorm.ErrRecordNotFound {
		newUser := &dao.UserEntity{
			Username: vo.Username,
			Password: vo.Password,
			Email:    vo.Email,
		}

		err = userService.userDAO.Save(newUser)
		if err == nil {
			return &UserInfoDTO{
				ID:       newUser.ID,
				Username: newUser.Username,
				Email:    newUser.Email,
			}, nil
		}
	}
	if err == nil {
		err = ErrUserExisted
	}
	return nil, err
}
