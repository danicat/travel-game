package main

import (
	"fmt"
	"log"
	"testing"
)

func TestPlay(t *testing.T) {
	LoadConfig("config.json")
	LoadCards("cards.json")
	tbl := []struct {
		name          string
		beforeCards   []string
		afterStatus   Status
		afterDistance int
		afterTerrain  string
	}{
		// No cards
		{
			"should be lost",
			[]string{},
			StatusLost,
			0,
			"",
		},
		// Green card effects
		{
			"should be oriented",
			[]string{"O"},
			StatusOriented,
			0,
			"",
		},
		{
			"should be oriented",
			[]string{"O", "P", "O"},
			StatusOriented,
			0,
			"",
		},
		{
			"should be escaping",
			[]string{"O", "PH", "F"},
			StatusEscaping,
			0,
			"",
		},
		{
			"should be working",
			[]string{"O", "FD", "T"},
			StatusWorking,
			0,
			"",
		},
		{
			"should be healing",
			[]string{"O", "E", "R"},
			StatusHealing,
			0,
			"",
		},
		// Red card effects
		{
			"should be lost",
			[]string{"O", "P"},
			StatusLost,
			0,
			"",
		},
		{
			"should be captive",
			[]string{"O", "PH"},
			StatusCaptive,
			0,
			"",
		},
		{
			"should be penniless",
			[]string{"O", "FD"},
			StatusPenniless,
			0,
			"",
		},
		{
			"should be sick",
			[]string{"O", "E"},
			StatusSick,
			0,
			"",
		},
		// Should be oriented (RA)
		{
			"should be oriented (RA)",
			[]string{"RA", "PH", "F"},
			StatusOriented,
			0,
			"",
		},
		{
			"should be oriented (RA)",
			[]string{"RA", "FD", "T"},
			StatusOriented,
			0,
			"",
		},
		{
			"should be oriented (RA)",
			[]string{"RA", "E", "R"},
			StatusOriented,
			0,
			"",
		},
		// Blue card status cancel
		{
			"should be oriented (RA)",
			[]string{"RA"},
			StatusOriented,
			0,
			"",
		},
		{
			"should be oriented (RA)",
			[]string{"O", "P", "RA"},
			StatusOriented,
			0,
			"",
		},
		// Blue card immunities
		{
			"should be immune to lost",
			[]string{"O", "RA", "P"},
			StatusOriented,
			0,
			"",
		},
		{
			"should be immune to captive",
			[]string{"O", "D", "PH"},
			StatusOriented,
			0,
			"",
		},
		{
			"should be immune to penniless",
			[]string{"O", "RQ", "FD"},
			StatusOriented,
			0,
			"",
		},
		{
			"should be immune to sick",
			[]string{"O", "S", "E"},
			StatusOriented,
			0,
			"",
		},
		// Yellow cards
		{
			"should be desert",
			[]string{"SR"},
			StatusLost,
			0,
			"desert",
		},
		{
			"should be desert",
			[]string{"TC", "SR"},
			StatusLost,
			0,
			"desert",
		},
		{
			"should be civilization",
			[]string{"SR", "TC"},
			StatusLost,
			0,
			"civilization",
		},
		{
			"should be savage land",
			[]string{"SR", "TC", "TS"},
			StatusLost,
			0,
			"savage_land",
		},
		{
			"should be sea",
			[]string{"SR", "TC", "TS", "M"},
			StatusLost,
			0,
			"sea",
		},
		{
			"should be no terrain (RA)",
			[]string{"SR", "TC", "TS", "M", "RA"},
			StatusOriented,
			0,
			"",
		},
		{
			"should be no terrain (RA)",
			[]string{"RA", "SR", "TC", "TS", "M"},
			StatusOriented,
			0,
			"",
		},
		// White cards
		{
			"1000",
			[]string{"O", "1000"},
			StatusOriented,
			1000,
			"",
		},
		{
			"1000",
			[]string{"O", "SR", "1000"},
			StatusOriented,
			1000,
			"",
		},
		{
			"1000",
			[]string{"O", "TS", "1000"},
			StatusOriented,
			1000,
			"",
		},
		{
			"1000",
			[]string{"O", "TC", "1000"},
			StatusOriented,
			1000,
			"",
		},
		{
			"0",
			[]string{"O", "M", "1000"},
			StatusOriented,
			0,
			"",
		},
		{
			"0",
			[]string{"O", "SR", "P", "1000"},
			StatusLost,
			0,
			"",
		},
		{
			"1000",
			[]string{"O", "M", "P", "RA", "1000"},
			StatusOriented,
			1000,
			"",
		},
	}

	for _, tc := range tbl {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			p := NewPlayer(0, "testplayer", 0, 0)
			for _, c := range tc.beforeCards {
				card := FindCardByID(c)
				if card == nil {
					t.Fatalf("card not found %s", c)
				}

				err := p.Receive(nil, *card)
				if err != nil {
					log.Println(err)
				}
			}

			if p.Status() != tc.afterStatus {
				t.Errorf("status should be %s, got %s", tc.afterStatus, p.Status())
			}

			if p.Distance != tc.afterDistance {
				t.Errorf("distance should be %d, got %d", tc.afterDistance, p.Distance)
			}
		})
	}
}
