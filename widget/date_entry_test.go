package widget

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_DateEntry(t *testing.T) {
	dateEntry := NewDateEntry()

	testDate := time.Date(2025, 07, 18, 0, 0, 0, 0, time.Local)
	testDateText := testDate.Format(dateEntry.DateFormat())

	inputTests := []struct {
		input    string
		wantErr  bool
		wantDate *time.Time
	}{
		{input: testDateText, wantErr: false, wantDate: &testDate},
		{input: "", wantErr: false, wantDate: nil},
		{input: testDateText, wantErr: false, wantDate: &testDate},
		{input: "not a valid date", wantErr: true, wantDate: nil},
	}

	for _, tt := range inputTests {
		dateEntry.SetText("")
		for _, r := range tt.input {
			dateEntry.TypedRune(r)
		}

		assert.Equal(t, tt.input, dateEntry.Text)
		if tt.wantErr {
			assert.Error(t, dateEntry.validationError)
		} else {
			assert.NoError(t, dateEntry.validationError)
		}

		if tt.wantDate == nil {
			assert.Nil(t, dateEntry.Date)
		} else {
			assert.NotNil(t, dateEntry.Date)
			assert.Equal(t, tt.wantDate.Year(), dateEntry.Date.Year())
			assert.Equal(t, tt.wantDate.Month(), dateEntry.Date.Month())
			assert.Equal(t, tt.wantDate.Day(), dateEntry.Date.Day())
		}
	}
}

func Test_DateEntry_CustomValidator(t *testing.T) {
	minDate := time.Date(1989, 9, 16, 0, 0, 0, 0, time.Local)

	noDateErr := fmt.Errorf("Please choose a date")
	invDateErr := fmt.Errorf("Invalid Date")
	toSoonErr := fmt.Errorf("Please eneter a date after 16 September 1989")

	dateEntry := NewDateEntry()
	dateEntry.Validator = func(s string) error {
		if s == "" {
			return noDateErr
		}

		if t, err := time.Parse(dateEntry.DateFormat(), s); err != nil {
			return invDateErr
		} else {
			if t.Before(minDate) {
				return toSoonErr
			}
		}
		return nil
	}

	testDate := time.Date(2025, 07, 18, 0, 0, 0, 0, time.Local)
	testDateText := testDate.Format(dateEntry.DateFormat())

	toSoonDate := minDate.Add(time.Duration(time.Hour * -24))
	toSoonDateText := toSoonDate.Format(dateEntry.DateFormat())

	inputTests := []struct {
		input    string
		wantErr  error
		wantDate *time.Time
	}{
		{input: testDateText, wantErr: nil, wantDate: &testDate},
		{input: "", wantErr: noDateErr, wantDate: nil},
		{input: testDateText, wantErr: nil, wantDate: &testDate},
		{input: "not a valid date", wantErr: invDateErr, wantDate: nil},
		{input: testDateText, wantErr: nil, wantDate: &testDate},
		{input: toSoonDateText, wantErr: toSoonErr, wantDate: &toSoonDate},
	}

	for _, tt := range inputTests {
		dateEntry.SetText("")
		for _, r := range tt.input {
			dateEntry.TypedRune(r)
		}

		assert.Equal(t, tt.input, dateEntry.Text)
		if tt.wantErr != nil {
			assert.EqualError(t, dateEntry.validationError, tt.wantErr.Error())
		} else {
			assert.NoError(t, dateEntry.validationError)
		}

		if tt.wantDate == nil {
			assert.Nil(t, dateEntry.Date)
		} else {
			assert.NotNil(t, dateEntry.Date)
			assert.Equal(t, tt.wantDate.Year(), dateEntry.Date.Year())
			assert.Equal(t, tt.wantDate.Month(), dateEntry.Date.Month())
			assert.Equal(t, tt.wantDate.Day(), dateEntry.Date.Day())
		}
	}
}

func Test_DateEntry_SetDate(t *testing.T) {
	dateEntry := NewDateEntry()

	assert.Nil(t, dateEntry.Date)
	assert.Equal(t, "", dateEntry.Text)
	assert.NoError(t, dateEntry.Validate())

	testDate := time.Date(2025, 07, 18, 0, 0, 0, 0, time.Local)
	testDateText := testDate.Format(dateEntry.DateFormat())
	dateEntry.SetDate(&testDate)

	assert.NotNil(t, dateEntry.Date)
	assert.Equal(t, testDate.Year(), dateEntry.Date.Year())
	assert.Equal(t, testDate.Month(), dateEntry.Date.Month())
	assert.Equal(t, testDate.Day(), dateEntry.Date.Day())
	assert.Equal(t, testDateText, dateEntry.Text)
	assert.NoError(t, dateEntry.Validate())

	dateEntry.SetDate(nil)

	assert.Nil(t, dateEntry.Date)
	assert.Equal(t, "", dateEntry.Text)
	assert.NoError(t, dateEntry.Validate())
}
