package tui

import "github.com/ivelsantos/cryptor/models"

func (m *model) deleteUser(user string) {
	err := models.DeleteUser(user)
	if err != nil {
		panic(err)
	}
}

func (m *model) updateUsers() {
	users := make([]string, 0, 5)
	accounts, err := models.GetAccounts()
	if err != nil {
		panic(err)
	}

	for _, account := range accounts {
		users = append(users, account.Name)
	}

	m.users = users
}
