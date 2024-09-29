package core

import (
	"image/color"
	"reflect"
	"testing"
)

func TestWorld_Continent(t *testing.T) {
	w := NewWorld()

	// invalid continent
	c := w.Continent("invalidTestContinent")
	if c == nil {
		t.Fatalf("is nul")
	}
	if c.Name != "invalidTestContinent" {
		t.Fatalf("invalid name")
	}
	if c.Points != 0 {
		t.Fatalf("invalid points")
	}
	if len(c.Countries) != 0 {
		t.Fatalf("invalid countries")
	}

	// valid continent
	c = w.Continent("Europe")
	if c == nil {
		t.Fatalf("is nul")
	}
	if c.Name != "Europe" {
		t.Fatalf("invalid name")
	}
	if c.Points != 6 {
		t.Fatalf("invalid points")
	}
	if len(c.Countries) != 7 {
		t.Fatalf("invalid countries")
	}
}

func TestWorld_Country(t *testing.T) {
	w := NewWorld()

	// invalid continent
	c := w.Country("invalidTestCountry")
	if c == nil {
		t.Fatalf("is nul")
	}
	if c.Name != "invalidTestCountry" {
		t.Fatalf("invalid name")
	}
	if len(c.Neighbors) != 0 {
		t.Fatalf("invalid Neighbors")
	}

	// valid Country
	c = w.Country("Alaska")
	if c == nil {
		t.Fatalf("is nul")
	}
	if c.Name != "Alaska" {
		t.Fatalf("invalid name")
	}
	if len(c.Neighbors) != 3 {
		t.Fatalf("invalid Neighbors")
	}
}

func TestWorld_RndCountryList(t *testing.T) {
	w := NewWorld()
	cl1 := w.RndCountryList()
	cl2 := w.RndCountryList()
	cl3 := w.RndCountryList()
	cl4 := w.RndCountryList()

	if 42 != len(cl1) || 42 != len(cl2) || 42 != len(cl3) || 42 != len(cl4) {
		t.Fatalf("countries: %d", len(cl1))
	}
	if cl1[0].Name == cl2[0].Name && cl1[0].Name == cl3[0].Name && cl1[0].Name == cl4[0].Name {
		t.Fatalf("not random")
	}
}

func TestWorld_Player(t *testing.T) {
	w := NewWorld()

	// not found
	if p := w.Player("Player1"); p == nil || p.Name != "Player1" || p.Color.R != 0 || p.Color.G != 0 || p.Color.B != 0 || p.Color.A != 0 {
		t.Fatalf("invalid player")
	}

	// add player
	if err := w.AddPlayer("Player1", color.RGBA{R: 255, G: 255, B: 255, A: 255}); err != nil {
		t.Fatalf("invalid player")
	}

	// found
	if p := w.Player("Player1"); p == nil || p.Name != "Player1" || p.Color.R != 255 || p.Color.G != 255 || p.Color.B != 255 || p.Color.A != 255 {
		t.Fatalf("invalid player")
	}
}

func TestWorld_CalcReinforcement(t *testing.T) {
	// init
	w := NewWorld()
	_ = w.AddPlayer("Player1", color.RGBA{R: 255, G: 255, B: 255, A: 255})
	_ = w.AddPlayer("Player2", color.RGBA{R: 0, G: 0, B: 0, A: 0})
	w.InitPopulation()

	// Act
	all, countries, continents, sackBonus := w.CalcReinforcement("Player1")
	if countries != 21 || continents != 0 || all != countries+continents+sackBonus {
		t.Fatalf("!!RANDOM TEST!!: invalid stats: %d, %d, %d", all, countries, continents)
	}
	all, countries, continents, sackBonus = w.CalcReinforcement("Player2")
	if countries != 21 || continents != 0 || all != countries+continents+sackBonus {
		t.Fatalf("!!RANDOM TEST!!: invalid stats: %d, %d, %d", all, countries, continents)
	}
	all, countries, continents, sackBonus = w.CalcReinforcement("Player3")
	if countries != 0 || continents != 0 || all != countries+continents+sackBonus {
		t.Fatalf("invalid stats: %d, %d, %d", all, countries, continents)
	}

	// add eu
	for _, c := range w.Continent("Europe").Countries {
		co := w.Country(c)
		co.Occupier.Player = "Player1"
	}
	all, countries, continents, sackBonus = w.CalcReinforcement("Player1")
	if countries <= 21 || continents == 0 {
		t.Fatalf("invalid stats: %d, %d, %d", all, countries, continents)
	}
}

