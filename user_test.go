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
