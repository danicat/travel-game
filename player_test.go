package main

import (
	"fmt"
	"log"
	"testing"
)

func TestPlay(t *testing.T) {
	LoadCards("cards.json")
	tbl := []struct {
		name            string
		beforeCards     []string
		afterStatus     Status
		afterDistance   int
		afterTerrain    string
		afterImmunities []Status
	}{
		{
			"should be lost",
			[]string{},
			StatusLost,
			0,
			"",
			nil,
		},
		// Green card effects
		{
			"should be oriented",
			[]string{"O"},
			StatusOriented,
			0,
			"",
			nil,
		},
		// Red card effects
		{
			"should be lost",
			[]string{"O", "P"},
			StatusLost,
			0,
			"",
			nil,
		},
		{
			"should be captive",
			[]string{"O", "PH"},
			StatusCaptive,
			0,
			"",
			nil,
		},
		{
			"should be penniless",
			[]string{"O", "FD"},
			StatusPenniless,
			0,
			"",
			nil,
		},
		{
			"should be sick",
			[]string{"O", "E"},
			StatusSick,
			0,
			"",
			nil,
		},
		// Blue card status cancel
		{
			"should be oriented (RA)",
			[]string{"RA"},
			StatusOriented,
			0,
			"",
			nil,
		},
		{
			"should be oriented (RA)",
			[]string{"O", "P", "RA"},
			StatusOriented,
			0,
			"",
			nil,
		},
		// Blue card immunities
		{
			"should be immune to lost",
			[]string{"O", "RA", "P"},
			StatusOriented,
			0,
			"",
			nil,
		},
		{
			"should be immune to captive",
			[]string{"O", "D", "PH"},
			StatusOriented,
			0,
			"",
			nil,
		},
		{
			"should be immune to penniless",
			[]string{"O", "RQ", "FD"},
			StatusOriented,
			0,
			"",
			nil,
		},
		{
			"should be immune to sick",
			[]string{"O", "S", "E"},
			StatusOriented,
			0,
			"",
			nil,
		},
	}

	for _, tc := range tbl {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			p := Player{Name: "testplayer"}
			for _, c := range tc.beforeCards {
				card := FindCardByID(c)
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
