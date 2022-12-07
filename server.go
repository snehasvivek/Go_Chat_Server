package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client)
		}
	}
}
func (s *server) newClient(conn net.Conn) *client {
	log.Printf("new client has connected: %s", conn.RemoteAddr().String())

	return &client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
	}

}

func (s *server) nick(c *client, args []string) {
	if len(args) < 2 {
		c.msg("nick is required. usage: /nick NAME")
		return
	}

	c.nick = args[1]
	c.msg(fmt.Sprintf("all right, I will call you %s", c.nick))
}

func (s *server) join(c *client, args []string) {
	if len(args) < 2 {
		c.msg("room name is required. usage: /join ROOM_NAME")
		return
	}

	// get desired room name from client
	roomName := args[1]

	// check if the room exists, creates new room if the room does not exist
	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}

	// adds client to new room
	r.members[c.conn.RemoteAddr()] = c

	// quits current room
	s.quitCurrentRoom(c)

	// sets the new room within the client structure
	c.room = r

	// broadcast message to the rest of the room
	r.broadcast(c, fmt.Sprintf("%s has joined the room", c.nick))

	// broadcast message to the client
	c.msg(fmt.Sprint("welcome to %s", roomName))
}

func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}
	c.msg(fmt.Sprintf("available rooms are %s", strings.Join(rooms, ", ")))

}

func (s *server) msg(c *client, args []string) {
	if len(args) < 2 {
		c.msg("message is required, usage: /msg MSG")
		return
	}
	msg := strings.Join(args[1:], " ")
	c.room.broadcast(c, c.nick+": "+msg)

}

func (s *server) quit(c *client) {
	log.Printf("client has disconnected: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.msg("sad to see you go :(")
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}
