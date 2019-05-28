// +build integration

package jgoweb

import (
	"testing"
	"database/sql"
)

//
func TestFetchUserByEmail(t *testing.T) {
	InitTestUser()

	user, err := FetchUserByEmail(testUser.Ctx, testUser.Email)

	if err != nil {
		t.Errorf("Failed to fetch user by email: %v", err)
	}

	if user == nil {
		t.Errorf("Failed to fetch user %v: ", testUser.Email)
	}

	if user.Email != testUser.Email {
		t.Errorf("Fetched wrong user? Expected: %v Got: %v", testUser.Email, user.Email)
	}

	testUser.Id = user.Id
	testUser.AccountId = user.AccountId
	testUser.RoleId = user.RoleId
	testUser.GivenName = user.GivenName
	testUser.FamilyName = user.FamilyName
	testUser.Email = user.Email
	testUser.CreatedAt = user.CreatedAt
	testUser.UpdatedAt = user.UpdatedAt
	testUser.DeletedAt = user.DeletedAt

	userName := "not_a_user"
	user, err = FetchUserByEmail(testUser.Ctx, userName)

	if err != nil {
		t.Errorf("Failed to fetch user by email: %v", err)
	}

	if user != nil {
		t.Errorf("Should have failed to find user: %v", userName)
	}
}

//
func TestFetchUserById(t *testing.T) {
	user, err := FetchUserById(testUser.Ctx, testUser.Id)

	if err != nil {
		t.Errorf("Failed to fetch user by Id: %v", err)
	}

	if user == nil {
		t.Errorf("Failed to fetch user %v: ", testUser.Email)
	}

	if user.Email != testUser.Email {
		t.Errorf("Fetched wrong user? Expected: %v Got: %v", testUser.Email, user.Email)
	}

	// force error
	userId := "not_a_user"
	_, err = FetchUserById(testUser.Ctx, userId)

	if err == nil {
		t.Errorf("FetchUserById should have failed with uuid of: %v", userId)
	}

	// force not found
	userId = "00000000-0000-0000-0000-000000000000"
	user, err = FetchUserById(testUser.Ctx, userId)

	if err != nil {
		t.Errorf("Failed to fetch user by id: %v", err)
	}

	if user != nil {
		t.Errorf("Should have failed to find user: %v", userId)
	}
}

//
func TestUserGetProjectList(t *testing.T) {
	_, err := testUser.GetProjectList(ProjectListParams{})

	if err != nil {
		t.Errorf("Failed to get project list %v: ", err)
	}
}

//
func TestUserIsValid(t *testing.T) {
	isValid, err := testUser.isValid()

	if !isValid {
		t.Errorf("User should be valid: %v", err)
	}
}

//
func TestUserBadId(t *testing.T) {
	id := testUser.Id
	testUser.Id = "bad_uuid"

	isValid, err := testUser.isValid()

	if isValid {
		t.Errorf("User.Id should be invalid: %v", testUser.Id)
	}

	if err == nil || err.Error() != "id: bad_uuid does not validate as uuid" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	testUser.Id = id
}

//
func TestUserOptionalId(t *testing.T) {
	id := testUser.Id
	testUser.Id = ""

	isValid, err := testUser.isValid()

	if !isValid {
		t.Errorf("User.Id should be valid (optional): %v %v", testUser.Id, err)
	}

	testUser.Id = id
}

//
func TestUserBadAccountId(t *testing.T) {
	id := testUser.AccountId
	testUser.AccountId = "bad_uuid"

	isValid, err := testUser.isValid()

	if isValid {
		t.Errorf("User.AccountId should be invalid: %v", testUser.AccountId)
	}

	if err == nil || err.Error() != "account_id: bad_uuid does not validate as uuid" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	testUser.AccountId = id
}

//
func TestUserBadRoleId(t *testing.T) {
	id := testUser.RoleId
	testUser.RoleId = 0

	isValid, err := testUser.isValid()

	if isValid {
		t.Errorf("User.RoleId should be invalid: %v", testUser.RoleId)
	}

	if err == nil || err.Error() != "role_id: non zero value required" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	//
	testUser.RoleId = 101

	isValid, err = testUser.isValid()

	if isValid {
		t.Errorf("User.RoleId should be invalid: %v", testUser.RoleId)
	}

	if err == nil || err.Error() != "role_id: 101 does not validate as range(1|100)" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	testUser.RoleId = id
}

