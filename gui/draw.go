package gui

import (
	"RISK-CodeConflict/core"
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"image/color"
	"math"
	"strings"
)

// drawAllMark draws marks on the screen for the selected country and its neighbors.
// It uses drawMark to draw a mark at the center of the selected country and for each neighbor.
func (g *GUI) drawAllMark(screen *ebiten.Image, bgImgWidth, bgImgHeight float64) {

	// Determine the active player from the queue
	activePlayer := ""
	if len(g.world.PlayerQueue) > 0 {
		activePlayer = g.world.PlayerQueue[0].Name
	}

	// Mark all countries occupied by the active player
	for _, c := range g.world.Countries {
		if c.Occupier != nil && c.Occupier.Player == activePlayer {
			// Draw a mark for the occupied country
			clr := color.White
			if !c.RecruitingRegion {
				clr = color.Black
			}
			g.drawMark(screen, bgImgWidth, bgImgHeight, c, 0.053, color.RGBA{R: 0, G: 0, B: 0, A: 0}, clr)
		}
	}

	// Check if there is a selected country
	if g.selectCountry != nil {
		clr := color.Black // Default color for the mark
		if g.selectCountry.Occupier != nil && g.selectCountry.Occupier.Player == activePlayer {
			clr = color.White // Change color if the selected country belongs to the active player
		}

		// Draw a mark for the selected country
		g.drawMark(screen, bgImgWidth, bgImgHeight, g.selectCountry, 0.053, color.RGBA{R: 255, G: 67, B: 7, A: 90}, clr)

		// Iterate over each neighbor of the selected country
		for _, ns := range g.selectCountry.Neighbors {
			nc := g.world.Countries[ns]

			clr = color.Black // Default color for neighbor marks
			if nc.Occupier != nil && nc.Occupier.Player == activePlayer {
				clr = color.White // Change color if the neighbor country belongs to the active player
			}

			// Draw a mark for the neighbor country
			g.drawMark(screen, bgImgWidth, bgImgHeight, nc, 0.053, color.RGBA{R: 77, G: 177, B: 7, A: 90}, clr)
		}
	}
}

// drawMark draws a mark (filled circle) on the screen at a specified position and size.
// The size of the mark is relative to the background image size.
func (g *GUI) drawMark(screen *ebiten.Image, bgImgWidth, bgImgHeight float64, country *core.Country, markSizeRelToBg float64, clr, clr2 color.Color) {

	// Calculate the radius of the mark based on the relative size
	radius := (bgImgWidth * markSizeRelToBg) / 2

	// Calculate the correct scaled position of the country on the screen
	countryPosX := float64(country.Position[0])*bgImgWidth/core.CountryPosScaleWidth - float64(g.viewport[0])
	countryPosY := float64(country.Position[1])*bgImgHeight/core.CountryPosScaleHeight - float64(g.viewport[1])

	// Draw a filled circle (mark) on the background image at the calculated position and size
	vector.DrawFilledCircle(screen, float32(countryPosX), float32(countryPosY), float32(radius), clr, false)
	drawCircle(screen, countryPosX, countryPosY, radius, clr2)
}

//--------------------------------------------------------------------------------------------------------------------//

// drawAllStats renders the military statistics (e.g., army strength) for all countries
// on the game map. It draws visual markers representing the occupier's and invader's army
// strengths at the respective positions of each country.
//
// Parameters:
// - screen: The *ebiten.Image object representing the game's display screen.
// - bgImgWidth: The width of the background image used as the game map.
// - bgImgHeight: The height of the background image used as the game map.
func (g *GUI) drawAllStats(screen *ebiten.Image, bgImgWidth, bgImgHeight float64) {
	// Countries
	for _, c := range g.world.Countries {
		countryPosX := float64(c.Position[0])
		countryPosY := float64(c.Position[1])
		// Invader
		if c.Invader != nil && c.Invader.Strength > 0 {
			// Invader movement
			if c.Name != c.Invader.HomeBase {
				homePosX := float64(c.Invader.HomeBaseObj().Position[0])
				homePosY := float64(c.Invader.HomeBaseObj().Position[1])
				g.drawMovement(screen, bgImgWidth, bgImgHeight, countryPosX-30, countryPosY-30, homePosX, homePosY, c.Invader.PlayerObj().Color)
			}
			// Invader stats
			g.drawStats(screen, bgImgWidth, bgImgHeight, countryPosX-30, countryPosY-30, 0.011, c.Invader.PlayerObj().Color, c.Invader.Strength)
		}
		// Occupier stats
		if c.Occupier != nil {
			g.drawStats(screen, bgImgWidth, bgImgHeight, countryPosX, countryPosY, 0.02, c.Occupier.PlayerObj().Color, c.Occupier.Strength)
		}
	}
}

