# docdb-api

Provides restful API access to the AWS DocumentDB service.

## Endpoints

```
GET /v1/docdb/ping
GET /v1/docdb/version
GET /v1/docdb/metrics

GET /v1/docdb/{account}
PUT /v1/docdb/{account}/{name}
DELETE /v1/docdb/{account}/{name}
```

## Authentication

Authentication is accomplished via an encrypted pre-shared key passed via the `X-Auth-Token` header.

## Usage

### Lists all docdbs

GET

/v1/docdb/<AWSAccountID>

### Get shows a single docdb

GET

/v1/docdb/<AWSAccountID>/<mytest-docdb>


### Create docdb cluster and instances

PUT

/v1/docdb/<AWSAccountID>/<mytest-docdb>

```JSON
{
  "AvailabilityZones": ["us-east-1a","us-east-1d", "us-east-1b"],
  "DBClusterIdentifier": "bestDocDBCluster",
  "DBSubnetGroupName": "a-subnetgroup-name",
  "Engine": "docdb",
  "MasterUsername": "foousername",
  "MasterUserPassword": "foobarbizbazboo",
  "Tags": [
    { "Key": "CreatedBy", "Value": "tom"},
    { "Key": "ChargingAccount", "Value": "xyz"}
  ]
}
```

### Delete docdb cluster and instances

DELETE

/v1/docdb/<AWSAccountID>/<mytest-docdb>

```JSON
{
  "ClusterName": "mytest-docdb",
  "InstanceNames":
    [
        "dkwrocks-1",
        "dkwrocks-2",
        "dkwrocks-3"
    ],
  "SkipFinalSnapshot": true
}
```

## Author

Darryl Wisneski <darryl.wisneski@yale.edu>

## License

GNU Affero General Public License v3.0 (GNU AGPLv3)  
Copyright © 2021 Yale University
