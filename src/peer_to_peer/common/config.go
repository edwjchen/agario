package common

import (
	"os"
	"encoding/json"
)

type Config struct {

	SCREEN_WIDTH      uint32
	SCREEN_HEIGHT     uint32
	MIN_FOOD_NUM      uint32

	STARTING_MASS     int32
	MAP_WIDTH         uint32
	MAP_HEIGHT        uint32
	REGION_MAP_WIDTH  uint32
	REGION_MAP_HEIGHT uint32

	NREGION_WIDTH     uint32
	NREGION_HEIGHT    uint32

	SPEED             float64
	EAT_RADIUS_DELTA  float64
	ZOOM              float64

	POISON_PROB       int

	VER_FILE          string

	PLAYER_PORT       string
	REGION_PORT       string

	ENTRY_SERVER      string 
	REGION_SERVERS    []string
}

type EntryServerConfig struct {
	ADDR        string 
	MIN_PLAYERS uint32
	MAX_PLAYERS uint32
	// MAP_LENGTH  int32
}

var Conf Config

func ReadConfig(filename string) error {
	file, err := os.Open(filename) 
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file) 
	err = decoder.Decode(&Conf) 
	if err != nil {
		return err 
	}
	Conf.NREGION_WIDTH = Conf.MAP_WIDTH/Conf.REGION_MAP_WIDTH
	Conf.NREGION_HEIGHT = Conf.MAP_HEIGHT/Conf.REGION_MAP_HEIGHT
	return nil
}

func ReadEntryServerConfig(filename string) (EntryServerConfig, error) {
	file, err := os.Open(filename) 
	if err != nil {
		return EntryServerConfig{}, err
	}
	decoder := json.NewDecoder(file) 
	esconfig := EntryServerConfig{}
	err = decoder.Decode(&esconfig) 
	if err != nil {
		return EntryServerConfig{}, err 
	}
	return esconfig, nil
}