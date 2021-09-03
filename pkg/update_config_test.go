package pkg

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateRegistryConfig(t *testing.T) {
	assert := assert.New(t)
	testConfigFile := "../testfiles/config_test.toml"

	var tests = []struct {
		username       string
		password       string
		registryURL    string
		clientCertName string
		configFile     string
		insecure       bool
		err            error
	}{
		{"testuser", "testpassword", "demo.registry.hpe.com",
			"test.pem", testConfigFile, true, nil},
		{"", "", "",
			"test.pem", testConfigFile, true,
			fmt.Errorf("argument registry" +
				" url value can't be empty, check program usage and rerun")},
	}

	for _, test := range tests {
		assert.Equal(UpdateRegistryConfig(test.configFile, test.username, test.password,
			test.registryURL, test.clientCertName, test.insecure), test.err)
	}
}

func TestAddNvidiaConfig(t *testing.T) {
	assert := assert.New(t)
	testConfigFile := "../testfiles/config_test.toml"
	fakeFile := ""

	_, err := os.Stat(fakeFile)
	var testErr = fmt.Errorf("specified config file %s not exists locally, error: %v", fakeFile, err)

	var tests2 = []struct {
		configFile string
		err error
	}{
		{testConfigFile, nil},
		{fakeFile, testErr},
	}

	for _, test := range tests2 {
		assert.Equal(AddNvidiaConfig(test.configFile), test.err)
	}
}

func TestDeleteNvidiaConfig(t *testing.T) {
	assert := assert.New(t)
	testFile := "../testfiles/config_test.toml"
	fakeFile := ""
	fakeNvidiaFile := "../testfiles/fake_config.toml"

	_, err := os.Stat(fakeFile)
	var fakeErr = fmt.Errorf("specified config file %s not exists locally, error: %v", fakeFile, err)
	var keyPathErr = fmt.Errorf("specified nvidia path not found in the config file")

	var tests = []struct {
		configFile string
		err error
	}{
		{testFile, nil},
		{fakeFile, fakeErr},
		{fakeNvidiaFile, keyPathErr},
	}

	for _, test := range tests {
		assert.Equal(DeleteNvidiaConfig(test.configFile), test.err)
	}
}

func TestDeleteRegistryConfig(t *testing.T) {
	assert := assert.New(t)
	testFile := "../testfiles/config_test.toml"
	fakeFile := ""
	fakeNvidiaFile := "../testfiles/fake_config.toml"

	_, err := os.Stat(fakeFile)
	var fakeErr = fmt.Errorf("specified config file %s not exists locally, error: %v", fakeFile, err)
	var keyPathErr = fmt.Errorf("specified registry path not found in the config file")

	var tests = []struct {
		configFile string
		registryPath string
		err error
	}{
		{testFile, "demo.registry.hpe.com",nil},
		{fakeFile, "demo.registry.hpe.com",fakeErr},
		{fakeNvidiaFile, "fake.registry.hpe.com",keyPathErr},
	}

	for _, test := range tests {
		assert.Equal(DeleteRegistryConfig(test.configFile, test.registryPath), test.err)
	}
}
