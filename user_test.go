// +build integration

package jgoweb

import (
	"testing"
	"database/sql"
)

//
func TestFetchUserByEmail(t *testing.T) {
	InitMockUser()

	user, err := FetchUserByEmail(MockUser.Ctx, MockUser.Email)

	if err != nil {
		t.Errorf("Failed to fetch user by email: %v", err)
	}

	if user == nil {
		t.Errorf("Failed to fetch user %v: ", MockUser.Email)
	}

	if user.Email != MockUser.Email {
		t.Errorf("Fetched wrong user? Expected: %v Got: %v", MockUser.Email, user.Email)
	}

	MockUser.Id = user.Id
	MockUser.AccountId = user.AccountId
	MockUser.RoleId = user.RoleId
	MockUser.GivenName = user.GivenName
	MockUser.FamilyName = user.FamilyName
	MockUser.Email = user.Email
	MockUser.CreatedAt = user.CreatedAt
	MockUser.UpdatedAt = user.UpdatedAt
	MockUser.DeletedAt = user.DeletedAt

	userName := "not_a_user"
	user, err = FetchUserByEmail(MockUser.Ctx, userName)

	if err != nil {
		t.Errorf("Failed to fetch user by email: %v", err)
	}

	if user != nil {
		t.Errorf("Should have failed to find user: %v", userName)
	}
}

//
func TestFetchUserById(t *testing.T) {
	user, err := FetchUserById(MockUser.Ctx, MockUser.Id)

	if err != nil {
		t.Errorf("Failed to fetch user by Id: %v", err)
	}

	if user == nil {
		t.Errorf("Failed to fetch user %v: ", MockUser.Email)
	}

	if user.Email != MockUser.Email {
		t.Errorf("Fetched wrong user? Expected: %v Got: %v", MockUser.Email, user.Email)
	}

	// force error
	userId := "not_a_user"
	_, err = FetchUserById(MockUser.Ctx, userId)

	if err == nil {
		t.Errorf("FetchUserById should have failed with uuid of: %v", userId)
	}

	// force not found
	userId = "00000000-0000-0000-0000-000000000000"
	user, err = FetchUserById(MockUser.Ctx, userId)

	if err != nil {
		t.Errorf("Failed to fetch user by id: %v", err)
	}

	if user != nil {
		t.Errorf("Should have failed to find user: %v", userId)
	}
}

//
func MockUserIsValid(t *testing.T) {
	isValid, err := MockUser.isValid()

	if !isValid {
		t.Errorf("User should be valid: %v", err)
	}
}

//
func MockUserBadId(t *testing.T) {
	id := MockUser.Id
	MockUser.Id = "bad_uuid"

	isValid, err := MockUser.isValid()

	if isValid {
		t.Errorf("User.Id should be invalid: %v", MockUser.Id)
	}

	if err == nil || err.Error() != "id: bad_uuid does not validate as uuid" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	MockUser.Id = id
}

//
func MockUserOptionalId(t *testing.T) {
	id := MockUser.Id
	MockUser.Id = ""

	isValid, err := MockUser.isValid()

	if !isValid {
		t.Errorf("User.Id should be valid (optional): %v %v", MockUser.Id, err)
	}

	MockUser.Id = id
}

//
func MockUserBadAccountId(t *testing.T) {
	id := MockUser.AccountId
	MockUser.AccountId = "bad_uuid"

	isValid, err := MockUser.isValid()

	if isValid {
		t.Errorf("User.AccountId should be invalid: %v", MockUser.AccountId)
	}

	if err == nil || err.Error() != "account_id: bad_uuid does not validate as uuid" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	MockUser.AccountId = id
}

//
func MockUserBadRoleId(t *testing.T) {
	id := MockUser.RoleId
	MockUser.RoleId = 0

	isValid, err := MockUser.isValid()

	if isValid {
		t.Errorf("User.RoleId should be invalid: %v", MockUser.RoleId)
	}

	if err == nil || err.Error() != "role_id: non zero value required" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	//
	MockUser.RoleId = 101

	isValid, err = MockUser.isValid()

	if isValid {
		t.Errorf("User.RoleId should be invalid: %v", MockUser.RoleId)
	}

	if err == nil || err.Error() != "role_id: 101 does not validate as range(1|100)" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	MockUser.RoleId = id
}