//
func TestUserBadGivenName(t *testing.T) {
	longStr := "IdViTf7vd9NezkhR8Ftvh4nTVe8re2RGvFGYkMN9alkUQxm7ZrEuHVVpr6CUgK5pdLTh8H4KVVTscwMlQL9ZJF0kcKJuyMizgKkorbNeblelBECnE2G6hxSecsL69dh9eXQGUSAbbdgB6BE3Q11Ffm4GjRbz4Z4cW6D2ZyZ4RIpFtHFlQcsiD2o8QwCitr7LAIRAwSW2DHenEUYh1OVDfpUFMMUtuvnRCYghwNj8iFwrdp3rxKeTBsZU5CdbVC3"
	name := testUser.GivenName
	testUser.GivenName = ""

	isValid, err := testUser.isValid()

	if isValid {
		t.Errorf("User.GivenName should be invalid: %v", testUser.GivenName)
	}

	if err == nil || err.Error() != "given_name: non zero value required" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	//

	testUser.GivenName = longStr

	isValid, err = testUser.isValid()

	if isValid {
		t.Errorf("User.GivenName should be invalid: %v", testUser.GivenName)
	}

	if err == nil || err.Error() != "given_name: IdViTf7vd9NezkhR8Ftvh4nTVe8re2RGvFGYkMN9alkUQxm7ZrEuHVVpr6CUgK5pdLTh8H4KVVTscwMlQL9ZJF0kcKJuyMizgKkorbNeblelBECnE2G6hxSecsL69dh9eXQGUSAbbdgB6BE3Q11Ffm4GjRbz4Z4cW6D2ZyZ4RIpFtHFlQcsiD2o8QwCitr7LAIRAwSW2DHenEUYh1OVDfpUFMMUtuvnRCYghwNj8iFwrdp3rxKeTBsZU5CdbVC3 does not validate as length(1|254)" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	testUser.GivenName = name
}

//
func TestUserBadFamilyName(t *testing.T) {
	longStr := "IdViTf7vd9NezkhR8Ftvh4nTVe8re2RGvFGYkMN9alkUQxm7ZrEuHVVpr6CUgK5pdLTh8H4KVVTscwMlQL9ZJF0kcKJuyMizgKkorbNeblelBECnE2G6hxSecsL69dh9eXQGUSAbbdgB6BE3Q11Ffm4GjRbz4Z4cW6D2ZyZ4RIpFtHFlQcsiD2o8QwCitr7LAIRAwSW2DHenEUYh1OVDfpUFMMUtuvnRCYghwNj8iFwrdp3rxKeTBsZU5CdbVC3"
	name := testUser.FamilyName
	testUser.FamilyName = ""

	isValid, err := testUser.isValid()

	if isValid {
		t.Errorf("User.FamilyName should be invalid: %v", testUser.FamilyName)
	}

	if err == nil || err.Error() != "family_name: non zero value required" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	//

	testUser.FamilyName = longStr

	isValid, err = testUser.isValid()

	if isValid {
		t.Errorf("User.FamilyName should be invalid: %v", testUser.FamilyName)
	}

	if err == nil || err.Error() != "family_name: IdViTf7vd9NezkhR8Ftvh4nTVe8re2RGvFGYkMN9alkUQxm7ZrEuHVVpr6CUgK5pdLTh8H4KVVTscwMlQL9ZJF0kcKJuyMizgKkorbNeblelBECnE2G6hxSecsL69dh9eXQGUSAbbdgB6BE3Q11Ffm4GjRbz4Z4cW6D2ZyZ4RIpFtHFlQcsiD2o8QwCitr7LAIRAwSW2DHenEUYh1OVDfpUFMMUtuvnRCYghwNj8iFwrdp3rxKeTBsZU5CdbVC3 does not validate as length(1|254)" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	testUser.FamilyName = name
}

//
func TestUserBadEmail(t *testing.T) {
	longStr := "IdViTf7vd9NezkhR8Ftvh4nTVe8re2RGvFGYkMN9alkUQxm7ZrEuHVVpr6CUgK5pdLTh8H4KVVTscwMlQL9ZJF0kcKJuyMizgKkorbNeblelBECnE2G6hxSecsL69dh9eXQGUSAbbdgB6BE3Q11Ffm4GjRbz4Z4cW6D2ZyZ4RIpFtHFlQcsiD2o8QwCitr7LAIRAwSW2DHenEUYh1OVDfpUFMMUtuvnRCYghwNj8iFwrdp3rxKeTBsZU5CdbVC3@gmail.com"
	email := testUser.Email
	testUser.Email = ""

	isValid, err := testUser.isValid()

	if isValid {
		t.Errorf("User.Email should be invalid: %v", testUser.Email)
	}

	if err == nil || err.Error() != "email: non zero value required" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	//

	testUser.Email = "not_email"

	isValid, err = testUser.isValid()

	if isValid {
		t.Errorf("User.Email should be invalid: %v", testUser.Email)
	}

	if err == nil || err.Error() != "email: not_email does not validate as email" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	//

	testUser.Email = longStr

	isValid, err = testUser.isValid()

	if isValid {
		t.Errorf("User.Email should be invalid: %v", testUser.Email)
	}

	if err == nil || err.Error() != "email: IdViTf7vd9NezkhR8Ftvh4nTVe8re2RGvFGYkMN9alkUQxm7ZrEuHVVpr6CUgK5pdLTh8H4KVVTscwMlQL9ZJF0kcKJuyMizgKkorbNeblelBECnE2G6hxSecsL69dh9eXQGUSAbbdgB6BE3Q11Ffm4GjRbz4Z4cW6D2ZyZ4RIpFtHFlQcsiD2o8QwCitr7LAIRAwSW2DHenEUYh1OVDfpUFMMUtuvnRCYghwNj8iFwrdp3rxKeTBsZU5CdbVC3@gmail.com does not validate as length(1|254)" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	testUser.Email = email
}

