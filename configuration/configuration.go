package configuration

import (
	"fmt"
	"strings"

	commoncfg "github.com/fabric8-services/fabric8-common/configuration"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

const (
	// Constants for viper variable names. Will be used to set
	// default values as well as to get each value

	varHTTPAddress          = "http.address"
	varMetricsHTTPAddress   = "metrics.http.address"
	varDeveloperModeEnabled = "developer.mode.enabled"
	varWITURL               = "wit.url"
	varAuthURL              = "auth.url"
	varMadrillAPIKey        = "mandrill.apikey"
	varLogLevel             = "log.level"
	varLogJSON              = "log.json"
	varServiceAccountID     = "service.account.id"
	varServiceAccountSecret = "service.account.secret"
)

// Data encapsulates the Viper configuration object which stores the configuration data in-memory.
type Data struct {
	v *viper.Viper
}

// NewData creates a configuration reader object using a configurable configuration file path
func NewData() (*Data, error) {
	c := Data{
		v: viper.New(),
	}
	c.v.SetEnvPrefix("F8")
	c.v.AutomaticEnv()
	c.v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	c.v.SetTypeByDefaultValue(true)
	c.setConfigDefaults()

	return &c, nil
}

// String returns the current configuration as a string
func (c *Data) String() string {
	allSettings := c.v.AllSettings()
	y, err := yaml.Marshal(&allSettings)
	if err != nil {
		log.WithFields(map[string]interface{}{
			"settings": allSettings,
			"err":      err,
		}).Panicln("Failed to marshall config to string")
	}
	return fmt.Sprintf("%s\n", y)
}

// GetData is a wrapper over NewData which reads configuration file path
// from the environment variable.
func GetData() (*Data, error) {
	cd, err := NewData()
	return cd, err
}

func (c *Data) setConfigDefaults() {
	c.v.SetTypeByDefaultValue(true)

	//-----
	// HTTP
	//-----
	c.v.SetDefault(varHTTPAddress, "0.0.0.0:8080")
	c.v.SetDefault(varMetricsHTTPAddress, "0.0.0.0:8080")

	c.v.SetDefault(varWITURL, defaultWITURL)
	c.v.SetDefault(varAuthURL, defaultAuthURL)

	//-----
	// Misc
	//-----

	// Enable development related features, e.g. token generation endpoint
	c.v.SetDefault(varDeveloperModeEnabled, false)
	c.v.SetDefault(varLogLevel, defaultLogLevel)

	c.v.SetDefault(varServiceAccountID, "4c34f6d4-f00b-487b-9a1f-e7d1adba6866")
	c.v.SetDefault(varServiceAccountSecret, "secret")

	// c.v.SetDefault(varMadrillAPIKey, "1234") // Enable for local testing.
}

// GetHTTPAddress returns the HTTP address (as set via default, config file, or environment variable)
// that the notification server binds to (e.g. "0.0.0.0:8080")
func (c *Data) GetHTTPAddress() string {
	return c.v.GetString(varHTTPAddress)
}

// GetMetricsHTTPAddress returns the address the /metrics endpoing will be mounted.
// By default GetMetricsHTTPAddress is the same as GetHTTPAddress
func (c *Data) GetMetricsHTTPAddress() string {
	return c.v.GetString(varMetricsHTTPAddress)
}

// IsDeveloperModeEnabled returns if development related features (as set via default, config file, or environment variable),
// e.g. token generation endpoint are enabled
func (c *Data) IsDeveloperModeEnabled() bool {
	return c.v.GetBool(varDeveloperModeEnabled)
}

// GetWITURL return the base WorkItemTracker API URL
func (c *Data) GetWITURL() string {
	return c.v.GetString(varWITURL)
}

// GetAuthServiceURL return the base Auth API URL
func (c *Data) GetAuthServiceURL() string {
	return c.v.GetString(varAuthURL)
}

// GetServiceAccountID returns service account ID for the notification service.
// This will be used by the notification service to request for a service account token
// from the Auth service.
func (c *Data) GetServiceAccountID() string {
	return c.v.GetString(varServiceAccountID)
}

// GetServiceAccountSecret returns service account secret for the notification service.
// This will be used by the notification service to request for a service account token
// from the Auth service.
func (c *Data) GetServiceAccountSecret() string {
	return c.v.GetString(varServiceAccountSecret)
}

// GetWebURL returns the base URL for the Web v
func (c *Data) GetWebURL() string {
	return strings.Replace(c.GetWITURL(), "api.", "", -1)
}

// GetMadrillAPIKey returns the API key used by the email sender
func (c *Data) GetMadrillAPIKey() string {
	return c.v.GetString(varMadrillAPIKey)
}

// GetLogLevel returns the loggging level (as set via config file or environment variable)
func (c *Data) GetLogLevel() string {
	return c.v.GetString(varLogLevel)
}

// IsLogJSON returns if we should log json format (as set via config file or environment variable)
func (c *Data) IsLogJSON() bool {
	if c.v.IsSet(varLogJSON) {
		return c.v.GetBool(varLogJSON)
	}
	if c.IsDeveloperModeEnabled() {
		return false
	}
	return true
}

func (c *Data) Validate() error {
	if c.GetMadrillAPIKey() == "" {
		return fmt.Errorf("Missing %v", varMadrillAPIKey)
	}
	return nil
}

// GetDevModePrivateKey returns the private key and its ID used in tests
func (c *Data) GetDevModePrivateKey() []byte {
	if c.IsDeveloperModeEnabled() {
		return []byte(commoncfg.DevModeRsaPrivateKey)
	}
	return nil
}

const (
	defaultWITURL   = "https://api.openshift.io/"
	defaultAuthURL  = "http://localhost:8089/"
	defaultLogLevel = "info"
)
