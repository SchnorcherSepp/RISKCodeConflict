package gui

import (
	"RISK-CodeConflict/core"
	"RISK-CodeConflict/gui/resources"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"image/color"
)

// preprocess generates a preprocessed image with the specified width and height,
// and draws country objects like fortresses and country names onto the background image.
// The countries parameter is a map where the key is not used and the value is a pointer to a core.Country object
// (see core.World.Countries).
//
// The native size is 2475x1392, which has an aspect ratio of 16:9 (see preprocessBg).
func preprocess(width, height float64, countries map[string]*core.Country) *ebiten.Image {

	// Generate a background image with continents
	img := preprocessBg(int(width), int(height))

	// Iterate over all countries in the map and draw their objects
	for _, country := range countries {

		// Calculate the correct scaled position of the country on the background
		posX := float64(country.Position[0]) * width / core.CountryPosScaleWidth
		posY := float64(country.Position[1]) * height / core.CountryPosScaleHeight

		// Draw the fortress if the country has a fortress region
		if country.FortressRegion {
			// Draw the fortress image at the calculated position
			// The fortress size is scaled relative to the background size
			preprocessObject(img, resources.Imgs.Fortress, posX, posY, 0.055) // Set object size here
		} else if country.RecruitingRegion {
			// Draw villages at recruiting regions
			preprocessObject(img, resources.Imgs.Village, posX, posY, 0.045) // Set object size here
		} else {
			// Draw border regions
			preprocessObject(img, resources.Imgs.Field, posX, posY, 0.045) // Set object size here
		}

		// Determine the color for the country name text
		txtClr := color.RGBA{R: 255, G: 255, B: 255, A: 255} // Default white color
		if country.BorderRegion {
			txtClr = color.RGBA{R: 255, G: 222, B: 3, A: 255} // Yellow color for border regions
		} else if country.FortressRegion {
			txtClr = color.RGBA{R: 255, G: 0, B: 0, A: 255} // red color for fortress regions
		}

		// Draw the country name text at the calculated position
		// The text size is scaled relative to the background size
		const heightOffset = 3.2
		preprocessText(img, country.Name, posX, posY, heightOffset, 0.007, txtClr) // Set text size here
	}

	// Return the preprocessed image
	return img
}

// preprocessBg combines multiple bg images into one and returns the preprocessed image.
// It creates a new image of the specified width and height, fills it with a light blue color,
// and draws the ocean and continent background images scaled to fit the new image.
//
// The native size is 2475x1392, which has an aspect ratio of 16:9.
func preprocessBg(width, height int) *ebiten.Image {

	// Create a new image with the specified width and height to hold the preprocessed result.
	img := ebiten.NewImage(width, height)

	// Clear the entire image with a light blue color to represent the ocean.
	img.Fill(color.RGBA{R: 173, G: 216, B: 230, A: 255})

	// --- ocean background ---

	// Get the dimensions of the ocean background image.
	bgOceanImgWidth := resources.Imgs.BgOcean.Bounds().Dx()
	bgOceanImgHeight := resources.Imgs.BgOcean.Bounds().Dy()

	// Create bgOceanImgOp for the ocean background image.
	// These options will scale the image to fit the specified width and height.
	bgOceanImgOp := new(ebiten.DrawImageOptions)
	bgOceanImgOp.GeoM.Scale(float64(width)/float64(bgOceanImgWidth), float64(height)/float64(bgOceanImgHeight))
	bgOceanImgOp.Filter = ebiten.FilterLinear

	// Draw the scaled ocean background image onto the new image.
	img.DrawImage(resources.Imgs.BgOcean, bgOceanImgOp)

	// --- continent background [16:9] ---

	// Get the dimensions of the continent background image.
	bgContinentImgWidth := resources.Imgs.BgContinent.Bounds().Dx()
	bgContinentImgHeight := resources.Imgs.BgContinent.Bounds().Dy()

	// Create bgContinentImgOp for the continent background image.
	// These options will scale the image to fit the specified width and height.
	bgContinentImgOp := new(ebiten.DrawImageOptions)
	bgContinentImgOp.GeoM.Scale(float64(width)/float64(bgContinentImgWidth), float64(height)/float64(bgContinentImgHeight))
	bgContinentImgOp.Filter = ebiten.FilterLinear

	// Draw the scaled continent background image onto the new image.
	img.DrawImage(resources.Imgs.BgContinent, bgContinentImgOp)

	// Return the preprocessed image, which now contains the combined backgrounds.
	return img
}

//--------------------------------------------------------------------------------------------------------------------//

// preprocessObject draws units or fortress images on a background image.
// It scales the object image based on the provided relative size (`objSizeRelToBg`) in relation to the background image.
// The object is positioned at coordinates (posX, posY) using its object image center as the reference point.
//
// Parameters:
// - bgImg: The background image on which the object will be drawn.
// - objImg: The object image to be drawn on the background.
// - posX, posY: The coordinates where the center of the object should be placed on the background image.
// - objSizeRelToBg: The relative size of the object image compared to the background image (as a percentage or scale factor).
func preprocessObject(bgImg, objImg *ebiten.Image, posX, posY, objSizeRelToBg float64) {

	// Get the dimensions of the background image
	bgWidth := bgImg.Bounds().Dx()

	// Get the dimensions of the object image
	objWidth := objImg.Bounds().Dx()
	objHeight := objImg.Bounds().Dy()

	// Calculate the scale of the object image relative to the background image
	objScale := (float64(bgWidth) * objSizeRelToBg) / float64(objWidth)
	objScaledWidth := float64(objWidth) * objScale
	objScaledHeight := float64(objHeight) * objScale

	// Configure the draw options for the object image
	op := new(ebiten.DrawImageOptions)
	op.GeoM.Scale(objScale, objScale)
	op.GeoM.Translate(posX-objScaledWidth/2, posY-objScaledHeight/2) // Center the object on (posX, posY)
	op.Filter = ebiten.FilterLinear

	// Draw the object image onto the background image
	bgImg.DrawImage(objImg, op)
}

// preprocessText draws text onto a background image at a specified position with a specified size and color.
// The text size is relative to the width of the background image, and a 3D effect is applied by drawing a shadow.
func preprocessText(bgImg *ebiten.Image, txt string, posX, posY, relOffY, txtSizeRelToBg float64, clr color.Color) {

	// Get the dimensions of the background image
	bgWidth := bgImg.Bounds().Dx()

	// Calculate the text size based on the background image width and the given relative size
	txtSize := (float64(bgWidth) * txtSizeRelToBg) / 0.9

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

	// Add vertically offset
	posY += relOffY * txtSize

	// Draw a shadow of the text to create a 3D effect
	text.Draw(bgImg, txt, fontFace, int(posX)+2, int(posY)+2, color.Black)

	// Draw the main text at the adjusted position with the specified color
	text.Draw(bgImg, txt, fontFace, int(posX), int(posY), clr)
}
