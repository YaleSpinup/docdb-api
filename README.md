# docdb-api

Provides RESTful API access to the AWS DocumentDB service.

## Endpoints

```
GET /v1/docdb/ping
GET /v1/docdb/version
GET /v1/docdb/metrics

GET /v1/docdb/flywheel?task=xxx[&task=yyy&task=zzz]

POST /v1/docdb/{account}
GET /v1/docdb/{account}
GET /v1/docdb/{account}/{name}
PUT /v1/docdb/{account}/{name}
DELETE /v1/docdb/{account}/{name}?snapshot=[true|false]
```

## Authentication

Authentication is accomplished via an encrypted pre-shared key passed in the `X-Auth-Token` header.

## Usage

### Create docdb cluster

Create requests are asynchronous and return a task ID in the header `X-Flywheel-Task`. This header can be used to get the task information and logs from the flywheel HTTP endpoint.

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
| **202 Accepted**              | success accepting docdb request |
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
    "Cluster": {
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
    },
    "Tags": [
        {
            "Key": "spinup:flavor",
            "Value": "docdb"
        },
        {
            "Key": "CreatedBy",
            "Value": "me"
        },
        {
            "Key": "spinup:org",
            "Value": "localdev"
        },
        {
            "Key": "spinup:type",
            "Value": "database"
        }
    ]
}
```

### Modify docdb cluster

The modify request can be used to change the master password for the DocumentDB cluster, or other parameters, such as `BackupRetentionPeriod`, `EngineVersion` or `DBInstanceClass`

PUT `/v1/docdb/{account}/{name}`

```json
{
  "BackupRetentionPeriod": 2,
  "DBInstanceClass": "db.r5.large",
  "MasterUserPassword": "newexamplepassword"
}```

| Response Code                 | Definition                      |
| ----------------------------- | --------------------------------|
| **200 OK**                    | success modifying docdb cluster |
| **400 Bad Request**           | badly formed request            |
| **403 Forbidden**             | bad token or fail to assume role|
| **404 Not Found**             | account not found               |
| **500 Internal Server Error** | a server error occurred         |

#### Example modify response
```json
{
    "Cluster": {
        "AssociatedRoles": null,
        "AvailabilityZones": [
            "us-east-1d",
            "us-east-1a",
            "us-east-1b"
        ],
        "BackupRetentionPeriod": 2,
        "ClusterCreateTime": "2021-10-08T15:01:01.073Z",
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
        "DBSubnetGroup": "spinup-spintg-docdb-subnetgroup",
        "DbClusterResourceId": "cluster-QX7X4WBDAQ3S44RXW226MLMYAM",
        "DeletionProtection": false,
        "EarliestRestorableTime": "2021-10-08T15:01:46.142Z",
        "EnabledCloudwatchLogsExports": null,
        "Endpoint": "mydocdb.cluster-c9ukc6s0rmbg.us-east-1.docdb.amazonaws.com",
        "Engine": "docdb",
        "EngineVersion": "4.0.0",
        "HostedZoneId": "ZNKXH85TT8WVW",
        "KmsKeyId": "arn:aws:kms:us-east-1:123456789012:key/92c75e09-8fcb-4e65-aba4-eed7f4a012f7",
        "LatestRestorableTime": "2021-10-08T15:31:03.749Z",
        "MasterUsername": "dadmin",
        "MultiAZ": false,
        "PercentProgress": null,
        "Port": 27017,
        "PreferredBackupWindow": "05:59-06:29",
        "PreferredMaintenanceWindow": "tue:07:43-tue:08:13",
        "ReadReplicaIdentifiers": null,
        "ReaderEndpoint": "mydocdb.cluster-ro-c9ukc6s0rmbg.us-east-1.docdb.amazonaws.com",
        "ReplicationSourceIdentifier": null,
        "Status": "available",
        "StorageEncrypted": true,
        "VpcSecurityGroups": [
            {
                "Status": "active",
                "VpcSecurityGroupId": "sg-0abcdef1234567890"
            }
        ]
    },
    "Instances": [
        {
            "AutoMinorVersionUpgrade": true,
            "AvailabilityZone": "us-east-1a",
            "BackupRetentionPeriod": 2,
            "CACertificateIdentifier": "rds-ca-2019",
            "DBClusterIdentifier": "mydocdb",
            "DBInstanceArn": "arn:aws:rds:us-east-1:123456789012:db:mydocdb-1",
            "DBInstanceClass": "db.r5.large",
            "DBInstanceIdentifier": "mydocdb-1",
            "DBInstanceStatus": "available",
            "DBSubnetGroup": {
                "DBSubnetGroupArn": null,
                "DBSubnetGroupDescription": "spinup-spintg-docdb-subnetgroup",
                "DBSubnetGroupName": "spinup-spintg-docdb-subnetgroup",
                "SubnetGroupStatus": "Complete",
                "Subnets": [
                    {
                        "SubnetAvailabilityZone": {
                            "Name": "us-east-1d"
                        },
                        "SubnetIdentifier": "subnet-01234567",
                        "SubnetStatus": "Active"
                    },
                    {
                        "SubnetAvailabilityZone": {
                            "Name": "us-east-1a"
                        },
                        "SubnetIdentifier": "subnet-01234568",
                        "SubnetStatus": "Active"
                    }
                ],
                "VpcId": "vpc-8bb612ec"
            },
            "DbiResourceId": "db-UUZG3SJ7BBBU4EJV5JW663DMIA",
            "EnabledCloudwatchLogsExports": null,
            "Endpoint": {
                "Address": "mydocdb-1.c9ukc6s0rmbg.us-east-1.docdb.amazonaws.com",
                "HostedZoneId": "ZNKXH85TT8WVW",
                "Port": 27017
            },
            "Engine": "docdb",
            "EngineVersion": "4.0.0",
            "InstanceCreateTime": "2021-10-08T15:06:00.146Z",
            "KmsKeyId": "arn:aws:kms:us-east-1:123456789012:key/92c73e09-9fcb-4e65-aba4-eed7f4a012f7",
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
            "PreferredBackupWindow": "05:59-06:29",
            "PreferredMaintenanceWindow": "mon:03:59-mon:04:29",
            "PromotionTier": 1,
            "PubliclyAccessible": false,
            "StatusInfos": null,
            "StorageEncrypted": true,
            "VpcSecurityGroups": [
                {
                    "Status": "active",
                    "VpcSecurityGroupId": "sg-0abcdef1234567890"
                }
            ]
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

### Get task information for asynchronous tasks

The status of a new task will initially be `running` and then change to either `failed` or `completed`

GET `/v1/docdb/flywheel?task=xxx[&task=yyy&task=zzz]`

```json
{
    "b403ea9a-a49e-4c4e-a05e-0743f0593c55": {
        "checkin_at": "2021-08-27T21:56:09.332745Z",
        "completed_at": "2021-08-27T21:56:09.539952Z",
        "created_at": "2021-08-27T21:55:05.771054Z",
        "id": "b403ea9a-a49e-4c4e-a05e-0743f0593c55",
        "status": "completed",
        "events": [
            "2021-08-27T21:55:07.77093Z starting task b403ea9a-a49e-4c4e-a05e-0743f0593c55",
            "2021-08-27T21:55:08.098951Z 2021-08-27T21:55:07.902583Z checkin task b403ea9a-a49e-4c4e-a05e-0743f0593c55",
            "2021-08-27T21:55:08.155465Z requested creation of docdb cluster myDocDA",
            "2021-08-27T21:55:08.326608Z 2021-08-27T21:55:08.227748Z checkin task b403ea9a-a49e-4c4e-a05e-0743f0593c55",
            "2021-08-27T21:55:08.386625Z checking if docdb cluster myDocDA is available before continuing",
            "2021-08-27T21:55:08.519555Z 2021-08-27T21:55:08.439783Z checkin task b403ea9a-a49e-4c4e-a05e-0743f0593c55",
            "2021-08-27T21:55:08.584438Z docdb cluster myDocDA is not yet available (creating)",
            "2021-08-27T21:55:10.823559Z 2021-08-27T21:55:10.723116Z checkin task b403ea9a-a49e-4c4e-a05e-0743f0593c55",
            "2021-08-27T21:55:10.879328Z checking if docdb cluster myDocDA is available before continuing",
            . . .
            "2021-08-27T21:56:09.214998Z 2021-08-27T21:56:09.120696Z checkin task b403ea9a-a49e-4c4e-a05e-0743f0593c55",
            "2021-08-27T21:56:09.26763Z checking if docdb cluster myDocDA is available before continuing",
            "2021-08-27T21:56:09.424279Z 2021-08-27T21:56:09.332745Z checkin task b403ea9a-a49e-4c4e-a05e-0743f0593c55",
            "2021-08-27T21:56:09.47998Z docdb cluster myDocDA is available",
            "2021-08-27T21:56:09.631862Z 2021-08-27T21:56:09.539952Z complete task b403ea9a-a49e-4c4e-a05e-0743f0593c55"
        ]
    }
}
```

## Author

Darryl Wisneski <darryl.wisneski@yale.edu>
Tenyo Grozev <tenyo.grozev@yale.edu>

## License

GNU Affero General Public License v3.0 (GNU AGPLv3)  
Copyright © 2021 Yale University
