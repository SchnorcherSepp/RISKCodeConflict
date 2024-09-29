package core

import (
	crnd "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"math/rand"
	"slices"
	"sort"
	"strings"
	"sync"
)

// World represents the entire game world, containing all continents, countries, and players.
// It acts as the main data structure managing the state of the game.
type World struct {
	rnd   *rand.Rand  // Random number generator used for various game mechanics.
	lock  *sync.Mutex // Mutex to handle concurrent access to the world state.
	NoLog bool

	// Freeze indicates whether the world state is locked. When set to true,
	// any SET-functions (such as AttackOrMove and EndTurn) have no effect,
	// effectively preventing any changes to the world.
	Freeze bool

	// Round keeps track of the current round number.
	// This value increments by 1 every time all players in the PlayerQueue have completed their turn.
	Round int

	// SubRound keeps track of the current player's turn within a round.
	// It increments after each player's turn and resets to 0 when all players have completed their turns in the round.
	SubRound int

	// Continents is a map of continent names to Continent structs.
	// The key is the name of the continent, and the value is a pointer to the Continent struct.
	// This map allows quick access to information about each continent in the game.
	Continents map[string]*Continent // Key: Continent.Name

	// Countries is a map of country names to Country structs.
	// The key is the name of the country, and the value is a pointer to the Country struct.
	// This map provides easy access to details about each country in the game.
	Countries map[string]*Country // Key: Country.Name

	// PlayerQueue is a slice that maintains the turn order of players during the game.
	// The first player in the queue is the active player. At the end of a turn,
	// the active player is moved to the end of the queue. When new players are added,
	// the queue is shuffled randomly to ensure a fair starting order.
	// The list managing all players participating in the game.
	PlayerQueue []*Player
}

//--------  GETTER  --------------------------------------------------------------------------------------------------//

// Continent retrieves a continent by its name from the world's Continents map.
// If the continent is not found, it returns an empty Continent struct with the given name.
// This ensures that the function always returns a valid Continent object.
func (w *World) Continent(name string) *Continent {
	ctt := w.Continents[name]
	if ctt != nil {
		return ctt
	} else {
		// Not found -> return empty Continent
		return &Continent{Name: name, Countries: []string{}}
	}
}

// Country retrieves a country by its name from the world's Countries map.
// If the country is not found, it returns an empty Country struct with the given name
// and an empty list of neighbors. This guarantees that the function always returns a valid Country object.
func (w *World) Country(name string) *Country {
	cnt := w.Countries[name]
	if cnt != nil {
		return cnt
	} else {
		// Not found -> return empty Country
		return &Country{world: w, Name: name, Neighbors: []string{}}
	}
}

// RndCountryList generates and returns a new, randomized list of all countries in the world.
// It creates a list of country pointers, shuffles it using the world's random number generator,
// and then returns the shuffled list.
func (w *World) RndCountryList() []*Country {
	// Create a list to hold the country pointers.
	list := make([]*Country, 0, len(w.Countries))

	// Append each country pointer from the world map to the list.
	for _, c := range w.Countries {
		list = append(list, c)
	}

	// Shuffle the list using the world's random number generator.
	w.rnd.Shuffle(len(list), func(i, j int) {
		list[i], list[j] = list[j], list[i]
	})

	// Return the shuffled list of countries.
	return list
}

// Player retrieves a player by their name from the world's PlayerQueue.
// If the player is not found, it returns an empty Player struct with the given name
// and a default color of black. This ensures that the function always returns a valid Player object.
func (w *World) Player(name string) *Player {
	// Create a default player to return if no match is found.
	ply := &Player{Name: name, Color: color.RGBA{}}

	// Search for the player in the PlayerQueue by name.
	for _, p := range w.PlayerQueue {
		if p != nil && p.Name == name {
			ply = p
			break // Player found, exit the loop.
		}
	}

	// If no matching player was found, return the default player.
	return ply
}

