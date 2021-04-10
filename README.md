# Spamhaus Response API

## Overview

This is an HTTP server that uses GraphQL as the messaging framework for communication. It sends DNS queries to Spamhaus to check if an IP address has been flagged as malicious.

### Spamhaus Response API Server Procedure

When the server is started the application first connects to the SQLite database included in the application. The GORM package performs an auto-migrate, so if any new columns have been added to the model it will create those columns in the table.

The HTTP server is the started and it listens on the port defined by the `PORT` environment variable.

At this point you can begin sending queries to the `graphql` endpoint. The graphql endpoint requires json data to be sent in the body of the HTTP message and the username and password described below to be sent using HTTP basic authorization.

When IP addresses are sent via the `enqueue` mutation the application reverses the IP and appends `.zen.spamhaus.org`. So, if 192.168.1.1 is provided then `1.1.168.192.zen.spamhaus.org` is queried. The query is performed using `dig` and the command is `dig +short 1.1.168.192.zen.spamhaus.org`. This returns only the response code, if there is a response. Once complete, the response, IP, a new UUID string, a the current date are saved. If a subsequent mutation is run against an existing IP then the updated_at column for that IP will be updated and if the response code is different then that will be updated also.

Lastly, in order to preclude duplicate IP's from being inserted into the database I've set a unique index on the column `ip_address`.

#### GQLGEN

Gqlgen is the primary library used in this application. The two main components used in this package are the `schema.graphql` file and the `gqlgen.yml` file. The schema file defines the Graphql types and inputs. The yaml file is used to define the application models and generate resolvers based on the schema file.

**For this application the schema defines the following**:
- One TYPE of IPDetails which contains attributes of uuid, response_code, ip_address, created_at, and updated_at of non-nullable scalars of ID, String, String, Time, and Time respectively.
- One Input of IP which contains an Address attribute which is a non-nullable String.
- One Query which defines a method of `getIPDetails` that takes a non-nullable IP type and returns a non-nullable IPDetails type. 
- One Mutation which defines a method `enqueue` that takes a non-nullable IP list of non-nullable IP types and returns a list of non-nullable IPDetails that are also non-nullable.

**When the binary `gqlgen` is executed gqlgen will look in gqlgen.yml and perform the following**:
1. Check the /schema directory for any .graphql files
2. Generate the models based on the types and inputs in the /model directory
3. Create resolver.go and schema.resolvers.go in /graph
4. Generate /graph/generated/generated.go
- Command: `gqlgen`

## HOW TO RUN

1. **Build the docker image**
    - From the root of the application directory enter the below.
    - `docker build -t spamhausresponseapi .`
    - The build copies in all of the directories from the application. This is coupled with the `.dockerignore` file which excludes some directories and files.
    - During the build `dig` is installed and `go mod vendor` and `go build` are run as well to speed up application startup during run.
2. **Run the docker image**
    - You must provide 3 environment variables in the runtime for the application to start.
    - `PORT`, `SpamhausUsername`, `SpamhausPassword`
    - Expose the port number you provided in the environment variable.
    - Mount the spamhausresponseapi.db during the run. The path for the mount be an absolute path for both the source and destination.
    - Command:
        ```docker
        docker run --rm -it -v $pwd/spamhausresponseapi.db:/src/spamhausresponseapi.db -e "PORT=8080" -e "SpamhausUsername=secureworks" -e "SpamhausPassword=supersecret" -p 8080:8080/tcp spamhausresponseapi go run server.go
        ```

### Dockerfile

```docker
FROM golang:1.16

WORKDIR /src
COPY ./ ./

RUN apt-get update && apt-get -y install dnsutils
RUN go mod vendor && go build
```

### .dockerignore

```docker
vendor/
README.md
go.sum
.dockerignore
spamhausresponseapi.db
```

## Features

### GraphQL API

The GraphQL API is reachable from the `/graphql` endpoint. It requires HTTP basic authentication with username `secureworks` and password `supersecret`. The endpoint accepts query `getIPDetails` which takes a parameter, `ip`, which is an INPUT of IP and returns an IPDetails type and a mutation `enqueue` which accepts a list of `ip` type and returns a list of IPDetails type.

#### Examples

* getIPDetails
    ```bash
    curl 'http://localhost:8080/graphql' \
    -u secureworks:supersecret \
    -H 'Content-Type: application/json' \
    --data-binary '{"query":"query {getIPDetails (ip: {address: \"46.102.177.99\"}){uuid, ip_address, response_code, created_at, updated_at}}"}' \
    | json_pp
    ```
* enqueue
    ```bash
    curl 'http://localhost:8080/graphql' \
    -u secureworks:supersecret \
    -H 'Content-Type: application/json' \
    --data-binary '{"query":"mutation {enqueue (ip: [{address: \"77.81.86.150\"},{address: \"77.36.62.11\"}]){uuid, ip_address, response_code, created_at, updated_at}}"}'\
    | json_pp
    ```

### Database

The sqlite database used is spamhausresponseapi.db. It is accessible using the GORM package and included the sqlite driver.

## External Libraries

### GORM

GORM is used for connecting to and interacting with the database. It allows queries to be run against the database and offers a measure of protection against SQL injection. The GORM package has an AutoMigrate method which will auto create the table and columns.

### GO-Mocket

SQL mocking library for use with tests. This library also works with GORM. It was used for the unit tests.

### UUID

Package that allows a uuid string to be generated for use in the database. It was used to create new UUID strings for the uuid column in the ip_details table.

### CHI

Lightweight HTTP router that was used to work as a middleware and provide HTTP basic authentication.

### GQLGEN

GraphQL generator package developed for Golang. Recommended by Secureworks for use with this application. Was very helpful in developing schemas, models, and methods for use with the application.

## Tests

The two main tests are unit tests for the `Enqueue` and `GetIPDetails` methods of `schema.resolvers.go`. They provide mock queries and return data to determine if the functions run against the database correctly.

Additionally, there are tests to check IP address validity. The first is a simple regular expression check to see if the IP address inputted conforms to the standard octect format of IPv4 and the second checks if any octet value exceeds 255.

Unit tests also exist to test the functionallity of IP reversal used to query against `.zen.spamhaus.org`.