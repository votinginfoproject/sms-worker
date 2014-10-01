# SMS-Worker

## Description
Grab messages from an AWS SQS queue, determine the action to take,
generate message(s) and send them to Twilio.

## Requirements
- Golang 1.3
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
2. Run `bundle`

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

- Build the binary
- Upload the binary to S3
- Upload all but the first THREE lines of the .env file to S3
- Restart the sms-web process on all instances

### Send Test Message
~~~~
rake test\[environment,number,message\]
~~~~

- Send a test SMS from the specified number
