package token

import (
	"crypto/rsa"
	"fmt"

	authjwk "github.com/fabric8-services/fabric8-auth/token/jwk"
	authservice "github.com/fabric8-services/fabric8-notification/auth/api"
	"github.com/fabric8-services/fabric8-wit/log"
	"github.com/pkg/errors"
)

// tokenManagerConfiguration represents configuration needed to construct a token manager
type tokenManagerConfiguration interface {
	GetAuthURL() string
}

type PublicKey struct {
	KeyID string
	Key   *rsa.PublicKey
}

// Manager generate and find auth token information
type Manager interface {
	PublicKeys() []*rsa.PublicKey
}

type tokenManager struct {
	publicKeysMap map[string]*rsa.PublicKey
	publicKeys    []*PublicKey
}

// NewManager returns a new token Manager for handling tokens
func NewManager(config tokenManagerConfiguration) (Manager, error) {
	// Load public keys from Auth service and add them to the manager
	tm := &tokenManager{
		publicKeysMap: map[string]*rsa.PublicKey{},
	}

	keysEndpoint := fmt.Sprintf("%s%s", config.GetAuthURL(), authservice.KeysTokenPath())
	remoteKeys, err := authjwk.FetchKeys(keysEndpoint)
	if err != nil {
		log.Error(nil, map[string]interface{}{
			"err":      err,
			"keys_url": keysEndpoint,
		}, "unable to load public keys from remote service")
		return nil, errors.New("unable to load public keys from remote service")
	}
	for _, remoteKey := range remoteKeys {
		tm.publicKeysMap[remoteKey.KeyID] = remoteKey.Key
		tm.publicKeys = append(tm.publicKeys, &PublicKey{KeyID: remoteKey.KeyID, Key: remoteKey.Key})
		log.Info(nil, map[string]interface{}{
			"kid": remoteKey.KeyID,
		}, "Public key added")
	}
	return tm, nil
}

// PublicKeys returns all the public keys
func (mgm *tokenManager) PublicKeys() []*rsa.PublicKey {
	keys := make([]*rsa.PublicKey, 0, len(mgm.publicKeysMap))
	for _, key := range mgm.publicKeys {
		keys = append(keys, key.Key)
	}
	return keys
}
