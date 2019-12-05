package jgoweb

import(
	"time"
	"database/sql"
	"github.com/gocraft/web"
	"github.com/jschneider98/jgoweb/util"
)
// Account
type Account struct {
	Id sql.NullString `json:"Id" validate:"omitempty,uuid"`
	Domain sql.NullString `json:"Domain" validate:"required"`
	CreatedAt sql.NullString `json:"CreatedAt" validate:"omitempty,rfc3339"`
	UpdatedAt sql.NullString `json:"UpdatedAt" validate:"omitempty,rfc3339"`
	DeletedAt sql.NullString `json:"DeletedAt" validate:"omitempty,rfc3339"`
	Ctx ContextInterface `json:"-" validate:"-"`
}


// Empty new model
func NewAccount(ctx ContextInterface) (*Account, error) {
	a := &Account{Ctx: ctx}
	a.SetDefaults()

	return a, nil
}

// Set defaults
func (a *Account) SetDefaults() {
	a.SetCreatedAt( time.Now().Format(time.RFC3339) )
	a.SetUpdatedAt( time.Now().Format(time.RFC3339) )
}

// New model with data
func NewAccountWithData(ctx ContextInterface, req *web.Request) (*Account, error) {
	a, err := NewAccount(ctx)

	if err != nil {
		return nil, err
	}

	err = a.Hydrate(req)

	if err != nil {
		return nil, err
	}

	return a, nil
}

// Factory Method
func FetchAccountById(ctx ContextInterface, id string) (*Account, error) {
	var a []Account

	stmt := ctx.Select("*").
	From("public.accounts").
	Where("id = ?", id).
	Limit(1)

	_, err := stmt.Load(&a)

	if err != nil {
		return nil, err
	}

	if (len(a) == 0) {
		return nil, nil
	}

	a[0].Ctx = ctx

	return &a[0], nil
}

//
func (a *Account) ProcessSubmit(req *web.Request) (string, bool, error) {
	err := a.Hydrate(req)

	if err != nil {
		return "", false, err
	}

	err = a.Ctx.GetValidator().Struct(a)

	if err != nil {
		return util.GetNiceErrorMessage(err, "</br>"), false, nil
	}
	
	err = a.Save()

	if err != nil {
		return "", false, err
	}

	return "Account saved.", true, nil
}

// Hydrate the model with data
func (a *Account) Hydrate(req *web.Request) error {
	err := req.ParseForm()

	if err != nil {
		return err
	}

	a.SetId(req.PostFormValue("Id"))
	a.SetDomain(req.PostFormValue("Domain"))
	a.SetCreatedAt(req.PostFormValue("CreatedAt"))
	a.SetUpdatedAt(req.PostFormValue("UpdatedAt"))
	a.SetDeletedAt(req.PostFormValue("DeletedAt"))

	return nil
}

// Validate the model
func (a *Account) IsValid() error {
	return a.Ctx.GetValidator().Struct(a)
}

// Insert/Update based on pkey value
func (a *Account) Save() error {
	err := a.IsValid()

	if err != nil {
		return err
	}

	if !a.Id.Valid {
		return a.Insert()
	} else {
		return a.Update()
	}
}

// Insert a new record
func (a *Account) Insert() error {
	tx, err := a.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	query := `
INSERT INTO
public.accounts (domain,
	deleted_at)
VALUES ($1,$2)
RETURNING id

`

	stmt, err := tx.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	err = stmt.QueryRow(a.Domain,
			a.DeletedAt).Scan(&a.Id)

	if err != nil {
		return err
	}

	return a.Ctx.OptionalCommit(tx)
}

// Update a record
func (a *Account) Update() error {
	if !a.Id.Valid {
		return nil
	}

	tx, err := a.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	a.SetUpdatedAt( time.Now().Format(time.RFC3339) )

	_, err = tx.Update("public.accounts").
		Set("id", a.Id).
		Set("domain", a.Domain).
		Set("updated_at", a.UpdatedAt).
		Set("deleted_at", a.DeletedAt).

		Where("id = ?", a.Id).
		Exec()

	if err != nil {
		return err
	}

	err = a.Ctx.OptionalCommit(tx)

	return err
}

// Soft delete a record
func (a *Account) Delete() error {

	if !a.Id.Valid {
		return nil
	}

	tx, err := a.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	a.SetDeletedAt( (time.Now()).Format(time.RFC3339) )

	_, err = tx.Update("public.accounts").
		Set("deleted_at", a.DeletedAt).
		Where("id = ?", a.Id).
		Exec()

	if err != nil {
		return err
	}

	return a.Ctx.OptionalCommit(tx)
}

// Soft undelete a record
func (a *Account) Undelete() error {

	if !a.Id.Valid {
		return nil
	}

	tx, err := a.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	a.SetDeletedAt("")

	_, err = tx.Update("public.accounts").
		Set("deleted_at", a.DeletedAt).
		Where("id = ?", a.Id).
		Exec()

	if err != nil {
		return err
	}

	return a.Ctx.OptionalCommit(tx)
}

//
func (a *Account) GetId() string {

	if a.Id.Valid {
		return a.Id.String
	}

	return ""
}

//
func (a *Account) SetId(val string) {

	if val == "" {
		a.Id.Valid = false
		a.Id.String = ""

		return
	}

	a.Id.Valid = true
	a.Id.String = val
}

//
func (a *Account) GetDomain() string {

	if a.Domain.Valid {
		return a.Domain.String
	}

	return ""
}

//
func (a *Account) SetDomain(val string) {

	if val == "" {
		a.Domain.Valid = false
		a.Domain.String = ""

		return
	}

	a.Domain.Valid = true
	a.Domain.String = val
}

//
func (a *Account) GetCreatedAt() string {

	if a.CreatedAt.Valid {
		return a.CreatedAt.String
	}

	return ""
}

//
func (a *Account) SetCreatedAt(val string) {

	if val == "" {
		a.CreatedAt.Valid = false
		a.CreatedAt.String = ""

		return
	}

	a.CreatedAt.Valid = true
	a.CreatedAt.String = val
}

//
func (a *Account) GetUpdatedAt() string {

	if a.UpdatedAt.Valid {
		return a.UpdatedAt.String
	}

	return ""
}

//
func (a *Account) SetUpdatedAt(val string) {

	if val == "" {
		a.UpdatedAt.Valid = false
		a.UpdatedAt.String = ""

		return
	}

	a.UpdatedAt.Valid = true
	a.UpdatedAt.String = val
}

//
func (a *Account) GetDeletedAt() string {

	if a.DeletedAt.Valid {
		return a.DeletedAt.String
	}

	return ""
}

//
func (a *Account) SetDeletedAt(val string) {

	if val == "" {
		a.DeletedAt.Valid = false
		a.DeletedAt.String = ""

		return
	}

	a.DeletedAt.Valid = true
	a.DeletedAt.String = val
}
