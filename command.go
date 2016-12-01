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
	"sort"
)

func cmdGet(c *cli.Context) {

	verbose = c.GlobalBool("verbose")

	path := c.Args().Get(0)
	key := c.String("key")
	files := strings.Split(c.String("file"), ",")
	envs := strings.Split(c.String("env"), ",")

	if len(path) < 1 {
		path = "." // default to current directory
	}

	if len(key) < 1 {
		key = strings.ToLower(filepath.Base(files[0]))
	}
	hash := checksum(files, envs)
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
	envs := strings.Split(c.String("env"), ",")

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

	hash := checksum(files, envs)
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
 * Calculates SHA256 checksum of an array of files and/or env vars
 *
 * Returns The String checksum of all files
 */
func checksum(files []string, envs []string) string {
	if len(files[0]) < 1 {
		printFatal("At least one dependency file is required.")
	}

	var buffer bytes.Buffer

	// first go thru files
	sort.Strings(files)
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			printFatal("Dependency file '%s' does not exist.", file)
		}
		contents, err := ioutil.ReadFile(file)
		checkError(err)
		buffer.Write(contents)
	}

	// then check any environment variables
	sort.Strings(envs)
	for _, v := range envs {
		envVar := os.Getenv(v)
		if len(envVar) > 0 {
			buffer.WriteString(envVar)
		}
	}

	hasher := sha256.New()
	hasher.Write(buffer.Bytes())
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
