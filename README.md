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
  "DBClusterIdentifier": "bestDocDB",
  "DBSubnetGroupName": "a-subnetgroup-name",
  "DBInstanceClass": "db.t3.medium",
  "Engine": "docdb",
  "InstanceCount": "1",
  "MaintenanceWindow": "Sun:04:00-Sun:04:30",
  "MasterUsername": "foousername",
  "MasterUserPassword": "foobarbizbazboo",
  "Tags": [
    { "Key": "CreatedBy", "Value": "tom"},
    { "Key": "MoneyMatters", "Value": "IT"}
  ]
}
```

Response:

```JSON
{
  "DBClusters": {
    "DBClusterArn": "arn:aws:rds:us-east-1:123456789012:cluster:bestDocDB",
    "DBClusterIdentifier": "bestDocDB",
    "Endpoint": "bestDocDB.cluster-cp7kklfeaq3g.us-east-1.docdb.amazonaws.com",
    "ReaderEndpoint": "bestDocDB.cluster-ro-cp7kklfeaq3g.us-east-1.docdb.amazonaws.com",
    "StorageEncrypted": false,
    "DBSubnetGroup": "default-vpc-0e7363e700630fab5",
    "DBInstances": [
      {
        "AvailabilityZone": "",
        "BackupRetentionPeriod": "",
        "DBInstanceArn": "arn:aws:rds:us-east-1:516855177326:db:dkwrocks-1",
        "DBInstanceClass": "",
        "DBInstanceStatus": "",
        "DBInstanceIdentifier": "dkwrocks-1",
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
        "DBInstanceArn": "arn:aws:rds:us-east-1:123456789012:db:dkwrocks-2",
        "DBInstanceClass": "",
        "DBInstanceStatus": "",
        "DBInstanceIdentifier": "dkwrocks-2",
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
        "DBInstanceArn": "arn:aws:rds:us-east-1:123456789012:db:dkwrocks-3",
        "DBInstanceClass": "",
        "DBInstanceStatus": "",
        "DBInstanceIdentifier": "dkwrocks-3",
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
