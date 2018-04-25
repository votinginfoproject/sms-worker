# SMS-Worker

## Description
Grab messages from an AWS SQS queue, determine the action to take,
generate message(s) and send them to Twilio.

## Requirements

- Docker

**OR**

- Golang 1.10
- A directory for Go development configured in $GOPATH, see below
- [Godep](https://github.com/tools/godep)
- [go-bindata](https://github.com/jteeuwen/go-bindata)
- A .env file with the following items...
    - The first three items in the example .env file below (AWS credentials and
      environment) MUST be in that order at the top of the file.

~~~~
ACCESS_KEY_ID=
SECRET_ACCESS_KEY=
ENVIRONMENT=development
QUEUE_PREFIX=vip-sms-app
DB_PREFIX=vip-sms-app-users
CIVIC_API_KEY=
CIVIC_API_ELECTION_ID=
CIVIC_API_OFFICIAL_ONLY=true
TWILIO_SID=
TWILIO_TOKEN=
TWILIO_NUMBER=
PROCS=24
ROUTINES=4
LOGGLY_TOKEN=
NEWRELIC_TOKEN=
~~~~

## Development System Setup
To set your system up to develop this application...

1. Make sure you have everything from the requirements section
2. If you do not have GOPATH configured, set one up pointing to a place you'll develop from. For example, `mkdir ~/dev/go` and then `export GOPATH=~/dev/go`
3. Use godeps to get the dependencies and this repository: `godep get github.com/votinginfoproject/sms-worker`
4. This repo should now be at $GOPATH/src/github.com/votinginfoproject/sms-worker, cd into it.
5. If you're using RBEnv to manage Ruby versions, make sure you have 2.1.2 installed and drop a .ruby-version file in this directory with the contents `2.1.2` to make sure you use the right version.
6. Run `bundle`

If you're not on Linux and AMD64, you'll need to compile Go for Linux since the deploy artifact needs to be this as that's the kind of server it is installed on in AWS.
1. Figure out where your Go install is: `which go` (/usr/local/bin/go if you installed with homebrew)
2. Find the go src path which should be nearby (/usr/local/Cellar/go/1.4.2/libexec/src if installed with homebrew)
4. Make go for linux/AMD64 by running `GOOS=linux GOARCH=amd64 ./make.bash --no-clean`

## Docker Development System Setup

It is potentially easier to develop with [Docker][docker].

To compile and run the project, the typical docker `build` and `run`
commands will work. When running, you will need to have the
environment variables above set, because the docker version does not
use the .env file.

With the environment variables set, the commands are:

```
$ docker build -t sms-worker .
$ docker run sms-worker
```

[docker]: https://www.docker.com/


## Commands
### Run Tests
~~~~
godep go test ./...
~~~~

To run the tests in docker:

- `docker build -t sms-worker .` to build the Docker image
- `docker run -ti --env-file .env sms-worker /bin/bash` to start the
container with an interactive terminal
- `godep go test ./...` to run the tests.


### Generate Go Data File From YAML
~~~~
go-bindata -prefix "data" -pkg "data" -o data/data.go data/raw
~~~~

- Generates a new data/data.go file

### Deploy

See the [sms-compose][sms-compose] repository for deployment
instructions.

[sms-compose]: https://github.com/votinginfoproject/sms-compose
