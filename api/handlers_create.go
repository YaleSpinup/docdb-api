package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/docdb"
	"github.com/gorilla/mux"

	//db "github.com/YaleSpinup/docdb-api/docdb"
	db "github.com/YaleSpinup/docdb-api/docdb"

	log "github.com/sirupsen/logrus"
)

// CreateDocumentDB creates a DocumentDB
func (s *server) CreateDocumentDB(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]

	log.Infoln("create documentBDs")

	role := fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName)

	sess, err := s.assumeRole(
		r.Context(),
		s.session.ExternalID,
		role,
		"",
		"arn:aws:iam::aws:policy/AmazonDocDBFullAccess",
	)
	if err != nil {
		msg := fmt.Sprintf("failed to assume role in account: %s", account)
		handleError(w, apierror.New(apierror.ErrForbidden, msg, err))
		return
	}

	client := db.New(db.WithSession(sess.Session))

	input := docdb.CreateDBClusterInput{
		AvailabilityZones: []*string{
			aws.String("us-east-1a"),
			aws.String("us-east-1d"),
		},
		DBClusterIdentifier: aws.String("foobar-cluster"),
		DBSubnetGroupName:   aws.String("default-vpc-0e7363e700630fab5"),
		Engine:              aws.String("docdb"),
		MasterUsername:      aws.String("foobarusername"),
		MasterUserPassword:  aws.String("foobarbazboo"),
		//Tags: "[]",
	}

	out, err := client.CreateDB(r.Context(), &input)
	if err != nil {
		handleError(w, err)
		return
	}

	j, err := json.Marshal(out)
	if err != nil {
		handleError(w, apierror.New(apierror.ErrBadRequest, "failed to marshal json", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// ListDocumentDb lists documentDBs
func (s *server) ListDocumentDB(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	log.Infof("account: %s\n", account)

	log.Infoln("list documentDBs")

	role := fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName)

	sess, err := s.assumeRole(
		r.Context(),
		s.session.ExternalID,
		role,
		"",
		"arn:aws:iam::aws:policy/AmazonDocDBFullAccess",
	)
	if err != nil {
		msg := fmt.Sprintf("failed to assume role in account: %s", account)
		handleError(w, apierror.New(apierror.ErrForbidden, msg, err))
		return
	}

	client := db.New(db.WithSession(sess.Session))

	/*
		input := docdb.DescribeDBClustersInput{
			Filters: []*docdb.Filter{
				{
					Name: aws.String("*"),
					Values: []*string{
						aws.String("us-east-1"),
						aws.String("us-east-2"),
						aws.String("us-east-3"),
					},
				},
			},
		}
	*/
	input := docdb.DescribeDBClustersInput{}

	out, err := client.ListDB(r.Context(), &input)
	if err != nil {
		handleError(w, err)
		return
	}

	j, err := json.Marshal(out)
	if err != nil {
		handleError(w, apierror.New(apierror.ErrBadRequest, "failed to marshal json", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}
