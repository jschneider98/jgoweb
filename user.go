package jgoweb

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gocraft/web"
	"github.com/jschneider98/jgoweb/util"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User
type User struct {
	Id                   sql.NullString   `json:"Id" validate:"omitempty,uuid"`
	AccountId            sql.NullString   `json:"AccountId" validate:"required,uuid"`
	RoleId               sql.NullString   `json:"RoleId" validate:"required,int"`
	FirstName            sql.NullString   `json:"FirstName" validate:"required"`
	LastName             sql.NullString   `json:"LastName" validate:"required"`
	Email                sql.NullString   `json:"Email" validate:"required"`
	Password             sql.NullString   `json:"Password" validate:"required,min=1,max=255"`
	CreatedAt            sql.NullString   `json:"CreatedAt" validate:"omitempty,rfc3339"`
	DeletedAt            sql.NullString   `json:"DeletedAt" validate:"omitempty,rfc3339"`
	UpdatedAt            sql.NullString   `json:"UpdatedAt" validate:"omitempty,rfc3339"`
	VerifiedAt           sql.NullString   `json:"VerifiedAt" validate:"omitempty,rfc3339"`
	Ctx                  ContextInterface `json:"-" validate:"-"`
	rawPassword          string           `json:"-", validate:"-"`
	verifyRawPassword    string           `json:"-", validate:"-"`
	currentPassword      string           `json:"-", validate:"-"`
	CurrentPasswordError string           `validate:"errorMsg"`
	UserUniqueError      string           `validate:"errorMsg"`
	RawPasswordError     string           `validate:"errorMsg"`
}

// Empty new model
func NewUser(ctx ContextInterface) (*User, error) {
	u := &User{Ctx: ctx}
	u.SetDefaults()

	return u, nil
}

// Set defaults
func (u *User) SetDefaults() {
	u.SetCreatedAt(time.Now().Format(time.RFC3339))
	u.SetUpdatedAt(time.Now().Format(time.RFC3339))
}

// New model with data
func NewUserWithData(ctx ContextInterface, req *web.Request) (*User, error) {
	u, err := NewUser(ctx)

	if err != nil {
		return nil, err
	}

	err = u.Hydrate(req)

	if err != nil {
		return nil, err
	}

	return u, nil
}

// Factory Method
func FetchUserById(ctx ContextInterface, id string) (*User, error) {
	var u []User

	stmt := ctx.Select("*").
		From("public.users").
		Where("id = ?", id).
		Limit(1)

	_, err := stmt.Load(&u)

	if err != nil {
		return nil, err
	}

	if len(u) == 0 {
		return nil, nil
	}

	u[0].Ctx = ctx

	return &u[0], nil
}

//
func (u *User) ProcessSubmit(req *web.Request) (string, bool, error) {
	err := u.Hydrate(req)

	if err != nil {
		return "", false, err
	}

	err = u.IsValid()

	if err != nil {
		return util.GetNiceErrorMessage(err, "</br>"), false, nil
	}

	err = u.Save()

	if err != nil {
		return "", false, err
	}

	return "User saved.", true, nil
}

// Hydrate the model with data
func (u *User) Hydrate(req *web.Request) error {
	err := req.ParseForm()

	if err != nil {
		return err
	}

	u.SetId(req.PostFormValue("Id"))
	u.SetAccountId(req.PostFormValue("AccountId"))
	u.SetRoleId(req.PostFormValue("RoleId"))
	u.SetFirstName(req.PostFormValue("FirstName"))
	u.SetLastName(req.PostFormValue("LastName"))
	u.SetEmail(req.PostFormValue("Email"))
	u.SetCreatedAt(req.PostFormValue("CreatedAt"))
	u.SetDeletedAt(req.PostFormValue("DeletedAt"))
	u.SetUpdatedAt(req.PostFormValue("UpdatedAt"))
	u.SetVerifiedAt(req.PostFormValue("VerifiedAt"))

	u.rawPassword = req.PostFormValue("rawPassword")
	u.verifyRawPassword = req.PostFormValue("verifyRawPassword")
	u.currentPassword = req.PostFormValue("currentPassword")

	if u.rawPassword != "" || u.verifyRawPassword != "" {
		u.SetPassword(u.rawPassword, u.verifyRawPassword)
	}

	return nil
}

// Validate the model
func (u *User) IsValid() error {
	u.UserUniqueError = ""

	valid, _, err := u.UserShardMapIsValid()

	// err != nil cases are kind of bad. Most likely a DB error occured, but IsValid only returns validation errors...
	if !valid || err != nil {
		u.UserUniqueError = "User already exists."

		fmt.Printf("ERROR: %v", err)
	}

	return u.Ctx.GetValidator().Struct(u)
}

