package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/YaleSpinup/apierror"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

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

	req := CreateDocDB{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		msg := fmt.Sprintf("cannot decode body into create documentdb input: %s", err)
		handleError(w, apierror.New(apierror.ErrBadRequest, msg, err))
		return
	}

	log.Debugf("GOOGLEY create data input: %s\n", req)

	orch := newDocDBOrchestrator(
		db.New(db.WithSession(sess.Session)),
		s.org,
	)

	resp, err := orch.createDocumentDB(r.Context(), &req)
	if err != nil {
		handleError(w, errors.Wrap(err, "failed to create documentDBs"))
		return
	}

	j, err := json.Marshal(resp)
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

	log.Infoln("delete documentBD cluster")

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

	//req := docdb.DeleteDBClusterInput{}
	req := DeleteDocDB{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		msg := fmt.Sprintf("cannot decode body into delete documentdb input: %s", err)
		handleError(w, apierror.New(apierror.ErrBadRequest, msg, err))
		return
	}

	orch := newDocDBOrchestrator(
		db.New(db.WithSession(sess.Session)),
		s.org,
	)

	resp, err := orch.deleteDocumentDB(r.Context(), &req)
	if err != nil {
		handleError(w, errors.Wrap(err, "failed to delete documentDBs"))
		return
	}

	j, err := json.Marshal(resp)
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

	orch := newDocDBOrchestrator(
		db.New(db.WithSession(sess.Session)),
		s.org,
	)

	//input := docdb.DescribeDBClustersInput{}

	resp, err := orch.listDocumentDB(r.Context())
	if err != nil {
		handleError(w, errors.Wrap(err, "failed to list documentDBs"))
		return
	}

	//client := db.New(db.WithSession(sess.Session))

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

	/*
		input := docdb.DescribeDBClustersInput{}

		out, err := client.ListDB(r.Context(), &input)
		if err != nil {
			handleError(w, err)
			return
		}
	*/

	j, err := json.Marshal(resp)
	if err != nil {
		handleError(w, apierror.New(apierror.ErrBadRequest, "failed to marshal json", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}
