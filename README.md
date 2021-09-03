### toml-configer-sample

This repo contains source to read and update toml file (like containerd config.toml). It uses https://github.com/pelletier/go-toml to deal with toml files.
If you are looking to tool manipulate toml files, you can take a look at this sample


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
