# docdb-api

Provides RESTful API access to the AWS DocumentDB service.

## Endpoints

```
GET /v1/docdb/ping
GET /v1/docdb/version
GET /v1/docdb/metrics

POST /v1/docdb/{account}
GET /v1/docdb/{account}
GET /v1/docdb/{account}/{name}
DELETE /v1/docdb/{account}/{name}?snapshot=[true|false]
```

## Authentication

Authentication is accomplished via an encrypted pre-shared key passed in the `X-Auth-Token` header.

## Usage

### Create docdb cluster

POST `/v1/docdb/{account}`

```json
{
  "BackupRetentionPeriod": 1,
  "DBClusterIdentifier": "myDocDB",
  "DBInstanceClass": "db.t3.medium",
  "EngineVersion": "4.0.0",
  "InstanceCount": 1,
  "MasterUsername": "dadmin",
  "MasterUserPassword": "examplepassword",
  "SubnetIds": ["subnet-12345678", "subnet-abcdef01"],
  "Tags": [
    { "Key": "CreatedBy", "Value": "me"}
  ],
  "VpcSecurityGroupIds": ["sg-0123456789abcdef0"]
}```

| Response Code                 | Definition                      |
| ----------------------------- | --------------------------------|
| **200 OK**                    | success creating docdb cluster  |
| **400 Bad Request**           | badly formed request            |
| **403 Forbidden**             | bad token or fail to assume role|
| **404 Not Found**             | account not found               |
| **500 Internal Server Error** | a server error occurred         |

#### Example create response
```json
{
    "Cluster": {
        "AssociatedRoles": null,
        "AvailabilityZones": [
            "us-east-1d",
            "us-east-1a",
            "us-east-1b"
        ],
        "BackupRetentionPeriod": 1,
        "ClusterCreateTime": "2021-08-05T14:16:58.019Z",
        "DBClusterArn": "arn:aws:rds:us-east-1:1234567890ab:cluster:mydocdb",
        "DBClusterIdentifier": "mydocdb",
        "DBClusterMembers": null,
        "DBClusterParameterGroup": "default.docdb4.0",
        "DBSubnetGroup": "spinup-local-docdb-subnetgroup",
        "DbClusterResourceId": "cluster-NERRHOWC4GF7JXQZQBCSCAQGCY",
        "DeletionProtection": false,
        "EarliestRestorableTime": null,
        "EnabledCloudwatchLogsExports": null,
        "Endpoint": "mydocdb.cluster-c9ukc6s0rmbg.us-east-1.docdb.amazonaws.com",
        "Engine": "docdb",
        "EngineVersion": "4.0.0",
        "HostedZoneId": "ZNKXH85TT8WVW",
        "KmsKeyId": "arn:aws:kms:us-east-1:1234567890ab:key/fffffff-8fcb-4e65-abb4-eed7f4a012f7",
        "LatestRestorableTime": null,
        "MasterUsername": "dadmin",
        "MultiAZ": false,
        "PercentProgress": null,
        "Port": 27017,
        "PreferredBackupWindow": "05:07-05:37",
        "PreferredMaintenanceWindow": "tue:07:05-tue:07:35",
        "ReaderEndpoint": "mydocdb.cluster-ro-c9ukc6s0rmbg.us-east-1.docdb.amazonaws.com",
        "Status": "creating",
        "StorageEncrypted": true,
        "VpcSecurityGroups": [
            {
                "Status": "active",
                "VpcSecurityGroupId": "sg-0123456789abcdef0"
            }
        ]
    },
    "Instances": [
        {
            "AutoMinorVersionUpgrade": true,
            "AvailabilityZone": null,
            "BackupRetentionPeriod": 1,
            "CACertificateIdentifier": "rds-ca-2019",
            "DBClusterIdentifier": "mydocdb",
            "DBInstanceArn": "arn:aws:rds:us-east-1:1234567890ab:db:mydocdb-1",
            "DBInstanceClass": "db.t3.medium",
            "DBInstanceIdentifier": "mydocdb-1",
            "DBInstanceStatus": "creating",
            "DBSubnetGroup": {
                "DBSubnetGroupArn": null,
                "DBSubnetGroupDescription": "spinup-local-docdb-subnetgroup",
                "DBSubnetGroupName": "spinup-local-docdb-subnetgroup",
                "SubnetGroupStatus": "Complete",
                "Subnets": [
                    {
                        "SubnetAvailabilityZone": {
                            "Name": "us-east-1d"
                        },
                        "SubnetIdentifier": "subnet-abcdef01",
                        "SubnetStatus": "Active"
                    },
                    {
                        "SubnetAvailabilityZone": {
                            "Name": "us-east-1a"
                        },
                        "SubnetIdentifier": "subnet-12345678",
                        "SubnetStatus": "Active"
                    }
                ],
                "VpcId": "vpc-12345678"
            },
            "DbiResourceId": "db-WWW2TBGBJKAINEFJOOL6YBP3LM",
            "EnabledCloudwatchLogsExports": null,
            "Endpoint": null,
            "Engine": "docdb",
            "EngineVersion": "4.0.0",
            "InstanceCreateTime": null,
            "KmsKeyId": "arn:aws:kms:us-east-1:1234567890ab:key/ffffffff-8fcb-4e65-abb4-eed7f4a012f7",
            "LatestRestorableTime": null,
            "PendingModifiedValues": {
                "AllocatedStorage": null,
                "BackupRetentionPeriod": null,
                "CACertificateIdentifier": null,
                "DBInstanceClass": null,
                "DBInstanceIdentifier": null,
                "DBSubnetGroupName": null,
                "EngineVersion": null,
                "Iops": null,
                "LicenseModel": null,
                "MasterUserPassword": null,
                "MultiAZ": null,
                "PendingCloudwatchLogsExports": null,
                "Port": null,
                "StorageType": null
            },
            "PreferredBackupWindow": "05:07-05:37",
            "PreferredMaintenanceWindow": "fri:07:04-fri:07:34",
            "PromotionTier": 1,
            "PubliclyAccessible": false,
            "StatusInfos": null,
            "StorageEncrypted": true,
            "VpcSecurityGroups": [
                {
                    "Status": "active",
                    "VpcSecurityGroupId": "sg-0123456789abcdef0"
                }
            ]
        }
    ]
}
```

