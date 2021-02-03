package service

import (
	"context"
	"testing"

	"github.com/gogocoding/micro-go-course/dao"
	"github.com/gogocoding/micro-go-course/redis"
)

func TestUserServiceImpl_Login(t *testing.T) {

	err := dao.InitMysql("39.107.51.59", "3306", "root", "123456", "user")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = redis.InitRedis("39.107.51.59", "6379", "hq123456")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	UserService := &UserServiceImpl{
		userDAO: &dao.UserDAOImpl{},
	}

	user, err := UserService.Login(context.Background(), "aoho@mail.com", "aoho")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("user id is %d", user.ID)
}

func TestUserServiceImpl_Register(t *testing.T) {
	err := dao.InitMysql("39.107.51.59", "3306", "root", "123456", "user")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = redis.InitRedis("39.107.51.59", "6379", "hq123456")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	UserService := &UserServiceImpl{
		userDAO: &dao.UserDAOImpl{},
	}

	user, err := UserService.Register(context.Background(),
		&RegisterUserVO{
			Username: "aoho",
			Password: "aoho",
			Email:    "aoho@mail.com",
		})

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("user id is %d", user.ID)
}
