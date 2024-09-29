package ai

import (
	"RISK-CodeConflict/core"
	"RISK-CodeConflict/remote"
	"image/color"
	"math/rand"
	"time"
)

// Play runs the AI logic for a specified player in the game world.
// The function continuously monitors if it's the player's turn to act.
// If it's the player's turn, the AI will reinforce its territories, send move/attack commands, and end the turn.
func Play(host, port string, player string, clr color.RGBA) {

	// init client
	client, err := remote.NewClient(host, port)
	if err != nil {
		println(err.Error())
		return // exit
	}

	// add player
	if err := client.AddPlayer(player, clr); err != nil {
		println(err.Error())
		return // exit
	}

	// Loop indefinitely, checking if it's the player's turn.
	for {
		// load world
		world := new(core.World)
		if err := client.Status(world); err != nil {
			println(err.Error())
		}

		// Check if it's the specified player's turn.
		if !world.Freeze && len(world.PlayerQueue) > 1 && world.PlayerQueue[0].Name == player {
			// --------- RUN AI ---------

			// Calculate distances of countries relative to enemy territories.
			// The countries are grouped into slices based on their distance from the nearest enemy.
			distance := countriesByDistance(world, player)

			// Randomize the order of countries within each distance group to create more unpredictable behavior.
			for _, d := range distance {
				rand.Shuffle(len(d), func(i, j int) { d[i], d[j] = d[j], d[i] })
			}

			// Reinforce phase: Add units to territories until the player's reinforcement points are exhausted.
			// A maximum of 600 reinforcement attempts will be made to avoid long loops.
			for i := 0; i < 600; i++ {
				// Select a random country to receive a reinforcement unit.
				for _, c := range world.RndCountryList() {
					// Try to reinforce one unit in the selected country.
					// If reinforcement is successful (no error), break out of the loop.
					if err := client.AttackOrMove(c.Name, c.Name, 1); err == nil {
						break
					}
				}

				// Check if the player has no reinforcement points left.
				if world.PlayerQueue[0].Reinforcement < 1 {
					break
				}
			}

			// Movement and attack phase: For each distance group, move or attack neighboring countries.
			for _, d := range distance {
				for _, c := range d {
					// Try to attack or move units to each neighboring country.
					for _, n := range c.Neighbors {
						var err error = nil

						// Continue sending units until an error occurs (e.g., no more units to move).
						for err == nil {
							err = client.AttackOrMove(n, c.Name, 1)
						}
					}
				}
			}

			// End the turn and wait briefly before continuing.
			time.Sleep(400 * time.Millisecond)
			if err := client.EndTurn(); err != nil {
				println(err.Error())
			}

		} else {
			// If it's not the player's turn, wait a short time before checking again.
			time.Sleep(200 * time.Millisecond)
		}
	}
}

//--------  HELPER  --------------------------------------------------------------------------------------------------//

// countriesByDistance returns a slice of slices of countries where each sub-slice contains countries
// that are a specific distance away from an enemy (i.e., c.Occupier.Player != player).
// The outer slice index represents the distance from the nearest enemy.
func countriesByDistance(w *core.World, player string) [][]*core.Country {

	// Map to store the distance of each country from the nearest enemy.
	distanceMap := make(map[*core.Country]int)

	// Queue for BFS
	queue := make([]*core.Country, 0)

	// Initialize queue with all enemy-occupied countries and set their distance to 0.
	for _, country := range w.Countries {
		if country.Occupier.Player != player {
			queue = append(queue, country)
			distanceMap[country] = 0 // Enemies have distance 0
		}
	}

	// Perform BFS to calculate the distance of each country from the nearest enemy.
	for len(queue) > 0 {
		// Dequeue the first country
		current := queue[0]
		queue = queue[1:]

		// Get the current country's distance from the nearest enemy
		currentDistance := distanceMap[current]

		// Iterate over neighboring countries
		for _, neighborName := range current.Neighbors {
			neighbor := w.Country(neighborName)

			// Check if the neighbor has already been visited
			if _, visited := distanceMap[neighbor]; !visited {
				// Set the neighbor's distance (current distance + 1)
				distanceMap[neighbor] = currentDistance + 1
				// Enqueue the neighbor
				queue = append(queue, neighbor)
			}
		}
	}

	// Build the result slice, where each index represents a distance from the enemy.
	maxDistance := 0
	for _, dist := range distanceMap {
		if dist > maxDistance {
			maxDistance = dist
		}
	}

	// Create a result slice with the size of the maximum distance + 1
	result := make([][]*core.Country, maxDistance+1)

	// Populate the result slice with countries at each distance
	for country, dist := range distanceMap {
		result[dist] = append(result[dist], country)
	}

	return result
}
