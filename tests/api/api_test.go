package api_tests

// Testing Refresh Token
// this method uses the TokenTest
func (s *TestSuite) TestLoginAndRefreshToken(c *C) {

	// Creating a new user
	login := &models.Login{Email: "TestLoginAndRefreshToken@email.com", Password: "123456"}
	postUser(login, c)

	// Trying to log in with this new user
	token := getToken(login, c)
	// Getting the brand new user from the session
	getMe(token, c)

	// Test Refresh Token!
	refreshToken(token, c)

}

// Testing User Creation
// this method also try to create
// another user using the same email
// to force the erro return
// this method also update users email and password
// and try to login with the new password
// it also checks if user has really changed in DB
func (s *TestSuite) TestCreateAndChangeUser(c *C) {

	// Creating a new user
	login := &models.Login{Email: "TestCreateAndChangeUser@email.com", Password: "123456"}
	postUser(login, c)

	// If we create a new user with same email
	// then server must return conflict status
	postUserMustConflict(login, c)

	//
	// Updating this new user email and password
	//

	token := getToken(login, c)
	user := getMe(token, c)

	// Attributes to be changed
	newEmail := "putuserchanged@email.com"
	newPassword := "123456changed"

	// Changing user attributes
	user.Email = newEmail
	user.Password = newPassword

	userModified := putUser(user, token, c)

	// Since Password isn't returned,
	// we will check all other fields
	user.Password = "" // Cleanning email to check all others
	c.Assert(userModified, DeepEquals, user)

	// Trying to login with the new Email and Password
	token = getToken(&models.Login{Email: newEmail, Password: newPassword}, c)

	// Testing if this token really works and our user has changed
	me := getMe(token, c)
	// Test if our changes has saved
	c.Assert(userModified, DeepEquals, me)
}