// Insert/Update based on pkey value
func (u *User) Save() error {
	err := u.IsValid()

	if err != nil {
		return err
	}

	if !u.Id.Valid {
		return u.Insert()
	} else {
		return u.Update()
	}
}

// Insert a new record
func (u *User) Insert() error {

	query := `
INSERT INTO
public.users (account_id,
	role_id,
	first_name,
	last_name,
	email,
	deleted_at,
	password,
	verified_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
RETURNING id
`

	stmt, err := u.Ctx.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	err = stmt.QueryRow(u.AccountId,
		u.RoleId,
		u.FirstName,
		u.LastName,
		u.Email,
		u.DeletedAt,
		u.Password,
		u.VerifiedAt).Scan(&u.Id)

	if err != nil {
		return err
	}

	return nil
}

// Update a record
func (u *User) Update() error {

	if !u.Id.Valid {
		return nil
	}

	u.SetUpdatedAt(time.Now().Format(time.RFC3339))

	_, err := u.Ctx.Update("public.users").
		Set("id", u.Id).
		Set("account_id", u.AccountId).
		Set("role_id", u.RoleId).
		Set("first_name", u.FirstName).
		Set("last_name", u.LastName).
		Set("email", u.Email).
		Set("password", u.Password).
		Set("deleted_at", u.DeletedAt).
		Set("updated_at", u.UpdatedAt).
		Set("verified_at", u.VerifiedAt).
		Where("id = ?", u.Id).
		Exec()

	if err != nil {
		return err
	}

	return nil
}

// Soft delete a record
func (u *User) Delete() error {

	if !u.Id.Valid {
		return nil
	}

	u.SetDeletedAt((time.Now()).Format(time.RFC3339))

	_, err := u.Ctx.Update("public.users").
		Set("deleted_at", u.DeletedAt).
		Where("id = ?", u.Id).
		Exec()

	if err != nil {
		return err
	}

	return nil
}

// Soft undelete a record
func (u *User) Undelete() error {

	if !u.Id.Valid {
		return nil
	}

	u.SetDeletedAt("")

	_, err := u.Ctx.Update("public.users").
		Set("deleted_at", u.DeletedAt).
		Where("id = ?", u.Id).
		Exec()

	if err != nil {
		return err
	}

	return nil
}

//
func (u *User) GetId() string {

	if u.Id.Valid {
		return u.Id.String
	}

	return ""
}

//
func (u *User) SetId(val string) {

	if val == "" {
		u.Id.Valid = false
		u.Id.String = ""

		return
	}

	u.Id.Valid = true
	u.Id.String = val
}

//
func (u *User) GetAccountId() string {

	if u.AccountId.Valid {
		return u.AccountId.String
	}

	return ""
}

//
func (u *User) SetAccountId(val string) {

	if val == "" {
		u.AccountId.Valid = false
		u.AccountId.String = ""

		return
	}

	u.AccountId.Valid = true
	u.AccountId.String = val
}

//
func (u *User) GetRoleId() string {

	if u.RoleId.Valid {
		return u.RoleId.String
	}

	return ""
}

//
func (u *User) SetRoleId(val string) {

	if val == "" {
		u.RoleId.Valid = false
		u.RoleId.String = ""

		return
	}

	u.RoleId.Valid = true
	u.RoleId.String = val
}

//
func (u *User) GetFirstName() string {

	if u.FirstName.Valid {
		return u.FirstName.String
	}

	return ""
}

//
func (u *User) SetFirstName(val string) {

	if val == "" {
		u.FirstName.Valid = false
		u.FirstName.String = ""

		return
	}

	u.FirstName.Valid = true
	u.FirstName.String = val
}

//
func (u *User) GetLastName() string {

	if u.LastName.Valid {
		return u.LastName.String
	}

	return ""
}

//
func (u *User) SetLastName(val string) {

	if val == "" {
		u.LastName.Valid = false
		u.LastName.String = ""

		return
	}

	u.LastName.Valid = true
	u.LastName.String = val
}

//
func (u *User) GetEmail() string {

	if u.Email.Valid {
		return u.Email.String
	}

	return ""
}

//
func (u *User) SetEmail(val string) {

	if val == "" {
		u.Email.Valid = false
		u.Email.String = ""

		return
	}

	u.Email.Valid = true
	u.Email.String = val
}

//
func (u *User) GetCreatedAt() string {

	if u.CreatedAt.Valid {
		return u.CreatedAt.String
	}

	return ""
}

//
func (u *User) SetCreatedAt(val string) {

	if val == "" {
		u.CreatedAt.Valid = false
		u.CreatedAt.String = ""

		return
	}

	u.CreatedAt.Valid = true
	u.CreatedAt.String = val
}

