package jgoweb

import(
	"time"
	"database/sql"
	"github.com/gocraft/web"
	"github.com/jschneider98/jgoweb/util"
)
// ShardMap
type ShardMap struct {
	Id sql.NullString `json:"Id" validate:"omitempty,int"`
	ShardId sql.NullString `json:"ShardId" validate:"required,int"`
	Domain sql.NullString `json:"Domain" validate:"required"`
	AccountId sql.NullString `json:"AccountId" validate:"required,uuid"`
	CreatedAt sql.NullString `json:"CreatedAt" validate:"omitempty,rfc3339"`
	UpdatedAt sql.NullString `json:"UpdatedAt" validate:"omitempty,rfc3339"`
	DeletedAt sql.NullString `json:"DeletedAt" validate:"omitempty,rfc3339"`
	Ctx ContextInterface `json:"-" validate:"-"`
}


// Empty new model
func NewShardMap(ctx ContextInterface) (*ShardMap, error) {
	sm := &ShardMap{Ctx: ctx}
	sm.SetDefaults()

	return sm, nil
}

// Set defaults
func (sm *ShardMap) SetDefaults() {
	sm.SetCreatedAt( time.Now().Format(time.RFC3339) )
	sm.SetUpdatedAt( time.Now().Format(time.RFC3339) )

}

// New model with data
func NewShardMapWithData(ctx ContextInterface, req *web.Request) (*ShardMap, error) {
	sm, err := NewShardMap(ctx)

	if err != nil {
		return nil, err
	}

	err = sm.Hydrate(req)

	if err != nil {
		return nil, err
	}

	return sm, nil
}

// Factory Method
func FetchShardMapById(ctx ContextInterface, id string) (*ShardMap, error) {
	var sm []ShardMap

	stmt := ctx.Select("*").
	From("public.shard_map").
	Where("id = ?", id).
	Limit(1)

	_, err := stmt.Load(&sm)

	if err != nil {
		return nil, err
	}

	if (len(sm) == 0) {
		return nil, nil
	}

	sm[0].Ctx = ctx

	return &sm[0], nil
}

//
func (sm *ShardMap) ProcessSubmit(req *web.Request) (string, bool, error) {
	err := sm.Hydrate(req)

	if err != nil {
		return "", false, err
	}

	err = sm.Ctx.GetValidator().Struct(sm)

	if err != nil {
		return util.GetNiceErrorMessage(err, "</br>"), false, nil
	}
	
	err = sm.Save()

	if err != nil {
		return "", false, err
	}

	return "Shard Map saved.", true, nil
}

// Hydrate the model with data
func (sm *ShardMap) Hydrate(req *web.Request) error {
	err := req.ParseForm()

	if err != nil {
		return err
	}

	sm.SetId(req.PostFormValue("Id"))
	sm.SetShardId(req.PostFormValue("ShardId"))
	sm.SetDomain(req.PostFormValue("Domain"))
	sm.SetAccountId(req.PostFormValue("AccountId"))
	sm.SetCreatedAt(req.PostFormValue("CreatedAt"))
	sm.SetUpdatedAt(req.PostFormValue("UpdatedAt"))
	sm.SetDeletedAt(req.PostFormValue("DeletedAt"))

	return nil
}

// Validate the model
func (sm *ShardMap) IsValid() error {
	return sm.Ctx.GetValidator().Struct(sm)
}

// Insert/Update based on pkey value
func (sm *ShardMap) Save() error {
	err := sm.IsValid()

	if err != nil {
		return err
	}

	if !sm.Id.Valid {
		return sm.Insert()
	} else {
		return sm.Update()
	}
}

// Insert a new record
func (sm *ShardMap) Insert() error {
	tx, err := sm.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	query := `
INSERT INTO
public.shard_map (shard_id,
	domain,
	account_id,
	deleted_at)
VALUES ($1,$2,$3,$4)
RETURNING id

`

	stmt, err := tx.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	err = stmt.QueryRow(sm.ShardId,
			sm.Domain,
			sm.AccountId,
			sm.DeletedAt).Scan(&sm.Id)

	if err != nil {
		return err
	}

	return sm.Ctx.OptionalCommit(tx)
}

// Update a record
func (sm *ShardMap) Update() error {
	if !sm.Id.Valid {
		return nil
	}

	tx, err := sm.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	sm.SetUpdatedAt( time.Now().Format(time.RFC3339) )

	_, err = tx.Update("public.shard_map").
		Set("id", sm.Id).
		Set("shard_id", sm.ShardId).
		Set("domain", sm.Domain).
		Set("account_id", sm.AccountId).
		Set("updated_at", sm.UpdatedAt).
		Set("deleted_at", sm.DeletedAt).

		Where("id = ?", sm.Id).
		Exec()

	if err != nil {
		return err
	}

	err = sm.Ctx.OptionalCommit(tx)

	return err
}

