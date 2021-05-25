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

// DocumentDBCreateHandler creates a documentDB cluster and instance(s)
func (s *server) DocumentDBCreateHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	name := vars["name"]

	log.Infof("creating documentBD cluster and instance(s): %s\n", name)

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

	// read the input against our struct in api/types.go
	req := CreateDocDB{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		msg := fmt.Sprintf("cannot decode body into create documentdb input: %s", err)
		handleError(w, apierror.New(apierror.ErrBadRequest, msg, err))
		return
	}

	orch := newDocDBOrchestrator(
		db.New(db.WithSession(sess.Session)),
		s.org,
	)

	resp, err := orch.createDocumentDB(r.Context(), name, &req)
	if err != nil {
		handleError(w, errors.Wrap(err, "failed to create documentDBs"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// DocumentDBDeleteHandler deletes a DocumentDB cluster and instance(s)
func (s *server) DocumentDBDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	name := vars["name"]

	log.Infof("delete documentBD: %s\n", name)

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

	// read the input against our struct in api/types.go
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

	resp, err := orch.deleteDocumentDB(r.Context(), name, &req)
	if err != nil {
		handleError(w, errors.Wrap(err, "failed to delete documentDBs"))
		return
	}

	j, err := json.Marshal(resp)
	if err != nil {
		handleError(w, apierror.New(apierror.ErrInternalError, "failed to marshal json", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// DocumentDbListHandler lists documentDBs
func (s *server) DocumentDBListHandler(w http.ResponseWriter, r *http.Request) {
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

	orch := newDocDBOrchestrator(
		db.New(db.WithSession(sess.Session)),
		s.org,
	)

	resp, err := orch.listDocumentDB(r.Context())
	if err != nil {
		handleError(w, errors.Wrap(err, "failed to list documentDBs"))
		return
	}

	j, err := json.Marshal(resp)
	if err != nil {
		handleError(w, apierror.New(apierror.ErrInternalError, "failed to marshal json", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}

// DocumentDbGetHandler gets a single named documentDB
func (s *server) DocumentDBGetHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	name := vars["name"]

	log.Infoln("get documentDB")

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

	orch := newDocDBOrchestrator(
		db.New(db.WithSession(sess.Session)),
		s.org,
	)

	resp, err := orch.getDocumentDB(r.Context(), name)
	if err != nil {
		msg := fmt.Sprintf("failed to get documentDB: %s\n", name)
		handleError(w, apierror.New(apierror.ErrBadRequest, msg, err))
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
