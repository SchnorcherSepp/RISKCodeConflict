package core

// Continent represents a continent in the game world. Each continent is a strategic region
// that groups multiple countries together and can be controlled by players.
// Controlling an entire continent provides additional points, reflecting its strategic importance in the game.
// The continent's name serves as its unique identifier, while its point value and the list of countries define its properties.
type Continent struct {

	// Name is the unique identifier for the continent.
	// It serves as a reference to the continent within the game world and can be used in game logic,
	// player interactions, or user interface elements. Typical values for Name might include:
	//  - "Europe"
	//  - "Asia"
	//  - "North America"
	Name string

	// Points is the number of points awarded to a player who controls all the countries within this continent.
	// The point value reflects the strategic significance of the continent. Larger continents or those with fewer
	// access points might provide a higher point value. For example:
	//  - Europe: 6 points
	//  - Asia: 8 points
	//  - Australia: 2 points
	// These points are added to the player's score, providing a tangible reward for controlling entire continents.
	Points int

	// Countries is a slice of strings representing the names of the countries that belong to this continent.
	// Each country name corresponds to a Country.Name value in the game (see World.Countries).
	//
	// Examples for the "Europe" continent could include:
	//  - "Great Britain"
	//  - "France"
	//  - "Germany"
	Countries []string // value: Country.Name
}
