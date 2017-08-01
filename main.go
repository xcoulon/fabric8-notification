package main

import (
	"net/http"

	"github.com/fabric8-services/fabric8-notification/app"
	"github.com/fabric8-services/fabric8-notification/collector"
	"github.com/fabric8-services/fabric8-notification/configuration"
	"github.com/fabric8-services/fabric8-notification/controller"
	"github.com/fabric8-services/fabric8-notification/email"
	"github.com/fabric8-services/fabric8-notification/jsonapi"
	"github.com/fabric8-services/fabric8-notification/keycloak"
	"github.com/fabric8-services/fabric8-notification/template"
	"github.com/fabric8-services/fabric8-notification/wit"

	"github.com/Sirupsen/logrus"
	witmiddleware "github.com/fabric8-services/fabric8-wit/goamiddleware"
	"github.com/fabric8-services/fabric8-wit/log"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/goadesign/goa/middleware/gzip"
	goajwt "github.com/goadesign/goa/middleware/security/jwt"

	goalogrus "github.com/goadesign/goa/logging/logrus"
)

func main() {

	// Initialized configuration
	config, err := configuration.GetData()
	if err != nil {
		logrus.Panic(nil, map[string]interface{}{
			"err": err,
		}, "failed to setup the configuration")
	}

	// Initialized developer mode flag for the logger
	log.InitializeLogger(config.IsLogJSON(), config.GetLogLevel())

	keycloakConfig := keycloak.Config{
		BaseURL: config.GetKeycloakURL(),
		Realm:   config.GetKeycloakRealm(),
	}

	publicKey, err := keycloak.GetPublicKey(keycloakConfig)
	if err != nil {
		log.Panic(nil, map[string]interface{}{
			"err": err,
		}, "failed to parse public token")
	}

	err = config.Validate()
	if err != nil {
		log.Panic(nil, map[string]interface{}{
			"err": err,
		}, "Missing required configuration")
	}

	witClient, err := wit.NewCachedClient(config.GetWITURL())
	if err != nil {
		log.Panic(nil, map[string]interface{}{
			"url": config.GetWITURL(),
			"err": err,
		}, "Could not create WIT client")
	}

	sender, err := email.NewMandrillSender(config.GetMadrillAPIKey())
	if err != nil {
		log.Panic(nil, map[string]interface{}{
			"err": err,
		}, "Could not create Madrill Sender")
	}

	notifier := email.NewAsyncWorkerNotifier(sender, 1)

	resolvers := &collector.LocalRegistry{}
	resolvers.Register("workitem.create", collector.ConfiguredVars(config, collector.NewWorkItemResolver(witClient)))
	resolvers.Register("workitem.update", collector.ConfiguredVars(config, collector.NewWorkItemResolver(witClient)))
	resolvers.Register("comment.create", collector.ConfiguredVars(config, collector.NewCommentResolver(witClient)))
	resolvers.Register("comment.update", collector.ConfiguredVars(config, collector.NewCommentResolver(witClient)))

	typeRegistry := &template.AssetRegistry{}
	service := goa.New("notification")

	// Mount middleware
	service.WithLogger(goalogrus.New(log.Logger()))
	service.Use(middleware.RequestID())
	service.Use(gzip.Middleware(9))
	service.Use(jsonapi.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	service.Use(witmiddleware.TokenContext(publicKey, nil, app.NewJWTSecurity()))
	service.Use(log.LogRequest(config.IsDeveloperModeEnabled()))
	app.UseJWTMiddleware(service, goajwt.New(publicKey, nil, app.NewJWTSecurity()))

	// Mount "status" controller
	statusCtrl := controller.NewStatusController(service)
	app.MountStatusController(service, statusCtrl)

	// Mount "notification" controller
	notifyCtrl := controller.NewNotifyController(service, resolvers, typeRegistry, notifier)
	app.MountNotifyController(service, notifyCtrl)

	log.Logger().Infoln("Git Commit SHA: ", controller.Commit)
	log.Logger().Infoln("UTC Build Time: ", controller.BuildTime)
	log.Logger().Infoln("UTC Start Time: ", controller.StartTime)
	log.Logger().Infoln("Dev mode:       ", config.IsDeveloperModeEnabled())

	http.Handle("/api/", service.Mux)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	// Start http
	if err := http.ListenAndServe(config.GetHTTPAddress(), nil); err != nil {
		log.Error(nil, map[string]interface{}{
			"addr": config.GetHTTPAddress(),
			"err":  err,
		}, "unable to connect to server")
		service.LogError("startup", "err", err)
	}
}
