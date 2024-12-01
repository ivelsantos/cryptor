package tui

import "fmt"

func home(user string) (string, error) {
	page := fmt.Sprintf("\nHello %s...\n", user)
	return page, nil
}
