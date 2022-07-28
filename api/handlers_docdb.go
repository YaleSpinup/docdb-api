package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/YaleSpinup/apierror"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// DocumentDBCreateHandler creates a documentDB cluster and instance(s)
func (s *server) DocumentDBCreateHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]

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

	orch, err := s.newDocDBOrchestrator(
		r.Context(),
		&sessionParams{
			role:       fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName),
			policyArns: []string{"arn:aws:iam::aws:policy/AmazonDocDBFullAccess"},
		},
	)
	if err != nil {
		handleError(w, errors.Wrap(err, "unable to create docdb orchestrator"))
		return
	}

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

	orch, err := s.newDocDBOrchestrator(
		r.Context(),
		&sessionParams{
			role:       fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName),
			policyArns: []string{"arn:aws:iam::aws:policy/AmazonDocDBFullAccess"},
		},
	)
	if err != nil {
		handleError(w, errors.Wrap(err, "unable to create docdb orchestrator"))
		return
	}

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

	orch, err := s.newDocDBOrchestrator(
		r.Context(),
		&sessionParams{
			role: fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName),
			policyArns: []string{
				"arn:aws:iam::aws:policy/AmazonDocDBReadOnlyAccess",
				"arn:aws:iam::aws:policy/ResourceGroupsandTagEditorReadOnlyAccess",
			},
		},
	)
	if err != nil {
		handleError(w, errors.Wrap(err, "unable to create docdb orchestrator"))
		return
	}

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

	orch, err := s.newDocDBOrchestrator(
		r.Context(),
		&sessionParams{
			role:       fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName),
			policyArns: []string{"arn:aws:iam::aws:policy/AmazonDocDBReadOnlyAccess"},
		},
	)
	if err != nil {
		handleError(w, errors.Wrap(err, "unable to create docdb orchestrator"))
		return
	}

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

	req := DocDBModifyRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		msg := fmt.Sprintf("cannot decode body into modify documentdb input: %s", err)
		handleError(w, apierror.New(apierror.ErrBadRequest, msg, err))
		return
	}

	orch, err := s.newDocDBOrchestrator(
		r.Context(),
		&sessionParams{
			role:       fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName),
			policyArns: []string{"arn:aws:iam::aws:policy/AmazonDocDBFullAccess"},
		},
	)
	if err != nil {
		handleError(w, errors.Wrap(err, "unable to create docdb orchestrator"))
		return
	}

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

// DocumentDBStateHandler Starts/Stops a DocumentDB cluster and instance(s)
func (s *server) DocumentDBStateHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	name := vars["name"]

	req := &docDBInstanceStateChangeRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		msg := fmt.Sprintf("cannot decode body into change power input: %s", err)
		handleError(w, apierror.New(apierror.ErrBadRequest, msg, err))
		return
	}

	if req.State == "" {
		handleError(w, apierror.New(apierror.ErrBadRequest, "missing required field: state", nil))
		return
	}

	orch, err := s.newDocDBOrchestrator(
		r.Context(),
		&sessionParams{
			role:       fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName),
			policyArns: []string{"arn:aws:iam::aws:policy/AmazonDocDBFullAccess"},
		},
	)
	if err != nil {
		handleError(w, errors.Wrap(err, "unable to create docdb orchestrator"))
		return
	}

	if err := orch.docDBState(r.Context(), req.State, name); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
