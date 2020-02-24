package jgoweb


import(
	"fmt"
	"errors"
	"time"
	"database/sql"
	"github.com/gocraft/web"
	"github.com/jschneider98/jgoweb/util"
)
// User
type User struct {
	Id sql.NullString `json:"Id" validate:"omitempty,uuid"`
	AccountId sql.NullString `json:"AccountId" validate:"required,uuid"`
	RoleId sql.NullString `json:"RoleId" validate:"required,int"`
	GivenName sql.NullString `json:"GivenName" validate:"required"`
	FamilyName sql.NullString `json:"FamilyName" validate:"required"`
	Email sql.NullString `json:"Email" validate:"required"`
	CreatedAt sql.NullString `json:"CreatedAt" validate:"omitempty,rfc3339"`
	DeletedAt sql.NullString `json:"DeletedAt" validate:"omitempty,rfc3339"`
	UpdatedAt sql.NullString `json:"UpdatedAt" validate:"omitempty,rfc3339"`
	Ctx ContextInterface `json:"-" validate:"-"`
}


// Empty new model
func NewUser(ctx ContextInterface) (*User, error) {
	u := &User{Ctx: ctx}
	u.SetDefaults()

	return u, nil
}

// Set defaults
func (u *User) SetDefaults() {
	u.SetCreatedAt( time.Now().Format(time.RFC3339) )
	u.SetUpdatedAt( time.Now().Format(time.RFC3339) )
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

	if (len(u) == 0) {
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

	err = u.Ctx.GetValidator().Struct(u)

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
	u.SetGivenName(req.PostFormValue("GivenName"))
	u.SetFamilyName(req.PostFormValue("FamilyName"))
	u.SetEmail(req.PostFormValue("Email"))
	u.SetCreatedAt(req.PostFormValue("CreatedAt"))
	u.SetDeletedAt(req.PostFormValue("DeletedAt"))
	u.SetUpdatedAt(req.PostFormValue("UpdatedAt"))

	return nil
}

// Validate the model
func (u *User) IsValid() error {
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
	tx, err := u.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	query := `
INSERT INTO
public.users (account_id,
	role_id,
	given_name,
	family_name,
	email,
	deleted_at)
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING id

`

	stmt, err := tx.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	err = stmt.QueryRow(u.AccountId,
			u.RoleId,
			u.GivenName,
			u.FamilyName,
			u.Email,
			u.DeletedAt).Scan(&u.Id)

	if err != nil {
		return err
	}

	return u.Ctx.OptionalCommit(tx)
}

// Update a record
func (u *User) Update() error {
	if !u.Id.Valid {
		return nil
	}

	tx, err := u.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	u.SetUpdatedAt( time.Now().Format(time.RFC3339) )

	_, err = tx.Update("public.users").
		Set("id", u.Id).
		Set("account_id", u.AccountId).
		Set("role_id", u.RoleId).
		Set("given_name", u.GivenName).
		Set("family_name", u.FamilyName).
		Set("email", u.Email).
		Set("deleted_at", u.DeletedAt).
		Set("updated_at", u.UpdatedAt).

		Where("id = ?", u.Id).
		Exec()

	if err != nil {
		return err
	}

	err = u.Ctx.OptionalCommit(tx)

	return err
}

// Soft delete a record
func (u *User) Delete() error {

	if !u.Id.Valid {
		return nil
	}

	tx, err := u.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	u.SetDeletedAt( (time.Now()).Format(time.RFC3339) )

	_, err = tx.Update("public.users").
		Set("deleted_at", u.DeletedAt).
		Where("id = ?", u.Id).
		Exec()

	if err != nil {
		return err
	}

	return u.Ctx.OptionalCommit(tx)
}

// Soft undelete a record
func (u *User) Undelete() error {

	if !u.Id.Valid {
		return nil
	}

	tx, err := u.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	u.SetDeletedAt("")

	_, err = tx.Update("public.users").
		Set("deleted_at", u.DeletedAt).
		Where("id = ?", u.Id).
		Exec()

	if err != nil {
		return err
	}

	return u.Ctx.OptionalCommit(tx)
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
func (u *User) GetGivenName() string {

	if u.GivenName.Valid {
		return u.GivenName.String
	}

	return ""
}

//
func (u *User) SetGivenName(val string) {

	if val == "" {
		u.GivenName.Valid = false
		u.GivenName.String = ""

		return
	}

	u.GivenName.Valid = true
	u.GivenName.String = val
}

//
func (u *User) GetFamilyName() string {

	if u.FamilyName.Valid {
		return u.FamilyName.String
	}

	return ""
}

//
func (u *User) SetFamilyName(val string) {

	if val == "" {
		u.FamilyName.Valid = false
		u.FamilyName.String = ""

		return
	}

	u.FamilyName.Valid = true
	u.FamilyName.String = val
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