//
func MockUserBadGivenName(t *testing.T) {
	longStr := "IdViTf7vd9NezkhR8Ftvh4nTVe8re2RGvFGYkMN9alkUQxm7ZrEuHVVpr6CUgK5pdLTh8H4KVVTscwMlQL9ZJF0kcKJuyMizgKkorbNeblelBECnE2G6hxSecsL69dh9eXQGUSAbbdgB6BE3Q11Ffm4GjRbz4Z4cW6D2ZyZ4RIpFtHFlQcsiD2o8QwCitr7LAIRAwSW2DHenEUYh1OVDfpUFMMUtuvnRCYghwNj8iFwrdp3rxKeTBsZU5CdbVC3"
	name := MockUser.GivenName
	MockUser.GivenName = ""

	isValid, err := MockUser.isValid()

	if isValid {
		t.Errorf("User.GivenName should be invalid: %v", MockUser.GivenName)
	}

	if err == nil || err.Error() != "given_name: non zero value required" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	//

	MockUser.GivenName = longStr

	isValid, err = MockUser.isValid()

	if isValid {
		t.Errorf("User.GivenName should be invalid: %v", MockUser.GivenName)
	}

	if err == nil || err.Error() != "given_name: IdViTf7vd9NezkhR8Ftvh4nTVe8re2RGvFGYkMN9alkUQxm7ZrEuHVVpr6CUgK5pdLTh8H4KVVTscwMlQL9ZJF0kcKJuyMizgKkorbNeblelBECnE2G6hxSecsL69dh9eXQGUSAbbdgB6BE3Q11Ffm4GjRbz4Z4cW6D2ZyZ4RIpFtHFlQcsiD2o8QwCitr7LAIRAwSW2DHenEUYh1OVDfpUFMMUtuvnRCYghwNj8iFwrdp3rxKeTBsZU5CdbVC3 does not validate as length(1|254)" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	MockUser.GivenName = name
}

//
func MockUserBadFamilyName(t *testing.T) {
	longStr := "IdViTf7vd9NezkhR8Ftvh4nTVe8re2RGvFGYkMN9alkUQxm7ZrEuHVVpr6CUgK5pdLTh8H4KVVTscwMlQL9ZJF0kcKJuyMizgKkorbNeblelBECnE2G6hxSecsL69dh9eXQGUSAbbdgB6BE3Q11Ffm4GjRbz4Z4cW6D2ZyZ4RIpFtHFlQcsiD2o8QwCitr7LAIRAwSW2DHenEUYh1OVDfpUFMMUtuvnRCYghwNj8iFwrdp3rxKeTBsZU5CdbVC3"
	name := MockUser.FamilyName
	MockUser.FamilyName = ""

	isValid, err := MockUser.isValid()

	if isValid {
		t.Errorf("User.FamilyName should be invalid: %v", MockUser.FamilyName)
	}

	if err == nil || err.Error() != "family_name: non zero value required" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	//

	MockUser.FamilyName = longStr

	isValid, err = MockUser.isValid()

	if isValid {
		t.Errorf("User.FamilyName should be invalid: %v", MockUser.FamilyName)
	}

	if err == nil || err.Error() != "family_name: IdViTf7vd9NezkhR8Ftvh4nTVe8re2RGvFGYkMN9alkUQxm7ZrEuHVVpr6CUgK5pdLTh8H4KVVTscwMlQL9ZJF0kcKJuyMizgKkorbNeblelBECnE2G6hxSecsL69dh9eXQGUSAbbdgB6BE3Q11Ffm4GjRbz4Z4cW6D2ZyZ4RIpFtHFlQcsiD2o8QwCitr7LAIRAwSW2DHenEUYh1OVDfpUFMMUtuvnRCYghwNj8iFwrdp3rxKeTBsZU5CdbVC3 does not validate as length(1|254)" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	MockUser.FamilyName = name
}

//
func MockUserBadEmail(t *testing.T) {
	longStr := "IdViTf7vd9NezkhR8Ftvh4nTVe8re2RGvFGYkMN9alkUQxm7ZrEuHVVpr6CUgK5pdLTh8H4KVVTscwMlQL9ZJF0kcKJuyMizgKkorbNeblelBECnE2G6hxSecsL69dh9eXQGUSAbbdgB6BE3Q11Ffm4GjRbz4Z4cW6D2ZyZ4RIpFtHFlQcsiD2o8QwCitr7LAIRAwSW2DHenEUYh1OVDfpUFMMUtuvnRCYghwNj8iFwrdp3rxKeTBsZU5CdbVC3@gmail.com"
	email := MockUser.Email
	MockUser.Email = ""

	isValid, err := MockUser.isValid()

	if isValid {
		t.Errorf("User.Email should be invalid: %v", MockUser.Email)
	}

	if err == nil || err.Error() != "email: non zero value required" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	//

	MockUser.Email = "not_email"

	isValid, err = MockUser.isValid()

	if isValid {
		t.Errorf("User.Email should be invalid: %v", MockUser.Email)
	}

	if err == nil || err.Error() != "email: not_email does not validate as email" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	//

	MockUser.Email = longStr

	isValid, err = MockUser.isValid()

	if isValid {
		t.Errorf("User.Email should be invalid: %v", MockUser.Email)
	}

	if err == nil || err.Error() != "email: IdViTf7vd9NezkhR8Ftvh4nTVe8re2RGvFGYkMN9alkUQxm7ZrEuHVVpr6CUgK5pdLTh8H4KVVTscwMlQL9ZJF0kcKJuyMizgKkorbNeblelBECnE2G6hxSecsL69dh9eXQGUSAbbdgB6BE3Q11Ffm4GjRbz4Z4cW6D2ZyZ4RIpFtHFlQcsiD2o8QwCitr7LAIRAwSW2DHenEUYh1OVDfpUFMMUtuvnRCYghwNj8iFwrdp3rxKeTBsZU5CdbVC3@gmail.com does not validate as length(1|254)" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	MockUser.Email = email
}

