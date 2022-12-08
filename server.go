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
		case NAME:
			s.name(cmd.user, cmd.args)
		case JOIN:
			s.join(cmd.user, cmd.args)
		case ROOMS:
			s.listRooms(cmd.user)
		case MSG:
			s.msg(cmd.user, cmd.args)
		case QUIT:
			s.quit(cmd.user)
		}
	}
}
func (s *server) newUser(conn net.Conn) *user {
	log.Printf("User Connected: %s", conn.RemoteAddr().String())

	return &user{
		conn:     conn,
		name:     "anonymous",
		commands: s.commands,
	}

}

func (s *server) name(u *user, args []string) {
	if len(args) < 2 {
		u.msg("Name Required. Format: *name NAME")
		return
	}

	u.name = args[1]
	u.msg(fmt.Sprintf("Name Set - %s", u.name))
}

func (s *server) join(u *user, args []string) {
	if len(args) < 2 {
		u.msg("Room Name Required. Format: *join ROOM_NAME")
		return
	}

	// get desired room name from client
	roomName := args[1]

	// check if the room exists, creates new room if the room does not exist
	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*user),
		}
		s.rooms[roomName] = r
	}

	// adds client to new room
	r.members[u.conn.RemoteAddr()] = u

	// quits current room
	s.quitCurrentRoom(u)

	// sets the new room within the client structure
	u.room = r

	// broadcast message to the rest of the room
	r.broadcast(u, fmt.Sprintf("%s has entered the room", u.name))

	// broadcast message to the client
	u.msg(fmt.Sprint("Welcome to ", roomName))
}

func (s *server) listRooms(u *user) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}
	u.msg(fmt.Sprintf("Current Rooms: %s", strings.Join(rooms, ", ")))

}

func (s *server) msg(u *user, args []string) {
	if len(args) < 2 {
		u.msg("Message is required. Format: *msg MSG")
		return
	}
	msg := strings.Join(args[1:], " ")
	u.room.broadcast(u, u.name+": "+msg)

}

func (s *server) quit(u *user) {
	log.Printf("User Disconncted: %s", u.conn.RemoteAddr().String())

	s.quitCurrentRoom(u)

	u.msg("Goodbye")
	u.conn.Close()
}

func (s *server) quitCurrentRoom(u *user) {
	if u.room != nil {
		oldRoom := s.rooms[u.room.name]
		delete(s.rooms[u.room.name].members, u.conn.RemoteAddr())
		oldRoom.broadcast(u, fmt.Sprintf("%s has left room", u.name))
	}
}
