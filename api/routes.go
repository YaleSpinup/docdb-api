/*
Copyright © 2021 Yale University

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package api

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *server) routes() {
	api := s.router.PathPrefix("/v1/docdb").Subrouter()
	api.HandleFunc("/ping", s.PingHandler).Methods(http.MethodGet)
	api.HandleFunc("/version", s.VersionHandler).Methods(http.MethodGet)
	api.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	api.Handle("/flywheel", s.flywheel.Handler())

	api.HandleFunc("/{account}", s.DocumentDBCreateHandler).Methods(http.MethodPost)
	api.HandleFunc("/{account}", s.DocumentDBListHandler).Methods(http.MethodGet)
	api.HandleFunc("/{account}/{name}", s.DocumentDBGetHandler).Methods(http.MethodGet)
	api.HandleFunc("/{account}/{name}", s.DocumentDBModifyHandler).Methods(http.MethodPut)
	api.HandleFunc("/{account}/{name}/power", s.DocumentDBStateHandler).Methods(http.MethodPut)
	api.HandleFunc("/{account}/{name}", s.DocumentDBDeleteHandler).Methods(http.MethodDelete)
}
