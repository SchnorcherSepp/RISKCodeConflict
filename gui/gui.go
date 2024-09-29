package gui

import (
	"RISK-CodeConflict/core"
	"RISK-CodeConflict/gui/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

//----------  GUI struct  --------------------------------------------------------------------------------------------//

// interface check: ebiten.Game
var _ ebiten.Game = (*GUI)(nil)

// GUI implements the ebiten.Game interface and manages the GUI.
type GUI struct {
	world        *core.World // A reference to the game world or state being managed by the GUI.
	screenWidth  int         // The width of the game window in pixels.
	screenHeight int         // The height of the game window in pixels.

	viewport   [2]int // The top left position (x, y) of the viewport in the game world, used for panning.
	isDragging bool   // A flag indicating whether the user is currently dragging the map.
	lastMouseX int    // The last recorded X position of the mouse, used to calculate dragging distance.
	lastMouseY int    // The last recorded Y position of the mouse, used to calculate dragging distance.

	zoom            float64       // The zoom level for the viewport, where 1.0 represents 100% zoom.
	lastZoom        float64       // lastZoom saves the last zoom level to detect changes.
	preprocessedImg *ebiten.Image // preprocessedImg saves a prepared basic image at the correct resolution.

	redraw     bool // A flag indicating whether the screen should be redrawn in the next frame.
	autoRedraw bool // draws every frame

	selectCountry *core.Country // saves a country selected via the GUI

	lastRound    int // save last round to detect changes
	lastSubRound int // save last sub-round to detect changes
}

// Update updates an img by one tick. The given argument represents a screen image.
//
// Update updates only the img logic and Draw draws the screen.
//
// In the first frame, it is ensured that Update is called at least once before Draw. You can use Update
// to initialize the img state.
//
// After the first frame, Update might not be called or might be called once
// or more for one frame. The frequency is determined by the current TPS (tick-per-second).
func (g *GUI) Update() error {

	// Call all update functions
	//----------------------------
	g.updateZoomAndViewport()
	g.updateActiveCountry()
	g.updateAttackCountry()
	g.updateTurn()
	//----------------------------

	// auto redraw on changes
	if g.world != nil && (g.world.Round != g.lastRound || g.world.SubRound != g.lastSubRound) {
		g.lastRound = g.world.Round
		g.lastSubRound = g.world.SubRound
		g.redraw = true // redraw
	}

	return nil
}

// Draw draws the img screen by one frame.
//
// The give argument represents a screen image. The updated content is adopted as the img screen.
func (g *GUI) Draw(screen *ebiten.Image) {

	// Check if redraw is needed.
	// If the redraw flag is false, skip the drawing process to avoid unnecessary updates.
	if !g.autoRedraw && !g.redraw {
		return // Skip the drawing process
	}
	// Reset the redraw flag to false after starting the drawing process
	g.redraw = false

	// Prepare the preprocessed image and update if the zoom level has changed
	if g.preprocessedImg == nil || g.lastZoom != g.zoom {
		// Debug output to track when preprocess is called
		//println("call preprocess", time.Now().String(), "zoom:", g.zoom, "lastZoom:", g.lastZoom) // DEBUG GUI
		// Call the preprocess function to create the basic image with updated parameters (zoom)
		g.preprocessedImg = preprocess(float64(g.screenWidth)*g.zoom, float64(g.screenHeight)*g.zoom, g.world.Countries)
		// Store the current zoom level as the last known zoom level
		g.lastZoom = g.zoom
	}

	// Draw the image onto the screen with the specified options (op).
	// Translate (move) the image based on the current viewport position,
	// effectively adjusting the position of the image on the screen.
	op := new(ebiten.DrawImageOptions)
	op.GeoM.Translate(float64(-g.viewport[0]), float64(-g.viewport[1]))
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(g.preprocessedImg, op)

	// Call all drawing functions to render the content on the screen
	//----------------------------------------------------------------
	bgImgWidth := float64(g.preprocessedImg.Bounds().Dx())
	bgImgHeight := float64(g.preprocessedImg.Bounds().Dy())
	g.drawAllMark(screen, bgImgWidth, bgImgHeight)
	g.drawAllStats(screen, bgImgWidth, bgImgHeight)
	g.drawControls(screen)
	//----------------------------------------------------------------

	// Debugging: Print a message indicating the Draw method has been called
	//println("call Draw", time.Now().String(), "zoom:", g.zoom, "viewport:", g.viewport[0], g.viewport[1])  // DEBUG GUI
}

// Layout accepts a native outside size in device-independent pixels and returns the img logical screen
// size.
//
// On desktops, the outside is a window or a monitor (fullscreen mode). On browsers, the outside is a body
// element. On mobiles, the outside is the view's size.
//
// Even though the outside size and the screen size differ, the rendering scale is automatically adjusted to
// fit with the outside.
//
// Layout is called almost every frame.
//
// It is ensured that Layout is invoked before Update is called in the first frame.
//
// If Layout returns non-positive numbers, the caller can panic.
//
// You can return a fixed screen size if you don't care, or you can also return a calculated screen size
// adjusted with the given outside size.
func (g *GUI) Layout(_, _ int) (int, int) {
	return g.screenWidth, g.screenHeight
}

//----------  Constructor  -------------------------------------------------------------------------------------------//

// RunGUI initializes the game window and starts the GUI loop.
// The Draw function is called with 30 Ticks per second.
//
// This function is blocking!
func RunGUI(screenWidth, screenHeight int, title string, world *core.World, autoRedraw bool) error {

	// Constants for the configuration
	const (
		tps             = 30    // Ticks per second for the game loop
		decorated       = true  // Whether the window should be decorated with title bar and borders
		floating        = false // Whether the window should always stay on top of other windows
		clearEveryFrame = false // Whether the screen should be cleared at the start of each frame
	)

	// Setup window properties
	ebiten.SetWindowTitle(title)                                   // Set the window title
	ebiten.SetWindowIcon([]image.Image{resources.Imgs.Icon})       // Set the window icon
	ebiten.SetWindowSize(screenWidth, screenHeight)                // Set the window size
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled) // Allow resizing of the window
	ebiten.SetTPS(tps)                                             // Set the ticks per second (game loop frequency)
	ebiten.SetWindowDecorated(decorated)                           // Enable window decoration (title bar, borders)
	ebiten.SetWindowFloating(floating)                             // Disable always-on-top behavior
	ebiten.SetScreenClearedEveryFrame(clearEveryFrame)             // Disable automatic clearing of the screen every frame

	// Create a new GUI instance
	gui := &GUI{
		world:        world,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		zoom:         1,
		redraw:       true,
		autoRedraw:   autoRedraw,
	}

	// Run the game loop (this call is blocking)
	return ebiten.RunGame(gui)
}

