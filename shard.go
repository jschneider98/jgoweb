package jgoweb

import(
	"time"
	"errors"
	"database/sql"
	"github.com/gocraft/web"
	"github.com/jschneider98/jgoweb/util"
)
// Shard
type Shard struct {
	Id sql.NullString `json:"Id" validate:"omitempty,int"`
	Name sql.NullString `json:"Name" validate:"required"`
	AccountCount sql.NullString `json:"AccountCount" validate:"required,int"`
	CreatedAt sql.NullString `json:"CreatedAt" validate:"omitempty,rfc3339"`
	UpdatedAt sql.NullString `json:"UpdatedAt" validate:"omitempty,rfc3339"`
	DeletedAt sql.NullString `json:"DeletedAt" validate:"omitempty,rfc3339"`
	Ctx ContextInterface `json:"-" validate:"-"`
}


// Empty new model
func NewShard(ctx ContextInterface) (*Shard, error) {
	s := &Shard{Ctx: ctx}
	s.SetDefaults()

	return s, nil
}

// Set defaults
func (s *Shard) SetDefaults() {
	s.SetCreatedAt( time.Now().Format(time.RFC3339) )
	s.SetUpdatedAt( time.Now().Format(time.RFC3339) )
}

// New model with data
func NewShardWithData(ctx ContextInterface, req *web.Request) (*Shard, error) {
	s, err := NewShard(ctx)

	if err != nil {
		return nil, err
	}

	err = s.Hydrate(req)

	if err != nil {
		return nil, err
	}

	return s, nil
}

// Factory Method
func FetchShardById(ctx ContextInterface, id string) (*Shard, error) {
	var s []Shard

	stmt := ctx.Select("*").
	From("public.shards").
	Where("id = ?", id).
	Limit(1)

	_, err := stmt.Load(&s)

	if err != nil {
		return nil, err
	}

	if (len(s) == 0) {
		return nil, nil
	}

	s[0].Ctx = ctx

	return &s[0], nil
}

//
func (s *Shard) ProcessSubmit(req *web.Request) (string, bool, error) {
	err := s.Hydrate(req)

	if err != nil {
		return "", false, err
	}

	err = s.Ctx.GetValidator().Struct(s)

	if err != nil {
		return util.GetNiceErrorMessage(err, "</br>"), false, nil
	}
	
	err = s.Save()

	if err != nil {
		return "", false, err
	}

	return "Shard saved.", true, nil
}

// Hydrate the model with data
func (s *Shard) Hydrate(req *web.Request) error {
	err := req.ParseForm()

	if err != nil {
		return err
	}

	s.SetId(req.PostFormValue("Id"))
	s.SetName(req.PostFormValue("Name"))
	s.SetAccountCount(req.PostFormValue("AccountCount"))
	s.SetCreatedAt(req.PostFormValue("CreatedAt"))
	s.SetUpdatedAt(req.PostFormValue("UpdatedAt"))
	s.SetDeletedAt(req.PostFormValue("DeletedAt"))

	return nil
}

// Validate the model
func (s *Shard) IsValid() error {
	return s.Ctx.GetValidator().Struct(s)
}

// Insert/Update based on pkey value
func (s *Shard) Save() error {
	err := s.IsValid()

	if err != nil {
		return err
	}

	if !s.Id.Valid {
		return s.Insert()
	} else {
		return s.Update()
	}
}

// Insert a new record
func (s *Shard) Insert() error {
	tx, err := s.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	query := `
INSERT INTO
public.shards (name,
	account_count,
	deleted_at)
VALUES ($1,$2,$3)
RETURNING id

`

	stmt, err := tx.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	err = stmt.QueryRow(s.Name,
			s.AccountCount,
			s.DeletedAt).Scan(&s.Id)

	if err != nil {
		return err
	}

	return s.Ctx.OptionalCommit(tx)
}

// Update a record
func (s *Shard) Update() error {
	if !s.Id.Valid {
		return nil
	}

	tx, err := s.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	s.SetUpdatedAt( time.Now().Format(time.RFC3339) )

	_, err = tx.Update("public.shards").
		Set("id", s.Id).
		Set("name", s.Name).
		Set("account_count", s.AccountCount).
		Set("updated_at", s.UpdatedAt).
		Set("deleted_at", s.DeletedAt).

		Where("id = ?", s.Id).
		Exec()

	if err != nil {
		return err
	}

	err = s.Ctx.OptionalCommit(tx)

	return err
}

