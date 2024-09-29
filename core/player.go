package core

import "image/color"

// Player represents a player in the game world. Each player has unique attributes, including a name, a color for visual
// representation on the map, and a pool of available reinforcements that they can deploy to strengthen their armies.
// The Player struct is used to track and manage each player's status, territory control, and in-game actions.
type Player struct {

	// Name is the unique identifier for the player within the game (see World.PlayerQueue). The Name value is often
	// used in logs, game commands, and user interfaces to identify which player is performing an action.
	//
	// Example values might include:
	//  - "Player1"
	//  - "AI_Opponent"
	//  - "Strategist_23"
	Name string

	// Color is the visual representation of the player on the game map.
	// Each player is assigned a specific color, which is used to paint their occupied territories,
	// display their units, and differentiate them from other players.
	//
	// The Color value can be set to any valid `color.Color` value. For instance:
	//  - color.RGBA{255, 0, 0, 255} for a red player.
	//  - color.RGBA{0, 0, 255, 255} for a blue player.
	Color color.RGBA // Use concrete type

	// Reinforcement represents the number of reinforcement units currently available to the player.
	// Reinforcements are typically awarded through various in-game mechanics, such as controlling entire continents,
	// capturing enemy territories, or meeting specific game objectives. These units can be deployed by the player
	// to strengthen their armies in occupied countries, giving them a strategic advantage.
	//
	// The Reinforcement value decreases as the player deploys units and increases as they earn new reinforcements
	// at the start of their turn or through special events.
	Reinforcement int

	// LastBattleWonRound indicates the most recent round in which the player won a battle.
	// This value is updated at the end of a turn by the `EndTurn()` function if the player has won any battles.
	// It is used for game mechanics such as granting bonuses or tracking player performance.
	//
	// Example:
	//  - If the player won a battle in round 5, this value would be set to 5.
	LastBattleWonRound int
}
