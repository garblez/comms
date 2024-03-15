package user

import (
	"net"
)

type User struct {
	Address net.TCPAddr `json:"address"`
	Alias   string      `json:"alias"`
}

func NewUser(alias string, address string, port int) User {
	return User{
		Alias:   alias,
		Address: net.TCPAddr{IP: net.ParseIP(address), Port: port},
	}
}
