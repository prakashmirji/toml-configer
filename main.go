package main

import (
	"flag"
	"strings"

	log "github.com/sirupsen/logrus"
	uc "github.hpe.com/hpe/containerd-migration/pkg"
)

func main() {
	log.Debug("containerd config file update hit")

	username := flag.String("username", "", "specify username of the registry server")
	password := flag.String("password", "", "specify username of the registry server")
	registryURL := flag.String("registryurl", "", "specify registry server url, include port if required")
	clientCertName := flag.String("clientcert", "", "specify client certificate file name")
	insecure := flag.Bool("insecure", false, "specify the value to skip insecure verification, default is set to false")
	op := flag.String("op", "", "specify the operation type to manipulate the config file, like add_nvidia, add_registry...etc")
	configFile := flag.String("configfile", "", "specify the containerd config file name")
	logl := flag.String("loglevel", "info", "specify the log level")

	flag.Parse()

	logLevel := strings.ToLower(*logl)
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatalf("unable to set log level: %v", err)
	}

	log.SetLevel(level)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)

	log.Debugf("registry username: %s", *username)
	log.Debugf("registry password: %s", *password)
	log.Debugf("registry URL: %s", *registryURL)
	log.Debugf("client certificate file name: %s", *clientCertName)
	log.Debugf("insecure skip verify: %v", *insecure)
	log.Debugf("op value: %v", *op)
	log.Debugf("config file name :%s", *configFile)

	switch {
	case *op == "add_nvidia":
		if err := uc.AddNvidiaConfig(*configFile); err != nil {
			log.Errorf("failed to add nvidia settings, error: %v", err)
		} else {
			log.Info("added nvidia configs successfully")
		}
	case *op == "add_registry":
		if err := uc.UpdateRegistryConfig(*configFile, *username, *password, *registryURL, *clientCertName, *insecure); err != nil {
			log.Errorf("failed to add registry details to the config file, error: %v", err)
		} else {
			log.Info("added registry details successfully")
		}
	case *op == "delete_nvidia":
		if err := uc.DeleteNvidiaConfig(*configFile); err != nil {
			log.Errorf("failed to delete nvidia settings, error: %v", err)
		} else {
			log.Info("deleted nvidia configs successfully")
		}
	case *op == "delete_registry":
		if err := uc.DeleteRegistryConfig(*configFile, *registryURL); err != nil {
			log.Errorf("failed to delete registry settings, error: %v", err)
		} else {
			log.Info("deleted registry configs successfully")
		}
	default:
		log.Infof("specified op value: {%s} is not supported", *op)
		log.Infof("supported op values are 'add_nvidia, add_registry, delete_nvidia, delete_registry'")
		log.Info("check out the program usage options like one of the below")
		prinfUsage()
	}
}

func prinfUsage() {
	log.Info("./bin/container_util -op=add_nvidia -configfile=<config filename>")
	log.Info("./bin/container_util -op=add_registry -username=<username> -password=<password> " +
		"-registryurl <registryurl> -clientcert=<ca file> -insecure=<true/false> -configfile=<config filename>")
	log.Info("./bin/container_util -op=delete_nvidia -configfile=<config filename>")
	log.Info("./bin/container_util -op=delete_registry -registryrul=<registry path> -configfile=<config filename>")
}
