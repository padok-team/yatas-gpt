package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/padok-team/yatas/plugins/commons"
	"gopkg.in/yaml.v2"
)

// generate a test that reads the file tests/tests.yml, unmarchall it into commons.Config.Tests

func TestUnmarshalGPT(t *testing.T) {
	// unmarchall the file tests/tests.yml into commons.Config.Tests

	var tests []commons.Tests
	// Read file tests/tests.yml
	filename, _ := filepath.Abs("./tests/tests.yml")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	yaml.Unmarshal(yamlFile, &tests)

	// Unmarshal the file tests/tests.yml into commons.Config.Tests
	var c commons.Config
	c.Tests = tests

	// Run the plugin
	gptCredentials := GPTPlugin{
		apiKey: "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",
		model:  "davinci",
		prompt: "This is a test",
	}

	generateReportChat(&YatasPlugin{}, gptCredentials, &c)
}
