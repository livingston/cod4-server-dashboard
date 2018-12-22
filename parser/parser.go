package parser

import (
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

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

type team struct {
	TeamName template.HTML
	Players  map[int]*Player
}

// Player properties
type Player struct {
	XMLName  xml.Name `xml:"Client"`
	ID       int      `xml:"CID,attr"`
	Name     string   `xml:"ColorName,attr"`
	IP       string   `xml:"IP,attr"`
	GUID     string   `xml:"PBID,attr"`
	Score    int      `xml:"Score,attr"`
	Kills    int      `xml:"Kills,attr"`
	Deaths   int      `xml:"Deaths,attr"`
	Assists  int      `xml:"Assists,attr"`
	Ping     int      `xml:"Ping,attr"`
	Team     int      `xml:"Team,attr"`
	TeamName string   `xml:"TeamName,attr"`
	Rank     int      `xml:"rank,attr"`
	Power    int      `xml:"power,attr"`
	Updated  string   `xml:"Updated,attr"`
}

// RankText - returns the player's rank title
func (p Player) RankText() string {
	return utils.GetRankText(p.Rank)
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

	for i := 0; i < len(server.Game.MetaData); i++ {
		parsedData[server.Game.MetaData[i].Key] = server.Game.MetaData[i].Value
	}

	server.Teams = make(map[int]*team)

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
		}

		currentPlayer := player
		server.Teams[player.Team].Players[player.ID] = &currentPlayer
	}

	server.FormattedName = template.HTML(utils.Colorize(parsedData["sv_hostname"]))
	parsedData["sv_hostname"] = utils.StripFormat(parsedData["sv_hostname"])

	return parsedData, server, err
}
