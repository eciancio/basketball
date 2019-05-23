package main

import (
    "os"
	"sync"
    "fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
    "golang.org/x/net/html"
	"net/http"
	"strings"
)

var waitGroup sync.WaitGroup
var sem = make(chan int, 50)

const (
	user     = "wsa@wsabasketball"
	database = "basketball"
	host     = "wsabasketball.mysql.database.azure.com"
	password = "LeBron>MJ!"
)

type PlayerPerformance struct {
	bbrefid string
	dateID string
	stats map[string]string
}

func getPlayerID(bbref string, db *sql.DB) string {
    var playerID string
    selectPlayerID := "select playerID from player_reference where bbrefID=?"
    rows, err:= db.Query(selectPlayerID, bbref)
    if err != nil {
        fmt.Println(err)
    }
    defer rows.Close()
    for rows.Next() {
        rows.Scan(&playerID)
    }
    return playerID
}

func newPlayerPerformance() *PlayerPerformance {
	var player PlayerPerformance
	player.stats = map[string]string {
		"mp": "0",
		"fg": "0",
		"fga": "0",
		"fg_pct": "0",
		"fg3": "0",
		"fg3a": "0",
		"fg3_pct": "0",
		"ft": "0",
		"fta": "0",
		"ft_pct": "0",
		"orb": "0",
		"drb": "0",
		"trb": "0",
		"blk": "-1",
		"tov": "0",
		"pf": "0",
		"pts": "0",
		"plus_minus": "0",
		"ts_pct": "0",
		"efg_pct": "0",
		"fg3a_per_fga_pct": "0",
		"fta_per_fga_pct": "0",
		"orb_pct": "0",
		"drb_pct": "0",
		"trb_pct": "0",
		"ast_pct": "0",
		"stl_pct": "0",
		"blk_pct": "0",
		"tov_pct": "0",
		"usg_pct": "0",
		"off_rtg": "0",
		"def_rtg": "0",
        "team": "",
        "opp": "",
        "home": "",
        "triple_double": "0",
        "double_double": "0",
	}
	return &player
}

func (p *PlayerPerformance) addToTable(db *sql.DB, dateID string) {
	defer waitGroup.Done()
    if _, ok := p.stats["reason"]; ok {
        // fmt.Println(p.stats)
        <-sem
        return
    }
    insertPerformance := "INSERT INTO performance (points, minutesPlayed, fieldGoals, fieldGoalsAttempted, fieldGoalPercent, 3PM, 3PA, 3PPercent, FT, FTA, FTPercent, offensiveRebounds, defensiveRebounds, totalRebounds, assists,  steals, blocks, turnovers, personalFouls, plusMinus, trueShootingPercent, effectiveFieldGoalPercent, 3pointAttemptRate, freeThrowAttemptRate, offensiveReboundPercent, defensiveReboundPercent, totalReboundPercent, assistPercent, stealPercent, blockPercent, turnoverPercent, usagePercent, offensiveRating, defensiveRating,  tripleDouble, doubleDouble, team, opponent, home, playerID, dateID) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
    playerID := getPlayerID(p.bbrefid,db)
    _, err := db.Exec(insertPerformance, p.stats["pts"], p.stats["mp"], p.stats["fg"], p.stats["fga"], p.stats["fg_pct"], p.stats["fg3"], p.stats["fg3a"], p.stats["fg3_pct"], p.stats["ft"], p.stats["fta"], p.stats["ft_pct"], p.stats["orb"], p.stats["drb"], p.stats["trb"], p.stats["ast"], p.stats["stl"], p.stats["blk"], p.stats["tov"], p.stats["pf"], p.stats["plus_minus"], p.stats["ts_pct"], p.stats["efg_pct"], p.stats["fg3a_per_fga_pct"], p.stats["fta_per_fga_pct"], p.stats["orb_pct"], p.stats["drb_pct"], p.stats["trb_pct"], p.stats["ast_pct"], p.stats["stl_pct"], p.stats["blk_pct"], p.stats["tov_pct"], p.stats["usg_pct"], p.stats["off_rtg"], p.stats["def_rtg"], p.stats["triple_double"], p.stats["double_double"], p.stats["team"], p.stats["opp"], p.stats["home"], playerID, dateID)
    if err != nil {
        fmt.Println(err)
    }
    <-sem
}