// CalcReinforcement calculates the total reinforcements a player receives based on:
//   - The number of countries they control.
//   - Any continent bonuses for fully controlled continents.
//   - A sack bonus for winning a battle in the last round.
//
// The function returns the total reinforcement points, as well as the individual contributions
// from countries, continents, and the sack bonus.
//
// Parameters:
//   - player: The name of the player for whom the reinforcement is being calculated.
//
// Returns:
//   - all: The total reinforcement points the player is awarded.
//   - countries: The reinforcement points from the number of countries controlled by the player.
//   - continents: The reinforcement points from any continents fully controlled by the player.
//   - sackBonus: The additional reinforcement points awarded if the player won a battle in the last round.
func (w *World) CalcReinforcement(player string) (all, countries, continents, sackBonus int) {

	//------  count controlled countries  ----------------------------//

	// Loop through all countries in the world and count how many are controlled by the player.
	for _, c := range w.Countries {
		if c.Occupier.Player == player {
			// Increment the count for controlled countries.
			countries++
		}
	}

	//------  check for continent control  ---------------------------//

	// For each continent, check if the player controls all countries within the continent.
	for _, continent := range w.Continents {
		// Assume the player controls the entire continent unless proven otherwise.
		totalControl := true

		// Check each country in the continent.
		for _, countryName := range continent.Countries {
			// Get the country object to check its occupier.
			countryObj := w.Country(countryName)
			if countryObj.Occupier.Player != player {
				// The player doesn't control this country, so they don't control the entire continent.
				totalControl = false
				break
			}
		}

		// If the player controls all countries in the continent, add the continent's bonus points.
		if totalControl {
			continents += continent.Points
		}
	}

	//------  get sack bonus  ----------------------------------------//

	// Check if the player won a battle in this last round.
	if w.Round == w.Player(player).LastBattleWonRound {
		sackBonus = minInt(w.Round, 20) // max bonus is 20
	}

	//------  calculate total reinforcements  ------------------------//

	// The total reinforcement is the sum of:
	//  - The number of countries controlled.
	//  - The continent control bonuses.
	//  - The sack bonus for winning a battle in this round.
	all = countries + continents + sackBonus

	//------  return values  -----------------------------------------//

	// Return the total reinforcements, along with the individual contributions from countries,
	// continents, and the sack bonus.
	return
}

// Clone creates a deep copy of the current World structure using JSON serialization and deserialization.
// This method utilizes the functions `Json()` and `FromJson()`.
//
// Returns:
//   - `*World`: A new World instance that is a deep copy of the original.
//   - `nil`: If an error occurs during the cloning process.
func (w *World) Clone() *World {
	clone := NewWorld()

	// Convert the current World to a JSON string.
	// The `Json()` function uses locking for thread safety.
	origJSON := w.Json()

	// Initialize the new World instance with the JSON string.
	// The `FromJson()` function also uses locking.
	err := clone.FromJson(origJSON)

	if err != nil {
		// Return `nil` in case of an error.
		println(err.Error())
		return nil
	} else {
		// Return the cloned World instance.
		return clone
	}
}

// Json converts the World object to a JSON-formatted string.
// This method uses locking to ensure thread safety.
//
// Returns:
//   - `string`: The JSON string representing the World object.
//     In case of an error, it returns the error message as a string.
func (w *World) Json() string {
	w.lock.Lock()         // Acquire lock for thread safety.
	defer w.lock.Unlock() // Release lock at the end of the function.

	// Serialize the World object to JSON format.
	b, err := json.Marshal(w)
	if err != nil {
		// Return the error message as a string in case of serialization failure.
		return err.Error()
	} else {
		// Return the JSON string.
		return string(b)
	}
}

//--------  SETTER  --------------------------------------------------------------------------------------------------//