### List all docdb clusters

GET `/v1/docdb/{account}`

| Response Code                 | Definition                       |
| ----------------------------- | ---------------------------------|
| **200 OK**                    | return list of docdb clusters    |
| **403 Forbidden**             | bad token or fail to assume role |
| **404 Not Found**             | account not found                |
| **500 Internal Server Error** | a server error occurred          |

#### Example list response
```json
[
    "mydocdb",
    "yourdocdb"
]
```

### Get details about a docdb cluster

GET `/v1/docdb/{account}/{name}`

| Response Code                 | Definition                       |
| ----------------------------- | ---------------------------------|
| **200 OK**                    | return details of docdb cluster  |
| **403 Forbidden**             | bad token or fail to assume role |
| **404 Not Found**             | account or docdb not found       |
| **500 Internal Server Error** | a server error occurred          |

#### Example get response
```json
{
    "AssociatedRoles": null,
    "AvailabilityZones": [
        "us-east-1f",
        "us-east-1d",
        "us-east-1a"
    ],
    "BackupRetentionPeriod": 1,
    "ClusterCreateTime": "2021-08-05T13:23:03.003Z",
    "DBClusterArn": "arn:aws:rds:us-east-1:123456789012:cluster:mydocdb",
    "DBClusterIdentifier": "mydocdb",
    "DBClusterMembers": [
        {
            "DBClusterParameterGroupStatus": "in-sync",
            "DBInstanceIdentifier": "mydocdb-1",
            "IsClusterWriter": true,
            "PromotionTier": 1
        }
    ],
    "DBClusterParameterGroup": "default.docdb4.0",
    "DBSubnetGroup": "spinup-example-docdb-subnetgroup",
    "DbClusterResourceId": "cluster-IBME365R7OUKGHZYEOHDJWLBSQ",
    "DeletionProtection": false,
    "EarliestRestorableTime": "2021-08-05T13:23:43.551Z",
    "EnabledCloudwatchLogsExports": null,
    "Endpoint": "mydocdb.cluster-z0ukc6s0rmbg.us-east-1.docdb.amazonaws.com",
    "Engine": "docdb",
    "EngineVersion": "4.0.0",
    "HostedZoneId": "ZZXXYY5TT8WVW",
    "KmsKeyId": "arn:aws:kms:us-east-1:123456789012:key/11aa0000-8fcb-4e65-abb4-eed7f4a012f7",
    "LatestRestorableTime": "2021-08-05T13:23:43.551Z",
    "MasterUsername": "dadmin",
    "MultiAZ": false,
    "PercentProgress": null,
    "Port": 27017,
    "PreferredBackupWindow": "08:41-09:11",
    "PreferredMaintenanceWindow": "mon:10:08-mon:10:38",
    "ReaderEndpoint": "mydocdb.cluster-ro-z0ukc6s0rmbg.us-east-1.docdb.amazonaws.com",
    "Status": "available",
    "StorageEncrypted": true,
    "VpcSecurityGroups": [
        {
            "Status": "active",
            "VpcSecurityGroupId": "sg-0abcdef1234567890"
        }
    ]
}
```

### Delete docdb cluster

Specify `snapshot=true` to create a final snapshot before deleting the cluster. By default, no snapshot will be created.

DELETE `/v1/docdb/{account}/{name}?snapshot=[true|false]`

| Response Code                 | Definition                               |
| ----------------------------- | -----------------------------------------|
| **204 Submitted**             | delete request is submitted              |
| **400 Bad Request**           | badly formed request                     |
| **403 Forbidden**             | bad token or fail to assume role         |
| **404 Not Found**             | account or docdb not found               |
| **500 Internal Server Error** | a server error occurred                  |

## Author

Darryl Wisneski <darryl.wisneski@yale.edu>
Tenyo Grozev <tenyo.grozev@yale.edu>

## License

GNU Affero General Public License v3.0 (GNU AGPLv3)  
Copyright © 2021 Yale University
