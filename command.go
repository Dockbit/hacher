package main

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"bytes"
)

func cmdGet(c *cli.Context) {

	verbose = c.GlobalBool("verbose")

	path := c.Args().Get(0)
	key := c.String("key")
	files := strings.Split(c.String("file"), ",")

	if len(path) < 1 {
		path = "." // default to current directory
	}

	if len(key) < 1 {
		key = strings.ToLower(filepath.Base(files[0]))
	}
	hash := checksum(files)
	fullPath := filepath.Join(CachePath, strings.Join([]string{key, hash}, "-")) + ".tar.gz"

	// get cache if exists
	if _, err := os.Stat(fullPath); err == nil {
		printInfo("Fetching cache '%s'. Please, wait...", key)

		args := []string{
			"-xzf",
			fullPath,
			"-C",
			path,
		}

		err := exec.Command("tar", args...).Run()
		checkError(err)
	}
}

func cmdSet(c *cli.Context) {

	verbose = c.GlobalBool("verbose")

	path := c.Args().Get(0)
	key := c.String("key")
	files := strings.Split(c.String("file"), ",")

	if len(path) < 1 {
		printFatal("Path to content is not provided as an argument.")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		printFatal("Content '%s' does not exist.", path)
	}

	if _, err := os.Stat(CachePath); os.IsNotExist(err) {
		if os.MkdirAll(CachePath, dirMode) != nil {
			printFatal("Couldn't create cache directory. "+
				"Is the %s directory writable?", CachePath)
		}
	}

	if len(key) < 1 {
		key = strings.ToLower(filepath.Base(files[0]))
	}

	hash := checksum(files)
	fullPath := filepath.Join(CachePath, strings.Join([]string{key, hash}, "-")) + ".tar.gz"

	// cache contents only if it doesn't exist already
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		printInfo("Caching '%s'. Please, wait...", key)

		args := []string{
			"-czf",
			fullPath,
			path,
		}

		err := exec.Command("tar", args...).Run()
		checkError(err)
	}
	go clean(key)
}

/*
 * Calculates SHA256 checksum of an array of files
 *
 * Returns The String checksum of all files
 */
func checksum(files []string) string {
	if len(files) < 1 {
		printFatal("At least one dependency file is not provided.")
	}

	var filesBuffer bytes.Buffer

	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			printFatal("Dependency file '%s' does not exist.", file)
		}
		contents, err := ioutil.ReadFile(file)
		checkError(err)
		filesBuffer.Write(contents)
	}
	hasher := sha256.New()
	hasher.Write(filesBuffer.Bytes())
	return hex.EncodeToString(hasher.Sum(nil))
}

/*
 * Cleans old caches
 */
func clean(key string) {
	files := fileSorter(CachePath, "^"+key+"-")

	for index, file := range files {
		if index+1 > CacheKeep {
			os.Remove(file.Path)
		}
	}
}
