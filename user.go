package jgoweb

import(
	"fmt"
	"errors"
	"database/sql"
)

type UserInterface interface {
	SetFromSession() error
}

// User is a retrieved and authentiacted user.
type User struct {
	Id string `json:"id" valid:"optional,uuid"`
	Email string `json:"email" valid:"required,email,length(1|254)"`
	CreatedAt string `json:"created_at" valid:"optional,rfc3339"`
	UpdatedAt string `json:"updated_at" valid:"optional,rfc3339"`
	DeletedAt sql.NullString `json:"deleted_at"`
	ctx *WebContext
}

// New User
func NewUser(ctx *WebContext) *User {
	return &User{ctx: ctx}
}

// set user from session
func (u *User) SetFromSession() error {
	var err error
	userId, err := u.ctx.Session.GetString("user_id")

	if err != nil {
		return err
	}

	if userId == "" {
		return nil
	}

	err = u.LoadById(userId)

	if err != nil {
		return  err
	}

	u.ctx.User = u

	return nil
}

// 
func (u *User) LoadById(id string) error {
	var user []User

	stmt := u.ctx.Select("*").
	From("users").
	Where("id = ?", id).
	Limit(1)

	_, err := stmt.Load(&user)

	if err != nil {
		return err
	}

	if (len(user) == 0) {
		return errors.New(fmt.Sprintf("User with Id (%v) not found", id))
	}

	user[0].ctx = u.ctx
	u = &user[0]

	return nil
}