// drawStats draws a visual marker representing the army strength for a country at its map position.
// It also displays the numerical strength of the army next to the visual marker.
//
// Parameters:
// - screen: The *ebiten.Image object where the stats will be drawn.
// - bgImgWidth: The width of the game map image.
// - bgImgHeight: The height of the game map image.
// - countryPosX: The X position of the country on the map (relative to the image).
// - countryPosY: The Y position of the country on the map (relative to the image).
// - markSizeRelToBg: The size of the visual marker (circle) as a relative proportion of the map width.
// - clr: The color used to draw the marker, representing the player.
// - strength: The strength of the army to display numerically near the marker.
func (g *GUI) drawStats(screen *ebiten.Image, bgImgWidth, bgImgHeight, countryPosX, countryPosY float64, markSizeRelToBg float64, clr color.Color, strength int) {

	// Calculate the radius of the mark based on the relative size
	radius := (bgImgWidth * markSizeRelToBg) / 2

	// Calculate the correct scaled position of the country on the screen
	posX := countryPosX*bgImgWidth/core.CountryPosScaleWidth - float64(g.viewport[0])
	posY := countryPosY*bgImgHeight/core.CountryPosScaleHeight - float64(g.viewport[1])

	// Draw a filled circle (mark) on the background image at the calculated position and size
	vector.DrawFilledCircle(screen, float32(posX), float32(posY), float32(radius), clr, false)

	//----------------------

	// stat text (army count)
	txt := fmt.Sprintf("%d", strength)
	txtSize := radius * 1.4
	// Parse the TrueType font from the provided font data
	// Create a new font face with the specified size and full hinting for better readability
	ttFont, _ := truetype.Parse(gomono.TTF)
	fontFace := truetype.NewFace(ttFont, &truetype.Options{
		Size:    txtSize,
		Hinting: font.HintingFull,
	})
	// Adjust the position to center the text horizontally and vertically relative to the given position
	posX -= float64(len(txt)) * txtSize * 0.31 // Adjust horizontally
	posY += txtSize * 0.35                     // Adjust vertically
	// Draw the main text at the adjusted position with the specified color
	text.Draw(screen, txt, fontFace, int(posX), int(posY), color.Black)
}

// TODO: description
func (g *GUI) drawMovement(screen *ebiten.Image, bgImgWidth, bgImgHeight, countryPosX, countryPosY, homePosX, homePosY float64, clr color.Color) {

	// Calculate the correct scaled position of the country on the screen
	posX := countryPosX*bgImgWidth/core.CountryPosScaleWidth - float64(g.viewport[0])
	posY := countryPosY*bgImgHeight/core.CountryPosScaleHeight - float64(g.viewport[1])
	homeX := homePosX*bgImgWidth/core.CountryPosScaleWidth - float64(g.viewport[0])
	homeY := homePosY*bgImgHeight/core.CountryPosScaleHeight - float64(g.viewport[1])

	// Draw a line
	vector.StrokeLine(screen, float32(posX), float32(posY), float32(homeX), float32(homeY), 3, clr, false)
}

//--------------------------------------------------------------------------------------------------------------------//

func (g *GUI) drawControls(screen *ebiten.Image) {
	// generate text
	sb := new(strings.Builder)
	sb.WriteString(fmt.Sprintf("Round: %d.%d\n", g.world.Round, g.world.SubRound+1))
	sb.WriteString("Press Enter to end the turn.\n\nPlayer queue:\n")
	for i, po := range g.world.PlayerQueue {
		if i == 0 {
			sb.WriteString(" > ")
		} else {
			sb.WriteString(" - ")
		}
		sb.WriteString(fmt.Sprintf("%s [%d]\n", po.Name, po.Reinforcement))
	}
	// print
	ebitenutil.DebugPrintAt(screen, sb.String(), 10, 10)
}

//--------------------------------------------------------------------------------------------------------------------//

// drawCircle draws a circle on the given image with the specified center (cx, cy), radius, and color.
func drawCircle(img *ebiten.Image, cx, cy, radius float64, col color.Color) {
	// Loop over all points in the bounding box of the circle
	for x := cx - radius; x <= cx+radius; x++ {
		for y := cy - radius; y <= cy+radius; y++ {
			// Calculate the distance from the center of the circle
			dx := x - cx
			dy := y - cy
			dist := math.Sqrt(dx*dx + dy*dy)

			// If the distance is approximately equal to the radius, set the pixel color
			if math.Abs(dist-radius) < 1 {
				img.Set(int(x-1), int(y), col)
				img.Set(int(x), int(y), col)
			}
		}
	}
}
