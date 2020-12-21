package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	nick     string
	conn     net.Conn
	room     *room
	commands chan<- command
}

func (c *client) readMessage() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")
		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/nick":
			c.commands <- command{
				id:     CMD_NICK,
				client: c,
				args:   args,
			}
		case "/join":
			c.commands <- command{
				id:     CMD_JOIN,
				client: c,
				args:   args,
			}
		case "/rooms":
			c.commands <- command{
				id:     CMD_ROOMS,
				client: c,
			}
		case "/msg":
			c.commands <- command{
				id:     CMD_MSG,
				client: c,
				args:   args,
			}
		case "/quit":
			c.commands <- command{
				id:     CMD_QUIT,
				client: c,
			}
		default:
			c.error(fmt.Errorf("Unknow command: %s", cmd))
		}
	}
}

func (c *client) error(err error) {
	c.conn.Write([]byte("err: " + err.Error() + "\n"))
}

func (c *client) sendMessage(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}
