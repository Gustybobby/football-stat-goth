package repos

import (
	"football-stat-goth/models"
	"maps"
	"slices"
	"sort"
)

func FindClubsWithNameAsc(repo *Repository) ([]models.Club, error) {
	var clubs []models.Club
	results := repo.DB.Select("ClubID", "Name", "Logo").Order("name ASC").Find(&clubs)
	if results.Error != nil {
		return nil, results.Error
	}
	return clubs, nil
}

func FindClubs(repo *Repository) ([]ClubStanding, error) {
	var clubs []models.Club
	clubResults := repo.DB.Find(&clubs)
	if clubResults.Error != nil {
		return nil, clubResults.Error
	}

	var matches []models.Match
	matchResults := repo.DB.Preload("HomeLineup").Preload("AwayLineup").Where(
		"is_finished = true").Find(&matches)
	if matchResults.Error != nil {
		return nil, matchResults.Error
	}

	var standingsMap = make(map[string]ClubStanding)
	for _, club := range clubs {
		standingsMap[club.ClubID] = ClubStanding{
			ClubID: club.ClubID,
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
		var homeID = match.HomeLineup.ClubID
		var awayID = match.AwayLineup.ClubID

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
		return a.Won*3+a.Drawn*1+a.Lost*0 > b.Won*3+b.Drawn*1+b.Lost*0
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
