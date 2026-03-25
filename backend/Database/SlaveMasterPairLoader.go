package database

import (
	"encoding/json"
	"os"
	"path"
	"syncjob/Logger"
)

type SlaveMasterMap map[string]SlaveMasterPair

// SlaveMasterMap
var SMMap SlaveMasterMap

func readFromJson(credsPath string, val any) error {
	bytes, err := os.ReadFile(credsPath)
	if err != nil {
		if os.IsNotExist(err) {
			if internErr := os.MkdirAll(path.Dir(credsPath), 0700); internErr != nil {
				return internErr
			}
		}
		return err
	}

	if err = json.Unmarshal(bytes, val); err != nil {
		return err
	}

	return nil
}

func saveToJson(jsonPath string, val any) error {
	f, err := os.Create(jsonPath)
	if err != nil {
		return err
	}
	defer f.Close();

	j := json.NewEncoder(f)
	j.SetIndent("", "\t")

	if err = j.Encode(val); err != nil {
		return err
	}

	return nil
}

func Close() {
	for k, v := range SMMap {
		if err := v.Close(); err != nil {
			logger.LogErr("Unable to close %s: %v", k, err)
		}
	}
}

// LoadSlaveMasterMap
func LoadSMPM() error {
	volPath := path.Join(os.Getenv("CREDS_SUBDIR"), "mappings.json")
	logger.LogInfo("Loading mappings from: %s", volPath)

	if err := readFromJson(path.Join(volPath, "mappings.json"), &SMMap); err != nil {
		SMMap = make(SlaveMasterMap)
		return err;
	}

	return nil
}

// SaveSlaveMasterMap
func SaveSMPM() error {
	volPath := path.Join(os.Getenv("CREDS_SUBDIR"), "mappings.json")
	logger.LogInfo("Loading mappings from: %s", volPath)

	if err := saveToJson(path.Join(volPath, "mappings.json"), SMMap); err != nil {
		return err;
	}

	return nil
}
