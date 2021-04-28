# docdb-api

Provides restful API access to the AWS DocumentDB service.

## Endpoints

```
GET /v1/docdb/ping
GET /v1/docdb/version
GET /v1/docdb/metrics

GET /v1/docdb/{account}
GET /v1/docdb/{account}/name
PUT /v1/docdb/{account}/{name}
DELETE /v1/docdb/{account}/{name}
```

## Authentication

Authentication is accomplished via an encrypted pre-shared key passed via the `X-Auth-Token` header.

## Usage

### Lists all docdbs

GET

/v1/docdb/{account}

### Get a single docdb

GET

/v1/docdb/{account}/{name}


### Create docdb cluster and instances

PUT

/v1/docdb/{account}/{name}

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
    { "Key": "MoneyMatters", "Value": "IT"}
  ]
}
```

### Delete docdb cluster and instances

DELETE

/v1/docdb/{account}/{name}

```JSON
{
  "ClusterName": "mytest-docdb",
  "InstanceNames":
    [
        "mytest-docdb-1",
        "mytest-docdb-2",
        "mytest-docdb-3"
    ],
  "SkipFinalSnapshot": true
}
```

## Author

Darryl Wisneski <darryl.wisneski@yale.edu>

## License

GNU Affero General Public License v3.0 (GNU AGPLv3)  
Copyright © 2021 Yale University
