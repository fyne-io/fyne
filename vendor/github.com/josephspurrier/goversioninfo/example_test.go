package goversioninfo

import (
	"io/ioutil"
	"log"
	"os"
)

// Example
func Example() {
	logic()
}

// Create the syso file
func logic() {
	// Read the config file
	jsonBytes, err := ioutil.ReadFile("versioninfo.json")
	if err != nil {
		log.Printf("Error reading versioninfo.json: %v", err)
		os.Exit(1)
	}

	// Create a new container
	vi := &VersionInfo{}

	// Parse the config
	if err := vi.ParseJSON(jsonBytes); err != nil {
		log.Printf("Could not parse the .json file: %v", err)
		os.Exit(2)
	}

	// Fill the structures with config data
	vi.Build()

	// Write the data to a buffer
	vi.Walk()

	// Optionally, embed an icon by path
	// If the icon has multiple sizes, all of the sizes will be embedded
	vi.IconPath = "icon.ico"

	// Create the file
	if err := vi.WriteSyso("resource.syso", "386"); err != nil {
		log.Printf("Error writing syso: %v", err)
		os.Exit(3)
	}
}
