package core

import (
	"slices"
	"testing"
)

func TestNewWorld(t *testing.T) {
	w := NewWorld()

	// test Players
	if w.PlayerQueue == nil || len(w.PlayerQueue) != 0 {
		t.Fatalf("invalid Players")
	}

	// test random
	if w.rnd == nil {
		t.Fatalf("invalid random")
	}

	// test lock
	w.lock.Lock()
	w.lock.Unlock()

	// test round
	if w.Round != 0 || w.SubRound != 0 {
		t.Fatalf("invalid round")
	}

	// test Continents
	if len(w.Continents) != 6 {
		t.Fatalf("invalid Continents")
	}
	for k, v := range w.Continents {
		// key/name link
		if k != v.Name {
			t.Fatalf("Continents: invalid key,value.Name")
		}
		// map key
		if len(k) < 4 || len(k) > 19 {
			t.Fatalf("Continents: invalid key")
		}
		// points
		if v.Points < 2 || v.Points > 10 {
			t.Fatalf("Continents: invalid Points")
		}
		// Countries
		if len(v.Countries) < 3 || len(v.Countries) > 12 {
			t.Fatalf("Continents: invalid Countries: %s=%d", v.Name, len(v.Countries))
		}
		// check Countries link
		for _, c := range v.Countries {
			if _, ok := w.Countries[c]; !ok {
				t.Fatalf("Continents: no link for country: %s", c)
			}
		}
	}

	// test Countries
	if len(w.Countries) != 42 {
		t.Fatalf("invalid Countries: %d", len(w.Countries))
	}
	for k, v := range w.Countries {
		// check world
		if v.world == nil {
			t.Fatalf("world is nil: %s", k)
		}
		// key/name link
		if k != v.Name {
			t.Fatalf("Countries: invalid key,value.Name")
		}
		// map key
		if len(k) < 4 || len(k) > 19 {
			t.Fatalf("Countries: invalid key: %d", len(k))
		}
		// fortress
		if v.FortressRegion && v.BorderRegion {
			t.Fatalf("Countries: invalid FortressRegion: BorderRegion is true!: %s", k)
		}
		// RecruitingRegion
		if v.FortressRegion && !v.RecruitingRegion {
			t.Fatalf("Countries: invalid RecruitingRegion: FortressRegion is true and RecruitingRegion not!: %s", k)
		}
		if v.BorderRegion && v.RecruitingRegion {
			t.Fatalf("Countries: invalid RecruitingRegion: BorderRegion and RecruitingRegion are true!: %s", k)
		}
		if !v.FortressRegion && !v.BorderRegion && !v.RecruitingRegion {
			t.Errorf("Countries: should be RecruitingRegion: %s", k)
		}
		// Occupier
		if v.Occupier != nil {
			t.Fatalf("Countries: invalid Occupier")
		}
		//Position
		if v.Position[0] < 1 || v.Position[1] < 1 {
			t.Fatalf("Countries: invalid Position")
		}
		// Neighbors
		if len(v.Neighbors) < 2 {
			t.Fatalf("Countries: invalid Neighbors")
		}
		// check Continent link
		if _, ok := w.Continents[v.Continent]; !ok {
			t.Fatalf("Countries: no link for Continent")
		}
		// check Neighbors link
		for _, n := range v.Neighbors {
			if _, ok := w.Countries[n]; !ok {
				t.Fatalf("Countries: no link for Neighbor: %s=%s", v.Name, n)
			}
		}
		// check Neighbors backlink
		for _, n := range v.Neighbors {
			c := w.Countries[n]
			if !slices.Contains(c.Neighbors, k) {
				t.Fatalf("Countries: no back link: 'test country'='%s', 'neighbor'='%s', 'neighbor country'='%s', 'neighbor country neighbors'=%v", k, n, c.Name, c.Neighbors)
			}
		}
	}

}
