package core

import (
	"fmt"
	"image/color"
	"testing"
)

func TestNewArmy_panic(t *testing.T) {
	// test nil -> panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	// The following is the code under test
	_ = NewArmy(nil, 12, "Player", "Base")
}

func TestNewArmy(t *testing.T) {
	// test A
	w := NewWorld()
	s := 12
	p := "Player"
	h := "Base"
	a := NewArmy(w, s, p, h)
	if a == nil || a.world != w || a.Strength != s || a.Player != p || a.HomeBase != h {
		t.Fatalf("wrong init")
	}

	// test B
	w = NewWorld()
	s = -1
	p = ""
	h = ""
	a = NewArmy(w, s, p, h)
	if a == nil || a.world != w || a.Strength != s || a.Player != p || a.HomeBase != h {
		t.Fatalf("wrong init")
	}
}

func TestPlayerObj(t *testing.T) {
	// world
	w := NewWorld()
	if err := w.AddPlayer("Player1", color.RGBA{R: 255, G: 255, B: 255, A: 255}); err != nil {
		t.Fatalf(err.Error())
	}

	// test A (not found -> black)
	a := NewArmy(w, 13, "Player99", "any")
	if po := a.PlayerObj(); po == nil || po.Name != "Player99" || po.Color.R != 0 || po.Color.G != 0 || po.Color.B != 0 || po.Color.A != 0 {
		t.Fatalf("wrong")
	}

	// test B (found -> white)
	a = NewArmy(w, 13, "Player1", "any")
	if po := a.PlayerObj(); po == nil || po.Name != "Player1" || po.Color.R != 255 || po.Color.G != 255 || po.Color.B != 255 || po.Color.A != 255 {
		t.Fatalf("wrong")
	}
}

func TestHomeBaseObj(t *testing.T) {
	// world
	w := NewWorld()
	if err := w.AddPlayer("Player1", color.RGBA{R: 255, G: 255, B: 255, A: 255}); err != nil {
		t.Fatalf(err.Error())
	}

	// test A (not found)
	a := NewArmy(w, 13, "Player99", "any")
	if c := a.HomeBaseObj(); c == nil || c.Name != "any" || len(c.Neighbors) != 0 {
		t.Fatalf("wrong")
	}

	// test B (found)
	a = NewArmy(w, 13, "Player1", "Alaska")
	if c := a.HomeBaseObj(); c == nil || c.Name != "Alaska" || len(c.Neighbors) == 0 {
		t.Fatalf("wrong")
	}
}

func TestDescription(t *testing.T) {
	a := NewArmy(NewWorld(), 1, "PLAYER", "BASE")
	d := a.Description()
	if d != "PLAYER's BASE Army with 1 men" {
		t.Fatalf("wrong: %s", d)
	}
}

func Test_rollDice(t *testing.T) {
	// get random from world
	world := NewWorld()
	rnd := world.rnd

	// check input
	if d := rollDice(nil, 1); d == nil || len(d) != 0 {
		t.Fatalf("ERROR: wrong dice")
	}
	if d := rollDice(rnd, 1); d == nil || len(d) != 1 || d[0] == 0 {
		t.Fatalf("ERROR: wrong dice")
	}
	if d := rollDice(rnd, 0); d == nil || len(d) != 0 {
		t.Fatalf("ERROR: wrong dice")
	}
	if d := rollDice(rnd, 2); d == nil || len(d) != 2 || d[1] == 0 {
		t.Fatalf("ERROR: wrong dice")
	}
	if d := rollDice(rnd, 6); d == nil || len(d) != 6 || d[5] == 0 || d[3] == 0 {
		t.Fatalf("ERROR: wrong dice")
	}

	// Akt
	results := make([]int, 6)
	for i := 0; i < 100000; i++ {
		d := rollDice(rnd, 1)
		results[d[0]-1]++
	}
	// check dice
	for i := 0; i < 6; i++ {
		r := results[i]
		t.Logf("LOG: dice %d: %d", i+1, r)
		if r < 100000/7 {
			t.Errorf("ERROR: wrong stat: %d=%d", i, r)
		}
	}
}

func Test_minInt(t *testing.T) {
	if x := minInt(3, 5); x != 3 {
		t.Fatalf("wrong min")
	}
	if x := minInt(5, 3); x != 3 {
		t.Fatalf("wrong min")
	}
	if x := minInt(3, 3); x != 3 {
		t.Fatalf("wrong min")
	}
	if x := minInt(-3, 3); x != -3 {
		t.Fatalf("wrong min")
	}
	if x := minInt(0, 0); x != 0 {
		t.Fatalf("wrong min")
	}
}

