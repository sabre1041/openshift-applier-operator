package handler

import (
	"log"
	"net/http"

	"github.com/redhat-cop/openshift-applier-operator/pkg/errors"
	"github.com/redhat-cop/openshift-applier-operator/pkg/manager"
	"github.com/redhat-cop/openshift-applier-operator/pkg/util"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request, applierManager *manager.ApplierManager) {

	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	// Parse input
	namespace, token, err := util.ParseQueryString(r.URL.Path)

	if err != nil {
		w.WriteHeader(500)
	}

	// First Attempt to locate a Applier Resource
	applier, err := applierManager.FindApplierResourceByToken(namespace, token)

	if err != nil {

		if util.IsErrorMessage(err, errors.NotFound) {
			w.WriteHeader(404)
			return
		}

		log.Printf("%v", err)
		w.WriteHeader(500)
		return
	}

	// Launch Applier Job
	err = applierManager.LaunchApplierJob(applier)

	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(201)

}
