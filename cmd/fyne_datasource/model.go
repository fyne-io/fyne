package main

// This whole section will be removed - its a test harness to try out databinding ideas
// Will move this to a separate repo shortly

import (
	"sort"

	"fyne.io/fyne/dataapi"
)

// DataModel - a single instance of this is used by all the views in the app
type DataModel struct {
	CountryAndState *CountryAndStateDB
	CountriesList   *dataapi.SliceDataSource
	StatesList      *dataapi.SliceDataSource
	SelectedCountry *dataapi.String
	SelectedState   *dataapi.String
	Clock           *dataapi.Clock
}

func NewDataModel() *DataModel {
	cs := NewCountryAndStateDB(
		map[string][]string{
			"UK":        {"City of London", "England", "Scotland", "Wales", "Nth Ireland", "Yorkshire", "Midlands"},
			"Zwazia":    {"North Zwazia", "Central", "Capitol City", "East Coast", "Mountains"},
			"USA":       {"NY", "WA", "DC", "MI", "TX", "CA", "NE", "BA", "AR", "AK", "OH", "MS", "OR", "ID", "MO", "FL", "VA", "NC"},
			"Australia": {"SA", "NSW", "QLD"},
			"USSR":      {"Byelorussia", "Ukrania", "Moskva", "Leningrad", "Novosibirsk"},
		},
	)
	return &DataModel{
		CountryAndState: cs,
		CountriesList:   dataapi.NewSliceDataSource().SetFromStringSlice(cs.Keys()),
		StatesList:      dataapi.NewSliceDataSource(),
		SelectedCountry: dataapi.NewString(""),
		SelectedState:   dataapi.NewString(""),
		Clock:           dataapi.NewClock(),
	}
}

type CountryAndStateDB struct {
	data map[string][]string
}

func NewCountryAndStateDB(data map[string][]string) *CountryAndStateDB {
	return &CountryAndStateDB{data: data}
}

// Keys returns the keys as a sorted slice
func (c *CountryAndStateDB) Keys() []string {
	data := make([]string, 0, len(c.data))
	for k := range c.data {
		data = append(data, k)
	}
	sort.Strings(data)
	return data
}

func (c *CountryAndStateDB) GetStates(country string) []string {
	data := c.data[country]
	sort.Strings(data)
	return data
}