func TestWorld_AddPlayer(t *testing.T) {
	w := NewWorld()

	// add user
	if err := w.AddPlayer("  user1  ", color.RGBA{R: 255, G: 0, B: 0, A: 255}); err != nil {
		t.Fatal(err)
	}
	if err := w.AddPlayer("user1", color.RGBA{R: 111, G: 111, B: 111, A: 111}); err == nil {
		t.Fatal("no error:", err) // user exist
	}
	if err := w.AddPlayer("user99", color.RGBA{R: 255, G: 0, B: 0, A: 255}); err == nil {
		t.Fatal("no error:", err) // color exist
	}
	if err := w.AddPlayer("     ", color.RGBA{R: 255, G: 0, B: 0, A: 255}); err == nil {
		t.Fatal("no error:", err) // player empty
	}
	if len(w.PlayerQueue) != 1 {
		t.Fatal("invalid player count")
	}

	// valid player
	p := w.Player("user1")
	if p == nil {
		t.Fatalf("is nul")
	}
	if p.Name != "user1" {
		t.Fatalf("invalid name")
	}
	if p.Color.R != 255 || p.Color.G != 0 || p.Color.B != 0 || p.Color.A != 255 {
		t.Fatalf("invalid color: %#v", p.Color)
	}
}

func TestWorld_InitPopulation(t *testing.T) {

	// no player
	w := NewWorld()
	w.InitPopulation()
	for _, c := range w.Countries {
		if c.Occupier != nil {
			t.Fatalf("wrong Occupier")
		}
	}

	// one player
	w = NewWorld()
	_ = w.AddPlayer("Player 1", color.RGBA{R: 1, G: 0, B: 0, A: 255})
	w.InitPopulation()
	for _, c := range w.Countries {
		if c.Occupier.Player != "Player 1" {
			t.Fatalf("wrong Occupier")
		}
		if c.Occupier.Strength != 1 {
			t.Fatalf("wrong Strength")
		}
	}
	if w.Player("Player 1").Reinforcement != 45-42 {
		t.Fatalf("wrong Reinforcement")
	}

	// six player
	w = NewWorld()
	_ = w.AddPlayer("Player 1", color.RGBA{R: 1, G: 0, B: 0, A: 255})
	_ = w.AddPlayer("Player 2", color.RGBA{R: 2, G: 0, B: 0, A: 255})
	_ = w.AddPlayer("Player 3", color.RGBA{R: 3, G: 0, B: 0, A: 255})
	_ = w.AddPlayer("Player 4", color.RGBA{R: 4, G: 0, B: 0, A: 255})
	_ = w.AddPlayer("Player 5", color.RGBA{R: 5, G: 0, B: 0, A: 255})
	_ = w.AddPlayer("Player 6", color.RGBA{R: 6, G: 0, B: 0, A: 255})
	w.InitPopulation()
	for _, c := range w.Countries {
		if c.Occupier == nil {
			t.Fatalf("wrong Occupier")
		}
		if c.Occupier.Strength != 1 {
			t.Fatalf("wrong Strength")
		}
	}
	if w.Player("Player 1").Reinforcement != 50-5*6-42/6 {
		t.Fatalf("wrong Reinforcement")
	}
	if w.Player("Player 6").Reinforcement != 50-5*6-42/6 {
		t.Fatalf("wrong Reinforcement")
	}

	// double init (nothing)
	w.InitPopulation()
	for _, c := range w.Countries {
		if c.Occupier == nil {
			t.Fatalf("wrong Occupier")
		}
		if c.Occupier.Strength != 1 {
			t.Fatalf("wrong Strength")
		}
	}
	if w.Player("Player 1").Reinforcement != 50-5*6-42/6 {
		t.Fatalf("wrong Reinforcement")
	}
	if w.Player("Player 6").Reinforcement != 50-5*6-42/6 {
		t.Fatalf("wrong Reinforcement")
	}
}

