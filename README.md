# docdb-api

Provides restful API access to the AWS DocumentDB service.

## Endpoints

```
GET /v1/docdb/ping
GET /v1/docdb/version
GET /v1/docdb/metrics

GET /v1/docdb/{account}/docdb
PUT /v1/docdb/{account}/docdb/{name}
DELETE /v1/docdb/{account}/docdb/{name}
```

## Authentication

Authentication is accomplished via an encrypted pre-shared key passed via the `X-Auth-Token` header.

## Usage

### List docdbs

GET

/v1/docdb/<AWSAccountID>

### Create docdb cluster and instances

PUT

/v1/docdb/<AWSAccountID>/mytest-docdb

### Delete docdb cluster and instances

DELETE

/v1/docdb/<AWSAccountID>/mytest-docdb

## Author

Darryl Wisneski <darryl.wisneski@yale.edu>

## License

GNU Affero General Public License v3.0 (GNU AGPLv3)  
Copyright © 2021 Yale University
