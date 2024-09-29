package gui

import (
	"RISK-CodeConflict/core"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// updateActiveCountry checks for mouse input to determine if a country has been clicked on the game map.
// It updates the currently selected country based on the user's mouse click and adjusts the screen accordingly.
//
// Functionality:
// - Checks if the world or its countries are properly initialized; if not, it exits early.
// - Checks if the left mouse button was clicked; if not, it exits early.
// - Retrieves the current mouse cursor position on the screen.
// - Iterates through all countries to determine if the cursor position falls within the bounds of any country.
// - Computes the position and dimensions of each country on the screen, considering the current zoom level and viewport offset.
// - If a country is clicked (i.e., the mouse cursor is within the country's visual bounds), it sets this country as the currently selected one.
// - Logs the name of the selected country or an "unselect" message if no country is selected.
// - Updates the `selectCountry` field with the newly selected country and triggers a screen redraw if the selection has changed.
func (g *GUI) updateActiveCountry() {
	// check input
	if g.world == nil || g.world.Countries == nil {
		return // skip
	}
	// no click, no check
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}

	//--------------------------------

	// get mouse position
	x, y := ebiten.CursorPosition()

	// find country
	var result *core.Country
	for _, country := range g.world.Countries {

		// basic image size
		var bgImgWidth int
		var bgImgHeight int
		if g.preprocessedImg != nil {
			bgImgWidth = g.preprocessedImg.Bounds().Dx()
			bgImgHeight = g.preprocessedImg.Bounds().Dy()
		}

		// Calculate the correct scaled position of the country on the screen
		countryPosX := country.Position[0]*bgImgWidth/core.CountryPosScaleWidth - g.viewport[0]
		countryPosY := country.Position[1]*bgImgHeight/core.CountryPosScaleHeight - g.viewport[1]

		// object dimension
		var dim = int(100 * g.zoom)
		x1 := countryPosX - dim/2
		y1 := countryPosY - dim/2
		x2 := x1 + dim
		y2 := y1 + dim

		// check position
		if x >= x1 && x <= x2 && y >= y1 && y <= y2 {
			result = country
			break
		}
	}

	// select country
	if g.selectCountry != result {
		// print action
		if result != nil {
			println("select", result.Name)
		} else {
			println("unselect")
		}
		// set new country
		g.selectCountry = result
		// update screen
		g.redraw = true
	}
}

// updateAttackCountry handles the logic for selecting a country to attack based on the user's right mouse click.
// It determines if the user clicks on a neighboring country of the currently selected country and triggers an attack.
func (g *GUI) updateAttackCountry() {
	// check input
	if g.world == nil || g.world.Countries == nil {
		return // skip
	}
	// no click, no check
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		return
	}
	// active country?
	selectCountry := g.selectCountry
	if selectCountry == nil {
		return
	}

	//--------------------------------

	// get mouse position
	x, y := ebiten.CursorPosition()

	// find neighbor countries
	var result *core.Country
	var list []*core.Country
	list = append(list, selectCountry.NeighborsObj()...)
	list = append(list, selectCountry)
	for _, country := range list {

		// basic image size
		var bgImgWidth int
		var bgImgHeight int
		if g.preprocessedImg != nil {
			bgImgWidth = g.preprocessedImg.Bounds().Dx()
			bgImgHeight = g.preprocessedImg.Bounds().Dy()
		}

		// Calculate the correct scaled position of the country on the screen
		countryPosX := country.Position[0]*bgImgWidth/core.CountryPosScaleWidth - g.viewport[0]
		countryPosY := country.Position[1]*bgImgHeight/core.CountryPosScaleHeight - g.viewport[1]

		// object dimension
		var dim = int(100 * g.zoom)
		x1 := countryPosX - dim/2
		y1 := countryPosY - dim/2
		x2 := x1 + dim
		y2 := y1 + dim

		// check position
		if x >= x1 && x <= x2 && y >= y1 && y <= y2 {
			result = country
			break
		}
	}

	// attack country
	if result != nil {
		// activePlayer
		activePlayer := ""
		if len(g.world.PlayerQueue) > 0 {
			activePlayer = g.world.PlayerQueue[0].Name
		}

		// ATTACK
		strength := 1
		if ebiten.IsKeyPressed(ebiten.KeyControl) {
			strength = 5
		}
		if err := g.world.AttackOrMove(selectCountry.Name, result.Name, strength, activePlayer); err != nil {
			println("ERROR:", err.Error())
		}

		// update screen
		g.redraw = true
	}
}

// updateTurn updates the game state for the active player's turn.
// It checks for input, processes the end of the turn, and triggers a screen redraw.
func (g *GUI) updateTurn() {
	// Check if the world and its countries are initialized.
	if g.world == nil || g.world.Countries == nil {
		return // Skip the turn update if the world data is not available.
	}

	// No input detected, skip processing.
	if !inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return // Skip the turn update if Enter key is not pressed.
	}

	// Retrieve the active player from the queue.
	activePlayer := ""
	if len(g.world.PlayerQueue) > 0 {
		activePlayer = g.world.PlayerQueue[0].Name // Get the first player in the queue.
	}

	// Process the end of the turn for the active player.
	if err := g.world.EndTurn(activePlayer); err != nil {
		println("ERROR:", err.Error()) // Print error message if ending the turn fails.
	}

	// Mark the screen for a redraw.
	g.redraw = true
}
