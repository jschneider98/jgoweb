package jgoweb

import(
	"database/sql"
	"github.com/gocraft/web"
	"github.com/jschneider98/jgoweb/util"
)

// ShardMap
type ShardMap struct {
	ShardId sql.NullString `json:"ShardId" validate:"required,int"`
	Domain sql.NullString `json:"Domain" validate:"required"`
	AccountId sql.NullString `json:"AccountId" validate:"required,uuid"`
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

	sm.SetShardId(req.PostFormValue("ShardId"))
	sm.SetDomain(req.PostFormValue("Domain"))
	sm.SetAccountId(req.PostFormValue("AccountId"))

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

	if !sm.AccountId.Valid {
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
	account_id)
VALUES ($1,$2,$3)
RETURNING account_id

`

	stmt, err := tx.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	err = stmt.QueryRow(sm.ShardId,
			sm.Domain,
			sm.AccountId).Scan(&sm.AccountId)

	if err != nil {
		return err
	}

	return sm.Ctx.OptionalCommit(tx)
}

// Update a record
func (sm *ShardMap) Update() error {
	if !sm.AccountId.Valid {
		return nil
	}

	tx, err := sm.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	

	_, err = tx.Update("public.shard_map").
		Set("shard_id", sm.ShardId).
		Set("domain", sm.Domain).
		Set("account_id", sm.AccountId).

		Where("account_id = ?", sm.AccountId).
		Exec()

	if err != nil {
		return err
	}

	err = sm.Ctx.OptionalCommit(tx)

	return err
}

// Hard delete a record
func (sm *ShardMap) Delete() error {

	if !sm.AccountId.Valid {
		return nil
	}

	tx, err := sm.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	_, err = tx.DeleteFrom("public.shard_map").
		Where("account_id = ?", sm.AccountId).
		Exec()

	if err != nil {
		return err
	}

	return sm.Ctx.OptionalCommit(tx)
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
