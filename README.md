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

```JSON
{
  "AvailabilityZones": ["us-east-1a", "us-east-1b"],
  "DBSubnetGroupName": "mysubnetgroup",
  "Engine": "docdb",
  "MasterUsername": "userfoo",
  "MasterUserPassword": "passbar",
  "Tags": {
    "ChargeOfAccounts": "123xyz",
    "Application": "goodstuff",
  },
}
```

### Delete docdb cluster and instances

DELETE

/v1/docdb/<AWSAccountID>/mytest-docdb

```JSON
{
  "SkipFinalSnapshot": "string",
}
```

## Author

Darryl Wisneski <darryl.wisneski@yale.edu>

## License

GNU Affero General Public License v3.0 (GNU AGPLv3)  
Copyright © 2021 Yale University
