package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/docdb"
	"github.com/gorilla/mux"

	db "github.com/YaleSpinup/docdb-api/docdb"

	log "github.com/sirupsen/logrus"
)

// CreateDocumentDB creates a DocumentDB
func (s *server) CreateDocumentDB(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	name := vars["name"]

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
		DBClusterIdentifier: aws.String(name),
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

// DeleteDocumentDB handles and organizes calls to deletes a DocumentDB cluster
func (s *server) DeleteDocumentDB(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	name := vars["name"]

	log.Infoln("delete documentBD cluster")

	// get request of body
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(w, err)
		return
	}
	defer r.Body.Close()

	data := DeleteDocDB{}
	if err := json.Unmarshal(raw, &data); err != nil {
		handleError(w, err)
	}

	// get assumerole
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

	input := docdb.DeleteDBClusterInput{
		DBClusterIdentifier: aws.String(name),
		SkipFinalSnapshot:   aws.Bool(data.SkipFinalSnapshot),
	}

	out, err := client.DeleteDB(r.Context(), &input)
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

/*
//DeleteDocumentDBCluster deletes documentDB instances and clusters
func (s *server) DeleteDocumentDBCluster(ctx context.Context, data *DeleteDocumentDBRequest) error {

}
*/

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
		"arn:aws:iam::aws:policy/AmazonDocDBReadOnlyAccess",
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
