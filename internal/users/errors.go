package users

type InvalidLogin struct{}

func (m *InvalidLogin) Error() string {
	return "Incorrect Credentials"
}
