package remote

import (
	"RISK-CodeConflict/core"
	"bufio"
	"fmt"
	"image/color"
	"log"
	"net"
	"net/textproto"
	"strconv"
	"strings"
)

// RunServer initializes and starts a TCP server that listens for incoming connections from clients.
// The server manages commands received from clients and executes them on the given World object.
// It remains BLOCKING until stopped manually.
//
// Parameters:
//   - host: The IP address or hostname on which the server should run (e.g., "0.0.0.0").
//   - port: The port on which the server should listen for connections (e.g., "1234").
//   - world: The World object representing the game state, shared between all connected clients.
//   - playerCount: The number of players required before the game starts (initializes population and unfreezes the world).
func RunServer(host, port string, world *core.World, maxPlayerCount int) {
	// Freeze the world state at the start to prevent any modifications before the game starts.
	world.Freeze = true

	// Set up the server to listen for incoming connections on the specified host and port.
	l, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Fatalf("Failed to start the server: %v\n", err)
	}

	// Ensure that the listener is closed when the server terminates.
	defer func(l net.Listener) {
		_ = l.Close()
	}(l)

	// Print the server start message to the console.
	fmt.Printf("Server started on [%s:%s]\n", host, port)

	// Track the number of connected players.
	count := 0
	for {
		// Wait for an incoming connection from a client.
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		// Handle each new connection in a separate goroutine.
		go handleRequest(conn, world, maxPlayerCount)
		count++

		// Print information about the connected player.
		fmt.Printf("Player connected from %v\n", conn.RemoteAddr())
	}
}

// handleRequest handles communication with a single client connection.
// It processes commands sent by the client and updates the shared World object accordingly.
//
// Parameters:
//   - conn: The network connection object representing the client connection.
//   - w: The World object representing the game state.
func handleRequest(conn net.Conn, w *core.World, maxPlayerCount int) {
	// Store the name of the player associated with this connection.
	var player string

	// Create a buffered reader to read client input line by line.
	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)

	// Ensure the connection is closed when the function exits.
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	// Continuously listen for commands from the client.
	for {
		// Read a line of input from the client.
		line, err := tp.ReadLine()
		if err != nil {
			break // Exit loop if an error occurs (e.g., client disconnect).
		}

		// Trim any leading/trailing whitespace and split the command into arguments.
		args := strings.Split(strings.TrimSpace(line), "|")

		// Extract the command keyword (the first element).
		var com string
		if len(args) > 0 {
			com = args[0]
		}

		// Handle different commands based on the command keyword.
		switch com {
		case "PLAYER":
			// Create or validate a player for the connection.
			if len(player) > 0 {
				// If the player is already set, send an error response.
				comResponse(conn, "err: player already created")
			} else {
				// Extract player information from command arguments and create the player.
				name, r, g, b := saveArgs(args)
				ri, _ := strconv.Atoi(r)
				gi, _ := strconv.Atoi(g)
				bi, _ := strconv.Atoi(b)
				col := color.RGBA{R: uint8(ri), G: uint8(gi), B: uint8(bi), A: 255}

				// Try adding the player to the world.
				e := w.AddPlayer(name, col)
				if e == nil {
					player = name // Set player name for this connection if successful.
					println("add player", name)
				}
				comResponseErr(conn, e)

				// Check if the number of players matches the required count.
				// If yes, initialize the world population and unfreeze the world to allow actions.
				if len(w.PlayerQueue) == maxPlayerCount {
					println("last player added")
					w.InitPopulation()
					w.Freeze = false
				}
			}
		case "STATUS":
			// Send the current world state as a JSON string.
			comResponse(conn, w.Json())
		case "END":
			// Handle the end of the turn for the player.
			comResponseErr(conn, w.EndTurn(player))
		case "MOVE":
			// Handle troop movements or attacks.
			attacker, defender, strength, _ := saveArgs(args)
			strengthInt, _ := strconv.Atoi(strength)
			comResponseErr(conn, w.AttackOrMove(attacker, defender, strengthInt, player))
		default:
			// If the command is invalid, send an error response.
			comResponse(conn, "err: invalid command")
		}
	}

	// Log the player's departure when the connection is closed.
	fmt.Printf("Player %s has disconnected\n", player)
}

// comResponse is a helper function that sends a formatted response message back to the client.
//
// Parameters:
//   - conn: The network connection object representing the client connection.
//   - s: The response message to send.
func comResponse(conn net.Conn, s string) {
	_, err := conn.Write([]byte(fmt.Sprintf("%s\r\n", s)))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}

// comResponseErr is a helper function that sends an error message (if any) or "OK" back to the client.
//
// Parameters:
//   - conn: The network connection object representing the client connection.
//   - err: The error object (if any) to send as a response.
func comResponseErr(conn net.Conn, err error) {
	if err != nil {
		comResponse(conn, err.Error())
	} else {
		comResponse(conn, "OK")
	}
}

// saveArgs is a helper function that extracts and returns up to four string arguments from a client command.
//
// Parameters:
//   - args: A list of command arguments received from the client.
//
// Returns:
//   - a1: The first argument as a string.
//   - a2: The second argument as a string.
//   - a3: The third argument as a string.
//   - a4: The fourth argument as a string.
func saveArgs(args []string) (a1, a2, a3, a4 string) {
	sArgs := make([]string, 5)
	copy(sArgs, args)
	return sArgs[1], sArgs[2], sArgs[3], sArgs[4]
}
