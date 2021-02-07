# Overview
This repo contains the results of the coding challenge for the mailgun interview. It is mostly a go program with some load testing using Node.js. This readme contains details of how to interact with the repository as well as the different components.

# Architecture
This problem lends itself well to horizontally scaling the application servers, and does not require any state on them. They are also fairly lightweight making them perfect for docker. With the containers sitting behind a load balancer. The key scaling bottleneck with this approach will be the database. The database choice I made was Postgres, partly due to familiarity and partly as it works pretty well with an upsert approach for updating records. For scaling Postgresql there are a few options. As 1 of the calls being handled is read-only, the use of a read-replica could help offload load from the write server. Scaling reads is a bit more challenging, but there are a few options. The first would be to batch updates to the database, specifically updates to a single domain. This works well if you the traffic pattern results in updates to a single domain coming in waves, but would result in a delay of data being available in the database. Alternatives include approaches like sharding, but are likely unnecessary with the load mentioned (100k events per second across 10 million domains).

# Running

## Local

### Setup
You will want two different shell windows for this.

```
docker-compose up database
migrate -path db/migrations -database postgres://postgres:postgres@localhost:15432/mailgun_dev?sslmode=disable up
go build
docker-compose up api
```

## Tests

This package contains unit tests as well as some load testing.

### Unit Tests
You will want two different shell windows for this.

Unit tests can be run by using:

```
docker-compose up test_database
migrate -path db/migrations -database postgres://postgres:postgres@localhost:25432/mailgun_test?sslmode=disable up
go test ./...
```

Manual testing can be done using commands like these:

```
curl -X PUT http://localhost:8080/events/example.com/delivered
curl -X PUT http://localhost:8080/events/example.com/bounced
curl http://localhost:8080/domains/example.com
```

### Load Testings

Load testing scripts can be seen under the `load_testing` directory. The README.md file in that directory contains details of how to run it. **NOTE:** Load testing will be greatly impacted by your setup. Running everything in docker locally will work, but at lower rates then running it on separate systems.


# Notes
* The service is not currently using https. This could be done by using a load balancer, though this would result in an unencrypted connection from the load balancer to the container.
* Bounced could be a bit (true/false). This would take up marginally less space but would not keep a count of occurences of bounced events for a domain. Storing the count could be useful depending on the use case (or if you wanted to change what is considred a catchall in the future).
* Could use the domain name for the primary key, but if you later wanted to join on this table, it would result in you having to store the domain name everywhere which isn't particularly efficient.
* Currently only minimal DB pool configuration is done. This could be optimized further.
* Errors are being returned, but are intentially devoid of many details. Most of the errors wouldn't help downstream services resolve the issue.
* To get docker-compose working as a single command, you could add a script that waits for the DB. The issue is always having that script in the docker container would not be ideal as it would run everywhere, creating one more potential point of fail during container startup in production.
* Tests could be optimized a lot with additional helper functions.
* Test coverage currently focuses on happy pathn, but should ideally be expanded to cover all validation and error conditions.
* Not a huge fan of having everything title case for JSON responses, but it was the default and works fine.
