package jgoweb

import (
	"fmt"
	"errors"
	"database/sql"
	"github.com/asaskevich/govalidator"
)

// User is a retrieved and authentiacted user.
type User struct {
	Id string `json:"id" valid:"optional,uuid"`
	AccountId string `json:"account_id" valid:"required,uuid"`
	RoleId int `json:"role_id" valid:"required,range(1|100)"`
	GivenName string `json:"given_name" valid:"required,length(1|254)"`
	FamilyName string `json:"family_name" valid:"required,length(1|254)"`
	Email string `json:"email" valid:"required,email,length(1|254)"`
	CreatedAt string `json:"created_at" valid:"optional,rfc3339"`
	UpdatedAt string `json:"updated_at" valid:"optional,rfc3339"`
	DeletedAt sql.NullString `json:"deleted_at"`
	Ctx ContextInterface `json:"-" valid:"-"`
}

//
func NewUser(ctx ContextInterface) *User {
	return &User{Ctx: ctx}
}

// 
func FetchUserById(ctx ContextInterface, id string) (*User, error) {
	var user []User

	stmt := ctx.Select("*").
	From("users").
	Where("id = ?", id).
	Limit(1)

	_, err := stmt.Load(&user)

	if err != nil {
		return nil, err
	}

	if (len(user) == 0) {
		return nil, nil
	}

	user[0].Ctx = ctx
	ctx.SetUser(&user[0])

	return &user[0], nil
}

// 
func FetchUserByEmail(ctx ContextInterface, email string) (*User, error) {
	var user []User

	stmt := ctx.Select("*").
	From("users").
	Where("email = ?", email).
	Limit(1)

	_, err := stmt.Load(&user)

	if err != nil {
		return nil, err
	}

	if (len(user) == 0) {
		return nil, nil
	}

	user[0].Ctx = ctx
	ctx.SetUser(&user[0])

	return &user[0], nil
}

// set user from session
func (u *User) SetFromSession() error {
	var err error
	userEmail, err := u.Ctx.SessionGetString("user_email")

	if err != nil {
		return err
	}

	if userEmail == "" {
		return nil
	}

	user, err := FetchUserByShardEmail(u.Ctx, userEmail)

	if err != nil {
		return  err
	}

	if user == nil {
		return errors.New(fmt.Sprintf("User with email (%v) not found", userEmail))
	} else {
		u = user
		u.Ctx.SetUser(user)
	}

	return nil
}

// 
func (u *User) Authenticate(password string) bool {

	// @TEMP: @TODO: @WIP: Hard coded universal password for now...
	if password != "letmein" {
		return false
	}

	return true
}

//
func (u *User) isValid() (bool, error) {
	isValid, err := govalidator.ValidateStruct(u)

	if !isValid {
		return isValid, err
	}

	// have to do DeletedAt validation manually because it's nullable...
	if !u.DeletedAt.Valid || u.DeletedAt.String == "" {
		return true, nil
	}

	isValid = govalidator.IsRFC3339(u.DeletedAt.String)

	if !isValid {
		// @TODO: Errors from govalidator suck. Need to make better...
		err = errors.New("deleted_at: is not a valid timestamp")
		return false, err
	}

	return true, nil
}

//
func (u *User) Save() error {
	isValid, err := u.isValid()

	if !isValid {
		return err
	}

	if u.Id == "" {
		return u.Insert()
	} else {
		return u.Update()
	}
}

//
func (u *User) Insert() error {
	tx, err := u.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	_, err = tx.InsertInto("users").
		Columns("account_id", "role_id", "given_name", "family_name", "email", "deleted_at").
		Record(u).
		Exec()

	if err != nil {
		return err
	}

	err = u.Ctx.OptionalCommit(tx)

	return err
}

//
func (u *User) Update() error {
	tx, err := u.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	_, err = tx.Update("users").
		Set("account_id", u.AccountId).
		Set("role_id", u.RoleId).
		Set("given_name", u.GivenName).
		Set("family_name", u.FamilyName).
		Set("email", u.Email).
		Set("deleted_at", u.DeletedAt).
		Where("id = ?", u.Id).
		Exec()

	if err != nil {
		return err
	}

	err = u.Ctx.OptionalCommit(tx)

	return err
}
