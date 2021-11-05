package docdb

import (
	"testing"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/docdb"
	"github.com/pkg/errors"
)

func TestErrCode(t *testing.T) {
	apiErrorTestCases := map[string]string{
		"": apierror.ErrBadRequest,

		docdb.ErrCodeDBClusterQuotaExceededFault:         apierror.ErrLimitExceeded,
		docdb.ErrCodeDBParameterGroupQuotaExceededFault:  apierror.ErrLimitExceeded,
		docdb.ErrCodeDBSubnetGroupQuotaExceededFault:     apierror.ErrLimitExceeded,
		docdb.ErrCodeDBSubnetQuotaExceededFault:          apierror.ErrLimitExceeded,
		docdb.ErrCodeEventSubscriptionQuotaExceededFault: apierror.ErrLimitExceeded,
		docdb.ErrCodeGlobalClusterQuotaExceededFault:     apierror.ErrLimitExceeded,
		docdb.ErrCodeInstanceQuotaExceededFault:          apierror.ErrLimitExceeded,
		docdb.ErrCodeSharedSnapshotQuotaExceededFault:    apierror.ErrLimitExceeded,
		docdb.ErrCodeSnapshotQuotaExceededFault:          apierror.ErrLimitExceeded,
		docdb.ErrCodeStorageQuotaExceededFault:           apierror.ErrLimitExceeded,

		docdb.ErrCodeDBClusterAlreadyExistsFault:         apierror.ErrConflict,
		docdb.ErrCodeDBClusterSnapshotAlreadyExistsFault: apierror.ErrConflict,
		docdb.ErrCodeDBInstanceAlreadyExistsFault:        apierror.ErrConflict,
		docdb.ErrCodeDBParameterGroupAlreadyExistsFault:  apierror.ErrConflict,
		docdb.ErrCodeDBSnapshotAlreadyExistsFault:        apierror.ErrConflict,
		docdb.ErrCodeDBSubnetGroupAlreadyExistsFault:     apierror.ErrConflict,
		docdb.ErrCodeGlobalClusterAlreadyExistsFault:     apierror.ErrConflict,
		docdb.ErrCodeSubnetAlreadyInUse:                  apierror.ErrConflict,
		docdb.ErrCodeSubscriptionAlreadyExistFault:       apierror.ErrConflict,
		docdb.ErrCodeDBUpgradeDependencyFailureFault:     apierror.ErrConflict,

		docdb.ErrCodeAuthorizationNotFoundFault:         apierror.ErrBadRequest,
		docdb.ErrCodeDBSubnetGroupDoesNotCoverEnoughAZs: apierror.ErrBadRequest,

		docdb.ErrCodeDBClusterNotFoundFault:         apierror.ErrNotFound,
		docdb.ErrCodeDBClusterSnapshotNotFoundFault: apierror.ErrNotFound,
		docdb.ErrCodeDBInstanceNotFoundFault:        apierror.ErrNotFound,
		docdb.ErrCodeDBSnapshotNotFoundFault:        apierror.ErrNotFound,
		docdb.ErrCodeDBSecurityGroupNotFoundFault:   apierror.ErrNotFound,
		docdb.ErrCodeDBSubnetGroupNotFoundFault:     apierror.ErrNotFound,
		docdb.ErrCodeGlobalClusterNotFoundFault:     apierror.ErrNotFound,
		docdb.ErrCodeResourceNotFoundFault:          apierror.ErrNotFound,

		docdb.ErrCodeInsufficientDBClusterCapacityFault:      apierror.ErrInternalError,
		docdb.ErrCodeInsufficientDBInstanceCapacityFault:     apierror.ErrInternalError,
		docdb.ErrCodeInsufficientStorageClusterCapacityFault: apierror.ErrInternalError,
		docdb.ErrCodeInvalidDBClusterSnapshotStateFault:      apierror.ErrInternalError,
		docdb.ErrCodeInvalidDBClusterStateFault:              apierror.ErrInternalError,
		docdb.ErrCodeInvalidDBInstanceStateFault:             apierror.ErrInternalError,
		docdb.ErrCodeInvalidDBParameterGroupStateFault:       apierror.ErrInternalError,
		docdb.ErrCodeInvalidDBSecurityGroupStateFault:        apierror.ErrInternalError,
		docdb.ErrCodeInvalidDBSnapshotStateFault:             apierror.ErrInternalError,
		docdb.ErrCodeInvalidDBSubnetGroupStateFault:          apierror.ErrInternalError,
		docdb.ErrCodeInvalidDBSubnetStateFault:               apierror.ErrInternalError,
		docdb.ErrCodeInvalidEventSubscriptionStateFault:      apierror.ErrInternalError,
		docdb.ErrCodeInvalidGlobalClusterStateFault:          apierror.ErrInternalError,
		docdb.ErrCodeInvalidRestoreFault:                     apierror.ErrInternalError,
		docdb.ErrCodeInvalidSubnet:                           apierror.ErrInternalError,
		docdb.ErrCodeInvalidVPCNetworkStateFault:             apierror.ErrInternalError,
	}

	for awsErr, apiErr := range apiErrorTestCases {
		expected := apierror.New(apiErr, "test error", awserr.New(awsErr, awsErr, nil))
		err := ErrCode("test error", awserr.New(awsErr, awsErr, nil))

		var aerr apierror.Error
		if !errors.As(err, &aerr) {
			t.Errorf("expected aws error %s to be an apierror.Error %s, got %s", awsErr, apiErr, err)
		}

		if aerr.String() != expected.String() {
			t.Errorf("expected error '%s', got '%s'", expected, aerr)
		}
	}

	err := ErrCode("test error", errors.New("Unknown"))
	if aerr, ok := errors.Cause(err).(apierror.Error); ok {
		t.Logf("got apierror '%s'", aerr)
	} else {
		t.Errorf("expected unknown error to be an apierror.ErrInternalError, got %s", err)
	}
}
