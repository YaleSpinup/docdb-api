# docdb-api

Provides RESTful API access to the AWS DocumentDB service.

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

POST

/v1/docdb/{account}/{name}

```JSON
{
  "AvailabilityZones": ["us-east-1a","us-east-1d", "us-east-1b"],
  "DBClusterIdentifier": "exampleDB",
  "DBSubnetGroupName": "some-subnet-group",
  "DBInstanceClass": "db.t3.medium",
  "Engine": "docdb",
  "InstanceCount": "3",
  "MaintenanceWindow": "Sun:04:00-Sun:04:30",
  "MasterUsername": "foousername",
  "MasterUserPassword": "foobarbizbazboo",
  "Tags": [
    { "Key": "CreatedBy", "Value": "tom"},
    { "Key": "MoneyMattersMost", "Value": "IT"}
  ]
}
```

Response:

```JSON
{
  "DBClusters": {
    "DBClusterArn": "arn:aws:rds:us-east-1:123456789012:cluster:exampleDB",
    "DBClusterIdentifier": "exampleDB",
    "Endpoint": "exampleDB.cluster-somestring.us-east-1.docdb.amazonaws.com",
    "ReaderEndpoint": "exampleDB.cluster-ro-somestring.us-east-1.docdb.amazonaws.com",
    "StorageEncrypted": false,
    "DBSubnetGroup": "some-subnet-group",
    "DBInstances": [
      {
        "AvailabilityZone": "",
        "BackupRetentionPeriod": "",
        "DBInstanceArn": "arn:aws:rds:us-east-1:123456789012:db:exampleDB-1",
        "DBInstanceClass": "",
        "DBInstanceStatus": "",
        "DBInstanceIdentifier": "exampleDB-1",
        "DBSubnetGroup": "",
        "Endpoint": "",
        "Engine": "",
        "EngineVersion": "",
        "InstanceCreateTime": "0001-01-01T00:00:00Z",
        "KmsKeyId": "",
        "ReaderEndpoint": "",
        "StorageEncrypted": false
      },
      {
        "AvailabilityZone": "",
        "BackupRetentionPeriod": "",
        "DBInstanceArn": "arn:aws:rds:us-east-1:123456789012:db:exampleDB-2",
        "DBInstanceClass": "",
        "DBInstanceStatus": "",
        "DBInstanceIdentifier": "exampleDB-2",
        "DBSubnetGroup": "",
        "Endpoint": "",
        "Engine": "",
        "EngineVersion": "",
        "InstanceCreateTime": "0001-01-01T00:00:00Z",
        "KmsKeyId": "",
        "ReaderEndpoint": "",
        "StorageEncrypted": false
      },
      {
        "AvailabilityZone": "",
        "BackupRetentionPeriod": "",
        "DBInstanceArn": "arn:aws:rds:us-east-1:123456789012:db:exampleDB-3",
        "DBInstanceClass": "",
        "DBInstanceStatus": "",
        "DBInstanceIdentifier": "exampleDB-3",
        "DBSubnetGroup": "",
        "Endpoint": "",
        "Engine": "",
        "EngineVersion": "",
        "InstanceCreateTime": "0001-01-01T00:00:00Z",
        "KmsKeyId": "",
        "ReaderEndpoint": "",
        "StorageEncrypted": false
      }
    ]
  }
}

```

### Delete docdb cluster and instances

DELETE

/v1/docdb/{account}/{name}

```JSON
{
  "ClusterName": "exampleDB",
  "InstanceNames":
    [
        "exampleDB-1",
        "exampleDB-2",
        "exampleDB-3"
    ],
  "SkipFinalSnapshot": true
}
```

## Author

Darryl Wisneski <darryl.wisneski@yale.edu>

## License

GNU Affero General Public License v3.0 (GNU AGPLv3)  
Copyright © 2021 Yale University
