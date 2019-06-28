package main

import "testing"

type MockStore struct {}

func (m *MockStore) GetAverage(model string) (int, error) {
	return 0, nil
}

func (m *MockStore) GetMedian(model string) (int, error) {
	return 200, nil
}

func (m *MockStore) GetEntryCount(model string) (int, error) {
	return 20000, nil
}

func (m *MockStore) GetGreaterCount(model string, moneyLine int) (int, error) {
	return 10000, nil
}

func TestGetAverage(t *testing.T) { 
	m := new(MockStore)
	k := Module{m}
	a := k.GetAverage("testing")
	if a != 0 {
		t.Fail()
	}
}

func TestGetAverageQuery(t *testing.T) { 
	model := "mlp"
	query := GetAverageQuery(model)
	if query != "select avg(actualPointsLineup) FROM basketball.historic_lineups where model like 'mlp_';" { 
		t.Fail()
	}
}

func TestGetEntryCountQuery(t *testing.T) { 
	model := "mlp"
	query := GetEntryCountQuery(model)
	if query != "select count(*) FROM basketball.historic_lineups where model like 'mlp_';" { 
		t.Fail()
	}
}
func TestGreaterEntryCountQuery(t *testing.T) { 
	model := "mlp"
	moneyLine := 300
	query := GetGreaterCountQuery(model, moneyLine)
	if query != "select count(*) FROM basketball.historic_lineups where model like 'mlp_' and actualPointsLineup > '300';" { 
		t.Fail()
	}
}
func TestGetMedianQuery(t *testing.T) { 
	model := "mlp"
	query := GetMedianQuery(model)
	if query != "SELECT AVG(dd.actualPointsLineup) as median_val FROM (SELECT d.actualPointsLineup, @rownum:=@rownum+1 as `row_number`, @total_rows:=@rownum FROM historic_lineups d, (SELECT @rownum:=0) r WHERE d.actualPointsLineup is NOT NULL and model like 'mlp_' ORDER BY d.actualPointsLineup) as dd WHERE dd.row_number IN ( FLOOR((@total_rows+1)/2), FLOOR((@total_rows+2)/2));" { 
		t.Fail()
	}
}

func TestGetMedian(t *testing.T) { 
	m := new(MockStore)
	k := Module{m}
	a := k.GetMedian("mlp")
	if a != 200 {
		t.Fail()
	}
}

func TestGetWinPercentage(t *testing.T) { 
	m := new(MockStore)
	k := Module{m}
	a := k.GetWinPercentage("mlp", 200)
	if a != 50 {
		t.Fail()
	}
}

