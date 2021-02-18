# Podistributor

A golang based program to distribute podcast episodes to audience.
[中文介绍](https://easonyang.com/2021/02/19/podistributor-cn-readme/)

## Features

- Process the request of the alias url to real episode resources so that we can easily handle cdn and failover.
- Asynchronously request analysis services to collect the information of the audience.
- Build-in local cache layer for the database to improve performance and minimize the risk of attacks.
- Monitor and metric mechanism based on Prometheus.

## Quick start

### Environment requirement

- MySQL: for data persistence.
- Golang: for compiling the project.
- Nginx(Recommended): reverse proxy the original podistributor server.
- Prometheus(Optional): monitor the server status with the build-in metric mechanism.

### Installation

Import the database and tables into MySQL. You can find the related sql here: [SQL example](https://github.com/MrEasonYang/podistributor/blob/main/podistributor.sql)

Compile and build the project:

```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o podistributor main.go
chmod 773 podistributor
./podistributor -decryptKey <AES decryption key> -configLocation <Config directory path>
```

The `decryptKey` and `configLocation` parameters are required, change [config file example](https://github.com/MrEasonYang/podistributor/blob/main/podistributor-config.yaml) to modify other configuration.

To run as a service, create a service configuration under `/usr/lib/systemd/system/` like this:

```
[Unit]
Description=podistributor
After=network.target

[Service]
Type=simple
User=root
ExecStart=/usr/sbin/podistributor/podistributor -decryptKey <AES decryption key> -configLocation <Config directory path>
StandardOutput=append:/var/log/podistributor.log
StandardError=append:/var/log/podistributor.log
ExecStop=/bin/kill -s QUIT $MAINPID
Restart=always

[Install]
WantedBy=multi-user.target
```

Then we can run podistributor as a service:

```shell
systemctl enable podistributor
systemctl start podistributor
```

### Usage

#### Request episode

Request in this format will be parsed and redirect to the real resource address:

```
<protocal http or https>://<domain or ip>:<Optional listen port if using nginx>/<listenPath>/<podcast unique name>/ep/<episode unique name>
```

The target url will be selected according to the level which actually is the index of the resource uri or url array set in the db.

Additionally, if the main uri list are unavailable, set the flag of the episode backup url to true and the backup resource url will be redirected to according to the backup url access level which is still the index of the url array.

If analysis urls are set, all the url will be requested asynchronously.

#### Access metrics

Podistributor expose the Prometheus client_golang metrics as a http service via the port configured in `podistributor-confi.yaml` with the name of `monitorPort` which has a default value of `11800` .

```
curl http://127.0.0.1:11800/metrics
```

As for Prometheus server, add a new pull job in `/etc/prometheus/prometheus.yml` :
```
...
    - job_name: 'podistributor'
      static_configs:
      - targets: ['127.0.0.1:11800']
        labels:
          instance: pod-instance
...
``` 

## License

[Apache License 2.0](https://github.com/MrEasonYang/podistributor/blob/main/LICENSE)
