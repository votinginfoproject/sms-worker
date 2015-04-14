# SMS-Worker

## Description
Grab messages from an AWS SQS queue, determine the action to take,
generate message(s) and send them to Twilio.

## Requirements
- Golang 1.3
- A directory for Go development configured in $GOPATH, see below
- Ruby 2.1.2
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

## Commands
### Run Tests
~~~~
godep go test ./...
~~~~

### Generate Go Data File From YAML
~~~~
rake gen_asset
~~~~

- Generates a new data/data.go file

### Deploy
~~~~
rake deploy\[environment\]
~~~~

Prerequisites
- The environment has been built with sms-infrastructure
- AWS console access
- The .env file has been built. You can retrieve the current one from s3 in the bucket "vip-sms-#{environment}" name "sms-worker-env".
- You have passwordless ssh access to the vip-sms-app-worker servers as the ubuntu user. To achieve this, obtain the private key and run `ssh-add <private-key-file>` and test it with `ssh ubuntu@<worker-public-ip>`

Deploy will rebuild the code, upload the new binary to S3, figure out which EC2 instances are the right workers for the environment, and restart them with the new binary. It also uploads the current .env to the s3 bucket "vip-sms-#{env}" as "sms-worker-env", so this can be a starting point for a new .env file. 

Steps of the deploy task:
- Build the binary
- Upload the binary to S3
- Upload all but the first THREE lines of the .env file to S3
- Restart the sms-worker process on all instances

### Send Test Message
~~~~
rake test\[environment,number,message\]
~~~~

- Send a test SMS from the specified number