func TestAttack(t *testing.T) {
	w := NewWorld()

	// no log
	att := NewArmy(w, 30, "Attacker", "AttBase")
	def := NewArmy(w, 10, "Defender", "DefBase")
	if log := att.Attack(def, true); len(log) != 0 {
		t.Fatalf("wrong log")
	}

	// log
	att = NewArmy(w, 30, "Attacker", "AttBase")
	def = NewArmy(w, 10, "Defender", "DefBase")
	if log := att.Attack(def, false); len(log) < 10 {
		t.Fatalf("wrong log")
	}

	// no battle
	att = NewArmy(w, 30, "Attacker", "AttBase")
	def = nil
	if log := att.Attack(def, false); log[0] != "Not all armies were ready to fight. There was no battle." {
		t.Fatalf("wrong log")
	}
	att = NewArmy(w, 30, "Attacker", "AttBase")
	def = NewArmy(w, 0, "Defender", "DefBase")
	if log := att.Attack(def, false); log[0] != "Not all armies were ready to fight. There was no battle." {
		t.Fatalf("wrong log")
	}
	att = NewArmy(w, 0, "Attacker", "AttBase")
	def = NewArmy(w, 10, "Defender", "DefBase")
	if log := att.Attack(def, false); log[0] != "Not all armies were ready to fight. There was no battle." {
		t.Fatalf("wrong log")
	}

	// too big battles
	att = NewArmy(w, 10000000, "Attacker", "AttBase")
	def = NewArmy(w, 10000000, "Defender", "DefBase")
	_ = att.Attack(def, true)
	if att.Strength > 0 && def.Strength > 0 {
		t.Fatalf("no figth to the end: %d vs %d", att.Strength, def.Strength)
	}

	// end loop attacker
	att = NewArmy(w, 10, "Attacker", "AttBase")
	def = NewArmy(w, 10000000, "Defender", "DefBase")
	_ = att.Attack(def, false)
	if att.Strength != 0 || def.Strength == 0 {
		t.Fatalf("wrong winner")
	}

	// end loop defender
	att = NewArmy(w, 10000000, "Attacker", "AttBase")
	def = NewArmy(w, 10, "Defender", "DefBase")
	_ = att.Attack(def, false)
	if def.Strength != 0 || att.Strength == 0 {
		t.Fatalf("wrong winner")
	}

	// fortress yes (def wins)
	att = NewArmy(w, 10, "Attacker", "AttBase")
	def = NewArmy(w, 10, "Defender", "Congo") // fortress
	_ = att.Attack(def, false)
	if att.Strength != 0 || def.Strength == 0 {
		t.Fatalf("!!RANDOM!!: wrong winner")
	}
}

func TestAttack_Table(t *testing.T) {
	w := NewWorld()

	// print table
	{ //-----------------------------------------

		// print header
		fmt.Printf("            \t")
		for defStr := 0; defStr < 22; defStr++ {
			fmt.Printf("%3.0f \t", float64(defStr))
		}
		fmt.Printf("\n")

		// print all tables
		defHomeBase := "CountryB"
		for nn := 0; nn < 2; nn++ {

			for attStr := 0; attStr < 33; attStr++ {
				fmt.Printf("Att. %d v. X:\t", attStr)
				for defStr := 0; defStr < 22; defStr++ {
					winA := 0
					winD := 0
					winN := 0
					for bat := 0; bat < 9999; bat++ {

						// config attacker
						attacker := NewArmy(w, attStr, "Attacker", "CountryA")
						if attacker == nil {
							t.Fatalf("attacker is nil")
						}
						// config defender
						defender := NewArmy(w, defStr, "Defender", defHomeBase)
						if defender == nil {
							t.Fatalf("defender is nil")
						}
						// simulate battle
						attacker.Attack(defender, true)

						// ----  checks -------------------------

						// invalid strength
						if attacker.Strength < 0 || defender.Strength < 0 {
							t.Fatalf("invalid strength")
						}
						if attacker.Strength > 0 && defender.Strength > 0 {
							t.Fatalf("invalid strength")
						}

						// stats
						if attacker.Strength <= 0 && defender.Strength <= 0 {
							winN++
						} else if attacker.Strength == 0 {
							winD++
						} else if defender.Strength == 0 {
							winA++
						}

					} // -----------------------------------------
					winARate := float64(winA) / float64(winA+winD+winN) * 100
					fmt.Printf("%3.0f%%\t", winARate)
				}
				fmt.Printf("\n")
			}
			defHomeBase = "Congo"
			fmt.Printf("\n\n\n")
		}
	} //-----------------------------------------
}