// FromJson initializes the world's state from a given JSON string.
// This function reads the JSON string and updates the World object accordingly.
// It uses locking to ensure thread safety.
//
// Parameters:
//   - `s`: The JSON string representing the world's state.
//
// Returns:
//   - `error`: Returns an error in case of failure; returns `nil` on success.
func (w *World) FromJson(s string) error {
	if w.lock != nil {
		w.lock.Lock()         // Acquire lock for thread safety.
		defer w.lock.Unlock() // Release lock at the end of the function.
	}

	// detect error string
	if strings.HasPrefix(s, "err") {
		return errors.New(s)
	}

	// Deserialize the JSON data and update the World object.
	if err := json.Unmarshal([]byte(s), &w); err != nil {
		return err // Return the error in case of failure.
	}

	// ----- not exported vars ----- ///

	// Reinitialize the random number generator.
	var seed int64
	_ = binary.Read(crnd.Reader, binary.LittleEndian, &seed)
	w.rnd = rand.New(rand.NewSource(seed))

	// Reinitialize the lock.
	w.lock = new(sync.Mutex)

	// add world link to countries & armies
	for _, c := range w.Countries {
		c.world = w
		if c.Occupier != nil {
			c.Occupier.world = w
		}
		if c.Invader != nil {
			c.Invader.world = w
		}
	}

	// Success; no error occurred.
	return nil
}

// AddPlayer adds a new player to the world with the specified name and color.
// Returns an error if the name is empty, already exists, or if the color is nil or already taken.
// Ensures player names are trimmed and unique, and colors are valid and unique.
func (w *World) AddPlayer(name string, clr color.RGBA) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	// Remove leading and trailing whitespace from the player's name.
	name = strings.TrimSpace(name)

	// Check if the player's name is empty after trimming.
	if len(name) == 0 {
		return errors.New("player name is empty")
	}

	// Check if a player with the same name already exists in the world.
	// Check if the specified color is already being used by another player.
	for _, p := range w.PlayerQueue {
		// name
		if p.Name == name {
			return errors.New("player already exists")
		}
		// color
		r0, g0, b0, _ := clr.RGBA()
		r1, g1, b1, _ := p.Color.RGBA()
		if r0 == r1 && g0 == g1 && b0 == b1 {
			return errors.New("player color already exists")
		}
	}

	// Add the new player to the world's player list with the provided name and color.
	newP := &Player{
		Name:  name,
		Color: clr,
	}
	w.PlayerQueue = append(w.PlayerQueue, newP)

	// Shuffle PlayerQueue using the world's random number generator.
	w.rnd.Shuffle(len(w.PlayerQueue), func(i, j int) {
		w.PlayerQueue[i], w.PlayerQueue[j] = w.PlayerQueue[j], w.PlayerQueue[i]
	})

	// Return nil to indicate that the player was added successfully.
	return nil
}

// InitPopulation distributes initial armies to each country in the world.
// It randomizes the order of countries and players, then assigns one army to each country,
// cycling through the players until all countries are occupied.
func (w *World) InitPopulation() {
	w.lock.Lock()
	defer w.lock.Unlock()

	// no player
	if len(w.PlayerQueue) < 1 {
		return // ERROR: no player
	}

	// Get a randomized list of all countries.
	list := w.RndCountryList()

	// Sorts Continents
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Continent > list[j].Continent
	})

	// Sorts FortressRegion
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].FortressRegion && !list[j].FortressRegion
	})

	// Sorts RecruitingRegion
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].RecruitingRegion && !list[j].RecruitingRegion
	})

	// set reinforcement
	for _, p := range w.PlayerQueue {
		p.Reinforcement = 50 - 5*len(w.PlayerQueue)
	}

	// Distribute one army per country, cycling through the players.
	for i := 0; ; {
		for _, p := range w.PlayerQueue {
			if len(list) > i {
				// Assign one army to the current country with the current player as the occupier.
				c := list[i]
				c.Occupier = NewArmy(w, 1, p.Name, c.Name)
				// Pay for the army with Reinforcement points
				p.Reinforcement--
			} else {
				// Return once all countries are occupied.
				return
			}
			i++
		}
	}
}

