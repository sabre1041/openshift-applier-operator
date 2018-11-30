package handler

import (
	"log"
	"net/http"

	"github.com/redhat-cop/openshift-applier-operator/pkg/manager"
	"github.com/redhat-cop/openshift-applier-operator/pkg/util"
	"k8s.io/apimachinery/pkg/api/errors"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request, applierManager *manager.ApplierManager) {

	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	// Parse input
	namespace, name, token, err := util.ParseQueryString(r.URL.Path)

	if err != nil {
		log.Printf("Error Parsing Webhook Query String: %v", err)
		w.WriteHeader(500)
		return
	}

	// First Attempt to locate a Applier Resource
	applier, err := applierManager.FindApplierResourceByToken(namespace, name, token)

	if err != nil {

		if errors.IsNotFound(err) {
			log.Printf("Resource with name '%s' in namespace '%s' not found", name, namespace)
			w.WriteHeader(404)
			return
		}

		log.Printf("Unexpected Error Finding Applier Resource: %v", err)
		w.WriteHeader(500)
		return
	}

	// Launch Applier Job
	err = applierManager.LaunchApplierJob(applier)

	if err != nil {
		log.Printf("Error Launching Job: %v", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(201)

}
