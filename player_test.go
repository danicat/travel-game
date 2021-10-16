package main

import (
	"fmt"
	"log"
	"testing"
)

func TestTravel(t *testing.T) {
	tbl := []struct {
		status           Status
		startDistance    int
		startTerrain     string
		cardName         string
		cardDistance     int
		terrains         []string
		expectedDistance int
	}{
		{
			StatusOriented,
			0,
			"desert",
			"1000",
			1000,
			[]string{"desert"},
			1000,
		},
		{
			StatusOriented,
			0,
			"sea",
			"1000",
			1000,
			[]string{"desert"},
			0,
		},
		{
			StatusEscaping,
			0,
			"desert",
			"1000",
			1000,
			[]string{"desert"},
			0,
		},
		{
			StatusOriented,
			0,
			"desert",
			"2000",
			2000,
			[]string{"savage_land", "civilization"},
			0,
		},
		{
			StatusLost,
			0,
			"civilization",
			"4000",
			4000,
			[]string{"civilization"},
			0,
		},
	}

	for _, tc := range tbl {
		t.Run(fmt.Sprintf("status %s terrain %s travel distance %d", tc.status, tc.startTerrain, tc.expectedDistance), func(t *testing.T) {
			p := Player{Status: tc.status}

			terrainCard := Card{
				Type: TypeYellow,
				Effects: Effects{
					Terrain: tc.startTerrain,
				},
			}

			travelCard := Card{
				Name: tc.cardName,
				Type: TypeWhite,
				Effects: Effects{
					Distance: tc.cardDistance,
				},
				Constraints: Constraints{
					Status:   []Status{StatusOriented},
					Terrains: tc.terrains,
				},
			}

			log.Println(tc.status.String())

			err := p.Receive(&p, terrainCard)
			if err != nil {
				t.Fatalf("expected no errors, got %s", err)
			}

			err = p.Receive(&p, travelCard)
			log.Println(err)

			if p.Distance != tc.expectedDistance {
				t.Fatalf("expected %d distance, got %d", tc.expectedDistance, p.Distance)
			}
		})
	}
}