// AttackOrMove processes an action where a player moves or attacks with troops from one country to a neighboring country.
// The function validates the input parameters, ensures that the player controls the attacking army, checks if the countries are neighbors,
// and then either moves troops or executes an attack. If a player attacks their own country, the function reinforces it
// using available reinforcements.
//
// Parameters:
//   - attacker: The name of the country initiating the attack or movement.
//   - defender: The name of the neighboring country being attacked or moved into. If it matches the attacker, reinforcements are deployed.
//   - strength: The number of troops to move or attack with. Must be greater than 0.
//   - player: The name of the player attempting the action. The player must control the attacking army.
//
// Returns:
//   - An error if any validation fails or the conditions for the attack or movement are not met.
//
// Error cases:
//   - Empty attacker or defender country names.
//   - Attacker army strength is less than 1.
//   - The player tries to command an army that doesn't belong to them.
//   - Not enough reinforcements when reinforcing.
//   - The attacker and defender countries are not neighbors.
func (w *World) AttackOrMove(attacker, defender string, strength int, player string) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	// check freeze
	if w.Freeze {
		return errors.New("world is frozen") // ERROR EXIT
	}

	//------  validate input  -----------------------------------------//

	// Validate that the attacker country name is not empty
	if attacker == "" {
		return errors.New("attacker is empty") // ERROR EXIT
	}

	// Validate that the defender country name is not empty
	if defender == "" {
		return errors.New("defender is empty") // ERROR EXIT
	}

	// Validate that the strength is positive and greater than 0
	if strength < 1 {
		return errors.New("attacker army strength must be greater than 0") // ERROR EXIT
	}

	//------  get objects  --------------------------------------------//

	// Player object
	playerObj := w.Player(player) // cannot be nil

	// Retrieve the attacker and defender country objects by name
	attackerObj := w.Country(attacker) // cannot be nil
	defenderObj := w.Country(defender) // cannot be nil

	// Retrieve the armies occupying the attacking and defending countries
	attackerArmy := attackerObj.Occupier // should not be null

	//------  second checks  ------------------------------------------//

	// Make sure that the player can only send orders on his own turn.
	// If 'player' is empty, commands can always be sent.
	if len(w.PlayerQueue) < 1 {
		return errors.New("no player found")
	}
	if player != "" && w.PlayerQueue[0].Name != player {
		return errors.New("not your turn")
	}

	// check attackerArmy
	if attackerArmy == nil {
		return errors.New("attacker army is nil or invalid")
	}

	// Make sure a player can only command his own armies.
	// An empty player can control all armies.
	if player != "" && attackerArmy.Player != player {
		return errors.New("cannot command enemy armies") // ERROR EXIT
	}

	// Ensure the attacking army has enough strength to leave at least one unit behind
	if attackerArmy.Strength-strength < 1 && attacker != defender {
		return errors.New("at least one man must stay behind") // ERROR EXIT
	}

	// Check if the countries are neighbors (i.e., they can interact with each other)
	if !slices.Contains(attackerObj.Neighbors, defender) && attacker != defender {
		return errors.New("attacker and defender are not neighbors") // ERROR EXIT
	}

	//------  EXIT  ---------------------------------------------------//

	// If the defender does not have an invader, create a new army for the invader
	if defenderObj.Invader == nil {
		// Create a new, empty army in the defender's territory representing the invader
		defenderObj.Invader = NewArmy(w, 0, attackerArmy.Player, attackerArmy.HomeBase)
	}

	// move, attack or reinforcement
	if attacker == defender {
		// MODE: Reinforcement
		//-----------------------

		// check RecruitingRegion flag
		if !defenderObj.RecruitingRegion {
			return errors.New("cannot recruit in this region")
		}
		// The attack on oneself is used to deploy reinforcement troops.
		if strength <= playerObj.Reinforcement {
			// The troops are withdrawn directly from the reinforcement pool.
			playerObj.Reinforcement -= strength
			defenderObj.Invader.Strength += strength
			return nil // SUCCESS EXIT
		} else {
			// not enough reinforcement
			return errors.New("not enough reinforcement") // ERROR EXIT
		}

	} else {
		// MODE: Move or Attack
		//-----------------------

		// Handle the attack or movement
		// Subtract the specified strength from the attacker army's strength
		attackerArmy.Strength -= strength

		// Add the moved or attacking units to the invader's strength
		defenderObj.Invader.Strength += strength

		// Return nil to indicate success with no errors
		return nil // SUCCESS EXIT
	}
}