// Soft delete a record
func (s *Shard) Delete() error {

	if !s.Id.Valid {
		return nil
	}

	tx, err := s.Ctx.OptionalBegin()

	if err != nil {
		return err
	}

	s.SetDeletedAt( (time.Now()).Format(time.RFC3339) )

	_, err = tx.Update("public.shards").
		Set("deleted_at", s.DeletedAt).
		Where("id = ?", s.Id).
		Exec()

	if err != nil {
		return err
	}

	return s.Ctx.OptionalCommit(tx)
}

//
func (s *Shard) GetId() string {

	if s.Id.Valid {
		return s.Id.String
	}

	return ""
}

//
func (s *Shard) SetId(val string) {

	if val == "" {
		s.Id.Valid = false
		s.Id.String = ""

		return
	}

	s.Id.Valid = true
	s.Id.String = val
}

//
func (s *Shard) GetName() string {

	if s.Name.Valid {
		return s.Name.String
	}

	return ""
}

//
func (s *Shard) SetName(val string) {

	if val == "" {
		s.Name.Valid = false
		s.Name.String = ""

		return
	}

	s.Name.Valid = true
	s.Name.String = val
}

//
func (s *Shard) GetAccountCount() string {

	if s.AccountCount.Valid {
		return s.AccountCount.String
	}

	return ""
}

//
func (s *Shard) SetAccountCount(val string) {

	if val == "" {
		s.AccountCount.Valid = false
		s.AccountCount.String = ""

		return
	}

	s.AccountCount.Valid = true
	s.AccountCount.String = val
}

//
func (s *Shard) GetCreatedAt() string {

	if s.CreatedAt.Valid {
		return s.CreatedAt.String
	}

	return ""
}

//
func (s *Shard) SetCreatedAt(val string) {

	if val == "" {
		s.CreatedAt.Valid = false
		s.CreatedAt.String = ""

		return
	}

	s.CreatedAt.Valid = true
	s.CreatedAt.String = val
}

//
func (s *Shard) GetUpdatedAt() string {

	if s.UpdatedAt.Valid {
		return s.UpdatedAt.String
	}

	return ""
}

//
func (s *Shard) SetUpdatedAt(val string) {

	if val == "" {
		s.UpdatedAt.Valid = false
		s.UpdatedAt.String = ""

		return
	}

	s.UpdatedAt.Valid = true
	s.UpdatedAt.String = val
}

//
func (s *Shard) GetDeletedAt() string {

	if s.DeletedAt.Valid {
		return s.DeletedAt.String
	}

	return ""
}

//
func (s *Shard) SetDeletedAt(val string) {

	if val == "" {
		s.DeletedAt.Valid = false
		s.DeletedAt.String = ""

		return
	}

	s.DeletedAt.Valid = true
	s.DeletedAt.String = val
}

// ******

// 
func FetchShardByAccountId(ctx ContextInterface, accountId string) (*Shard, error) {
	var shards []Shard

	// Shard data is stored on every DB
	dbConn, err := ctx.GetDb().GetRandomConn()

	if err != nil {
		return nil, err
	}

	dbSess := dbConn.NewSession(nil)

	stmt := dbSess.SelectBySql(`
	SELECT
		s.*
	FROM public.shard_map sm
	JOIN public.shards s ON s.id = sm.shard_id
	WHERE sm.account_id = ?
	LIMIT 1`,
	accountId)

	_, err = stmt.Load(&shards)

	if err != nil {
		return nil, err
	}

	if (len(shards) == 0) {
		return nil, nil
	}

	// Set db session for this shard
	dbSess, err = ctx.GetDb().GetSessionByName(shards[0].GetName())

	if err != nil {
		return nil, err
	}

	ctx.SetDbSession(dbSess)
	shards[0].Ctx = ctx

	return &shards[0], nil
}

