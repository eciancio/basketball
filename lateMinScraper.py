from bs4 import BeautifulSoup
import demjson
import mysql.connector
import constants
import urllib2
'''
This scrapes rotogrinders for daily Minutes projections in order to insert and use in our regression model

'''

def rotogrinders_season_scraper():
	url = "https://rotogrinders.com/game-stats/nba-player?site=fanduel&range=season"
	stats = ["Player", "Team", "Pos", "Salary", "GP", "Min", "Reb", "Ast", "Stl", "Blk", "TO", "Pts", "USG", "FPTS"]
	csv_filename = 'rotogrinders_season_basketball.csv'
	rotogrindersBasketball(url, csv_filename, stats)

def rotogrinders_fourweeks_scraper():
	url = "https://rotogrinders.com/game-stats/nba-player?site=fanduel&range=4weeks"
	stats = ["Player", "Team", "Pos", "Salary", "GP", "Min", "Reb", "Ast", "Stl", "Blk", "TO", "Pts", "USG", "FPTS"]
	csv_filename = 'rotogrinders_fourweeks_basketball.csv'
	rotogrindersBasketball(url, csv_filename, stats)

def rotogrinders_minutes_scraper(dateID):
	url = "https://rotogrinders.com/projected-stats/nba-player"
	stats = ["Player", "Team", "Pos", "Salary", "GP", "Min", "Reb", "Ast", "Stl", "Blk", "TO", "Pts", "USG", "FPTS"]
	csv_filename = 'rotogrinders_minutes_basketball.csv'
	rotogrindersBasketball(url, csv_filename, stats, dateID)

# use rotogrinders to find the projected minutes for individuals 
def rotogrindersBasketball(url, csv_filename, stats, dateID):
	page = urllib2.urlopen(url).read()
	soup = BeautifulSoup(page, "html.parser")

	file = open(csv_filename, 'w')
	for stat in stats:
		file.write(stat + ", ")
	file.write("\n")

	jsonData = soup.find_all('script')
	jsonData = jsonData[len(jsonData) -1]

	jsonData = jsonData.text

	junk,jsonData = jsonData.split("data =")
	jsonData = jsonData.lstrip()
	jsonData = jsonData.rstrip()

	jsonData = jsonData.split(";")
	jsonData = jsonData[0]


        cnx = mysql.connector.connect(user=constants.databaseUser,
                                  host=constants.databaseHost,
                                  database=constants.databaseName,
                                  password=constants.databasePassword)
        cursor = cnx.cursor(buffered=True)


	basketballData = demjson.decode(jsonData)

        selec_id = "select playerID from player_reference where nickName=\""

        # this is date id for today
        date = dateID  
        update = "UPDATE futures SET projMinutes = %s where dateID = %s and playerID = %s"
        getPrevProj = "SELECT projMinutes from performance where dateID = %s and playerID = %s and projMinutes > 0"
        deleteFutures = "DELETE from futures where dateID = %s and playerID = %s"

        for playerData in basketballData:
		nickName = playerData['player']['first_name'] + " " + playerData['player']['last_name']
                
                get_id = selec_id + nickName + "\""
                try:
                    # attempt to insert into performance table
                    cursor.execute(get_id)
                    player_id = cursor.fetchall()[0][0]
                    inserts = (playerData['minutes'],date, player_id)

                    # delete row from futures if the player's old projMinutes > 0 and their new projMinutes == 0
                    getNonZeros = (date, player_id)
                    cursor.execute(getPrevProj, getNonZeros)
                    oldProjMinutes = cursor.fetchall()
                    if playerData['minutes'] == 0
                    	cursor.execute(deleteFutures, getNonZeros)

                    # execute insert
                    cursor.execute(update, inserts)
                    cnx.commit()
                except:
                    # if it fails print out who failed
                    player_id = nickName + "Failed"


		print nickName, playerData['minutes'], player_id
      

        cursor.close()
        cnx.commit()
        cnx.close()

def auto(dateID):
	rotogrinders_minutes_scraper(dateID)

if __name__ == "__main__":
	rotogrinders_minutes_scraper(dateID)
