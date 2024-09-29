package remote

import (
	"RISK-CodeConflict/core"
	"bufio"
	"errors"
	"fmt"
	"image/color"
	"net"
	"net/textproto"
	"strings"
	"sync"
)

// Client represents a remote connection to the game server, allowing communication and interaction with the game world.
type Client struct {
	conn *net.TCPConn      // TCP connection to the game server
	tp   *textproto.Reader // Text protocol reader for the connection
	mux  *sync.Mutex       // Mutex for thread-safe operations
}

// NewClient creates a new Client instance and establishes a connection to the game server at the provided host and port.
// It initializes the TCP connection.
func NewClient(host, port string) (*Client, error) {

	// Resolve TCP address
	tcpAddr, err := net.ResolveTCPAddr("tcp", host+":"+port)
	if err != nil {
		return nil, err
	}

	// Establish TCP connection
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	// Create a new client instance
	c := &Client{
		conn: conn,
		tp:   textproto.NewReader(bufio.NewReader(conn)),
		mux:  new(sync.Mutex),
	}

	// Return the client instance
	return c, nil
}

// AddPlayer registers or identifies the player with the given name on the server.
func (c *Client) AddPlayer(name string, clr color.RGBA) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	resp := c.command(fmt.Sprintf("PLAYER|%s|%d|%d|%d", name, clr.R, clr.G, clr.B))

	if strings.HasPrefix(resp, "OK") {
		return nil // Operation successful
	} else {
		return errors.New(resp)
	}
}

// Status retrieves the current world status from the server and updates the provided World instance.
func (c *Client) Status(update *core.World) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	resp := c.command("STATUS")

	if update == nil {
		return errors.New("world is nil")
	} else {
		return update.FromJson(resp)
	}
}

// EndTurn signals the server that the player has finished their turn.
func (c *Client) EndTurn() error {
	c.mux.Lock()
	defer c.mux.Unlock()

	resp := c.command("END")

	if strings.HasPrefix(resp, "OK") {
		return nil // Operation successful
	} else {
		return errors.New(resp)
	}
}

// AttackOrMove sends a command to the server to attack or move from one country to another with a specified strength.
func (c *Client) AttackOrMove(attacker, defender string, strength int) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	resp := c.command(fmt.Sprintf("MOVE|%s|%s|%d", attacker, defender, strength))

	if strings.HasPrefix(resp, "OK") {
		return nil // Operation successful
	} else {
		return errors.New(resp)
	}
}

// Reinforcement sends a command to reinforce a country with additional strength.
func (c *Client) Reinforcement(country string, strength int) error {
	return c.AttackOrMove(country, country, strength)
}

//---------------- HELPER --------------------------------------------------------------------------------------------//

// command sends the command string to the server and returns the response.
func (c *Client) command(cmd string) string {
	if c == nil || c.conn == nil || c.tp == nil {
		return "err: TcpClient connection closed."
	}

	// Remove protocol breaks
	cmd = strings.ReplaceAll(cmd, "\n", "")
	cmd = strings.ReplaceAll(cmd, "\r", "")
	cmd = strings.ReplaceAll(cmd, "  ", " ")

	// Send command
	_, err := c.conn.Write([]byte(fmt.Sprintf("%s\r\n", cmd)))
	if err != nil {
		return fmt.Sprintf("err: TcpClient write: %v", err)
	}

	// Read response
	resp, err := c.tp.ReadLine()
	if err != nil {
		return fmt.Sprintf("err: TcpClient read: %v", err)
	}

	// Return server response
	return resp
}
