package models

import (
	"errors"
	"strings"

	"flexoo-service-order/src/config"
)

type Client struct {
	ClientID int64  `json:"client_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"-"`
}

func AuthenticateClient(username, password string) (*Client, error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return nil, errors.New("invalid credentials")
	}

	if config.DB == nil {
		storeMu.RLock()
		defer storeMu.RUnlock()
		for _, c := range fallbackClients {
			if (strings.EqualFold(c.Username, username) || strings.EqualFold(c.Email, username)) && c.Password == password {
				copyClient := c
				return &copyClient, nil
			}
		}
		return nil, errors.New("invalid username or password")
	}

	var client Client
	err := config.DB.QueryRow(
		`SELECT client_id, full_name, email, COALESCE(username, ''), COALESCE(password, '')
		 FROM clients
		 WHERE (LOWER(username) = LOWER($1) OR LOWER(email) = LOWER($1)) AND password = $2
		 LIMIT 1`,
		username, password,
	).Scan(&client.ClientID, &client.FullName, &client.Email, &client.Username, &client.Password)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}
	return &client, nil
}