//----------  zoom and scroll  ---------------------------------------------------------------------------------------//

// updateZoomAndViewport adjusts the zoom level and the position of the viewport based on user input.
// It handles mouse wheel input for zooming in and out, arrow/WASD key presses for panning the viewport,
// and mouse dragging for moving the map.
//
// Manipulates GUI.viewport and GUI.zoom
func (g *GUI) updateZoomAndViewport() {

	//-------------------------------------//
	// --- Handle zoom via mouse wheel --- //
	//-------------------------------------//

	// max zoom level
	// FIX for "panic: atlas: the image being put on an atlas is too big: width: 17512, height: 9849"
	const maxZoomLvl = 5
	const zoomSpeed = 10

	// Handle mouse wheel input for zooming
	_, dy := ebiten.Wheel()
	if dy != 0 {
		// Get the current mouse position relative to the window
		mouseX, mouseY := ebiten.CursorPosition()

		// Calculate the relative position of the mouse within the window
		mouseRelX := float64(mouseX) / float64(g.screenWidth)
		mouseRelY := float64(mouseY) / float64(g.screenHeight)

		// Calculate the position of the mouse within the game world
		// considering the current zoom level and viewport offset
		currentViewX := float64(g.viewport[0]) + mouseRelX*float64(g.screenWidth)/g.zoom
		currentViewY := float64(g.viewport[1]) + mouseRelY*float64(g.screenHeight)/g.zoom

		// Calculate the new zoom level by adjusting the current zoom
		newZoom := g.zoom + dy/zoomSpeed*g.zoom

		// Ensure the zoom level stays within the bounds [1, 10]
		if newZoom < 1 {
			newZoom = 1 // Minimum zoom level
		}
		if newZoom > maxZoomLvl {
			newZoom = maxZoomLvl // Maximum zoom level
		}

		// Calculate the zoom factor as the ratio of the new zoom to the current zoom
		zoomFactor := newZoom / g.zoom

		// Adjust the viewport to keep the mouse position centered
		// This shifts the viewport to compensate for the change in zoom level
		g.viewport[0] = int(currentViewX*zoomFactor - mouseRelX*float64(g.screenWidth)/newZoom)
		g.viewport[1] = int(currentViewY*zoomFactor - mouseRelY*float64(g.screenHeight)/newZoom)

		// Update the zoom level and set the redraw flag to true
		g.zoom = newZoom
		g.redraw = true
	}

	//--------------------------------------------------------//
	// --- Handle mouse dragging for panning the viewport --- //
	//--------------------------------------------------------//

	// Check if the left mouse button is pressed
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// Get the current mouse position
		mouseX, mouseY := ebiten.CursorPosition()

		// If dragging has just started, initialize lastMouseX and lastMouseY
		if !g.isDragging {
			g.isDragging = true
			g.lastMouseX, g.lastMouseY = mouseX, mouseY
		}

		// Calculate the difference in mouse position since the last frame
		dragX := mouseX - g.lastMouseX
		dragY := mouseY - g.lastMouseY

		// Check if there was any movement
		if dragX != 0 || dragY != 0 {
			// Update the viewport position based on the mouse movement
			g.viewport[0] -= dragX
			g.viewport[1] -= dragY

			// Update lastMouseX and lastMouseY for the next frame
			g.lastMouseX, g.lastMouseY = mouseX, mouseY

			// Mark the screen for redraw
			g.redraw = true
		}
	} else {
		// Reset dragging state if the left mouse button is released
		g.isDragging = false
	}

	//----------------------------------------------------------------//
	// --- Handle arrow/WASD key presses for panning the viewport --- //
	//----------------------------------------------------------------//

	// movement speed
	const scrollSpeed = 30

	// Move the viewport up
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		// Calculate the amount to move based on screen height
		g.viewport[1] -= int(float64(g.screenWidth) / scrollSpeed) // always use the width for the calculation!
		// Mark the screen for redraw
		g.redraw = true
	}

	// Move the viewport down
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		// Calculate the amount to move based on screen height
		g.viewport[1] += int(float64(g.screenWidth) / scrollSpeed) // always use the width for the calculation!
		// Mark the screen for redraw
		g.redraw = true
	}

	// Move the viewport left
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		// Calculate the amount to move based on screen width
		g.viewport[0] -= int(float64(g.screenWidth) / scrollSpeed)
		// Mark the screen for redraw
		g.redraw = true
	}

	// Move the viewport right
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		// Calculate the amount to move based on screen width
		g.viewport[0] += int(float64(g.screenWidth) / scrollSpeed)
		// Mark the screen for redraw
		g.redraw = true
	}

	//-------------------------------------//
	// --- check boundary for viewport --- //
	//-------------------------------------//

	// Ensure the viewport does not go out of the top boundary
	if g.viewport[1] < 0 {
		g.viewport[1] = 0 // Top boundary
		g.redraw = true
	}

	// Ensure the viewport does not go out of the bottom boundary
	var maxY = int(float64(g.screenHeight) * (g.zoom - 1))
	if g.viewport[1] > maxY {
		g.viewport[1] = maxY // Bottom boundary
		g.redraw = true
	}

	// Ensure the viewport does not go out of the left boundary
	if g.viewport[0] < 0 {
		g.viewport[0] = 0 // Left boundary
		g.redraw = true
	}

	// Ensure the viewport does not go out of the right boundary
	var maxX = int(float64(g.screenWidth) * (g.zoom - 1))
	if g.viewport[0] > maxX {
		g.viewport[0] = maxX // Right boundary
		g.redraw = true
	}

}
