package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/YaleSpinup/apierror"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	db "github.com/YaleSpinup/docdb-api/docdb"
	"github.com/YaleSpinup/docdb-api/resourcegroupstaggingapi"
)

// DocumentDBCreateHandler creates a documentDB cluster and instance(s)
func (s *server) DocumentDBCreateHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]

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
	req := DocDBCreateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		msg := fmt.Sprintf("cannot decode body into create documentdb input: %s", err)
		handleError(w, apierror.New(apierror.ErrBadRequest, msg, err))
		return
	}

	if req.SubnetIds == nil {
		handleError(w, apierror.New(apierror.ErrBadRequest, "SubnetIds is a required field", nil))
		return
	}

	if len(req.SubnetIds) < 2 {
		handleError(w, apierror.New(apierror.ErrBadRequest, "At least 2 SubnetIds are required", nil))
		return
	}

	orch := newDocDBOrchestrator(
		db.New(db.WithSession(sess.Session)),
		nil,
		s.flywheel,
		s.org,
	)

	resp, task, err := orch.documentDBCreate(r.Context(), &req)
	if err != nil {
		handleError(w, err)
		return
	}

	j, err := json.Marshal(resp)
	if err != nil {
		handleError(w, errors.Wrap(err, "unable to marshal response from the docdb service"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Flywheel-Task", task.ID)
	w.WriteHeader(http.StatusAccepted)
	w.Write(j)
}

// DocumentDBDeleteHandler deletes a DocumentDB cluster and instance(s)
func (s *server) DocumentDBDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	name := vars["name"]

	queries := r.URL.Query()
	snapshot := false
	if len(queries["snapshot"]) > 0 {
		if b, err := strconv.ParseBool(queries["snapshot"][0]); err == nil {
			snapshot = b
		}
	}

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
		nil,
		nil,
		s.org,
	)

	if err := orch.documentDBDelete(r.Context(), name, snapshot); err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// DocumentDbListHandler lists documentDBs
func (s *server) DocumentDBListHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]

	role := fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName)

	sess, err := s.assumeRole(
		r.Context(),
		s.session.ExternalID,
		role,
		"",
		"arn:aws:iam::aws:policy/AmazonDocDBReadOnlyAccess",
		"arn:aws:iam::aws:policy/ResourceGroupsandTagEditorReadOnlyAccess",
	)
	if err != nil {
		msg := fmt.Sprintf("failed to assume role in account: %s", account)
		handleError(w, apierror.New(apierror.ErrForbidden, msg, err))
		return
	}

	orch := newDocDBOrchestrator(
		db.New(db.WithSession(sess.Session)),
		resourcegroupstaggingapi.New(resourcegroupstaggingapi.WithSession(sess.Session)),
		nil,
		s.org,
	)

	resp, err := orch.documentDBList(r.Context())
	if err != nil {
		handleError(w, errors.Wrap(err, "failed to list documentDBs"))
		return
	}

	j, err := json.Marshal(resp)
	if err != nil {
		handleError(w, apierror.New(apierror.ErrInternalError, "failed to marshal json", err))
		return
	}

	w.Header().Set("X-Items", strconv.Itoa(len(resp)))
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
		nil,
		nil,
		s.org,
	)

	resp, err := orch.documentDBDetails(r.Context(), name)
	if err != nil {
		handleError(w, err)
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

// DocumentDBModifyHandler modifies parameters for a documentDB
func (s *server) DocumentDBModifyHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	name := vars["name"]

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

	req := DocDBModifyRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		msg := fmt.Sprintf("cannot decode body into modify documentdb input: %s", err)
		handleError(w, apierror.New(apierror.ErrBadRequest, msg, err))
		return
	}

	orch := newDocDBOrchestrator(
		db.New(db.WithSession(sess.Session)),
		nil,
		nil,
		s.org,
	)

	resp, err := orch.documentDBModify(r.Context(), name, &req)
	if err != nil {
		handleError(w, err)
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
