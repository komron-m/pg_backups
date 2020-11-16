package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"time"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	// MainBackupDir - folder where all backups will be stored after dumping postgres db
	MainBackupDir string
	// RemoveDailyBackupFolderAfterNDay - removes the backup folder after `N` days
	// e.g. one may want to delete entire backup dir after 3 days (only deletes from `MainBackupDir` root)
	RemoveDailyBackupFolderAfterNDay int
	// SecondaryBackupDir - optionally specifying directory where backups will be copied
	// on specified times in `MakeSecondaryBackupsAt`...
	SecondaryBackupDir string
	// MakeSecondaryBackupsAt - ... format hour:minute, ex: {"08:30", "13:30", "21:00"}
	MakeSecondaryBackupsAt []string
}

func main() {
	// initialize configs
	c := NewConfig()

	// create command for taking backup
	dumpFilePath := makeFileNameWithPath(c.MainBackupDir)
	cmd := exec.Command(
		"pg_dump",
		"--compress=9",
		"--file="+dumpFilePath,
		"-U"+c.User,
		"-h"+c.Host,
		"-p"+c.Port,
		"-d"+c.Database,
		"-w",
	)

	// set PGPASSWORD for current command
	// this is easiest way to pass postgres password, without handling stdin
	if c.Password != "" {
		if err := os.Setenv("PGPASSWORD", c.Password); err != nil {
			panic(err)
		}
	}
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	// ------------------------------------------------------------------

	// copy the dumpFile into SecondaryBackupDir
	if shouldCreateSecondaryBackup(c) {
		mainDumpFile, err := ioutil.ReadFile(dumpFilePath)
		if err != nil {
			panic(err)
		}

		pathForSecondaryDump := makeFileNameWithPath(c.SecondaryBackupDir)
		if err := ioutil.WriteFile(pathForSecondaryDump, mainDumpFile, 0666); err != nil {
			panic(err)
		}
	}

	// ------------------------------------------------------------------

	// remove backups that live more that N days from MainBackupDir
	if c.RemoveDailyBackupFolderAfterNDay != 0 {
		removeFolder := createPathBeforeNDay(c.MainBackupDir, c.RemoveDailyBackupFolderAfterNDay)
		if err := os.RemoveAll(removeFolder); err != nil {
			panic(err)
		}
	}
}

func shouldCreateSecondaryBackup(c *Config) bool {
	if c.SecondaryBackupDir == "" {
		return false
	}

	hm := time.Now().Format("15:04")
	for _, v := range c.MakeSecondaryBackupsAt {
		if v == hm {
			return true
		}
	}

	return false
}

func createPathBeforeNDay(basePath string, daysBefore int) string {
	if daysBefore > 0 {
		daysBefore = -daysBefore
	}
	t := time.Now().AddDate(0, 0, daysBefore)
	return path.Join(basePath, t.Format("2006"), t.Format("01"), t.Format("02"))
}

func createFilePath(basePath string) string {
	t := time.Now()
	ymd := path.Join(basePath, t.Format("2006"), t.Format("01"), t.Format("02"))
	if err := os.MkdirAll(ymd, os.ModeDir); err != nil {
		panic(err)
	}
	return ymd
}

func makeFileNameWithPath(basePath string) string {
	fullPath := createFilePath(basePath)
	return path.Join(fullPath, time.Now().Format("2006-01-02T15_04_05.sql.gz"))
}