func TestWorld_AttackOrMove(t *testing.T) {
	w := NewWorld() // empty world

	// freeze
	w.Freeze = true
	if err := w.AttackOrMove("", "", 0, "Player 3"); err == nil || err.Error() != "world is frozen" {
		t.Fatal(err)
	}
	w.Freeze = false

	// input
	if err := w.AttackOrMove("", "", 0, "Player 3"); err == nil || err.Error() != "attacker is empty" {
		t.Fatal(err)
	}
	if err := w.AttackOrMove("test", "", 0, "Player 3"); err == nil || err.Error() != "defender is empty" {
		t.Fatal(err)
	}
	if err := w.AttackOrMove("test", "test", 0, "Player 3"); err == nil || err.Error() != "attacker army strength must be greater than 0" {
		t.Fatal(err)
	}

	// second checks
	if err := w.AttackOrMove("Alaska", "test", 1, "Player 3"); err == nil || err.Error() != "no player found" {
		t.Fatal(err)
	}
	_ = w.AddPlayer("Player 3", color.RGBA{R: 255, G: 255, B: 255, A: 255})
	if err := w.AttackOrMove("Alaska", "test", 1, "Player 4"); err == nil || err.Error() != "not your turn" {
		t.Fatal(err)
	}
	if err := w.AttackOrMove("Alaska", "test", 1, "Player 3"); err == nil || err.Error() != "attacker army is nil or invalid" {
		t.Fatal(err)
	}

	// init world
	_ = w.AddPlayer("Player2", color.RGBA{R: 0, G: 0, B: 0, A: 0})
	w.PlayerQueue[0].Name = "Player1"
	w.PlayerQueue[1].Name = "Player2"
	w.InitPopulation()
	w.Country("Alaska").Occupier.Player = "Player2"

	// second checks (#2)
	if err := w.AttackOrMove("Alaska", "test", 1, "Player1"); err == nil || err.Error() != "cannot command enemy armies" {
		t.Fatal(err)
	}

	if err := w.AttackOrMove("Alaska", "test", 1, ""); err == nil || err.Error() != "at least one man must stay behind" {
		t.Fatal(err)
	}

	// add men power
	w.Country("Alaska").Occupier.Strength += 1

	// second checks (#3)
	if err := w.AttackOrMove("Alaska", "test", 1, ""); err == nil || err.Error() != "attacker and defender are not neighbors" {
		t.Fatal(err)
	}

	// EXIT: reinforcement
	w.Country("Alaska").Occupier.Player = "Player2"
	if err := w.AttackOrMove("Alaska", "Alaska", 10000, "Player1"); err == nil || err.Error() != "cannot command enemy armies" {
		t.Fatal(err)
	}
	w.Country("Alaska").Occupier.Player = "Player1"
	if err := w.AttackOrMove("Alaska", "Alaska", 1, ""); err == nil || err.Error() != "cannot recruit in this region" {
		t.Fatal(err)
	}

	// change RecruitingRegion flag
	w.Country("Alaska").RecruitingRegion = true

	if err := w.AttackOrMove("Alaska", "Alaska", 10000, "Player1"); err == nil || err.Error() != "not enough reinforcement" {
		t.Fatal(err)
	}
	if w.Player("Player1").Reinforcement != 19 {
		t.Fatal("wrong Reinforcement")
	}
	if err := w.AttackOrMove("Alaska", "Alaska", 1, "Player1"); err != nil || w.Country("Alaska").Invader.Strength != 1 {
		t.Fatal(err)
	}
	if w.Player("Player1").Reinforcement != 18 {
		t.Fatal("wrong Reinforcement")
	}

	// EXIT: move or attack
	w.Country("Alaska").Occupier.Player = "Player2"
	if err := w.AttackOrMove("Alaska", "Kamchatka", 1, "Player1"); err == nil || err.Error() != "cannot command enemy armies" {
		t.Fatal(err)
	}
	w.Country("Alaska").Occupier.Player = "Player1"
	if w.Country("Alaska").Occupier.Strength != 2 {
		t.Fatal("wrong Strength:")
	}
	if err := w.AttackOrMove("Alaska", "Kamchatka", 1, "Player1"); err != nil || w.Country("Kamchatka").Invader.Strength != 1 {
		t.Fatal(err)
	}
	if w.Country("Alaska").Occupier.Strength != 1 {
		t.Fatal("wrong Strength")
	}
}