func getBoxScoreUrls(dateID string, db *sql.DB) []string {
    var boxScore string
    urls := make([]string, 0)
    rows, err:= db.Query("Select url from box_score_urls WHERE dateID >= ? AND dateID <= ?", dateID, dateID)
    if err != nil {
        fmt.Println(err)
    }
    defer rows.Close()
    for rows.Next() {
        rows.Scan(&boxScore)
		urls = append(urls, boxScore)
    }
    fmt.Println(urls)
    return urls // retunrs urls 
}

func getBasicStats(z *html.Tokenizer, playerMap map[string]*PlayerPerformance, team string, home string, opp string) {
	for {
		tt:= z.Next()
		if tt == html.ErrorToken {
			return
		}
		if tt == html.EndTagToken {
			t := z.Token()
			isTable := t.Data == "table"
			if isTable {
				return
			}
		}
        if tt == html.StartTagToken {
            t := z.Token()
			isHeader := t.Data == "th"
			if isHeader {
				if(len(t.Attr) == 5){
                    var player *PlayerPerformance
                    if val, ok := playerMap[t.Attr[2].Val]; ok {
                        player = val
                    } else {
                        player = newPlayerPerformance()
                        player.bbrefid = t.Attr[2].Val
                        playerMap[t.Attr[2].Val] = player
                    }
                    player.stats["team"] = team
                    player.stats["opp"] = opp
                    player.stats["home"] = home
					// fmt.Println(t.Attr[2].Val)
			        z.Next()
                    t := z.Token()
                    for t.Data != "tr" {
                        if(t.Data == "td" && tt == html.StartTagToken) {
                            tt = z.Next()
                            text := (string)(z.Text())
                            stats := strings.TrimSpace(text)
                            // fmt.Println(t.Attr[1].Val, stats)
                            player.stats[t.Attr[1].Val] = stats
                        }
                        tt = z.Next()
                        t = z.Token()
			        }
				}

			}
		}
	}
}


func getTables(url string, db *sql.DB, dateID string) {
    playerMap := make(map[string]*PlayerPerformance)
	resp, _ := http.Get(url)
	z := html.NewTokenizer(resp.Body)
    i := 0
    home := "1"
    awayTeam := ""
    homeTeam := ""
    opp := ""
	for {
        tt := z.Next()
		if tt == html.ErrorToken {
			break
		}
        if tt == html.StartTagToken {
            t := z.Token()
			isTable := t.Data == "table"
			if isTable {
                team := t.Attr[1].Val[4:7]
                if i < 2 {
                    home = "0"
                    awayTeam = team
                } else {
                    homeTeam = team
                    opp = awayTeam
                }
				getBasicStats(z, playerMap, team, home, opp)
                i += 1
			}
		}
	}
    for _, val :=range( playerMap) {
        if val.stats["opp"] == ""{
            val.stats["opp"] = homeTeam
        }
        sem <-1
		waitGroup.Add(1)
        go val.addToTable(db, dateID)
    }
	waitGroup.Done()

}
func updateAndInsertPlayerRef(dateID string, db *sql.DB) {
    urls := getBoxScoreUrls(dateID, db)
	for _, url := range urls {
		waitGroup.Add(1)
		go getTables(url, db, dateID)
	}
}
func main() {
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", user, password, host, database)
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	arguements := os.Args
    dateID := arguements[1]

    updateAndInsertPlayerRef(dateID, db)
	waitGroup.Wait()
}
