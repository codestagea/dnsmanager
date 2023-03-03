package store

import (
	"errors"
	"github.com/codestagea/bindmgr/internal/model"
	"gorm.io/gorm"
	"strings"
)

type User struct {
	Model
	Username string `json:"username" gorm:"username"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
	Status   int    `json:"status" gorm:"status"`
}

type UserQuery struct {
	Search string
	Status int
}

func (User) TableName() string {
	return "sys_user"
}

type UserStore interface {
	GetByUsername(username string) (*User, error)
	ListUser(q UserQuery, page *model.PageQuery) ([]User, int64, error)
	AddUser(u *User) error
	UpdateUser(u *User) error
	ChangeStatus(username string, status int) error
	ChangePwd(username, password string) error
}

type userStore struct {
	db *gorm.DB
}

var _ UserStore = &userStore{}

func (s *userStore) GetByUsername(username string) (*User, error) {
	var user User
	err := s.db.Table(user.TableName()).Where("username = ?", username).Take(&user).Error
	return &user, err
}

func (s *userStore) ListUser(q UserQuery, page *model.PageQuery) ([]User, int64, error) {
	fields := make([]string, 0)
	values := make([]interface{}, 0)
	if q.Search != "" {
		fields = append(fields, `username = like ('%?%')`)
		values = append(values, q.Search)
	}

	if q.Status >= 0 {
		fields = append(fields, `status = ?`)
		values = append(values, q.Status)
	}

	sql := s.db.Table(User{}.TableName())

	if len(fields) > 0 {
		sql = sql.Where(strings.Join(fields, " and "), values...)
	}

	var total int64
	if err := sql.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []User
	if err := sql.Limit(page.PageSize).Offset(page.Offset).Order("id desc").Find(&users).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, err
		}
	}

	return users, total, nil
}

func (s *userStore) AddUser(u *User) error {
	err := s.db.Table(u.TableName()).Save(u).Error
	return err
}

func (s *userStore) UpdateUser(u *User) error {
	err := s.db.Table(u.TableName()).Where("id = ?", u.ID).Updates(u).Error
	return err
}

func (s *userStore) ChangeStatus(username string, status int) error {

	return s.db.Table(User{}.TableName()).Where("username = ?", username).Updates(map[string]interface{}{
		"status": status,
	}).Error
}

func (s *userStore) ChangePwd(username, password string) error {
	user := User{
		Password: password,
	}
	return s.db.Table(user.TableName()).Where("username = ?", username).Updates(user).Error
}
