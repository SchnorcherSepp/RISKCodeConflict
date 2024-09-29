package core

import (
	crnd "crypto/rand"
	"encoding/binary"
	"math/rand"
	"sync"
)

// CountryPosScaleWidth and CountryPosScaleHeight define the dimensions of the image
// that represents the game world, where the positions of the countries are scaled.
// These constants are used to ensure that the coordinates for the countries fit within
// the specified image dimensions, which helps maintain visual accuracy when rendering.
const (
	CountryPosScaleWidth  = 1778 // Width of the game world image.
	CountryPosScaleHeight = 1000 // Height of the game world image.
)

// NewWorld initializes and returns a new instance of the World struct.
// It sets up the initial state of the game world, including the continents and countries,
// along with their respective properties such as positions and neighboring countries.
func NewWorld() *World {
	world := &World{
		Continents: map[string]*Continent{
			"North America": {Name: "North America", Points: 5, Countries: []string{"Alaska", "Alberta", "Central America", "Eastern US", "Greenland", "Northwest Territory", "Ontario", "Quebec", "Western US"}},
			"Europe":        {Name: "Europe", Points: 6, Countries: []string{"Great Britain", "Iceland", "Northern Europe", "Scandinavia", "Southern Europe", "Ukraine", "Western Europe"}},
			"South America": {Name: "South America", Points: 2, Countries: []string{"Argentina", "Brazil", "Venezuela", "Peru"}},
			"Africa":        {Name: "Africa", Points: 3, Countries: []string{"Congo", "East Africa", "Egypt", "Madagascar", "North Africa", "South Africa"}},
			"Australia":     {Name: "Australia", Points: 2, Countries: []string{"Eastern Australia", "New Guinea", "Indonesia", "Western Australia"}},
			"Asia":          {Name: "Asia", Points: 8, Countries: []string{"Afghanistan", "China", "India", "Irkutsk", "Japan", "Kamchatka", "Middle East", "Mongolia", "Siam", "Siberia", "Ural", "Yakutsk"}},
		},

		Countries: map[string]*Country{
			// North America
			"Alaska": {
				Name:         "Alaska",
				Position:     [2]int{110, 134},
				Neighbors:    []string{"Northwest Territory", "Alberta", "Kamchatka"},
				Continent:    "North America",
				BorderRegion: true,
			},
			"Alberta": {
				Name:             "Alberta",
				Position:         [2]int{258, 218},
				Neighbors:        []string{"Alaska", "Northwest Territory", "Ontario", "Western US"},
				Continent:        "North America",
				RecruitingRegion: true,
			},
			"Central America": {
				Name:         "Central America",
				Position:     [2]int{276, 469},
				Neighbors:    []string{"Western US", "Eastern US", "Venezuela"},
				Continent:    "North America",
				BorderRegion: true,
			},
			"Eastern US": {
				Name:             "Eastern US",
				Position:         [2]int{404, 364},
				Neighbors:        []string{"Central America", "Western US", "Ontario", "Quebec"},
				Continent:        "North America",
				RecruitingRegion: true,
			},
			"Greenland": {
				Name:         "Greenland",
				Position:     [2]int{626, 112},
				Neighbors:    []string{"Iceland", "Northwest Territory", "Ontario", "Quebec"},
				Continent:    "North America",
				BorderRegion: true,
			},
			"Northwest Territory": {
				Name:             "Northwest Territory",
				Position:         [2]int{280, 115},
				Neighbors:        []string{"Alaska", "Alberta", "Ontario", "Greenland"},
				Continent:        "North America",
				RecruitingRegion: true,
			},
			"Ontario": {
				Name:             "Ontario",
				Position:         [2]int{390, 240},
				Neighbors:        []string{"Northwest Territory", "Alberta", "Western US", "Eastern US", "Quebec", "Greenland"},
				Continent:        "North America",
				FortressRegion:   true,
				RecruitingRegion: true,
			},
			"Quebec": {
				Name:             "Quebec",
				Position:         [2]int{505, 249},
				Neighbors:        []string{"Greenland", "Ontario", "Eastern US"},
				Continent:        "North America",
				RecruitingRegion: true,
			},
			"Western US": {
				Name:             "Western US",
				Position:         [2]int{263, 334},
				Neighbors:        []string{"Alberta", "Eastern US", "Ontario", "Central America"},
				Continent:        "North America",
				RecruitingRegion: true,
			},

			// South America
			"Argentina": {
				Name:             "Argentina",
				Position:         [2]int{441, 831},
				Neighbors:        []string{"Peru", "Brazil"},
				Continent:        "South America",
				RecruitingRegion: true,
			},
			"Brazil": {
				Name:         "Brazil",
				Position:     [2]int{546, 656},
				Neighbors:    []string{"Argentina", "Peru", "Venezuela", "North Africa"},
				Continent:    "South America",
				BorderRegion: true,
			},
			"Venezuela": {
				Name:         "Venezuela",
				Position:     [2]int{405, 551},
				Neighbors:    []string{"Brazil", "Peru", "Central America"},
				Continent:    "South America",
				BorderRegion: true,
			},
			"Peru": {
				Name:             "Peru",
				Position:         [2]int{419, 687},
				Neighbors:        []string{"Venezuela", "Brazil", "Argentina"},
				Continent:        "South America",
				FortressRegion:   true,
				RecruitingRegion: true,
			},

			// Europe
			"Great Britain": {
				Name:             "Great Britain",
				Position:         [2]int{719, 300},
				Neighbors:        []string{"Iceland", "Scandinavia", "Northern Europe", "Western Europe"},
				Continent:        "Europe",
				RecruitingRegion: true,
			},
			"Iceland": {
				Name:         "Iceland",
				Position:     [2]int{758, 184},
				Neighbors:    []string{"Great Britain", "Scandinavia", "Greenland"},
				Continent:    "Europe",
				BorderRegion: true,
			},
			"Northern Europe": {
				Name:             "Northern Europe",
				Position:         [2]int{891, 318},
				Neighbors:        []string{"Great Britain", "Scandinavia", "Ukraine", "Southern Europe", "Western Europe"},
				Continent:        "Europe",
				FortressRegion:   true,
				RecruitingRegion: true,
			},
			"Scandinavia": {
				Name:             "Scandinavia",
				Position:         [2]int{896, 174},
				Neighbors:        []string{"Iceland", "Great Britain", "Northern Europe", "Ukraine"},
				Continent:        "Europe",
				RecruitingRegion: true,
			},
			"Southern Europe": {
				Name:         "Southern Europe",
				Position:     [2]int{906, 419},
				Neighbors:    []string{"Western Europe", "Northern Europe", "Ukraine", "Middle East", "Egypt", "North Africa"},
				Continent:    "Europe",
				BorderRegion: true,
			},
			"Ukraine": {
				Name:         "Ukraine",
				Position:     [2]int{1058, 253},
				Neighbors:    []string{"Southern Europe", "Northern Europe", "Scandinavia", "Ural", "Afghanistan", "Middle East"},
				Continent:    "Europe",
				BorderRegion: true,
			},
			"Western Europe": {
				Name:         "Western Europe",
				Position:     [2]int{769, 440},
				Neighbors:    []string{"Great Britain", "Northern Europe", "Southern Europe", "North Africa"},
				Continent:    "Europe",
				BorderRegion: true,
			},

			// Africa
			"Congo": {
				Name:             "Congo",
				Position:         [2]int{963, 745},
				Neighbors:        []string{"North Africa", "East Africa", "South Africa"},
				Continent:        "Africa",
				FortressRegion:   true,
				RecruitingRegion: true,
			},
			"East Africa": {
				Name:         "East Africa",
				Position:     [2]int{1072, 696},
				Neighbors:    []string{"Egypt", "Congo", "South Africa", "Madagascar", "North Africa", "Middle East"},
				Continent:    "Africa",
				BorderRegion: true,
			},
			"Egypt": {
				Name:         "Egypt",
				Position:     [2]int{965, 569},
				Neighbors:    []string{"North Africa", "Southern Europe", "Middle East", "East Africa"},
				Continent:    "Africa",
				BorderRegion: true,
			},
			"Madagascar": {
				Name:             "Madagascar",
				Position:         [2]int{1138, 888},
				Neighbors:        []string{"East Africa", "South Africa"},
				Continent:        "Africa",
				RecruitingRegion: true,
			},
			"North Africa": {
				Name:         "North Africa",
				Position:     [2]int{828, 614},
				Neighbors:    []string{"Egypt", "East Africa", "Congo", "Western Europe", "Southern Europe", "Brazil"},
				Continent:    "Africa",
				BorderRegion: true,
			},
			"South Africa": {
				Name:             "South Africa",
				Position:         [2]int{983, 887},
				Neighbors:        []string{"Congo", "East Africa", "Madagascar"},
				Continent:        "Africa",
				RecruitingRegion: true,
			},

			// Asia
			"Afghanistan": {
				Name:         "Afghanistan",
				Position:     [2]int{1208, 362},
				Neighbors:    []string{"Ural", "China", "India", "Middle East", "Ukraine"},
				Continent:    "Asia",
				BorderRegion: true,
			},
			"China": {
				Name:             "China",
				Position:         [2]int{1411, 430},
				Neighbors:        []string{"Ural", "Siberia", "Mongolia", "Siam", "India", "Afghanistan"},
				Continent:        "Asia",
				RecruitingRegion: true,
			},
			"India": {
				Name:             "India",
				Position:         [2]int{1315, 525},
				Neighbors:        []string{"Middle East", "Siam", "China", "Afghanistan"},
				Continent:        "Asia",
				FortressRegion:   true,
				RecruitingRegion: true,
			},
			"Irkutsk": {
				Name:             "Irkutsk",
				Position:         [2]int{1446, 232},
				Neighbors:        []string{"Siberia", "Yakutsk", "Kamchatka", "Mongolia"},
				Continent:        "Asia",
				RecruitingRegion: true,
			},
			"Japan": {
				Name:             "Japan",
				Position:         [2]int{1649, 353},
				Neighbors:        []string{"Mongolia", "Kamchatka"},
				Continent:        "Asia",
				RecruitingRegion: true,
			},
			"Kamchatka": {
				Name:         "Kamchatka",
				Position:     [2]int{1583, 140},
				Neighbors:    []string{"Yakutsk", "Irkutsk", "Mongolia", "Alaska", "Japan"},
				Continent:    "Asia",
				BorderRegion: true,
			},
			"Middle East": {
				Name:         "Middle East",
				Position:     [2]int{1108, 526},
				Neighbors:    []string{"Ukraine", "Afghanistan", "India", "Egypt", "Southern Europe", "East Africa"},
				Continent:    "Asia",
				BorderRegion: true,
			},
			"Mongolia": {
				Name:             "Mongolia",
				Position:         [2]int{1472, 333},
				Neighbors:        []string{"Japan", "China", "Siberia", "Irkutsk", "Kamchatka"},
				Continent:        "Asia",
				RecruitingRegion: true,
			},
			"Siam": {
				Name:         "Siam",
				Position:     [2]int{1453, 558},
				Neighbors:    []string{"Indonesia", "India", "China"},
				Continent:    "Asia",
				BorderRegion: true,
			},
			"Siberia": {
				Name:             "Siberia",
				Position:         [2]int{1327, 148},
				Neighbors:        []string{"Ural", "China", "Mongolia", "Yakutsk", "Irkutsk"},
				Continent:        "Asia",
				FortressRegion:   true,
				RecruitingRegion: true,
			},
			"Ural": {
				Name:         "Ural",
				Position:     [2]int{1239, 224},
				Neighbors:    []string{"Ukraine", "Siberia", "Afghanistan", "China"},
				Continent:    "Asia",
				BorderRegion: true,
			},
			"Yakutsk": {
				Name:             "Yakutsk",
				Position:         [2]int{1465, 109},
				Neighbors:        []string{"Kamchatka", "Irkutsk", "Siberia"},
				Continent:        "Asia",
				RecruitingRegion: true,
			},

			// Australia
			"Eastern Australia": {
				Name:             "Eastern Australia",
				Position:         [2]int{1684, 852},
				Neighbors:        []string{"Western Australia", "New Guinea"},
				Continent:        "Australia",
				RecruitingRegion: true,
			},
			"New Guinea": {
				Name:             "New Guinea",
				Position:         [2]int{1618, 692},
				Neighbors:        []string{"Eastern Australia", "Western Australia", "Indonesia"},
				Continent:        "Australia",
				RecruitingRegion: true,
			},
			"Indonesia": {
				Name:         "Indonesia",
				Position:     [2]int{1464, 744},
				Neighbors:    []string{"New Guinea", "Western Australia", "Siam"},
				Continent:    "Australia",
				BorderRegion: true,
			},
			"Western Australia": {
				Name:             "Western Australia",
				Position:         [2]int{1551, 886},
				Neighbors:        []string{"Eastern Australia", "New Guinea", "Indonesia"},
				Continent:        "Australia",
				FortressRegion:   true,
				RecruitingRegion: true,
			},
		},
	}

	// init random
	var seed int64
	_ = binary.Read(crnd.Reader, binary.LittleEndian, &seed)
	world.rnd = rand.New(rand.NewSource(seed))

	// init lock
	world.lock = new(sync.Mutex)

	// init player list
	world.PlayerQueue = make([]*Player, 0, 12)

	// add world link to countries
	for _, c := range world.Countries {
		c.world = world
	}

	// return
	return world
}
