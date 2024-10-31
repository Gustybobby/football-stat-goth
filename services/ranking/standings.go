package ranking

import (
	"context"
	"football-stat-goth/queries"
	"maps"
	"slices"
	"sort"
)

func FindClubStandings(db *queries.Queries, ctx context.Context) ([]ClubStanding, error) {
	clubs, err := db.ListClubsOrderByNameAsc(ctx)
	if err != nil {
		return nil, err
	}

	matches, err := db.ListMatchesWithClubsAndGoals(ctx, queries.ListMatchesWithClubsAndGoalsParams{
		FilterClubID: false,
		ClubID:       "",
		IsFinished:   true,
		Order:        "DESC",
	})
	if err != nil {
		return nil, err
	}

	var standingsMap = make(map[string]ClubStanding)
	for _, club := range clubs {
		standingsMap[club.ID] = ClubStanding{
			ClubID: club.ID,
			Name:   club.Name,
			Logo:   club.Logo,
			Games:  0,
			Won:    0,
			Drawn:  0,
			Lost:   0,
			GF:     0,
			GA:     0,
		}
	}
	for _, match := range matches {
		var homeID = match.HomeClubID
		var awayID = match.AwayClubID

		homeStanding, homeOK := standingsMap[homeID]
		awayStanding, awayOK := standingsMap[awayID]

		if homeOK && awayOK {
			homeStanding.Games = homeStanding.Games + 1
			awayStanding.Games += 1

			if match.HomeGoals > match.AwayGoals {
				homeStanding.Won += 1
				awayStanding.Lost += 1
			} else if match.HomeGoals < match.AwayGoals {
				homeStanding.Lost += 1
				awayStanding.Won += 1
			} else {
				homeStanding.Drawn += 1
				awayStanding.Drawn += 1
			}

			homeStanding.GF += int(match.HomeGoals)
			homeStanding.GA += int(match.AwayGoals)

			awayStanding.GF += int(match.AwayGoals)
			awayStanding.GA += int(match.HomeGoals)

			standingsMap[homeID] = homeStanding
			standingsMap[awayID] = awayStanding
		}
	}

	standings := slices.Collect(maps.Values(standingsMap))
	sort.SliceStable(standings, func(i int, j int) bool {
		var a = standings[i]
		var b = standings[j]
		var aPoints = a.Won*3 + a.Drawn*1 + a.Lost*0
		var bPoints = b.Won*3 + b.Drawn*1 + b.Lost*0
		if aPoints == bPoints {
			var aGD = a.GF - a.GA
			var bGD = b.GF - b.GA
			if aGD == bGD {
				return a.GF > b.GF
			}
			return aGD > bGD
		}
		return aPoints > bPoints
	})

	return standings, nil
}

type ClubStanding struct {
	ClubID string
	Name   string
	Logo   string
	Games  int
	Won    int
	Drawn  int
	Lost   int
	GF     int
	GA     int
}