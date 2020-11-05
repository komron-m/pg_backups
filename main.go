package main

import (
	"log"
	"os"
	"os/exec"
	"path"
	"time"
)

type PostgresCredentials struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	DumpsDir string
}

func main() {
	// initialize credentials
	c := NewPostgresCredentials()

	// create command for taking backup
	cmd := exec.Command(
		"pg_dump",
		"--compress=9",
		"--file="+makeFileNameWithPath(c.DumpsDir),
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
			log.Fatal(err)
		}
	}

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func createFilePath(basePath string) string {
	t := time.Now()
	ymd := path.Join(basePath, t.Format("2006"), t.Format("01"), t.Format("02"))
	if err := os.MkdirAll(ymd, os.ModeDir); err != nil {
		log.Fatal(err)
	}
	return ymd
}

func makeFileNameWithPath(basePath string) string {
	fullPath := createFilePath(basePath)
	return path.Join(fullPath, time.Now().Format("2006-01-02T15_04_05.sql.gz"))
}
