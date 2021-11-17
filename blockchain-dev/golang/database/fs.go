package database

import (
	"os"
	"path/filepath"
)

func InitDataDirIfNotExist(dataDir string, genesis []byte) error {
	if fileExist(getGenesisJsonFilePath(dataDir)) {
		return nil
	}
	if err := os.MkdirAll(getDatabaseDirPath(dataDir), os.ModePerm); err != nil {
		return err
	}

	if err := writeGenesisToDisk(getGenesisJsonFilePath(dataDir), genesis); err != nil {
		return err
	}

	return nil

}

func getDatabaseDirPath(dataDir string) string {
	return filepath.Join(dataDir, "Database")
}

func getGenesisJsonFilePath(dataDir string) string {
	return filepath.Join(getDatabaseDirPath(dataDir), "genesis.json")
}

func getBlocksDbFilePath(dataDir string) string {
	return filepath.Join(getDatabaseDirPath(dataDir), "block.db")
}

func fileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	return true
}
