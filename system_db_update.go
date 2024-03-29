package jgoweb

import (
	"database/sql"
	"errors"
	"github.com/gocraft/web"
	"github.com/jschneider98/jgoweb/util"
	"time"
)

//
type SystemDbUpdateInterface interface {
	NeedsToRun() (bool, error)
	Run() error
	SetContext(ctx ContextInterface)
	SetComplete() error
	GetUpdateName() string
	GetDescription() string
	Clone() SystemDbUpdateInterface
}

// SystemDbUpdate
type SystemDbUpdate struct {
	ApplyUpdate func(ctx ContextInterface) error `json:"-" validate:"-"`
	Id          sql.NullString                   `json:"Id" validate:"omitempty,int"`
	UpdateName  sql.NullString                   `json:"UpdateName" validate:"required,min=1,max=255"`
	Description sql.NullString                   `json:"Description" validate:"required,min=1,max=255"`
	CreatedAt   sql.NullString                   `json:"CreatedAt" validate:"omitempty,rfc3339"`
	Ctx         ContextInterface                 `json:"-" validate:"-"`
}

// Empty new model
func NewSystemDbUpdate(ctx ContextInterface) (*SystemDbUpdate, error) {
	sdu := &SystemDbUpdate{Ctx: ctx}
	sdu.SetDefaults()

	return sdu, nil
}

// Set defaults
func (sdu *SystemDbUpdate) SetDefaults() {
	sdu.SetCreatedAt(time.Now().Format(time.RFC3339))

}

// New model with data
func NewSystemDbUpdateWithData(ctx ContextInterface, req *web.Request) (*SystemDbUpdate, error) {
	sdu, err := NewSystemDbUpdate(ctx)

	if err != nil {
		return nil, err
	}

	err = sdu.Hydrate(req)

	if err != nil {
		return nil, err
	}

	return sdu, nil
}

// Factory Method
func FetchSystemDbUpdateById(ctx ContextInterface, id string) (*SystemDbUpdate, error) {
	var sdu []SystemDbUpdate

	stmt := ctx.Select("*").
		From("system.db_updates").
		Where("id = ?", id).
		Limit(1)

	_, err := stmt.Load(&sdu)

	if err != nil {
		return nil, err
	}

	if len(sdu) == 0 {
		return nil, nil
	}

	sdu[0].Ctx = ctx

	return &sdu[0], nil
}

// Factory Method
func FetchSystemDbUpdateByUpdateName(ctx ContextInterface, updateName string) (*SystemDbUpdate, error) {
	var sdu []SystemDbUpdate

	stmt := ctx.Select("*").
		From("system.db_updates").
		Where("update_name = ?", updateName).
		Limit(1)

	_, err := stmt.Load(&sdu)

	if err != nil {
		return nil, err
	}

	if len(sdu) == 0 {
		return nil, nil
	}

	sdu[0].Ctx = ctx

	return &sdu[0], nil
}

// Factory Method
func CreateSystemDbUpdateByUpdateName(ctx ContextInterface, updateName string) (*SystemDbUpdate, error) {
	sdu, err := FetchSystemDbUpdateByUpdateName(ctx, updateName)

	if err != nil {
		return nil, err
	}

	if sdu != nil {
		return sdu, nil
	}

	sdu, err = NewSystemDbUpdate(ctx)

	if err != nil {
		return nil, err
	}

	sdu.SetUpdateName(updateName)

	return sdu, nil
}

//
func (sdu *SystemDbUpdate) ProcessSubmit(req *web.Request) (string, bool, error) {
	err := sdu.Hydrate(req)

	if err != nil {
		return "", false, err
	}

	err = sdu.Ctx.GetValidator().Struct(sdu)

	if err != nil {
		return util.GetNiceErrorMessage(err, "</br>"), false, nil
	}

	err = sdu.Save()

	if err != nil {
		return "", false, err
	}

	return "System Db Update saved.", true, nil
}

// Hydrate the model with data
func (sdu *SystemDbUpdate) Hydrate(req *web.Request) error {
	err := req.ParseForm()

	if err != nil {
		return err
	}

	sdu.SetId(req.PostFormValue("Id"))
	sdu.SetUpdateName(req.PostFormValue("UpdateName"))
	sdu.SetDescription(req.PostFormValue("Description"))
	sdu.SetCreatedAt(req.PostFormValue("CreatedAt"))

	return nil
}

// Validate the model
func (sdu *SystemDbUpdate) IsValid() error {
	return sdu.Ctx.GetValidator().Struct(sdu)
}

