package docdb

import (
	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/docdb"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func ErrCode(msg string, err error) error {
	if aerr, ok := errors.Cause(err).(awserr.Error); ok {
		switch aerr.Code() {
		case
			"Forbidden":

			return apierror.New(apierror.ErrForbidden, msg, aerr)
		case
			// ErrCodeDBClusterQuotaExceededFault for service response error code
			// "DBClusterQuotaExceededFault".
			//
			// The cluster can't be created because you have reached the maximum allowed
			// quota of clusters.
			docdb.ErrCodeDBClusterQuotaExceededFault,

			// ErrCodeDBParameterGroupQuotaExceededFault for service response error code
			// "DBParameterGroupQuotaExceeded".
			//
			// This request would cause you to exceed the allowed number of parameter groups.
			docdb.ErrCodeDBParameterGroupQuotaExceededFault,

			// ErrCodeDBSubnetGroupQuotaExceededFault for service response error code
			// "DBSubnetGroupQuotaExceeded".
			//
			// The request would cause you to exceed the allowed number of subnet groups.
			docdb.ErrCodeDBSubnetGroupQuotaExceededFault,

			// ErrCodeDBSubnetQuotaExceededFault for service response error code
			// "DBSubnetQuotaExceededFault".
			//
			// The request would cause you to exceed the allowed number of subnets in a
			// subnet group.
			docdb.ErrCodeDBSubnetQuotaExceededFault,

			// ErrCodeEventSubscriptionQuotaExceededFault for service response error code
			// "EventSubscriptionQuotaExceeded".
			//
			// You have reached the maximum number of event subscriptions.
			docdb.ErrCodeEventSubscriptionQuotaExceededFault,

			// ErrCodeGlobalClusterQuotaExceededFault for service response error code
			// "GlobalClusterQuotaExceededFault".
			//
			// The number of global clusters for this account is already at the maximum
			// allowed.
			docdb.ErrCodeGlobalClusterQuotaExceededFault,

			// ErrCodeInstanceQuotaExceededFault for service response error code
			// "InstanceQuotaExceeded".
			//
			// The request would cause you to exceed the allowed number of instances.
			docdb.ErrCodeInstanceQuotaExceededFault,

			// ErrCodeSharedSnapshotQuotaExceededFault for service response error code
			// "SharedSnapshotQuotaExceeded".
			//
			// You have exceeded the maximum number of accounts that you can share a manual
			// DB snapshot with.
			docdb.ErrCodeSharedSnapshotQuotaExceededFault,

			// ErrCodeSnapshotQuotaExceededFault for service response error code
			// "SnapshotQuotaExceeded".
			//
			// The request would cause you to exceed the allowed number of snapshots.
			docdb.ErrCodeSnapshotQuotaExceededFault,

			// ErrCodeStorageQuotaExceededFault for service response error code
			// "StorageQuotaExceeded".
			//
			// The request would cause you to exceed the allowed amount of storage available
			// across all instances.
			docdb.ErrCodeStorageQuotaExceededFault,

			// Limit Exceeded
			"LimitExceeded":

			return apierror.New(apierror.ErrLimitExceeded, msg, aerr)
		case
			// ErrCodeDBClusterAlreadyExistsFault for service response error code
			// "DBClusterAlreadyExistsFault".
			//
			// You already have a cluster with the given identifier.
			docdb.ErrCodeDBClusterAlreadyExistsFault,

			// ErrCodeDBClusterSnapshotAlreadyExistsFault for service response error code
			// "DBClusterSnapshotAlreadyExistsFault".
			//
			// You already have a cluster snapshot with the given identifier.
			docdb.ErrCodeDBClusterSnapshotAlreadyExistsFault,

			// ErrCodeDBInstanceAlreadyExistsFault for service response error code
			// "DBInstanceAlreadyExists".
			//
			// You already have a instance with the given identifier.
			docdb.ErrCodeDBInstanceAlreadyExistsFault,

			// ErrCodeDBParameterGroupAlreadyExistsFault for service response error code
			// "DBParameterGroupAlreadyExists".
			//
			// A parameter group with the same name already exists.
			docdb.ErrCodeDBParameterGroupAlreadyExistsFault,

			// ErrCodeDBSnapshotAlreadyExistsFault for service response error code
			// "DBSnapshotAlreadyExists".
			//
			// DBSnapshotIdentifier is already being used by an existing snapshot.
			docdb.ErrCodeDBSnapshotAlreadyExistsFault,

			// ErrCodeDBSubnetGroupAlreadyExistsFault for service response error code
			// "DBSubnetGroupAlreadyExists".
			//
			// DBSubnetGroupName is already being used by an existing subnet group.
			docdb.ErrCodeDBSubnetGroupAlreadyExistsFault,

			// ErrCodeGlobalClusterAlreadyExistsFault for service response error code
			// "GlobalClusterAlreadyExistsFault".
			//
			// The GlobalClusterIdentifier already exists. Choose a new global cluster identifier
			// (unique name) to create a new global cluster.
			docdb.ErrCodeGlobalClusterAlreadyExistsFault,

			// ErrCodeSubnetAlreadyInUse for service response error code
			// "SubnetAlreadyInUse".
			//
			// The subnet is already in use in the Availability Zone.
			docdb.ErrCodeSubnetAlreadyInUse,

			// ErrCodeSubscriptionAlreadyExistFault for service response error code
			// "SubscriptionAlreadyExist".
			//
			// The provided subscription name already exists.
			docdb.ErrCodeSubscriptionAlreadyExistFault,

			// ErrCodeDBUpgradeDependencyFailureFault for service response error code
			// "DBUpgradeDependencyFailure".
			//
			// The upgrade failed because a resource that the depends on can't be modified.
			docdb.ErrCodeDBUpgradeDependencyFailureFault:

			return apierror.New(apierror.ErrConflict, msg, aerr)
		case
			// ErrCodeDBClusterNotFoundFault for service response error code
			// "DBClusterNotFoundFault".
			//
			// DBClusterIdentifier doesn't refer to an existing cluster.
			docdb.ErrCodeDBClusterNotFoundFault,

			// ErrCodeDBClusterSnapshotNotFoundFault for service response error code
			// "DBClusterSnapshotNotFoundFault".
			//
			// DBClusterSnapshotIdentifier doesn't refer to an existing cluster snapshot.
			docdb.ErrCodeDBClusterSnapshotNotFoundFault,

			// ErrCodeDBInstanceNotFoundFault for service response error code
			// "DBInstanceNotFound".
			//
			// DBInstanceIdentifier doesn't refer to an existing instance.
			docdb.ErrCodeDBInstanceNotFoundFault,

			// ErrCodeDBSnapshotNotFoundFault for service response error code
			// "DBSnapshotNotFound".
			//
			// DBSnapshotIdentifier doesn't refer to an existing snapshot.
			docdb.ErrCodeDBSnapshotNotFoundFault,

			// ErrCodeGlobalClusterNotFoundFault for service response error code
			// "GlobalClusterNotFoundFault".
			//
			// The GlobalClusterIdentifier doesn't refer to an existing global cluster.
			docdb.ErrCodeGlobalClusterNotFoundFault,

			// ErrCodeResourceNotFoundFault for service response error code
			// "ResourceNotFoundFault".
			//
			// The specified resource ID was not found.
			docdb.ErrCodeResourceNotFoundFault,

			// Not found.
			"NotFound":

			return apierror.New(apierror.ErrNotFound, msg, aerr)
		case
			// ErrCodeInsufficientDBClusterCapacityFault for service response error code
			// "InsufficientDBClusterCapacityFault".
			//
			// The cluster doesn't have enough capacity for the current operation.
			docdb.ErrCodeInsufficientDBClusterCapacityFault,

			// ErrCodeInsufficientDBInstanceCapacityFault for service response error code
			// "InsufficientDBInstanceCapacity".
			//
			// The specified instance class isn't available in the specified Availability
			// Zone.
			docdb.ErrCodeInsufficientDBInstanceCapacityFault,

			// ErrCodeInsufficientStorageClusterCapacityFault for service response error code
			// "InsufficientStorageClusterCapacity".
			//
			// There is not enough storage available for the current action. You might be
			// able to resolve this error by updating your subnet group to use different
			// Availability Zones that have more storage available.
			docdb.ErrCodeInsufficientStorageClusterCapacityFault,

			// ErrCodeInvalidDBClusterSnapshotStateFault for service response error code
			// "InvalidDBClusterSnapshotStateFault".
			//
			// The provided value isn't a valid cluster snapshot state.
			docdb.ErrCodeInvalidDBClusterSnapshotStateFault,

			// ErrCodeInvalidDBClusterStateFault for service response error code
			// "InvalidDBClusterStateFault".
			//
			// The cluster isn't in a valid state.
			docdb.ErrCodeInvalidDBClusterStateFault,

			// ErrCodeInvalidDBInstanceStateFault for service response error code
			// "InvalidDBInstanceState".
			//
			// The specified instance isn't in the available state.
			docdb.ErrCodeInvalidDBInstanceStateFault,

			// ErrCodeInvalidDBParameterGroupStateFault for service response error code
			// "InvalidDBParameterGroupState".
			//
			// The parameter group is in use, or it is in a state that is not valid. If
			// you are trying to delete the parameter group, you can't delete it when the
			// parameter group is in this state.
			docdb.ErrCodeInvalidDBParameterGroupStateFault,

			// ErrCodeInvalidDBSecurityGroupStateFault for service response error code
			// "InvalidDBSecurityGroupState".
			//
			// The state of the security group doesn't allow deletion.
			docdb.ErrCodeInvalidDBSecurityGroupStateFault,

			// ErrCodeInvalidDBSnapshotStateFault for service response error code
			// "InvalidDBSnapshotState".
			//
			// The state of the snapshot doesn't allow deletion.
			docdb.ErrCodeInvalidDBSnapshotStateFault,

			// ErrCodeInvalidDBSubnetGroupStateFault for service response error code
			// "InvalidDBSubnetGroupStateFault".
			//
			// The subnet group can't be deleted because it's in use.
			docdb.ErrCodeInvalidDBSubnetGroupStateFault,

			// ErrCodeInvalidDBSubnetStateFault for service response error code
			// "InvalidDBSubnetStateFault".
			//
			// The subnet isn't in the available state.
			docdb.ErrCodeInvalidDBSubnetStateFault,

			// ErrCodeInvalidEventSubscriptionStateFault for service response error code
			// "InvalidEventSubscriptionState".
			//
			// Someone else might be modifying a subscription. Wait a few seconds, and try
			// again.
			docdb.ErrCodeInvalidEventSubscriptionStateFault,

			// ErrCodeInvalidGlobalClusterStateFault for service response error code
			// "InvalidGlobalClusterStateFault".
			//
			// The requested operation can't be performed while the cluster is in this state.
			docdb.ErrCodeInvalidGlobalClusterStateFault,

			// ErrCodeInvalidRestoreFault for service response error code
			// "InvalidRestoreFault".
			//
			// You cannot restore from a virtual private cloud (VPC) backup to a non-VPC
			// DB instance.
			docdb.ErrCodeInvalidRestoreFault,

			// ErrCodeInvalidSubnet for service response error code
			// "InvalidSubnet".
			//
			// The requested subnet is not valid, or multiple subnets were requested that
			// are not all in a common virtual private cloud (VPC).
			docdb.ErrCodeInvalidSubnet,

			// ErrCodeInvalidVPCNetworkStateFault for service response error code
			// "InvalidVPCNetworkStateFault".
			//
			// The subnet group doesn't cover all Availability Zones after it is created
			// because of changes that were made.
			docdb.ErrCodeInvalidVPCNetworkStateFault:

			return apierror.New(apierror.ErrInternalError, msg, err)
		case
			// Service Unavailable
			"ServiceUnavailable":

			return apierror.New(apierror.ErrServiceUnavailable, msg, aerr)
		default:
			return apierror.New(apierror.ErrBadRequest, msg, aerr)
		}
	}

	log.Warnf("uncaught error: %s, returning Internal Server Error", err)
	return apierror.New(apierror.ErrInternalError, msg, err)
}
