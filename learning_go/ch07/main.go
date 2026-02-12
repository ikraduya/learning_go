package main

import (
	"cmp"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
)

type Team struct {
	teamName    string
	playerNames []string
}

type League struct {
	Teams []Team
	Wins  map[string]int
}

func (l *League) MatchResult(team1 string, team1Score int, team2 string, team2Score int) {
	if team1Score > team2Score {
		l.Wins[team1] += 1
	} else if team2Score > team1Score {
		l.Wins[team2] += 1
	}
}

func (l *League) Ranking() []string {
	rankings := slices.SortedStableFunc(maps.Keys(l.Wins), func(a, b string) int {
		return cmp.Compare(l.Wins[b], l.Wins[a]) // descending
	})
	return rankings
}

type Ranker interface {
	Ranking() []string
}

func RankPrinter(r Ranker, writer io.Writer) {
	for _, v := range r.Ranking() {
		io.WriteString(writer, v+"\n")
	}
}

func main() {
	teams := []Team{
		{teamName: "Gophers", playerNames: []string{"Ava", "Ben", "Chloe"}},
		{teamName: "Rockets", playerNames: []string{"Diego", "Ella", "Finn"}},
		{teamName: "Tigers", playerNames: []string{"Grace", "Hugo", "Ivy"}},
		{teamName: "Wolves", playerNames: []string{"Jade", "Kai", "Liam"}},
	}

	league := League{
		Teams: teams,
		Wins:  map[string]int{},
	}

	for _, t := range teams {
		league.Wins[t.teamName] = 0
	}

	league.MatchResult("Gophers", 3, "Rockets", 1)
	league.MatchResult("Tigers", 2, "Wolves", 2) // draw
	league.MatchResult("Rockets", 4, "Tigers", 2)
	league.MatchResult("Wolves", 1, "Gophers", 2)
	league.MatchResult("Wolves", 3, "Rockets", 0)

	rankings := league.Ranking()
	for i, name := range rankings {
		fmt.Println(i+1, name, league.Wins[name])
	}

	var writer io.Writer = os.Stdout
	RankPrinter(&league, writer)
}
