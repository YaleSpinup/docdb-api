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

	api.HandleFunc("/{account}", s.ListDocumentDB).Methods(http.MethodGet)
	//api.HandleFunc("/{account}/docdb/{name}", s.ShowDocumentDB).Methods(http.MethodGet)
	api.HandleFunc("/{account}/{name}", s.CreateDocumentDB).Methods(http.MethodPut)
	//api.HandleFunc("/{account}/docdb/{name}", s.ModifyDocumentDB).Methods(http.MethodPut)
	//api.HandleFunc("/{account}/docdb/{name}", s.DeleteDocumentDB).Methods(http.MethodDelete)
}

/*
	// cost endpoints for a space
	api.HandleFunc("/{account}/spaces/{space}", s.SpaceGetHandler).
		Queries("start", "{start}", "end", "{end}").Methods(http.MethodGet)
	api.HandleFunc("/{account}/spaces/{space}", s.SpaceGetHandler).Methods(http.MethodGet)
	api.HandleFunc("/{account}/spaces/{space}/{resourcename}", s.SpaceResourceGetHandler).
		Queries("start", "{start}", "end", "{end}").Methods(http.MethodGet)
	api.HandleFunc("/{account}/spaces/{space}/{resourcename}", s.SpaceResourceGetHandler).Methods(http.MethodGet)

	// metrics endpoints for EC2 instances
	// TODO: deprecated but left for backwards compatability, remove me once the UI is updated
	api.HandleFunc("/{account}/instances/{id}/metrics/graph", s.GetEC2MetricsURLHandler).Methods(http.MethodGet)

	// metrics subrouter - /v1/metrics
	metricsApi := s.router.PathPrefix("/v1/metrics").Subrouter()
	metricsApi.HandleFunc("/ping", s.PingHandler).Methods(http.MethodGet)
	metricsApi.HandleFunc("/version", s.VersionHandler).Methods(http.MethodGet)
	metricsApi.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	// metrics endpoints for EC2 instances
	metricsApi.HandleFunc("/{account}/instances/{id}/graph", s.GetEC2MetricsURLHandler).Methods(http.MethodGet)
	// metrics endpoints for ECS services
	metricsApi.HandleFunc("/{account}/clusters/{cluster}/services/{service}/graph", s.GetECSMetricsURLHandler).Methods(http.MethodGet)
	// metrics endpoints for S3 buckets
	metricsApi.HandleFunc("/{account}/buckets/{bucket}/graph", s.GetS3MetricsURLHandler).Queries("metric", "{metric:(?:BucketSizeBytes|NumberOfObjects)}").Methods(http.MethodGet)
	// metrics endpoints for RDS services
	metricsApi.HandleFunc("/{account}/rds/{type}/{id}/graph", s.GetRDSMetricsURLHandler).Methods(http.MethodGet)
*/

/*
	// elastigroup handlers
	api.HandleFunc("/{account}/elastigroups", s.ElastigroupsListHandler).Methods(http.MethodGet)
	api.HandleFunc("/{account}/elastigroups", s.ElastigroupCreateHandler).Methods(http.MethodPost)
	api.HandleFunc("/{account}/elastigroups/{elastigroup}", s.ElastigroupUpdateHandler).Methods(http.MethodPut)
	api.HandleFunc("/{account}/elastigroups/{elastigroup}", s.ElastigroupShowHandler).Methods(http.MethodGet)
	api.HandleFunc("/{account}/elastigroups/{elastigroup}", s.ElastigroupDeleteHandler).Methods(http.MethodDelete)

	// managedinstance handlers
	api.HandleFunc("/{account}/instances", s.ManagedInstanceListHandler).Methods(http.MethodGet)
	api.HandleFunc("/{account}/instances", s.ManagedInstanceCreateHandler).Methods(http.MethodPost)
	api.HandleFunc("/{account}/instances/{instance}", s.ManagedInstanceUpdateHandler).Methods(http.MethodPut)
	api.HandleFunc("/{account}/instances/{instance}", s.ManagedInstanceShowHandler).Methods(http.MethodGet)
	api.HandleFunc("/{account}/instances/{instance}", s.ManagedInstanceDeleteHandler).Methods(http.MethodDelete)
	api.HandleFunc("/{account}/instances/{instance}/status", s.ManagedInstanceStatusHandler).Methods(http.MethodGet)
	api.HandleFunc("/{account}/volumes", s.ManagedVolumesListHandler).Methods(http.MethodGet)
*/

// f5-api
// api.HandleFunc("/{host}/clientssl", s.ListClientSSLProfiles).Methods(http.MethodGet)
// api.HandleFunc("/{host}/clientssl/{name}", s.ShowClientSSLProfile).Methods(http.MethodGet)
// api.HandleFunc("/{host}/createclientssl/{name}", s.CreateClientSSLProfile).Methods(http.MethodPut)
// api.HandleFunc("/{host}/updateclientssl/{name}", s.ModifyClientSSLProfile).Methods(http.MethodPut)
