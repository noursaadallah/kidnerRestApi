package settings

import (
	"fmt"
	"testing"
)

func TestSettings_GetFabricSetup(t *testing.T) {
	// Load config file containing parameters for FabricSetup initialization
	fs, err := GetFabricSettings()
	if err != nil {
		fmt.Println("error loading config for FabricSetup init")
		t.FailNow()
	}
	fmt.Println(fs)
}

func TestSettings_GetWebSetup(t *testing.T) {
	// Load config file containing parameters for the web server
	ws, err := GetWebSettings()
	if err != nil {
		fmt.Println("error loading config for FabricSetup init")
		t.FailNow()
	}
	fmt.Println(ws)
}
