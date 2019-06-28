package main

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Module struct { 
	Store Store
}

type Store interface {
	GetAverage(model string) (int, error) 
	GetMedian(model string) (int, error) 
	GetEntryCount(model string) (int, error) 
	GetGreaterCount(model string, moneyLine int) (int, error) 
}

type store struct {
	db *sql.DB
}

func (d *store) GetRowInt(query string) (int, error){
	var val float64
	row := d.db.QueryRow(query)
	err := row.Scan(&val)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return int(val), nil
}

func (m *Module) GetAverage(model string) int {
	avg, err := m.Store.GetAverage(model)
	if err != nil {
		return 0
	}
	return avg
}

func (m *Module) GetMedian(model string) int {
	median, err := m.Store.GetMedian(model)
	if err != nil {
		return 0
	}
	return median
}

func (m * Module) GetWinPercentage(model string, moneyLine int) int {
	total, err := m.Store.GetEntryCount(model)
	if err != nil {
		return 0
	}
	greater, err := m.Store.GetGreaterCount(model, moneyLine)
	if err != nil {
		return 0
	}
	return (100 * greater)/total

}


func GetAverageQuery( model string )string {
	base_string := "select avg(actualPointsLineup) FROM basketball.historic_lineups where model like '%s_';"
	return fmt.Sprintf(base_string, model)
}

func GetEntryCountQuery( model string )string {
	base_string := "select count(*) FROM basketball.historic_lineups where model like '%s_';"
	return fmt.Sprintf(base_string, model)
}

func GetGreaterCountQuery( model string, moneyLine int) string {
	base_string := "select count(*) FROM basketball.historic_lineups where model like '%s_' and actualPointsLineup > '%d';"
	return fmt.Sprintf(base_string, model, moneyLine)
}

func GetMedianQuery(model string ) (string) {
	medianQuery := "SELECT AVG(dd.actualPointsLineup) as median_val FROM (SELECT d.actualPointsLineup, @rownum:=@rownum+1 as `row_number`, @total_rows:=@rownum FROM historic_lineups d, (SELECT @rownum:=0) r WHERE d.actualPointsLineup is NOT NULL and model like '%s_' ORDER BY d.actualPointsLineup) as dd WHERE dd.row_number IN ( FLOOR((@total_rows+1)/2), FLOOR((@total_rows+2)/2));"
	return fmt.Sprintf(medianQuery, model)
}

func (d *store) GetAverage(model string) (int, error) {
	query := GetAverageQuery(model)
	return d.GetRowInt(query)
}

func (d *store) GetEntryCount(model string) (int, error) {
	query := GetEntryCountQuery(model)
	return d.GetRowInt(query)
}

func (d *store) GetGreaterCount(model string, moneyLine int) (int, error) {
	query := GetGreaterCountQuery(model, moneyLine)
	return d.GetRowInt(query)
}
func (d *store) GetMedian(model string) (int, error) {
	query := GetMedianQuery(model)
	return d.GetRowInt(query)
}