//
func TestUserBadCreatedAt(t *testing.T) {
	createdAt := testUser.CreatedAt
	testUser.CreatedAt = "bad_timestamp"

	isValid, err := testUser.isValid()

	if isValid {
		t.Errorf("User.CreatedAt should be invalid: %v", testUser.CreatedAt)
	}

	if err == nil || err.Error() != "created_at: bad_timestamp does not validate as rfc3339" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	testUser.CreatedAt = createdAt
}

//
func TestUserOptionalCreatedAt(t *testing.T) {
	createdAt := testUser.CreatedAt
	testUser.CreatedAt = ""

	isValid, err := testUser.isValid()

	if !isValid {
		t.Errorf("User.CreatedAt should be valid (optional): %v %v", testUser.CreatedAt, err)
	}

	testUser.CreatedAt = createdAt
}

//
func TestUserBadUpdatedAt(t *testing.T) {
	updatedAt := testUser.UpdatedAt
	testUser.UpdatedAt = "bad_timestamp"

	isValid, err := testUser.isValid()

	if isValid {
		t.Errorf("User.UpdatedAt should be invalid: %v", testUser.UpdatedAt)
	}

	if err == nil || err.Error() != "updated_at: bad_timestamp does not validate as rfc3339" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	testUser.UpdatedAt = updatedAt
}

//
func TestUserOptionalUpdatedAt(t *testing.T) {
	updatedAt := testUser.UpdatedAt
	testUser.UpdatedAt = ""

	isValid, err := testUser.isValid()

	if !isValid {
		t.Errorf("User.UpdatedAt should be valid (optional): %v %v", testUser.UpdatedAt, err)
	}

	testUser.UpdatedAt = updatedAt
}

//
func TestUserBadDeletedAt(t *testing.T) {
	deletedAt := testUser.DeletedAt
	testUser.DeletedAt = sql.NullString{"bad_timestamp", true}

	isValid, err := testUser.isValid()

	if isValid {
		t.Errorf("User.DeletedAt should be invalid: %v", testUser.DeletedAt)
	}

	if err == nil || err.Error() != "deleted_at: is not a valid timestamp" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	testUser.DeletedAt = deletedAt
}

//
func TestUserOptionalDeletedAt(t *testing.T) {
	deletedAt := testUser.DeletedAt
	testUser.DeletedAt = sql.NullString{"", true}

	isValid, err := testUser.isValid()

	if !isValid {
		t.Errorf("User.DeletedAt should be valid (optional): %v %v", testUser.DeletedAt, err)
	}

	//
	testUser.DeletedAt = sql.NullString{testUser.CreatedAt, true}

	isValid, err = testUser.isValid()

	if !isValid {
		t.Errorf("User.DeletedAt should be valid: %v %v", testUser.DeletedAt, err)
	}

	testUser.DeletedAt = deletedAt
}

//
func TestUserUpdate(t *testing.T) {
	var err error
	_, err = testUser.Ctx.Begin()

	if err != nil {
		t.Errorf("Failed to start transaction %v", err)
	}

	testUser.GivenName = "Bob"

	err = testUser.Save()

	if err != nil {
		t.Errorf("Failed to update %v", err)
	}

	//

	user, err := FetchUserByEmail(testUser.Ctx, testUser.Email)

	if err != nil {
		t.Errorf("Failed to fetch user by email: %v", err)
	}

	if user == nil {
		t.Errorf("Failed to fetch user %v: ", testUser.Email)
	}

	if user.GivenName != "Bob" {
		t.Errorf("Given name not updated: Expected %v Got: %v", "Bob", user.GivenName)
	}

	err = testUser.Ctx.Rollback()

	if err != nil {
		t.Errorf("Failed to rollback transaction %v", err)
	}
}

//
func TestUserInsert(t *testing.T) {
	var err error
	_, err = testUser.Ctx.Begin()

	if err != nil {
		t.Errorf("Failed to start transaction %v", err)
	}

	var user *User
	user = &User{}
	user.Ctx = testUser.Ctx

	user.AccountId = testUser.AccountId
	user.RoleId = 1
	user.GivenName = "Test"
	user.FamilyName = "User"
	user.Email = "testuser@uxt.com"

	err = user.Save()

	if err != nil {
		t.Errorf("Failed to insert %v", err)
	}

	user, err = FetchUserByEmail(testUser.Ctx, user.Email)

	if err != nil {
		t.Errorf("Failed to fetch user by email: %v", err)
	}

	if user == nil {
		t.Errorf("Failed to fetch user %v: ", testUser.Email)
	}

	if user.GivenName != "Test" {
		t.Errorf("Given name not inserted: Expected %v Got: %v", "Test", user.GivenName)
	}

	err = testUser.Ctx.Rollback()

	if err != nil {
		t.Errorf("Failed to rollback transaction %v", err)
	}
}
