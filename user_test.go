// +build integration

package jgoweb

import (
	"testing"
)

//
func TestNewUser(t *testing.T) {
	user := &User{}

	InitMockCtx()
	user, err := NewUser(MockCtx)

	if err != nil {
		t.Errorf("ERROR: %v", err)
	}

	user.SetEmail("test_valid_email@sure_not_a_valid_domain.com")
	user.SetFirstName("test")
	user.SetLastName("user1")
	user.SetAccountId(MockUser.GetAccountId())
	user.SetPassword("test", "test")
	user.SetRoleId("1")

	err = user.Save()

	if err != nil {
		t.Errorf("ERROR: %v", err)
	}
}

//
func TestFetchUserByEmail(t *testing.T) {
	InitMockUser()

	user, err := FetchUserByEmail(MockUser.Ctx, MockUser.GetEmail())

	if err != nil {
		t.Errorf("Failed to fetch user by email: %v", err)
	}

	if user == nil {
		t.Errorf("Failed to fetch user %v: ", MockUser.GetEmail())
	}

	if user.Email != MockUser.Email {
		t.Errorf("Fetched wrong user? Expected: %v Got: %v", MockUser.Email, user.Email)
	}

	MockUser.Id = user.Id
	MockUser.AccountId = user.AccountId
	MockUser.RoleId = user.RoleId
	MockUser.FirstName = user.FirstName
	MockUser.LastName = user.LastName
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
	var userId string
	user, err := FetchUserById(MockUser.Ctx, MockUser.GetId())

	if err != nil {
		t.Errorf("Failed to fetch user by Id: %v", err)
	}

	if user == nil {
		t.Errorf("Failed to fetch user %v: ", MockUser.Email)
	}

	if user.Email != MockUser.Email {
		t.Errorf("Fetched wrong user? Expected: %v Got: %v", MockUser.Email, user.Email)
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
func TestUserId(t *testing.T) {
	InitMockUser()
	origVal := MockUser.GetId()
	testVal := "test"

	MockUser.SetId("")

	if MockUser.Id.Valid {
		t.Errorf("ERROR: Id should be invalid.\n")
	}

	if MockUser.GetId() != "" {
		t.Errorf("ERROR: Set Id failed. Should have a blank value. Got: %s", MockUser.GetId())
	}

	MockUser.SetId(testVal)

	if !MockUser.Id.Valid {
		t.Errorf("ERROR: Id should be valid.\n")
	}

	if MockUser.GetId() != testVal {
		t.Errorf("ERROR: Set Id failed. Expected: %s, Got: %s", testVal, MockUser.GetId())
	}

	MockUser.SetId(origVal)
}

//
func TestUserAccountId(t *testing.T) {
	InitMockUser()
	origVal := MockUser.GetAccountId()
	testVal := "test"

	MockUser.SetAccountId("")

	if MockUser.AccountId.Valid {
		t.Errorf("ERROR: AccountId should be invalid.\n")
	}

	if MockUser.GetAccountId() != "" {
		t.Errorf("ERROR: Set AccountId failed. Should have a blank value. Got: %s", MockUser.GetAccountId())
	}

	MockUser.SetAccountId(testVal)

	if !MockUser.AccountId.Valid {
		t.Errorf("ERROR: AccountId should be valid.\n")
	}

	if MockUser.GetAccountId() != testVal {
		t.Errorf("ERROR: Set AccountId failed. Expected: %s, Got: %s", testVal, MockUser.GetAccountId())
	}

	MockUser.SetAccountId(origVal)
}

//
func TestUserRoleId(t *testing.T) {
	InitMockUser()
	origVal := MockUser.GetRoleId()
	testVal := "test"

	MockUser.SetRoleId("")

	if MockUser.RoleId.Valid {
		t.Errorf("ERROR: RoleId should be invalid.\n")
	}

	if MockUser.GetRoleId() != "" {
		t.Errorf("ERROR: Set RoleId failed. Should have a blank value. Got: %s", MockUser.GetRoleId())
	}

	MockUser.SetRoleId(testVal)

	if !MockUser.RoleId.Valid {
		t.Errorf("ERROR: RoleId should be valid.\n")
	}

	if MockUser.GetRoleId() != testVal {
		t.Errorf("ERROR: Set RoleId failed. Expected: %s, Got: %s", testVal, MockUser.GetRoleId())
	}

	MockUser.SetRoleId(origVal)
}

//
func TestUserFirstName(t *testing.T) {
	InitMockUser()
	origVal := MockUser.GetFirstName()
	testVal := "test"

	MockUser.SetFirstName("")

	if MockUser.FirstName.Valid {
		t.Errorf("ERROR: FirstName should be invalid.\n")
	}

	if MockUser.GetFirstName() != "" {
		t.Errorf("ERROR: Set FirstName failed. Should have a blank value. Got: %s", MockUser.GetFirstName())
	}

	MockUser.SetFirstName(testVal)

	if !MockUser.FirstName.Valid {
		t.Errorf("ERROR: FirstName should be valid.\n")
	}

	if MockUser.GetFirstName() != testVal {
		t.Errorf("ERROR: Set FirstName failed. Expected: %s, Got: %s", testVal, MockUser.GetFirstName())
	}

	MockUser.SetFirstName(origVal)
}

//
func TestUserLastName(t *testing.T) {
	InitMockUser()
	origVal := MockUser.GetLastName()
	testVal := "test"

	MockUser.SetLastName("")

	if MockUser.LastName.Valid {
		t.Errorf("ERROR: LastName should be invalid.\n")
	}

	if MockUser.GetLastName() != "" {
		t.Errorf("ERROR: Set LastName failed. Should have a blank value. Got: %s", MockUser.GetLastName())
	}

	MockUser.SetLastName(testVal)

	if !MockUser.LastName.Valid {
		t.Errorf("ERROR: LastName should be valid.\n")
	}

	if MockUser.GetLastName() != testVal {
		t.Errorf("ERROR: Set LastName failed. Expected: %s, Got: %s", testVal, MockUser.GetLastName())
	}

	MockUser.SetLastName(origVal)
}

//
func TestUserEmail(t *testing.T) {
	InitMockUser()
	origVal := MockUser.GetEmail()
	testVal := "test"

	MockUser.SetEmail("")

	if MockUser.Email.Valid {
		t.Errorf("ERROR: Email should be invalid.\n")
	}

	if MockUser.GetEmail() != "" {
		t.Errorf("ERROR: Set Email failed. Should have a blank value. Got: %s", MockUser.GetEmail())
	}

	MockUser.SetEmail(testVal)

	if !MockUser.Email.Valid {
		t.Errorf("ERROR: Email should be valid.\n")
	}

	if MockUser.GetEmail() != testVal {
		t.Errorf("ERROR: Set Email failed. Expected: %s, Got: %s", testVal, MockUser.GetEmail())
	}

	MockUser.SetEmail(origVal)
}

//
func TestUserCreatedAt(t *testing.T) {
	InitMockUser()
	origVal := MockUser.GetCreatedAt()
	testVal := "test"

	MockUser.SetCreatedAt("")

	if MockUser.CreatedAt.Valid {
		t.Errorf("ERROR: CreatedAt should be invalid.\n")
	}

	if MockUser.GetCreatedAt() != "" {
		t.Errorf("ERROR: Set CreatedAt failed. Should have a blank value. Got: %s", MockUser.GetCreatedAt())
	}

	MockUser.SetCreatedAt(testVal)

	if !MockUser.CreatedAt.Valid {
		t.Errorf("ERROR: CreatedAt should be valid.\n")
	}

	if MockUser.GetCreatedAt() != testVal {
		t.Errorf("ERROR: Set CreatedAt failed. Expected: %s, Got: %s", testVal, MockUser.GetCreatedAt())
	}

	MockUser.SetCreatedAt(origVal)
}

//
func TestUserDeletedAt(t *testing.T) {
	InitMockUser()
	origVal := MockUser.GetDeletedAt()
	testVal := "test"

	MockUser.SetDeletedAt("")

	if MockUser.DeletedAt.Valid {
		t.Errorf("ERROR: DeletedAt should be invalid.\n")
	}

	if MockUser.GetDeletedAt() != "" {
		t.Errorf("ERROR: Set DeletedAt failed. Should have a blank value. Got: %s", MockUser.GetDeletedAt())
	}

	MockUser.SetDeletedAt(testVal)

	if !MockUser.DeletedAt.Valid {
		t.Errorf("ERROR: DeletedAt should be valid.\n")
	}

	if MockUser.GetDeletedAt() != testVal {
		t.Errorf("ERROR: Set DeletedAt failed. Expected: %s, Got: %s", testVal, MockUser.GetDeletedAt())
	}

	MockUser.SetDeletedAt(origVal)
}

//
func TestUserUpdatedAt(t *testing.T) {
	InitMockUser()
	origVal := MockUser.GetUpdatedAt()
	testVal := "test"

	MockUser.SetUpdatedAt("")

	if MockUser.UpdatedAt.Valid {
		t.Errorf("ERROR: UpdatedAt should be invalid.\n")
	}

	if MockUser.GetUpdatedAt() != "" {
		t.Errorf("ERROR: Set UpdatedAt failed. Should have a blank value. Got: %s", MockUser.GetUpdatedAt())
	}

	MockUser.SetUpdatedAt(testVal)

	if !MockUser.UpdatedAt.Valid {
		t.Errorf("ERROR: UpdatedAt should be valid.\n")
	}

	if MockUser.GetUpdatedAt() != testVal {
		t.Errorf("ERROR: Set UpdatedAt failed. Expected: %s, Got: %s", testVal, MockUser.GetUpdatedAt())
	}

	MockUser.SetUpdatedAt(origVal)
}

//
func TestUserPassword(t *testing.T) {
	InitMockUser()
	testVal := "test"

	MockUser.SetPassword("", "")

	if MockUser.RawPasswordError == "" {
		t.Errorf("ERROR: Password should be invalid.\n")
	}

	MockUser.SetPassword(testVal, testVal)

	if !MockUser.Password.Valid {
		t.Errorf("ERROR: Password should be valid.\n")
	}
}

// @TODO: will need to be updated when authenticate is implemented properly
func TestAuthenticate(t *testing.T) {
	result := MockUser.Authenticate("bad_password")

	if result == true {
		t.Errorf("User authentication failed. Incorrect password returned true.")
	}

	// result = MockUser.Authenticate("letmein")

	// if result == false {
	// 	t.Errorf("User authentication failed. Correct password returned false.")
	// }
}

//
func TestUserIsValid(t *testing.T) {
	err := MockUser.IsValid()

	if err != nil {
		t.Errorf("User should be valid: %v", err)
	}
}

//
// func TestUserBadId(t *testing.T) {
// 	id := MockUser.Id
// 	MockUser.SetId("bad_uuid")

// 	err := MockUser.IsValid()

// 	if err == nil {
// 		t.Errorf("Wrong error returned by validator: %v", err)
// 	}

// 	MockUser.Id = id
// }

//
// func TestUserOptionalId(t *testing.T) {
// 	id := MockUser.GetId()
// 	MockUser.SetId("")

// 	err := MockUser.IsValid()

// 	if err != nil {
// 		t.Errorf("User.Id should be valid (optional): %v %v", MockUser.Id, err)
// 	}

// 	MockUser.SetId(id)
// }

//
// func TestUserBadAccountId(t *testing.T) {
// 	id := MockUser.GetAccountId()
// 	MockUser.SetAccountId("bad_uuid")

// 	err := MockUser.IsValid()

// 	if err == nil {
// 		t.Errorf("User.AccountId should be invalid: %v", MockUser.GetAccountId())
// 	}

// 	MockUser.SetAccountId(id)
// }

//
func TestUserUpdate(t *testing.T) {
	var err error

	MockUser.SetFirstName("Bob")

	err = MockUser.Save()

	if err != nil {
		t.Errorf("Failed to update %v, email: %v, accountID: %v", err, MockUser.GetEmail(), MockUser.GetAccountId())
	}

	//
	user, err := FetchUserByEmail(MockUser.Ctx, MockUser.GetEmail())

	if err != nil {
		t.Errorf("Failed to fetch user by email: %v", err)
	}

	if user == nil {
		t.Errorf("Failed to fetch user %v: ", MockUser.GetEmail())
	}

	if user.GetFirstName() != "Bob" {
		t.Errorf("First name not updated: Expected %v Got: %v", "Bob", user.GetFirstName())
	}
}

//
func TestUserInsert(t *testing.T) {
	var err error

	var user *User
	user = &User{}
	user.Ctx = MockUser.Ctx

	user.AccountId = MockUser.AccountId
	user.SetRoleId("1")
	user.SetFirstName("Test")
	user.SetLastName("User")
	user.SetEmail("MockUser@uxt.com")
	user.SetPassword("test", "test")

	err = user.Save()

	if err != nil {
		t.Errorf("Failed to insert %v", err)
	}

	user, err = FetchUserByEmail(MockUser.Ctx, user.GetEmail())

	if err != nil {
		t.Errorf("Failed to fetch user by email: %v", err)
	}

	if user == nil {
		t.Errorf("Failed to fetch user %v: ", MockUser.GetEmail())
	}

	if user.GetFirstName() != "Test" {
		t.Errorf("First name not inserted: Expected %v Got: %v", "Test", user.GetFirstName())
	}
}
