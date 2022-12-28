package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v3"
)

var (
	DefaultConfFilePath = "string-checker/config.yml"
)

type conf struct {
	Rules []struct {
		Name          string   `yaml:"name"`
		Message       string   `yaml:"message"`
		Level         level    `yaml:"level"`
		CaseSensitive bool     `yaml:"caseSensitive"`
		Targets       []string `yaml:"targets"`
	} `yaml:"rules"`
}

type level string

var (
	levelWarning level = "warning" //lint:ignore U1000 ignore
	levelError   level = "error"   //lint:ignore U1000 ignore
)

func main() {
	checkEnv()

	filePattern := os.Getenv("INPUT_FILEPATTERN")

	targetFiles, err := filepath.Glob(filePattern)
	checkError(err)

	confFilePath := os.Getenv("INPUT_CONFFILEPATH")
	if confFilePath == "" {
		confFilePath = DefaultConfFilePath
	}

	c, err := getConf(confFilePath)
	checkError(err)

	err = checkLevel(c)
	checkError(err)

	merr := new(multierror.Error)
	for _, targetFile := range targetFiles {
		err := scan(targetFile, c)
		merr = multierror.Append(merr, err)
	}

	checkError(merr.ErrorOrNil())
}

func checkEnv() {
	envs := []string{
		"INPUT_FILEPATTERN",
	}

	for _, env := range envs {
		if os.Getenv(env) == "" {
			log.Fatalf("env: %s is required", env)
		}
	}
}

func checkError(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
	}
}

func checkLevel(c *conf) error {
	for _, r := range c.Rules {
		if r.Level != levelError && r.Level != levelWarning {
			return fmt.Errorf("invalid level: %s", r.Level)
		}
	}

	return nil
}

func getConf(fileName string) (*conf, error) {
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	c := &conf{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func contains(a, b string, caseSensitive bool) bool {
	if caseSensitive {
		return strings.Contains(a, b)
	}

	return strings.Contains(strings.ToLower(a), strings.ToLower(b))
}

func scan(fileName string, c *conf) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	merr := new(multierror.Error)
	line := 1

	for scanner.Scan() {
		for _, r := range c.Rules {
			for _, t := range r.Targets {
				if contains(scanner.Text(), t, r.CaseSensitive) {
					msg := fmt.Sprintf(`::%s file=%s,line=%v,title=%s::%s`, r.Level, fileName, line, r.Name, r.Message)

					if r.Level == levelError {
						merr = multierror.Append(merr, fmt.Errorf(msg))
					}

					fmt.Println(msg)
				}
			}
		}
		line++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return merr.ErrorOrNil()
}
