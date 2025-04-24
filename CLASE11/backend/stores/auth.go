package stores

type AuthStore struct {
	IsLoggedIn  bool
	Username    string
	Password    string
	PartitionID string
}

var Auth = &AuthStore{
	IsLoggedIn:  false,
	Username:    "",
	Password:    "",
	PartitionID: "",
}

func (a *AuthStore) Login(username, password, partitionID string) {
	a.IsLoggedIn = true
	a.Username = username
	a.Password = password
	a.PartitionID = partitionID
}

func (a *AuthStore) Logout() {
	a.IsLoggedIn = false
	a.Username = ""
	a.Password = ""
	a.PartitionID = ""
}

func (a *AuthStore) IsAuthenticated() bool {
	return a.IsLoggedIn
}

func (a *AuthStore) GetCurrentUser() (string, string, string) {
	return a.Username, a.Password, a.PartitionID
}

func (a *AuthStore) GetPartitionID() string {
	return a.PartitionID
}
