/**
Copyright (c) 2016 The ConnectorDB Contributors
Licensed under the MIT license.
**/
package feed

import (
	"connectordb"
	"connectordb/authoperator"
	"connectordb/datastream"
	"connectordb/users"
	"net/http"
	"server/restapi/restcore"

	log "github.com/Sirupsen/logrus"

	"github.com/gorilla/mux"
)

var (
	//EntryLimit is the maximum number of entries that the feeds will display at one time
	EntryLimit = int64(500)
)

//Get the last week's data
func getFeedData(o *authoperator.AuthOperator, writer http.ResponseWriter, request *http.Request, logger *log.Entry) (*users.Stream, datastream.DataRange, error) {
	_, _, _, streampath := restcore.GetStreamPath(request)
	transform := request.URL.Query().Get("transform")

	s, err := o.ReadStream(streampath)
	if err != nil {
		restcore.WriteError(writer, logger, http.StatusForbidden, err, false)
		return nil, nil, err
	}

	dr, err := o.GetStreamIndexRange(streampath, -EntryLimit, 0, transform)
	if err != nil {
		restcore.WriteError(writer, logger, http.StatusInternalServerError, err, true)
		return nil, nil, err
	}

	return s, dr, err
}

//Router returns a fully formed Gorilla router given an optional prefix
func Router(db *connectordb.Database, prefix *mux.Router) *mux.Router {
	if prefix == nil {
		prefix = mux.NewRouter()
	}

	//Allow for the application to match /path and /path/ to the same place.
	prefix.StrictSlash(true)

	prefix.HandleFunc("/{user}/{device}/{stream}.atom", restcore.Authenticator(GetAtom, db)).Methods("GET")
	prefix.HandleFunc("/{user}/{device}/{stream}", restcore.Authenticator(GetAtom, db)).Methods("GET")

	return prefix
}
