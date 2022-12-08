package main

type commandID int

const (
	NAME commandID = iota
	JOIN
	ROOMS
	MSG
	QUIT
)

type command struct {
	id   commandID
	user *user
	args []string
}
