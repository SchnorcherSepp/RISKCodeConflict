package remote

import (
	"RISK-CodeConflict/core"
	"image/color"
	"testing"
	"time"
)

func TestClient_AddPlayer(t *testing.T) {
	world := core.NewWorld()

	go RunServer("127.0.0.1", "2222", world, 3)

	client, err := NewClient("127.0.0.1", "2222")
	if err != nil {
		t.Fatal(err)
	}
	//------------------------------------------

	// add user
	if err := client.AddPlayer("     ", color.RGBA{R: 255, A: 255}); err == nil || err.Error() != "player name is empty" {
		t.Fatal(err)
	}
	if err := client.AddPlayer("  user1  ", color.RGBA{R: 255, A: 255}); err != nil {
		t.Fatal(err)
	}
	if err := client.AddPlayer("1234567", color.RGBA{R: 255, A: 255}); err == nil || err.Error() != "err: player already created" {
		t.Fatal(err)
	}
}

func TestClient_AttackOrMove_EndTurn(t *testing.T) {
	world := core.NewWorld()

	go RunServer("127.0.0.1", "3333", world, 2)

	client, err := NewClient("127.0.0.1", "3333")
	if err != nil {
		t.Fatal(err)
	}
	client2, err := NewClient("127.0.0.1", "3333")
	if err != nil {
		t.Fatal(err)
	}
	//------------------------------------------

	if err := client.AddPlayer("Player1", color.RGBA{R: 255, A: 255}); err != nil {
		t.Fatal(err)
	}
	if err := client2.AddPlayer("Player2", color.RGBA{G: 255, A: 255}); err != nil {
		t.Fatal(err)
	}

	time.Sleep(600 * time.Millisecond)
	world.Country("Argentina").Occupier.Player = "Player1"
	world.PlayerQueue[0].Name = "Player1"
	world.PlayerQueue[1].Name = "Player2"
	time.Sleep(600 * time.Millisecond)

	if err := client.Reinforcement("Argentina", 1); err != nil {
		t.Fatal(err)
	}

	if err := client.EndTurn(); err != nil {
		t.Fatal(err)
	}
	if err := client2.EndTurn(); err != nil {
		t.Fatal(err)
	}

	if err := client.AttackOrMove("Argentina", "Peru", 1); err != nil {
		t.Fatal(err)
	}
}