func TestWorld_EndTurn(t *testing.T) {
	w := NewWorld()

	// freeze
	w.Freeze = true
	if err := w.EndTurn("Player1"); err == nil || err.Error() != "world is frozen" {
		t.Fatal(err)
	}
	w.Freeze = false

	// error
	if err := w.EndTurn("Player1"); err == nil || err.Error() != "no other player found" {
		t.Fatal(err)
	}
	_ = w.AddPlayer("Player1", color.RGBA{R: 255, G: 255, B: 255, A: 255})
	if err := w.EndTurn("Player1"); err == nil || err.Error() != "no other player found" {
		t.Fatal(err)
	}
	_ = w.AddPlayer("Player2", color.RGBA{R: 0, G: 0, B: 0, A: 0})
	w.PlayerQueue[0].Name = "PlayerA"
	w.PlayerQueue[1].Name = "PlayerB"
	if err := w.EndTurn("PlayerB"); err == nil || err.Error() != "cannot end enemy turn" {
		t.Fatal(err)
	}

	// init
	w.InitPopulation()
	for _, c := range w.Countries {
		c.Invader = NewArmy(w, 1, "PlayerA", c.Name)
	}

	// success
	if err := w.EndTurn("PlayerA"); err != nil {
		t.Fatal(err)
	}
	if w.PlayerQueue[0].Name != "PlayerB" || w.PlayerQueue[1].Name != "PlayerA" {
		t.Fatal("wrong users", w.PlayerQueue[0].Name, w.PlayerQueue[1].Name)
	}

	//
	if err := w.EndTurn("PlayerB"); err != nil {
		t.Fatal(err)
	}
	if err := w.EndTurn("PlayerA"); err != nil {
		t.Fatal(err)
	}
	if err := w.EndTurn("PlayerB"); err != nil {
		t.Fatal(err)
	}
}

func TestWorldClone(t *testing.T) {
	// Create an initial world instance and modify its state
	originalWorld := NewWorld()
	_ = originalWorld.AddPlayer("Fritz", color.RGBA{R: 0, G: 0, B: 0, A: 0})
	_ = originalWorld.AddPlayer("Bob", color.RGBA{R: 255, G: 255, B: 255, A: 255})
	originalWorld.InitPopulation()
	originalWorld.Round = 5
	originalWorld.Freeze = true
	originalWorld.RndCountryList()[0].Invader = NewArmy(originalWorld, 19, "inv", "home")

	// Clone the world
	clonedWorld := originalWorld.Clone()
	if clonedWorld == nil {
		t.Fatal("Clone returned nil")
	}

	// Check that the cloned world is not the same instance as the original
	if originalWorld == clonedWorld {
		t.Error("Cloned world should be a different instance")
	}

	// check random
	if originalWorld.rnd == nil || clonedWorld.rnd == nil {
		t.Error("Cloned world random should not be nil")
	}

	// check lock
	if originalWorld.lock == nil || clonedWorld.lock == nil {
		t.Error("Cloned world lock should not be nil")
	}

	// remove lock and random
	clonedWorld.lock = originalWorld.lock
	clonedWorld.rnd = originalWorld.rnd

	// check links
	if clonedWorld.RndCountryList()[0].world == nil || clonedWorld.RndCountryList()[0].Occupier.world == nil {
		t.Error("world link is nil")
	}

	// remove world links
	for _, c := range originalWorld.Countries {
		c.world = nil
		c.Occupier.world = nil
		if c.Invader != nil {
			c.Invader.world = nil
		}
	}
	for _, c := range clonedWorld.Countries {
		c.world = nil
		c.Occupier.world = nil
		if c.Invader != nil {
			c.Invader.world = nil
		}
	}

	// Compare the original and cloned world
	if !reflect.DeepEqual(originalWorld, clonedWorld) {
		t.Error("Cloned world does not match the original world")
	}

	// Modify the cloned world and ensure it does not affect the original
	clonedWorld.RndCountryList()[0].Occupier.Player = "Invalid"
	if reflect.DeepEqual(originalWorld, clonedWorld) {
		t.Error("Modifying the cloned world should not affect the original world")
	}
}
