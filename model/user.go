package model

// User is a retrieved and authentiacted user.
type User struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

var UsersCache = make([]User, 0)

func UserIsCached(user User) bool {
	for _, v := range UsersCache {
		if v.Email == user.Email && v.EmailVerified {
			return true
		}
	}
	return false
}

func AddUser(user User) {
	if !UserIsCached(user) {
		UsersCache = append(UsersCache, user)
	}
}
