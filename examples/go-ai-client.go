package main

import (
	"RISK-CodeConflict/core"
	"RISK-CodeConflict/remote"
	"fmt"
	"image/color"
	"time"
)

//--------  AI (main)  -----------------------------------------------------------------------------------------------//

func main() {

	// config
	const name = "My Mega AI 9000"
	var clr = color.RGBA{R: 252, G: 3, B: 236, A: 255}
	const host = "localhost"
	const port = "1234"

	// init client
	client, err := remote.NewClient(host, port)
	if err != nil {
		panic(err)
	}

	// add player
	if err := client.AddPlayer(name, clr); err != nil {
		panic(err)
	}

	// Loop indefinitely, checking if it's the player's turn.
	for {
		// load world
		world := new(core.World)
		if err := client.Status(world); err != nil {
			fmt.Printf("load word error: %v\n", err)
		}

		// Check if it's the specified player's turn.
		if !world.Freeze && len(world.PlayerQueue) > 1 && world.PlayerQueue[0].Name == name {
			// --------- RUN AI -------------------------------------------------------------------

			// TODO: implement your ai here

			// EXAMPLE: list of your countries
			myCountries := make([]*core.Country, 0, 41)
			for _, c := range world.RndCountryList() {
				if c.Occupier != nil && c.Occupier.Player == name {
					myCountries = append(myCountries, c)
				}
			}
			fmt.Printf("myCountries: %#v\n", myCountries)

			// EXAMPLE: list of your recruiting regions
			myRecruitingRegion := make([]*core.Country, 0, 41)
			for _, c := range world.RndCountryList() {
				if c.Occupier != nil && c.Occupier.Player == name && c.RecruitingRegion {
					myRecruitingRegion = append(myRecruitingRegion, c)
				}
			}
			fmt.Printf("myRecruitingRegion: %#v\n", myRecruitingRegion)

			// EXAMPLE: recruiting
			err := client.Reinforcement(myRecruitingRegion[0].Name, 3)
			fmt.Printf("recruiting error: %v\n", err)

			// EXAMPLE: move or attack
			country := myRecruitingRegion[0]
			neighbor := country.Neighbors[0]
			err = client.AttackOrMove(country.Name, neighbor, 1)
			fmt.Printf("move error: %v\n", err)

			// EXAMPLE: End the turn and wait briefly before continuing.
			time.Sleep(400 * time.Millisecond)
			err = client.EndTurn()
			fmt.Printf("end turn error: %v\n\n", err)

			// ------------------------------------------------------------------------------------
		} else {
			// If it's not the player's turn, wait a short time before checking again.
			time.Sleep(200 * time.Millisecond)
		}
	}
}