//
func (u *User) GetDeletedAt() string {

	if u.DeletedAt.Valid {
		return u.DeletedAt.String
	}

	return ""
}

//
func (u *User) SetDeletedAt(val string) {

	if val == "" {
		u.DeletedAt.Valid = false
		u.DeletedAt.String = ""

		return
	}

	u.DeletedAt.Valid = true
	u.DeletedAt.String = val
}

//
func (u *User) GetUpdatedAt() string {

	if u.UpdatedAt.Valid {
		return u.UpdatedAt.String
	}

	return ""
}

//
func (u *User) SetUpdatedAt(val string) {

	if val == "" {
		u.UpdatedAt.Valid = false
		u.UpdatedAt.String = ""

		return
	}

	u.UpdatedAt.Valid = true
	u.UpdatedAt.String = val
}

//
func (u *User) GetPassword() string {

	if u.Password.Valid {
		return u.Password.String
	}

	return ""
}

//
func (u *User) SetPassword(password string, verifyPassword string) {
	u.RawPasswordError = ""
	u.CurrentPasswordError = ""
	u.rawPassword = password
	u.verifyRawPassword = verifyPassword

	if u.currentPassword != "" && !u.Authenticate(u.currentPassword) {
		u.CurrentPasswordError = "Invalid password. You must supply your current password in order to change your password."
		u.Password.Valid = false
		return
	}

	if u.rawPassword == "" || u.rawPassword != u.verifyRawPassword {
		u.RawPasswordError = "Password and Verify Password do not match."
		u.Password.Valid = false
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		u.Password.Valid = false
		u.Password.String = ""
	}

	u.Password.Valid = true
	u.Password.String = string(hash)
}

//
func (u *User) GetVerifiedAt() string {

	if u.VerifiedAt.Valid {
		return u.VerifiedAt.String
	}

	return ""
}

//
func (u *User) SetVerifiedAt(val string) {

	if val == "" {
		u.VerifiedAt.Valid = false
		u.VerifiedAt.String = ""

		return
	}

	u.VerifiedAt.Valid = true
	u.VerifiedAt.String = val
}

// ***

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

	if len(user) == 0 {
		return nil, nil
	}

	user[0].Ctx = ctx
	ctx.SetUser(&user[0])

	return &user[0], nil
}

//
func FetchAllUserByAccountId(ctx ContextInterface, accountId string) ([]User, error) {
	var u []User

	stmt := ctx.Select("*").
		From("public.users").
		Where("account_id = ?", accountId).
		OrderBy("last_name, first_name")

	_, err := stmt.Load(&u)

	if err != nil {
		return nil, err
	}

	if len(u) == 0 {
		return nil, nil
	}

	return u, nil
}

// set user from session
func (u *User) SetFromSession() error {
	var err error
	var shard *Shard
	var user *User
	userEmail, err := u.Ctx.SessionGetString("user_email")

	if err != nil {
		return err
	}

	if userEmail == "" {
		err = errors.New("User not in session.")
		return err
	}

	accountId, _ := u.Ctx.SessionGetString("account_id")

	if accountId != "" {
		shard, _ = FetchShardByAccountId(u.Ctx, accountId)
	}

	if shard == nil {
		user, err = FetchUserByShardEmail(u.Ctx, userEmail)
	} else {
		user, err = FetchUserByEmail(u.Ctx, userEmail)
	}

	if err != nil {
		return err
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
	hash := u.GetPassword()

	if !u.Password.Valid || hash == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		return false
	}

	return true
}

// Make sure user is connected to one account id in the shard_map
// @TODO: logic needs to be better... should not be able to create a 2nd user with the same email etc.
func (u *User) UserShardMapIsValid() (bool, string, error) {
	var accountId string

	if !u.Email.Valid || !u.AccountId.Valid {
		return true, "", nil
	}

	stmt := u.Ctx.Select("account_id").
		From("system.shard_map").
		Where("domain = ?", u.GetEmail()).
		Where("account_id <> ?", u.GetAccountId()).
		Limit(1)

	_, err := stmt.Load(&accountId)

	return accountId == "", accountId, err
}

//
func (u *User) EmailPreviouslyVerified() bool {
	return u.VerifiedAt.Valid
}

//
func (u *User) VerifyEmail(token string) (bool, error) {

	if u.EmailPreviouslyVerified() || !u.Id.Valid {
		return false, nil
	}

	if u.GetId() != token {
		return false, nil
	}

	u.SetVerifiedAt(time.Now().Format(time.RFC3339))

	err := u.Save()

	if err != nil {
		return false, err
	}

	return true, nil
}
