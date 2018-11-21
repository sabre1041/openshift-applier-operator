package main

import (
	"flag"
	"log"
	"runtime"

	"net/http"

	sdkVersion "github.com/operator-framework/operator-sdk/version"
	"github.com/redhat-cop/openshift-applier-operator/pkg/apis"
	"github.com/redhat-cop/openshift-applier-operator/pkg/controller"
	"github.com/redhat-cop/openshift-applier-operator/pkg/handler"
	"github.com/redhat-cop/openshift-applier-operator/pkg/manager"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	k8sManager "sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

func printVersion() {
	log.Printf("Go Version: %s", runtime.Version())
	log.Printf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	log.Printf("operator-sdk Version: %v", sdkVersion.Version)
}

func main() {
	printVersion()
	flag.Parse()

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := k8sManager.New(cfg, k8sManager.Options{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Registering Components.")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Fatal(err)
	}

	// Setup all Controllers
	if err := controller.AddToManager(mgr); err != nil {
		log.Fatal(err)
	}

	applierMgr, err := manager.New(mgr)

	// Set up Applier Manager

	mux := http.NewServeMux()
	mux.HandleFunc("/webhook/", func(w http.ResponseWriter, r *http.Request) { handler.WebhookHandler(w, r, applierMgr) })

	log.Printf("Starting the Web Server")
	go http.ListenAndServe(":8080", mux)

	log.Printf("Starting the Cmd.")

	// Start the Cmd
	log.Fatal(mgr.Start(signals.SetupSignalHandler()))
}
