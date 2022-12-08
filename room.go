package main

import "net"

type room struct {
	name    string
	members map[net.Addr]*user
}

func (r *room) broadcast(sender *user, msg string) {
	for addr, m := range r.members {
		if addr != sender.conn.RemoteAddr() {
			m.msg(msg)
		}
	}
}
