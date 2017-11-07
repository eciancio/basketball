import re
import mysql.connector
from datetime import timedelta, date
import constants
from bs4 import BeautifulSoup, Comment
import urllib2
import requests


def updateAndInsertTeamRef(startDay, startMonth, startYear, endDay, endMonth, endYear, cursor):
    # set range of dates
    # urls = generateURLs(startDay, startMonth, startYear, endDay, endMonth, endYear)
    urls = ["https://www.basketball-reference.com/boxscores/201704070DAL.html"]
    start_date = date(startYear, startMonth, startDay)
    end_date = date(endYear, endMonth, endDay)

    # loop through all url's
    for url in urls:
        page = requests.get(url)
        soup = BeautifulSoup(page.text, 'html.parser')

        soup.findAll(text=lambda text: isinstance(text, Comment))
        comments = soup.findAll(text=lambda text: isinstance(text, Comment))
        comment = comments[len(comments) - 23]

        soup1 = BeautifulSoup(comment, "html.parser")
        for tr in soup1.find_all('tr')[2:]:
            # team name
            team = tr.find_all('th')[0].text

            # table columns
            tds = tr.find_all('td')
            pace = tds[0].text
            eShooting = tds[1].text
            turnoverPercent = tds[2].text
            oReboundPercent = tds[3].text
            FTOverFQAttempts = tds[4].text
            ORTG = tds[5].text

            print pace, eShooting, turnoverPercent, oReboundPercent, FTOverFQAttempts, ORTG

        for tr in soup.find_all('tr')[2:]:
            #	print tr
            # first find, then updated
            tds = tr.find_all('td')
            #       print tds[1].a.text

if __name__ == "__main__":
    cnx = mysql.connector.connect(user=constants.databaseUser,
                                  host=constants.databaseHost,
                                  database=constants.databaseName,
                                  password=constants.databasePassword)
    cursor = cnx.cursor(buffered=True)

    updateAndInsertTeamRef(constants.startDayP, constants.startMonthP, constants.startYearP, constants.endDayP,
                             constants.endMonthP, constants.endYearP, cursor)

    cursor.close()
    cnx.commit()
    cnx.close()