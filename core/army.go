package core

import (
	"fmt"
	"math/rand"
	"sort"
)

// Army represents a military unit in the game. Each army is associated with a specific player
// and has a certain strength, which determines its combat effectiveness. Armies are stationed
// in a particular country (HomeBase), from where they can launch attacks or defend against enemy forces.
type Army struct {
	world *World

	// Strength indicates the combat power of the army.
	// A higher strength value means the army has more units and is stronger in battle.
	Strength int

	// Player is the name of the player who controls this army.
	// This should correspond to a Player.Name value in the game, identifying the owner of the army (see World.PlayerQueue).
	Player string // value: Player.Name

	// HomeBase is the name of the country where the army is currently stationed.
	// This should match a Country.Name value in the game, indicating the army's current location (see World.Countries).
	HomeBase string // value: Country.Name
}

// NewArmy creates and returns a new Army instance with the specified strength, player name, and home base country.
// This function initializes the army with the provided values, allowing for the creation of new military units
// in the game. It also associates the army with the game world to access relevant game data.
//
// Parameters:
//   - world: A pointer to the game world (`*World`) that the army is part of.
//   - strength: The number of units in the army, indicating its initial combat power.
//   - player: The name of the player controlling the army.
//   - homeBase: The name of the country where the army is stationed.
//
// Returns:
//   - A pointer to the newly created `Army` instance.
//
// Panics:
//   - If `world` is nil, the function will panic with an error message indicating that a valid world must be provided.
func NewArmy(world *World, strength int, player, homeBase string) *Army {
	if world == nil {
		panic("NewArmy: world is nil")
	}
	return &Army{
		world:    world,
		Strength: strength,
		Player:   player,
		HomeBase: homeBase,
	}
}

//--------  GETTER  --------------------------------------------------------------------------------------------------//

// PlayerObj retrieves the `Player` object associated with this army.
// It uses the Player name stored in the army to look up the corresponding Player in the World.
// This function provides an easy way to access the full Player struct controlling the army.
//
// Returns:
//   - A pointer to the `Player` object representing the player who owns this army.
func (a *Army) PlayerObj() *Player {
	return a.world.Player(a.Player)
}

// HomeBaseObj retrieves the `Country` object representing the army's home base.
// It uses the HomeBase name stored in the army to look up the corresponding Country in the World.
// This function provides an easy way to access the full Country struct where the army is stationed.
//
// Returns:
//   - A pointer to the `Country` object representing the country where the army is currently located.
func (a *Army) HomeBaseObj() *Country {
	return a.world.Country(a.HomeBase)
}

// Description provides a textual description of the army, including the player's name,
// the home base country, and the army's current strength. This function is useful for debugging
// or displaying army information in the game UI.
//
// Returns:
//   - A string describing the army in the format: "<Player's Name>'s <HomeBase> Army with <Strength> men".
func (a *Army) Description() string {
	return fmt.Sprintf("%s's %s Army with %d men", a.Player, a.HomeBase, a.Strength)
}

//--------  SETTER  --------------------------------------------------------------------------------------------------//

