package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"

	soccer "github.com/stein-f/oink-soccer-common"
)

// go run cmd/simulate/main.go
// my_team_y
// my_team_diamond
// my_team_pyramid
//
//go:embed other_team_saturn.json
var homeTeamConfig []byte

//go:embed my_team_diamond.json
var awayTeamConfig []byte

func main() {
	var homeWins, awayWins, draws, goals, homeChances, awayChances int
	gameCount := 10000

	scorerByPosition := make(map[soccer.PlayerPosition]int)

	homeLineup := loadConfig(homeTeamConfig)
	// homeLineup.ItemBoosts = []soccer.Boost{
	// 	{BoostType: soccer.BoostTypeTeam, MinBoost: 1.01, MaxBoost: 1.05}, // apply a 1-5% boost to the team
	// }
	awayLineup := loadConfig(awayTeamConfig)
	// awayLineup.ItemBoosts = []soccer.Boost{
	// 	{BoostType: soccer.BoostTypeTeam, MinBoost: 1.03, MaxBoost: 1.07}, // apply a 3-7% boost to the team
	// }

	for i := 0; i < gameCount; i++ {
		gameEvents, err := soccer.RunGame(homeLineup, awayLineup)
		if err != nil {
			panic(err)
		}

		gameStats := soccer.CreateGameStats(gameEvents)

		if gameStats.HomeTeamStats.Goals > gameStats.AwayTeamStats.Goals {
			homeWins++
		} else if gameStats.HomeTeamStats.Goals < gameStats.AwayTeamStats.Goals {
			awayWins++
		} else {
			draws++
		}

		for _, event := range gameEvents {
			if event.Type == soccer.GameEventTypeGoal {
				scorerID := event.Event.(soccer.GoalEvent).PlayerID
				homeScorer, homeFound := homeLineup.FindPlayer(scorerID)
				awayScorer, awayFound := awayLineup.FindPlayer(scorerID)
				if !homeFound && !awayFound {
					panic(fmt.Sprintf("scorer %s not found", scorerID))
				}
				if homeFound {
					scorerByPosition[homeScorer.SelectedPosition]++
					continue
				}
				scorerByPosition[awayScorer.SelectedPosition]++
			}
		}

		goals += gameStats.HomeTeamStats.Goals + gameStats.AwayTeamStats.Goals
		homeChances += gameStats.HomeTeamStats.Shots
		awayChances += gameStats.AwayTeamStats.Shots
	}

	goalsPerGame := float64(goals) / float64(gameCount)
	homeTeamChancePerGame := float64(homeChances) / float64(gameCount)
	awayTeamChancePerGame := float64(awayChances) / float64(gameCount)

	fmt.Printf("\nGame summary:\n")
	fmt.Printf("Games played: %d\n", gameCount)
	fmt.Printf("Home Team wins: %d\n", homeWins)
	fmt.Printf("Home Team chances/game: %f\n", homeTeamChancePerGame)
	fmt.Printf("Away Team wins: %d\n", awayWins)
	fmt.Printf("Away Team chances/game: %f\n", awayTeamChancePerGame)
	fmt.Printf("Draws: %d\n", draws)
	fmt.Printf("Goals/game: %f\n", goalsPerGame)

	homeWinPercent := float64(homeWins) / float64(gameCount) * 100
	awayWinPercent := float64(awayWins) / float64(gameCount) * 100
	drawPercent := float64(draws) / float64(gameCount) * 100
	fmt.Printf("Home Team Win Percentage: %f%% \n", homeWinPercent)
	fmt.Printf("Away Team Win Percentage: %f%% \n", awayWinPercent)
	fmt.Printf("Draw Percentage: %f%% \n", drawPercent)

	attackerGoals := scorerByPosition[soccer.PlayerPositionAttack]
	totalGoals := scorerByPosition[soccer.PlayerPositionAttack] + scorerByPosition[soccer.PlayerPositionMidfield] + scorerByPosition[soccer.PlayerPositionDefense] + scorerByPosition[soccer.PlayerPositionGoalkeeper]
	attackerGoalsPercentage := float64(attackerGoals) / float64(totalGoals) * 100
	fmt.Printf("Attacker goals: %d (%f%%)\n", attackerGoals, attackerGoalsPercentage)

	midfielderGoals := scorerByPosition[soccer.PlayerPositionMidfield]
	midfielderGoalsPercentage := float64(midfielderGoals) / float64(totalGoals) * 100
	fmt.Printf("Midfielder goals: %d (%f%%)\n", midfielderGoals, midfielderGoalsPercentage)

	defenderGoals := scorerByPosition[soccer.PlayerPositionDefense]
	defenderGoalsPercentage := float64(defenderGoals) / float64(totalGoals) * 100
	fmt.Printf("Defender goals: %d (%f%%)\n", defenderGoals, defenderGoalsPercentage)

	goalkeeperGoals := scorerByPosition[soccer.PlayerPositionGoalkeeper]
	goalkeeperGoalsPercentage := float64(goalkeeperGoals) / float64(totalGoals) * 100
	fmt.Printf("Goalkeeper goals: %d (%f%%)\n", goalkeeperGoals, goalkeeperGoalsPercentage)
}

func loadConfig(config []byte) soccer.GameLineup {
	var lineup soccer.GameLineup
	if err := json.Unmarshal(config, &lineup); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Lineup: ", lineup)
	return lineup
}
