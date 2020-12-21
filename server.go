package main

import (
	"errors"
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
			s.setNick(&cmd)
		case CMD_JOIN:
			s.joinRoom(&cmd)
		case CMD_ROOMS:
			s.listRooms(&cmd)
		case CMD_MSG:
			s.sendMsg(&cmd)
		case CMD_QUIT:
			s.quit(&cmd)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("new client has joined: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
	}

	c.readMessage()
}

func (s *server) setNick(cmd *command) {
	cmd.client.nick = cmd.args[1]
	cmd.client.sendMessage(fmt.Sprintf("all right, I will call you %s", cmd.client.nick))
}

func (s *server) joinRoom(cmd *command) {
	roomName := cmd.args[1]

	r, ok := s.rooms[roomName]

	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}

	r.members[cmd.client.conn.RemoteAddr()] = cmd.client
	s.quitCurrentRoomt(cmd.client)

	cmd.client.room = r

	r.broadcast(cmd.client, fmt.Sprintf("%s has joined the room", cmd.client.nick))
	cmd.client.sendMessage(fmt.Sprintf("Welcome to the room: %s", roomName))
}

func (s *server) listRooms(cmd *command) {
	var rooms []string

	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	cmd.client.sendMessage(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ", ")))
}

func (s *server) sendMsg(cmd *command) {
	if cmd.client.room == nil {
		cmd.client.error(errors.New("You must to join the room first."))
		return
	}
	cmd.client.room.broadcast(
		cmd.client,
		cmd.client.nick+": "+strings.Join(cmd.args[1:len(cmd.args)], " "))

}

func (s *server) quit(cmd *command) {
	log.Printf("Client has disconnected: %s", cmd.client.conn.RemoteAddr().String())
	s.quitCurrentRoomt(cmd.client)
	cmd.client.sendMessage("Se you soon...")
	cmd.client.conn.Close()
}

func (s *server) quitCurrentRoomt(c *client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}