// Attack simulates a battle between two armies, determining the outcome based on their respective strengths,
// random dice rolls, and potential terrain bonuses (e.g., fortress regions).
// The function uses a dice-based combat mechanic, where each army rolls a number of dice proportional to their strength.
// Additional rules apply for defending armies in fortified regions, giving them a strategic advantage.
//
// Combat Mechanism:
//   - The attacker rolls up to 3 dice, while the defender rolls up to 2 dice (or 3 dice if in a fortress region).
//   - Dice are rolled for each side, sorted in descending order, and compared pairwise.
//   - For each pair of dice, the side with the lower roll loses one unit of strength.
//   - If the dice values are equal in a comparison, the defender always wins the tie.
//   - The battle continues in rounds until one army is defeated (i.e., its strength reaches 0).
//
// Parameters:
//   - `defender`: A pointer to the `Army` instance representing the defending army.
//   - `noLog`: A boolean flag to control logging. If `noLog` is true, logging is disabled to avoid performance overhead.
//
// Returns:
//   - A slice of strings (if `noLog` is false) that contains a step-by-step log of the battle, including dice rolls,
//     losses, and the final outcome of the battle. If `noLog` is true, the function returns an empty slice.
func (a *Army) Attack(defender *Army, noLog bool) (log []string) {
	// Initialize the log slice if logging is enabled.
	if !noLog {
		log = make([]string, 0, 20)
	}
	attacker := a

	// Check if both armies are ready for battle.
	if attacker == nil || defender == nil || attacker.Strength <= 0 || defender.Strength <= 0 {
		if !noLog {
			log = append(log, "Not all armies were ready to fight. There was no battle.")
		}
		return
	}

	// Log initial army details.
	if !noLog {
		log = append(log, attacker.Description()+"  attacks  "+defender.Description())
	}

	// Conduct battle rounds until one army is defeated.
	for round := 1; true; round++ {
		// Log the current round number.
		if !noLog {
			log = append(log, fmt.Sprintf("--- ROUND %d ---", round))
		}

		// Determine the number of dice each army will roll based on their strengths.
		attackDiceCount := minInt(3, attacker.Strength)
		defendDiceCount := minInt(2, defender.Strength)

		// Check if the defender is in a fortified region and adjust their dice count.
		if defender.HomeBaseObj().FortressRegion {
			defendDiceCount = minInt(3, defender.Strength) // Defender receives a bonus.

			// Log the defender's advantage if in a fortress region.
			if !noLog {
				log = append(log, fmt.Sprintf("%s is a fortress region!", defender.HomeBase))
			}
		}

		// Roll dice for both armies.
		attackDice := rollDice(a.world.rnd, attackDiceCount)
		defendDice := rollDice(a.world.rnd, defendDiceCount)

		// Sort dice rolls in descending order for comparison.
		sort.Sort(sort.Reverse(sort.IntSlice(attackDice)))
		sort.Sort(sort.Reverse(sort.IntSlice(defendDice)))

		// Log the dice rolls.
		if !noLog {
			log = append(log, fmt.Sprintf("Attacker dice: %v", attackDice))
			log = append(log, fmt.Sprintf("Defender dice: %v", defendDice))
		}

		// Compare the highest dice rolls and determine unit losses.
		oldAttackerStr := attacker.Strength
		oldDefenderStr := defender.Strength
		for i := 0; i < minInt(attackDiceCount, defendDiceCount); i++ {
			if attackDice[i] > defendDice[i] {
				defender.Strength-- // Defender loses a unit.
			} else {
				attacker.Strength-- // Attacker loses a unit.
			}
		}

		// Log the losses.
		if !noLog {
			log = append(log, fmt.Sprintf("The attacker lost %d units.", oldAttackerStr-attacker.Strength))
			log = append(log, fmt.Sprintf("The defender lost %d units.", oldDefenderStr-defender.Strength))
		}

		// Determine if the battle should end based on remaining strengths.
		if attacker.Strength <= 0 {
			if !noLog {
				log = append(log, fmt.Sprintf("The defender was victorious with %d men left.", defender.Strength))
			}
			break
		}
		if defender.Strength <= 0 {
			if !noLog {
				log = append(log, fmt.Sprintf("The attacker was victorious with %d men left.", attacker.Strength))
			}
			break
		}
	}
	return
}

//--------  HELPER  --------------------------------------------------------------------------------------------------//

// rollDice simulates rolling a specified number of dice and returns a slice of integers representing the results.
// Each die rolled produces a random number between 1 and 6 (inclusive).
//
// Parameters:
//   - rnd: A pointer to a random number generator (`*rand.Rand`). If `rnd` is nil, the function returns an empty slice.
//   - count: The number of dice to roll. If `count` is less than 1, the function returns an empty slice.
//
// Returns:
//   - A slice of integers, each representing the result of a die roll. The slice has a length equal to `count`.
//     If invalid input is provided (e.g., `rnd` is nil or `count` < 1), an empty slice is returned.
func rollDice(rnd *rand.Rand, count int) []int {
	// Validate input parameters.
	if rnd == nil || count < 1 {
		return make([]int, 0) // Invalid input: Return an empty slice.
	}

	// Create a slice to store the results of the dice rolls.
	dice := make([]int, count)

	// Roll the specified number of dice, generating a random value between 1 and 6 for each die.
	for i := 0; i < count; i++ {
		dice[i] = rnd.Intn(6) + 1
	}
	return dice
}

// minInt returns the smaller of the two provided integers.
//
// Parameters:
//   - a: The first integer to compare.
//   - b: The second integer to compare.
//
// Returns:
//   - The smaller of the two integers (`a` or `b`). If both integers are equal, it returns that value.
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