// Does not set db session. Mainly used for logical replication of shards.
func FetchShardByName(ctx ContextInterface, shardName string) (*Shard, error) {
	var shards []Shard
	var err error

	stmt := ctx.SelectBySql(`
	SELECT
		s.*
	FROM public.shards s
	WHERE s.name = ?
	LIMIT 1`,
	shardName)

	_, err = stmt.Load(&shards)

	if err != nil {
		return nil, err
	}

	if (len(shards) == 0) {
		return nil, nil
	}

	shards[0].Ctx = ctx

	return &shards[0], nil
}

//
func CreateShardByName(ctx ContextInterface, shardName string) (*Shard, error) {
	shard, err := FetchShardByName(ctx, shardName)

	if err != nil {
		return nil, err
	}

	if shard != nil {
		return shard, nil
	}

	shard, err = NewShard(ctx)

	if err != nil {
		return nil, err
	}

	shard.SetName(shardName)
	shard.SetAccountCount("0")

	return shard, nil
}

// 
func FetchShardByEmail(ctx ContextInterface, email string) (*Shard, error) {
	var shards []Shard

	// Shard data is stored on every DB
	dbConn, err := ctx.GetDb().GetRandomConn()

	if err != nil {
		return nil, err
	}

	dbSess := dbConn.NewSession(nil)

	// @TODO: domain shouldn't be full email for non-personal accounts
	stmt := dbSess.SelectBySql(`
	SELECT
		s.*
	FROM public.shard_map sm
	JOIN public.shards s ON s.id = sm.shard_id
	WHERE sm.domain = ?
	LIMIT 1`,
	email)

	_, err = stmt.Load(&shards)

	if err != nil {
		return nil, err
	}

	if (len(shards) == 0) {
		return nil, nil
	}

	// Set db session for this shard
	dbSess, err = ctx.GetDb().GetSessionByName(shards[0].GetName())

	if err != nil {
		return nil, err
	}

	ctx.SetDbSession(dbSess)
	shards[0].Ctx = ctx

	return &shards[0], nil
}

// 
func FetchUserByShardEmail(ctx ContextInterface, email string) (*User, error) {
	shard, err := FetchShardByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	if shard == nil {
		return nil, nil
	}
	
	user, err := FetchUserByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// 
func GetAllShards(ctx ContextInterface) ([]Shard, error) {
	var s []Shard

	stmt := ctx.Select("*").
	From("public.shards").
	OrderBy("account_count, name")

	_, err := stmt.Load(&s)

	if err != nil {
		return nil, err
	}

	if s == nil {
		s = make([]Shard, 0)
	}

	return s, nil
}

// Get shard data for all DBs
func ClusterGetShards(ctx ContextInterface) (map[string][]Shard, error) {
	var err error
	shards := make(map[string][]Shard)

	for dbName, dbConn := range ctx.GetDb().GetConns() {
		curCtx := NewContext(ctx.GetDb())
		curCtx.SetDbSession(dbConn.NewSession(nil))
		
		shards[dbName], err = GetAllShards(curCtx)

		if err != nil {
			return nil, err
		}
	}

	return shards, nil
}

// 
func ClusterAddShard(ctx ContextInterface, shardName string) error {
	var err error
	var shard *Shard
	var msg string

	for dbName, dbConn := range ctx.GetDb().GetConns() {
		curCtx := NewContext(ctx.GetDb())
		curCtx.SetDbSession(dbConn.NewSession(nil))
		
		shard, err = CreateShardByName(curCtx, shardName)

		if err != nil {
			msg += dbName + ": " + err.Error() + "\n"
		} else {
			err = shard.Save()

			if err != nil {
				msg += dbName + ": " + err.Error() + "\n"
			}
		}
	}

	if msg != "" {
		err = errors.New(msg)

		return err
	}

	return nil
}

// 
func ClusterDeleteShard(ctx ContextInterface, shardName string) error {
	var err error
	var shard *Shard
	var msg string

	for dbName, dbConn := range ctx.GetDb().GetConns() {
		curCtx := NewContext(ctx.GetDb())
		curCtx.SetDbSession(dbConn.NewSession(nil))
		
		shard, err = FetchShardByName(curCtx, shardName)

		if err != nil {
			msg += dbName + ": " + err.Error() + "\n"
		} else if shard != nil {
			err = shard.Delete()

			if err != nil {
				msg += dbName + ": " + err.Error() + "\n"
			}
		}
	}

	if msg != "" {
		err = errors.New(msg)

		return err
	}

	return nil
}
