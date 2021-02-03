package dao

import "testing"

func TestUserDAOImpl_Save(t *testing.T) {
	userDAO := &UserDAOImpl{}

	err := InitMysql("39.107.51.59", "3306", "root", "123456", "user")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	user := &UserEntity{
		Username: "aoho",
		Password: "aoho",
		Email:    "aoho@mail.com",
	}

	err = userDAO.Save(user)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("new User ID is %d", user.ID)
}

func TestUserDAPImpl_SelectByEmail(t *testing.T) {
	UserDAO := &UserDAOImpl{}

	err := InitMysql("39.107.51.59", "3306", "root", "123456", "user")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	user, err := UserDAO.SelectByEmail("aoho@mail.com")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("result username is %s", user.Username)
}
