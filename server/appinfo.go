package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/newity/glog"
	"github.com/pkg/errors"
)

// AppInfo is a struct with app info
type AppInfo struct {
	Ver          string
	VerSdkFabric string
}

func appInfoHandler(ctx context.Context, appInfo *AppInfo) (http.HandlerFunc, error) {
	log := glog.FromContext(ctx)
	resp, err := json.Marshal(appInfo)
	if err != nil {
		return nil, errors.Wrapf(err, "marshal app info error")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		if _, err := w.Write(resp); err != nil {
			log.Errorf("write appinfo response error: %s", err)
		}
	}, nil
}
