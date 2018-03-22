+++
title = "Installing using Docker"
description = "Installing Grafana using Docker guide"
keywords = ["grafana", "configuration", "documentation", "docker"]
type = "docs"
[menu.docs]
name = "Installing using Docker"
identifier = "docker"
parent = "installation"
weight = 4
+++

# Installing using Docker

Grafana is very easy to install and run using the offical docker container.

```bash
$ docker run -d -p 3000:3000 grafana/grafana
```

## Configuration

All options defined in conf/grafana.ini can be overridden using environment
variables by using the syntax `GF_<SectionName>_<KeyName>`.
For example:

```bash
$ docker run \
  -d \
  -p 3000:3000 \
  --name=grafana \
  -e "GF_SERVER_ROOT_URL=http://grafana.server.name" \
  -e "GF_SECURITY_ADMIN_PASSWORD=secret" \
  grafana/grafana
```

The back-end web server has a number of configuration options. Go to the
[Configuration]({{< relref "configuration.md" >}}) page for details on all
those options.

## Installing Plugins for Grafana

Pass the plugins you want installed to docker with the `GF_INSTALL_PLUGINS` environment variable as a comma separated list. This will pass each plugin name to `grafana-cli plugins install ${plugin}`.

```bash
docker run \
  -d \
  -p 3000:3000 \
  --name=grafana \
  -e "GF_INSTALL_PLUGINS=grafana-clock-panel,grafana-simple-json-datasource" \
  grafana/grafana
```

## Building a custom Grafana image with pre-installed plugins

Dockerfile:
```Dockerfile
FROM grafana/grafana:5.0.0
ENV GF_PATHS_PLUGINS=/opt/grafana-plugins
RUN mkdir -p $GF_PATHS_PLUGINS
RUN grafana-cli --pluginsDir $GF_PATHS_PLUGINS plugins install grafana-clock-panel
```

Add lines with `RUN grafana-cli ...` for each plugin you wish to install in your custom image. Don't forget to specify what version of Grafana you wish to build from (replace 5.0.0 in the example).

Example of how to build and run:
```bash
docker build -t grafana:5.0.0-custom . 
docker run \
  -d \
  -p 3000:3000 \
  --name=grafana \
  grafana:5.0.0-custom
```

// TODO: mention that GF?PARTHSPLUG is an official env vars

## Running a Specific Version of Grafana

```bash
# specify right tag, e.g. 4.5.2 - see Docker Hub for available tags
$ docker run \
  -d \
  -p 3000:3000 \
  --name grafana \
  grafana/grafana:5.0.3
```

## Configuring AWS Credentials for CloudWatch Support

```bash
$ docker run \
  -d \
  -p 3000:3000 \
  --name=grafana \
  -e "GF_AWS_PROFILES=default" \
  -e "GF_AWS_default_ACCESS_KEY_ID=YOUR_ACCESS_KEY" \
  -e "GF_AWS_default_SECRET_ACCESS_KEY=YOUR_SECRET_KEY" \
  -e "GF_AWS_default_REGION=us-east-1" \
  grafana/grafana
```

You may also specify multiple profiles to `GF_AWS_PROFILES` (e.g.
`GF_AWS_PROFILES=default another`).

Supported variables:

- `GF_AWS_${profile}_ACCESS_KEY_ID`: AWS access key ID (required).
- `GF_AWS_${profile}_SECRET_ACCESS_KEY`: AWS secret access  key (required).
- `GF_AWS_${profile}_REGION`: AWS region (optional).

## Grafana container with persistent storage (recommended)

 // TODO: implicit vs explicit volumes

```
# create /var/lib/grafana as persistent volume storage
docker run -d -v /var/lib/grafana --name grafana-storage busybox:latest

# start grafana
docker run \
  -d \
  -p 3000:3000 \
  --name=grafana \
  --volumes-from grafana-storage \
  grafana/grafana
```

// TODO: how to map a config from your local filesystem
// TODO: how to point to a different folder, is that really neccessary?
// TODO: mention the different env vars that exists specific to the container (in run.sh)

## Grafana container with host binding and running as a different user

...



## Migration from a previous version of the docker container to 5.1+

In Grafana docker containers prior to 5.1 Grafana was run as the Grafana user (id = `104`). In 5.1 we switched over to the `nobody` user (id = `65534`) instead and also made it possible to change user more easily. Unfortunately this may cause issues when upgrading to 5.1 or later if you have files created by previous versions.

There are two possible solutions to this problem. Either you start the new container as the root user and changes ownership from `104` to `65534` or you start the upgraded container as user `104`.

Examples:

### running docker as a different user

`docker run --user 104 grafana/grafana:5.1`


### docker-compose.yml

```Dockerfile
version: "2"

services:
  grafana:
    image: grafana/grafana:5.1
    ports:
      - 3000:3000
    user: "104"
```

### starting the container and changing ownership

```bash
$ docker run -ti --user root --entrypoint bash grafana/grafana:5.1
$ chown -R root:root /etc/grafana && \
  chmod -R a+r /etc/grafana && \
  chown -R nobody:nogroup /var/lib/grafana && \
  chown -R nobody:nogroup /usr/share/grafana
```