package rcon

import (
	"fmt"

	mcnet "github.com/Tnze/go-mc/net"
	prompt "github.com/c-bata/go-prompt"
)

type rcon struct {
	server   string
	port     int
	password string
	prompt   *prompt.Prompt
}

type RCONer interface {
	RunPrompt()
}

func (r *rcon) RunPrompt() {
	fmt.Println("🔌 Connected to RCON (control-D to exit)\n🆘 Type help for list of commands")
	r.prompt.Run()
}

func (r *rcon) executor(t string) {
	conn, err := mcnet.DialRCON(fmt.Sprintf("%s:%d", r.server, r.port), r.password)
	if err != nil {
		fmt.Printf("Error dialing to rcon: %s\n", err.Error())
		return
	}
	defer func(conn mcnet.RCONClientConn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("Error closing RCON client: %s\n", err.Error())
		}
	}(conn)
	err = conn.Cmd(t)
	if err != nil {
		fmt.Printf("Error executing command: %s\n", err.Error())
	}
	resp, err := conn.Resp()
	if err != nil {
		fmt.Printf("Error retrieving rcon response: %s\n", err.Error())
		return
	}
	fmt.Println(resp)
}

func completer(_ prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{}
}

func NewRCON(server, passwort string, port int) *rcon {
	var r = &rcon{
		server:   server,
		password: passwort,
		port:     port,
	}

	p := prompt.New(
		r.executor,
		completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("minectl🗺 RCON"),
	)
	r.prompt = p
	return r
}