// EndTurn processes the end of a player's turn, simulates any ongoing battles or troop movements,
// and transitions the game to the next player's turn. If all players have completed a turn, the game
// round is incremented.
//
// Parameters:
//   - player: The name of the player ending their turn. The player must be the one whose turn it is.
//
// Returns:
//   - An error if the player is attempting to end another player's turn or if no players are found.
//
// Error cases:
//   - No players found in the queue.
//   - Player tries to end the turn of another player.
func (w *World) EndTurn(player string) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	// check freeze
	if w.Freeze {
		return errors.New("world is frozen") // ERROR EXIT
	}

	//------  validate input  -----------------------------------------//

	// Ensure that the player can only end their own turn.
	// If 'player' is empty, all turns can be ended (for debug or admin purposes).
	// If the player does not match the current active player in PlayerQueue, return an error.
	if len(w.PlayerQueue) <= 1 {
		return errors.New("no other player found") // ERROR: No or one player in the queue.
	}
	if player != "" && w.PlayerQueue[0].Name != player {
		return errors.New("cannot end enemy turn") // ERROR: The player tries to end another player's turn.
	}

	//------  simulate battles  ---------------------------------------//

	// Simulate battles or movements for all countries with an invader army.
	// The invader either moves into the country (if they belong to the same player) or attacks the occupier
	// (if different players).
	for _, c := range w.Countries {
		if c.Invader != nil {

			// Check if the invader belongs to the same player as the occupier.
			if c.Invader.Player == c.Occupier.Player {
				// MODE: Move
				//-------------

				// Troop movement: Add the invader's strength to the occupier's.
				c.Occupier.Strength += c.Invader.Strength

			} else {
				// MODE: Attack
				//---------------

				// Battle: If the players differ, an attack occurs.
				log := c.Invader.Attack(c.Occupier, w.NoLog)

				// Print the battle log to show the results of each battle.
				for i, l := range log {
					if i > 0 {
						print(" | ")
					}
					println(l) // Outputs each event in the battle log.
				}

				// If the occupier's strength drops below 1, he loses the battle.
				if c.Occupier.Strength < 1 {
					// Replace the occupier with the invader (the invader now controls the country).
					c.Occupier = c.Invader
					c.Occupier.HomeBase = c.Name
					// The attacker has won a battle.
					c.Invader.PlayerObj().LastBattleWonRound = w.Round
				}
			}

			// Clear the invader (either they merged with the occupier or their attack was resolved).
			c.Invader = nil
		}
	}

	//------  end turn and go to next player  -------------------------//

	// Move the current player to the end of the queue and update the queue order.
	old := w.PlayerQueue[0]
	for i := 1; i < len(w.PlayerQueue); i++ {
		w.PlayerQueue[i-1] = w.PlayerQueue[i]
	}
	w.PlayerQueue[len(w.PlayerQueue)-1] = old // The old player is placed at the end of the queue.

	//------  round increment & reinforcement logic  ------------------//

	// Increment the SubRound counter, which tracks the turns of individual players within a round.
	w.SubRound++

	// Check if all players have completed their turns in the current round.
	if w.SubRound%len(w.PlayerQueue) == 0 {
		// A new round begins as all players have completed their turns.

		// Calculate and distribute reinforcements for all players.
		var livingPlayers = make([]*Player, 0, len(w.PlayerQueue))
		for _, p := range w.PlayerQueue {
			// calc reinforcement
			all, countries, continents, sackBonus := w.CalcReinforcement(p.Name)
			p.Reinforcement += all
			println(fmt.Sprintf("Reinforcements %s: countries=%d, continents=%d, sackBonus=%d", p.Name, countries, continents, sackBonus))

			// save living players
			if countries > 0 {
				livingPlayers = append(livingPlayers, p)
			}
		}
		w.PlayerQueue = livingPlayers

		// Go to next Round and reset the SubRound
		w.Round++
		w.SubRound = 0

		// print new turn
		println(fmt.Sprintf("\n==========  Round %d  ==========", w.Round))
	}

	// Return nil to indicate that the turn ended successfully without errors.
	return nil
}
