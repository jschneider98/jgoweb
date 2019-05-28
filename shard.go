package jgoweb

import (
	"database/sql"
)

//
type Shard struct {
	Id string `json:"id"`
	Name string `json:"name"`
	AccountCount int `json:"account_count"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt sql.NullString `json:"deleted_at"`
	Ctx ContextInterface
}

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
	dbSess, err = ctx.GetDb().GetSessionByName(shards[0].Name)

	if err != nil {
		return nil, err
	}

	ctx.SetDbSession(dbSess)
	shards[0].Ctx = ctx

	return &shards[0], nil
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
	dbSess, err = ctx.GetDb().GetSessionByName(shards[0].Name)

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
