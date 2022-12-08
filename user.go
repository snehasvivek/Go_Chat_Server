package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type user struct {
	conn     net.Conn
	name     string
	room     *room
	commands chan<- command
}

func (u *user) readInput() {
	for {
		msg, err := bufio.NewReader(u.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "*name":
			u.commands <- command{
				id:   NAME,
				user: u,
				args: args,
			}
		case "*join":
			u.commands <- command{
				id:   JOIN,
				user: u,
				args: args,
			}
		case "*rooms":
			u.commands <- command{
				id:   ROOMS,
				user: u,
			}
		case "*msg":
			u.commands <- command{
				id:   MSG,
				user: u,
				args: args,
			}
		case "*quit":
			u.commands <- command{
				id:   QUIT,
				user: u,
			}
		default:
			u.err(fmt.Errorf("Unknown Command: %s", cmd))
		}
	}
}

func (u *user) err(err error) {
	u.conn.Write([]byte("ERR: " + err.Error() + "\n"))
}

func (u *user) msg(msg string) {
	u.conn.Write([]byte("~ " + msg + "\n"))
}
