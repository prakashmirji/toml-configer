### containerd-migration

This repo contains source and design doc for containerd migration project

### How to update containerd config.toml with airgap settings
Follow these commands
```
git clone <this repo>
cd to cmd directory
make -f Makefile build-mac or
make -f Makefile build-linux
```

### How to test
Here is the sample execution command

````
<binary name> -username=user1 -password=testpw -clientcert="abc.pem" -insecure=true -registryurl="demo.hpe.com -loglevel=debug"
  
or run go test
go test ./...
````