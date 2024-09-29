package core_test

import (
	"RISK-CodeConflict/core"
	"testing"
)

func TestNeighborsObj(t *testing.T) {
	// Setup world and countries
	world := core.NewWorld()
	country := world.Country("Alaska")

	// Act
	neighbors := country.NeighborsObj()

	// Assert
	// []string{"Northwest Territory", "Alberta", "Kamchatka"},
	if len(neighbors) != 3 {
		t.Error("neighbors length should be 3")
	}
	if neighbors[0] == nil || neighbors[0].Name != "Northwest Territory" || len(neighbors[0].Neighbors) == 0 {
		t.Error("wrong neighbor")
	}
	if neighbors[1] == nil || neighbors[1].Name != "Alberta" || len(neighbors[1].Neighbors) == 0 {
		t.Error("wrong neighbor")
	}
	if neighbors[2] == nil || neighbors[2].Name != "Kamchatka" || len(neighbors[2].Neighbors) == 0 {
		t.Error("wrong neighbor")
	}
}

func TestContinentObj(t *testing.T) {

	// Setup world and countries
	world := core.NewWorld()
	country := world.Country("Alaska")

	// Act
	continent := country.ContinentObj()

	// test
	if continent == nil || continent.Name != "North America" || len(continent.Countries) == 0 {
		t.Error("wrong continent")
	}
}