//
func MockUserBadCreatedAt(t *testing.T) {
	createdAt := MockUser.CreatedAt
	MockUser.CreatedAt = "bad_timestamp"

	isValid, err := MockUser.isValid()

	if isValid {
		t.Errorf("User.CreatedAt should be invalid: %v", MockUser.CreatedAt)
	}

	if err == nil || err.Error() != "created_at: bad_timestamp does not validate as rfc3339" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	MockUser.CreatedAt = createdAt
}

//
func MockUserOptionalCreatedAt(t *testing.T) {
	createdAt := MockUser.CreatedAt
	MockUser.CreatedAt = ""

	isValid, err := MockUser.isValid()

	if !isValid {
		t.Errorf("User.CreatedAt should be valid (optional): %v %v", MockUser.CreatedAt, err)
	}

	MockUser.CreatedAt = createdAt
}

//
func MockUserBadUpdatedAt(t *testing.T) {
	updatedAt := MockUser.UpdatedAt
	MockUser.UpdatedAt = "bad_timestamp"

	isValid, err := MockUser.isValid()

	if isValid {
		t.Errorf("User.UpdatedAt should be invalid: %v", MockUser.UpdatedAt)
	}

	if err == nil || err.Error() != "updated_at: bad_timestamp does not validate as rfc3339" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	MockUser.UpdatedAt = updatedAt
}

//
func MockUserOptionalUpdatedAt(t *testing.T) {
	updatedAt := MockUser.UpdatedAt
	MockUser.UpdatedAt = ""

	isValid, err := MockUser.isValid()

	if !isValid {
		t.Errorf("User.UpdatedAt should be valid (optional): %v %v", MockUser.UpdatedAt, err)
	}

	MockUser.UpdatedAt = updatedAt
}

//
func MockUserBadDeletedAt(t *testing.T) {
	deletedAt := MockUser.DeletedAt
	MockUser.DeletedAt = sql.NullString{"bad_timestamp", true}

	isValid, err := MockUser.isValid()

	if isValid {
		t.Errorf("User.DeletedAt should be invalid: %v", MockUser.DeletedAt)
	}

	if err == nil || err.Error() != "deleted_at: is not a valid timestamp" {
		t.Errorf("Wrong error returned by validator: %v", err)
	}

	MockUser.DeletedAt = deletedAt
}

//
func MockUserOptionalDeletedAt(t *testing.T) {
	deletedAt := MockUser.DeletedAt
	MockUser.DeletedAt = sql.NullString{"", true}

	isValid, err := MockUser.isValid()

	if !isValid {
		t.Errorf("User.DeletedAt should be valid (optional): %v %v", MockUser.DeletedAt, err)
	}

	//
	MockUser.DeletedAt = sql.NullString{MockUser.CreatedAt, true}

	isValid, err = MockUser.isValid()

	if !isValid {
		t.Errorf("User.DeletedAt should be valid: %v %v", MockUser.DeletedAt, err)
	}

	MockUser.DeletedAt = deletedAt
}

//
func MockUserUpdate(t *testing.T) {
	var err error
	_, err = MockUser.Ctx.Begin()

	if err != nil {
		t.Errorf("Failed to start transaction %v", err)
	}

	MockUser.GivenName = "Bob"

	err = MockUser.Save()

	if err != nil {
		t.Errorf("Failed to update %v", err)
	}

	//

	user, err := FetchUserByEmail(MockUser.Ctx, MockUser.Email)

	if err != nil {
		t.Errorf("Failed to fetch user by email: %v", err)
	}

	if user == nil {
		t.Errorf("Failed to fetch user %v: ", MockUser.Email)
	}

	if user.GivenName != "Bob" {
		t.Errorf("Given name not updated: Expected %v Got: %v", "Bob", user.GivenName)
	}

	err = MockUser.Ctx.Rollback()

	if err != nil {
		t.Errorf("Failed to rollback transaction %v", err)
	}
}

//
func MockUserInsert(t *testing.T) {
	var err error
	_, err = MockUser.Ctx.Begin()

	if err != nil {
		t.Errorf("Failed to start transaction %v", err)
	}

	var user *User
	user = &User{}
	user.Ctx = MockUser.Ctx

	user.AccountId = MockUser.AccountId
	user.RoleId = 1
	user.GivenName = "Test"
	user.FamilyName = "User"
	user.Email = "MockUser@uxt.com"

	err = user.Save()

	if err != nil {
		t.Errorf("Failed to insert %v", err)
	}

	user, err = FetchUserByEmail(MockUser.Ctx, user.Email)

	if err != nil {
		t.Errorf("Failed to fetch user by email: %v", err)
	}

	if user == nil {
		t.Errorf("Failed to fetch user %v: ", MockUser.Email)
	}

	if user.GivenName != "Test" {
		t.Errorf("Given name not inserted: Expected %v Got: %v", "Test", user.GivenName)
	}

	err = MockUser.Ctx.Rollback()

	if err != nil {
		t.Errorf("Failed to rollback transaction %v", err)
	}
}
