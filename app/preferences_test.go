package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func loadPreferences(id string) *preferences {
	p := newPreferences(&fyneApp{uniqueID: id})
	p.load()

	return p
}

func TestPreferences_Remove(t *testing.T) {
	p := loadPreferences("dummy")
	p.WriteValues(func(val map[string]interface{}) {
		val["keyString"] = "value"
		val["keyInt"] = 4
	})

	p.RemoveValue("keyString")
	err := p.saveToFile(p.storagePath())
	assert.Nil(t, err)

	// check it doesn't write values that were removed
	p = loadPreferences("dummy")
	assert.Equal(t, 4, p.Int("keyInt"))
	assert.Equal(t, "missing", p.StringWithFallback("keyString", "missing"))
}

func TestPreferences_Save(t *testing.T) {
	p := loadPreferences("dummy")
	p.WriteValues(func(val map[string]interface{}) {
		val["keyString"] = "value"
		val["keyStringList"] = []string{"1", "2", "3"}
		val["keyInt"] = 4
		val["keyIntList"] = []int{1, 2, 3}
		val["keyFloat"] = 3.5
		val["keyFloatList"] = []float64{1.1, 2.2, 3.3}
		val["keyBool"] = true
		val["keyBoolList"] = []bool{true, false, true}
	})

	path := p.storagePath()
	defer os.Remove(path)
	err := p.saveToFile(path)
	assert.Nil(t, err)

	expected, err := os.ReadFile(filepath.Join("testdata", "preferences.json"))
	if err != nil {
		assert.Fail(t, "Failed to load, %v", err)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		assert.Fail(t, "Failed to load, %v", err)
	}
	assert.JSONEq(t, string(expected), string(content))

	// check it reads the saved output
	p = loadPreferences("dummy")
	assert.Equal(t, "value", p.String("keyString"))
	assert.Equal(t, 3, len(p.StringList("keyStringList")))
}

func TestPreferences_Save_OverwriteFast(t *testing.T) {
	p := loadPreferences("dummy2")
	p.WriteValues(func(val map[string]interface{}) {
		val["key"] = "value"
	})

	path := filepath.Join(os.TempDir(), "fynePrefs2.json")
	defer os.Remove(path)
	p.saveToFile(path)

	p.WriteValues(func(val map[string]interface{}) {
		val["key2"] = "value2"
	})
	p.saveToFile(path)

	p2 := loadPreferences("dummy")
	p2.loadFromFile(path)
	assert.Equal(t, "value2", p2.String("key2"))
}

func TestPreferences_Load(t *testing.T) {
	p := loadPreferences("dummy")
	p.loadFromFile(filepath.Join("testdata", "preferences.json"))

	assert.Equal(t, "value", p.String("keyString"))
	assert.Equal(t, []string{"1", "2", "3"}, p.StringList("keyStringList"))
	assert.Equal(t, 4, p.Int("keyInt"))
	assert.Equal(t, []int{1, 2, 3}, p.IntList("keyIntList"))
	assert.Equal(t, 3.5, p.Float("keyFloat"))
	assert.Equal(t, []float64{1.1, 2.2, 3.3}, p.FloatList("keyFloatList"))
	assert.Equal(t, true, p.Bool("keyBool"))
	assert.Equal(t, []bool{true, false, true}, p.BoolList("keyBoolList"))
}
