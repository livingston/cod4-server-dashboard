package parser

import (
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"

	"github.com/livingston/cod4-server-dashboard/utils"
	"golang.org/x/net/html/charset"
)

// Server properties
type Server struct {
	XMLName       xml.Name `xml:"B3Status"`
	Game          metaData `xml:"Game"`
	Time          string   `xml:"Time,attr"`
	Players       []Player `xml:"Clients>Client"`
	FormattedName template.HTML
	Teams         map[int]*team
}

type metaData struct {
	XMLName  xml.Name    `xml:"Game"`
	MetaData []metaDatum `xml:"Data"`
}

type metaDatum struct {
	XMLName xml.Name `xml:"Data"`
	Key     string   `xml:"Name,attr"`
	Value   string   `xml:"Value,attr"`
}

// Players is a List of players
type Players map[int]*Player

type team struct {
	TeamName     template.HTML
	Players      Players
	TotalPlayers int
}

// Player properties
type Player struct {
	XMLName   xml.Name `xml:"Client"`
	ID        int      `xml:"CID,attr"`
	ColorName string   `xml:"ColorName,attr"`
	IP        string   `xml:"IP,attr"`
	GUID      string   `xml:"PBID,attr"`
	Score     int      `xml:"Score,attr"`
	Kills     int      `xml:"Kills,attr"`
	Deaths    int      `xml:"Deaths,attr"`
	Assists   int      `xml:"Assists,attr"`
	Ping      int      `xml:"Ping,attr"`
	Team      int      `xml:"Team,attr"`
	TeamName  string   `xml:"TeamName,attr"`
	Rank      int      `xml:"rank,attr"`
	Power     int      `xml:"power,attr"`
	Updated   string   `xml:"Updated,attr"`
}

// Len - Interface method for sort
func (players Players) Len() int {
	return len(players)
}

// Less - Interface method for sort
func (players Players) Less(i, j int) bool {
	return players[i].Score < players[j].Score
}

// Swap - Interface method for sort
func (players Players) Swap(i, j int) {
	players[i], players[j] = players[j], players[i]
}

// RankTitle - returns the player's rank title
func (p Player) RankTitle() string {
	return utils.GetRankTitle(p.Rank)
}

// Name of the player colorised
func (p Player) Name() template.HTML {
	return template.HTML(utils.Colorize(p.ColorName))
}

// Parse the serverstatus.xml file
func Parse(file string) (map[string]string, Server, error) {
	var server Server

	xmlFileFilePath, err := filepath.Abs(file)
	if err != nil {
		fmt.Println(err)

		return nil, server, errors.New("Missing " + file)
	}

	xmlFile, err := os.Open(xmlFileFilePath)
	if err != nil {
		return nil, server, errors.New("Unable to open " + file)
	}

	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	decoder.CharsetReader = charset.NewReaderLabel

	if err := decoder.Decode(&server); err != nil {
		return nil, server, err
	}

	parsedData := make(map[string]string)

	for _, metaData := range server.Game.MetaData {
		parsedData[metaData.Key] = metaData.Value
	}

	server.Teams = make(map[int]*team)
	playerIndices := make(map[int]int)

	for _, player := range server.Players {
		switch player.TeamName {
		case "Connecting...":
			player.Team = 4
		case "Loading...":
			player.Team = 5
		case "Free":
			player.Team = 6
			player.TeamName = "No team"
		}

		_, f := server.Teams[player.Team]
		if !f {
			server.Teams[player.Team] = &team{}
			server.Teams[player.Team].TeamName = template.HTML(utils.Colorize(player.TeamName))
			server.Teams[player.Team].Players = make(map[int]*Player)

			playerIndices[player.Team] = 0
		}

		currentPlayer := player
		currentPlayerIndex := playerIndices[player.Team]

		server.Teams[player.Team].Players[currentPlayerIndex] = &currentPlayer

		playerIndices[player.Team]++
	}

	for _, team := range server.Teams {
		team.TotalPlayers = len(team.Players)

		sort.Sort(sort.Reverse(team.Players))
	}

	server.FormattedName = template.HTML(utils.Colorize(parsedData["sv_hostname"]))
	parsedData["sv_hostname"] = utils.StripFormat(parsedData["sv_hostname"])

	return parsedData, server, err
}
