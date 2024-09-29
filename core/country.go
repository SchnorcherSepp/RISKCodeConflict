package core

// Country represents a country in the game world. Each country is a distinct region with its own unique name,
// geographical position, neighboring countries, and continent affiliation. It serves as a strategic point of control
// within the game, as each country can be occupied by one player's army and may host battles between opposing armies.
// Countries are connected to their neighbors and continents, influencing both movement and control dynamics.
//
// Additionally, countries can have special statuses, such as being a border region, a fortress region, or a recruiting region,
// which affect gameplay and strategies. Border regions often form strategic choke points, while fortress regions grant
// defensive bonuses, and recruiting regions allow the creation of new units.
type Country struct {
	world *World

	// Name is the unique identifier for this country within the game world.
	// It is used to refer to the country in various game mechanics, such as when identifying
	// ownership, checking borders, or issuing commands. Examples of country names might include:
	//  - "France"
	//  - "Germany"
	//  - "Brazil"
	Name string

	// Position represents the geographical coordinates of the country on the game map.
	// This position is stored as an array of two integers [x, y], corresponding to a scaled map of
	// CountryPosScaleWidth x CountryPosScaleHeight. The coordinates can be used to draw the country
	// on a visual map. Example values might include:
	//  - [50, 75] indicating a position at x=50 and y=75 on the game map.
	Position [2]int

	// Neighbors is a list of names of the countries that share a border with this country.
	// These neighboring countries are directly adjacent to the current country and can be moved to or attacked.
	// The names in this list correspond to the Name values of other Country structs in the game.
	Neighbors []string // value: Country.Name

	// Continent is the name of the continent to which this country belongs.
	// This value corresponds to a Continent.Name value in the game (see World.Continents), linking the country to its continent.
	// Continent affiliation influences game mechanics such as scoring and continent control bonuses.
	// Example values might include:
	//  - "Europe"
	//  - "Africa"
	//  - "South America"
	Continent string // value: Continent.Name

	// BorderRegion indicates whether this country is considered a border region.
	// A border region is a country that shares borders with countries from different continents or holds
	// significant strategic importance due to its location. Controlling border regions often influences
	// the flow of the game, as they serve as gateways between continents or key defensive points.
	BorderRegion bool

	// FortressRegion indicates whether this country is a designated fortress region.
	// Fortresses provide defensive bonuses to the controlling player, making it harder for opposing armies to capture them.
	// A country that is a fortress region may not be a border region, as fortresses are often located away from immediate threats.
	FortressRegion bool

	// RecruitingRegion indicates whether new units can be recruited or raised in this country.
	// Typically, recruiting regions represent strongholds or capitals within the game. A country must be a fortress region
	// to be a recruiting region, while border regions cannot be designated recruiting regions to prevent overpowered troop generation.
	RecruitingRegion bool

	// Occupier is a pointer to the army currently occupying and controlling this country.
	// This value indicates which player owns the country and can defend it against attacks.
	// There must always be an occupier.
	Occupier *Army

	// Invader is a pointer to an attacking army currently attempting to take control of the country.
	// During a battle, the invader will either defeat the occupier and take control, or be destroyed in the attempt.
	// If Invader and Occupier are controlled by the same player, it is only a troop transfer.
	// If Invader is nil, no army is currently attacking this country.
	Invader *Army
}

//--------  GETTER  --------------------------------------------------------------------------------------------------//

// NeighborsObj retrieves a list of Country objects representing the neighboring countries.
// It uses the list of neighbor names stored in the Country struct to look up the corresponding
// Country objects in the World. This function provides an easy way to access the full Country
// structs of all neighboring countries. This allows the game logic to interact with neighboring
// countries for movement, attacks, or other strategic decisions.
func (c *Country) NeighborsObj() []*Country {
	ret := make([]*Country, 0, len(c.Neighbors))
	for _, n := range c.Neighbors {
		ret = append(ret, c.world.Country(n))
	}
	return ret
}

// ContinentObj retrieves the Continent object associated with this country.
// It uses the Continent name stored in the Country struct to look up the corresponding
// Continent object in the World. This function provides an easy way to access the full
// Continent struct to which this country belongs.
func (c *Country) ContinentObj() *Continent {
	return c.world.Continent(c.Continent)
}