// Insert/Update based on pkey value
func (sdu *SystemDbUpdate) Save() error {
	err := sdu.IsValid()

	if err != nil {
		return err
	}

	if !sdu.Id.Valid {
		return sdu.Insert()
	} else {
		return sdu.Update()
	}
}

// Insert a new record
func (sdu *SystemDbUpdate) Insert() error {

	query := `
INSERT INTO
system.db_updates (update_name,
	description)
VALUES ($1,$2)
RETURNING id
`

	stmt, err := sdu.Ctx.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	err = stmt.QueryRow(sdu.UpdateName,
		sdu.Description).Scan(&sdu.Id)

	if err != nil {
		return err
	}

	return nil
}

// Update a record
func (sdu *SystemDbUpdate) Update() error {

	if !sdu.Id.Valid {
		return nil
	}

	_, err := sdu.Ctx.Update("system.db_updates").
		Set("id", sdu.Id).
		Set("update_name", sdu.UpdateName).
		Set("description", sdu.Description).
		Where("id = ?", sdu.Id).
		Exec()

	if err != nil {
		return err
	}

	return nil
}

// Hard delete a record
func (sdu *SystemDbUpdate) Delete() error {

	if !sdu.Id.Valid {
		return nil
	}

	_, err := sdu.Ctx.DeleteFrom("system.db_updates").
		Where("id = ?", sdu.Id).
		Exec()

	if err != nil {
		return err
	}

	return nil
}

//
func (sdu *SystemDbUpdate) GetId() string {

	if sdu.Id.Valid {
		return sdu.Id.String
	}

	return ""
}

//
func (sdu *SystemDbUpdate) SetId(val string) {

	if val == "" {
		sdu.Id.Valid = false
		sdu.Id.String = ""

		return
	}

	sdu.Id.Valid = true
	sdu.Id.String = val
}

//
func (sdu *SystemDbUpdate) GetUpdateName() string {

	if sdu.UpdateName.Valid {
		return sdu.UpdateName.String
	}

	return ""
}

//
func (sdu *SystemDbUpdate) SetUpdateName(val string) {

	if val == "" {
		sdu.UpdateName.Valid = false
		sdu.UpdateName.String = ""

		return
	}

	sdu.UpdateName.Valid = true
	sdu.UpdateName.String = val
}

//
func (sdu *SystemDbUpdate) GetDescription() string {

	if sdu.Description.Valid {
		return sdu.Description.String
	}

	return ""
}

//
func (sdu *SystemDbUpdate) SetDescription(val string) {

	if val == "" {
		sdu.Description.Valid = false
		sdu.Description.String = ""

		return
	}

	sdu.Description.Valid = true
	sdu.Description.String = val
}

//
func (sdu *SystemDbUpdate) GetCreatedAt() string {

	if sdu.CreatedAt.Valid {
		return sdu.CreatedAt.String
	}

	return ""
}

//
func (sdu *SystemDbUpdate) SetCreatedAt(val string) {

	if val == "" {
		sdu.CreatedAt.Valid = false
		sdu.CreatedAt.String = ""

		return
	}

	sdu.CreatedAt.Valid = true
	sdu.CreatedAt.String = val
}

// ****** Interface Methods ******

// Empty new update
func CreateSystemDbUpdateNoContext(updateName string, desc string) *SystemDbUpdate {
	sdu := &SystemDbUpdate{}

	sdu.SetUpdateName(updateName)
	sdu.SetDescription(desc)

	return sdu
}

//
func (sdu *SystemDbUpdate) SetContext(ctx ContextInterface) {
	sdu.Ctx = ctx
}

//
func (sdu *SystemDbUpdate) NeedsToRun() (bool, error) {
	var err error
	updateName := sdu.UpdateName.String
	description := sdu.Description

	if sdu.Ctx == nil {
		err := errors.New("Context not set in SystemDbUpdate.NeedsToRun()")
		return false, err
	}

	sdu, err = CreateSystemDbUpdateByUpdateName(sdu.Ctx, updateName)
	sdu.Description = description

	if err != nil {
		return false, err
	}

	return !sdu.Id.Valid, nil
}

//
func (sdu *SystemDbUpdate) Run() error {
	return sdu.ApplyUpdate(sdu.Ctx)
}

//
func (sdu *SystemDbUpdate) SetComplete() error {
	return sdu.Save()
}

//
func (sdu *SystemDbUpdate) Clone() SystemDbUpdateInterface {
	clone := *sdu
	clone.Ctx = nil

	return &clone
}