// Soft delete a record
func (sm *ShardMap) Delete() error {

	if !sm.Id.Valid {
		return nil
	}

	tx, err := sm.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	sm.SetDeletedAt( (time.Now()).Format(time.RFC3339) )

	_, err = tx.Update("public.shard_map").
		Set("deleted_at", sm.DeletedAt).
		Where("id = ?", sm.Id).
		Exec()

	if err != nil {
		return err
	}

	return sm.Ctx.OptionalCommit(tx)
}

//
func (sm *ShardMap) GetId() string {

	if sm.Id.Valid {
		return sm.Id.String
	}

	return ""
}

//
func (sm *ShardMap) SetId(val string) {

	if val == "" {
		sm.Id.Valid = false
		sm.Id.String = ""

		return
	}

	sm.Id.Valid = true
	sm.Id.String = val
}

//
func (sm *ShardMap) GetShardId() string {

	if sm.ShardId.Valid {
		return sm.ShardId.String
	}

	return ""
}

//
func (sm *ShardMap) SetShardId(val string) {

	if val == "" {
		sm.ShardId.Valid = false
		sm.ShardId.String = ""

		return
	}

	sm.ShardId.Valid = true
	sm.ShardId.String = val
}

//
func (sm *ShardMap) GetDomain() string {

	if sm.Domain.Valid {
		return sm.Domain.String
	}

	return ""
}

//
func (sm *ShardMap) SetDomain(val string) {

	if val == "" {
		sm.Domain.Valid = false
		sm.Domain.String = ""

		return
	}

	sm.Domain.Valid = true
	sm.Domain.String = val
}

//
func (sm *ShardMap) GetAccountId() string {

	if sm.AccountId.Valid {
		return sm.AccountId.String
	}

	return ""
}

//
func (sm *ShardMap) SetAccountId(val string) {

	if val == "" {
		sm.AccountId.Valid = false
		sm.AccountId.String = ""

		return
	}

	sm.AccountId.Valid = true
	sm.AccountId.String = val
}

//
func (sm *ShardMap) GetCreatedAt() string {

	if sm.CreatedAt.Valid {
		return sm.CreatedAt.String
	}

	return ""
}

//
func (sm *ShardMap) SetCreatedAt(val string) {

	if val == "" {
		sm.CreatedAt.Valid = false
		sm.CreatedAt.String = ""

		return
	}

	sm.CreatedAt.Valid = true
	sm.CreatedAt.String = val
}

//
func (sm *ShardMap) GetUpdatedAt() string {

	if sm.UpdatedAt.Valid {
		return sm.UpdatedAt.String
	}

	return ""
}

//
func (sm *ShardMap) SetUpdatedAt(val string) {

	if val == "" {
		sm.UpdatedAt.Valid = false
		sm.UpdatedAt.String = ""

		return
	}

	sm.UpdatedAt.Valid = true
	sm.UpdatedAt.String = val
}

//
func (sm *ShardMap) GetDeletedAt() string {

	if sm.DeletedAt.Valid {
		return sm.DeletedAt.String
	}

	return ""
}

//
func (sm *ShardMap) SetDeletedAt(val string) {

	if val == "" {
		sm.DeletedAt.Valid = false
		sm.DeletedAt.String = ""

		return
	}

	sm.DeletedAt.Valid = true
	sm.DeletedAt.String = val
}

// ******

// Factory Method
func FetchShardMapByAccountId(ctx ContextInterface, accountId string) (*ShardMap, error) {
	var sm []ShardMap

	stmt := ctx.Select("*").
	From("public.shard_map").
	Where("account_id = ?", accountId).
	Limit(1)

	_, err := stmt.Load(&sm)

	if err != nil {
		return nil, err
	}

	if (len(sm) == 0) {
		return nil, nil
	}

	sm[0].Ctx = ctx

	return &sm[0], nil
}

// 
func GetAllShardMaps(ctx ContextInterface) ([]ShardMap, error) {
	var sm []ShardMap

	stmt := ctx.Select("*").
	From("public.shard_map").
	OrderBy("domain")

	_, err := stmt.Load(&sm)

	if err != nil {
		return nil, err
	}

	if sm == nil {
		sm = make([]ShardMap, 0)
	}

	return sm, nil
}
