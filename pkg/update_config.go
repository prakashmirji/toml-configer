package pkg

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

// UpdateRegistryConfig - updates registry settings in the config file
func UpdateRegistryConfig(configFile, username, password, registryURL, clientCertName string, insecure bool) error {
	log.Debugf("UpdateRegistryConfig hit, config file: %s", configFile)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("specified config file %s not exists locally, error: %v", configFile, err)
	}

	if registryURL == "" {
		return fmt.Errorf("argument registry" +
			" url value can't be empty, check program usage and rerun")
	}
	if !insecure && clientCertName == "" {
		return fmt.Errorf("argument client certificate file can't be " +
			"empty when insecure is false, check program usage and rerun")
	}

	tree, err := getTomlTree(configFile)
	if err != nil {
		return fmt.Errorf("failed to load the config file as toml tree, error: %v", err)
	}

	// auth
	tree.SetPath([]string{"plugins", "io.containerd.grpc.v1.cri", "registry", "configs", registryURL, "auth", "username"}, username)
	tree.SetPath([]string{"plugins", "io.containerd.grpc.v1.cri", "registry", "configs", registryURL, "auth", "password"}, password)
	tree.SetPath([]string{"plugins", "io.containerd.grpc.v1.cri", "registry", "configs", registryURL, "auth", "auth"}, "")
	tree.SetPath([]string{"plugins", "io.containerd.grpc.v1.cri", "registry", "configs", registryURL, "auth", "identitytoken"}, "")

	// tls
	tree.SetPath([]string{"plugins", "io.containerd.grpc.v1.cri", "registry", "configs", registryURL, "tls", "ca_file"}, clientCertName)
	tree.SetPath([]string{"plugins", "io.containerd.grpc.v1.cri", "registry", "configs", registryURL, "tls", "insecure_skip_verify"}, insecure)

	if err := persistTomlTree(configFile, tree); err != nil {
		return err
	}

	log.Debug("registry settings added successfully")
	return nil
}

// AddNvidiaConfig - adds nvidia settings in config file
func AddNvidiaConfig(configFile string) error {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("specified config file %s not exists locally, error: %v", configFile, err)
	}

	tree, err := getTomlTree(configFile)
	if err != nil {
		return fmt.Errorf("failed to load the config file as toml tree, error: %v", err)
	}

	tree.SetPath([]string{"plugins", "io.containerd.grpc.v1.cri", "containerd", "default_runtime_name"}, "nvidia")
	tree.SetPath([]string{"plugins", "io.containerd.grpc.v1.cri", "containerd", "runtimes", "runc", "options", "SystemdCgroup"}, true)
	tree.SetPath([]string{"plugins", "io.containerd.grpc.v1.cri", "containerd", "runtimes", "nvidia", "privileged_without_host_devices"}, false)
	tree.SetPath([]string{"plugins", "io.containerd.grpc.v1.cri", "containerd", "runtimes", "nvidia", "runtime_engine"}, "")
	tree.SetPath([]string{"plugins", "io.containerd.grpc.v1.cri", "containerd", "runtimes", "nvidia", "runtime_root"}, "")
	tree.SetPath([]string{"plugins", "io.containerd.grpc.v1.cri", "containerd", "runtimes", "nvidia", "runtime_type"}, "io.containerd.runc.v2")
	tree.SetPath([]string{"plugins", "io.containerd.grpc.v1.cri", "containerd", "runtimes", "nvidia", "options","BinaryName"}, "/usr/bin/nvidia-container-runtime")
	tree.SetPath([]string{"plugins", "io.containerd.grpc.v1.cri", "containerd", "runtimes", "nvidia", "options","SystemdCgroup"}, true)

	if err := persistTomlTree(configFile, tree); err != nil {
		return err
	}

	log.Debug("nvidia settings added successfully")
	return nil
}

// DeleteRegistryConfig - deletes registry settings from config file
func DeleteRegistryConfig(configFile, registryURL string) error {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("specified config file %s not exists locally, error: %v", configFile, err)
	}
	if registryURL == "" {
		return fmt.Errorf("argument registry" +
			" url value can't be empty for deleting registry from config file, check program usage and rerun")
	}

	tree, err := getTomlTree(configFile)
	if err != nil {
		return fmt.Errorf("failed to load the config file as toml tree, error: %v", err)
	}

	// check if key path exists
	if !tree.HasPath([]string{"plugins", "io.containerd.grpc.v1.cri", "registry", "configs", registryURL}) {
		return fmt.Errorf("specified registry path not found in the config file")
	}
	// delete
	tree.DeletePath([]string{"plugins", "io.containerd.grpc.v1.cri", "registry", "configs", registryURL})

	if err := persistTomlTree(configFile, tree); err != nil {
		return err
	}

	log.Debug("nvidia settings deleted successfully")
	return nil
}

// DeleteNvidiaConfig - deletes nvidia settings from config file
func DeleteNvidiaConfig(configFile string) error {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("specified config file %s not exists locally, error: %v", configFile, err)
	}

	tree, err := getTomlTree(configFile)
	if err != nil {
		return fmt.Errorf("failed to load the config file as toml tree, error: %v", err)
	}

	// revert these two
	tree.SetPath([]string{"plugins", "io.containerd.grpc.v1.cri", "containerd", "default_runtime_name"}, "runc")
	tree.SetPath([]string{"plugins", "io.containerd.grpc.v1.cri", "containerd", "runtimes", "runc", "options", "SystemdCgroup"}, false)

	// check if key path exists
	if !tree.HasPath([]string{"plugins", "io.containerd.grpc.v1.cri", "containerd", "runtimes", "nvidia"}) {
		return fmt.Errorf("specified nvidia path not found in the config file")
	}
	// delete
	tree.DeletePath([]string{"plugins", "io.containerd.grpc.v1.cri", "containerd", "runtimes", "nvidia"})

	if err := persistTomlTree(configFile, tree); err != nil {
		return err
	}

	log.Debug("nvidia settings deleted successfully")
	return nil
}

func getTomlTree(configFile string) (*toml.Tree, error) {
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file, error: %v", err)
	}

	content := string(bytes)
	tree, err := toml.Load(content)
	if err != nil {
		return nil, fmt.Errorf("toml library failed to load config from the file content, error: %v", err)
	}
	return tree, nil
}

func persistTomlTree(configFile string, tree *toml.Tree) error {
	str, err := tree.ToTomlString()
	if err != nil {
		return fmt.Errorf("toml library failed to convert config to a string, error: %v", err)
	}

	data := []byte(str)
	if err := ioutil.WriteFile(configFile, data, 0666); err != nil {
		return fmt.Errorf("failed to write config to a file, error: %v", err)
	}

	log.Info("updates written to file successfully")
	return nil
}