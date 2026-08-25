package configuration

// KeychainConfigLoader reads the credentials stored in the keychain.
// Load returns an empty Config along with the error when the keychain cannot be read.
type KeychainConfigLoader struct {
	Keychain Keychain
}

func (e KeychainConfigLoader) Load() (Config, error) {
	user, err := e.Keychain.Get(KeychainUserKey)
	if err != nil {
		return Config{}, err
	}

	hostname, err := e.Keychain.Get(KeychainHostNameKey)
	if err != nil {
		return Config{}, err
	}

	apiKey, err := e.Keychain.Get(KeychainApiKeyKey)
	if err != nil {
		return Config{}, err
	}

	return Config{
		User:     user,
		HostName: hostname,
		ApiKey:   apiKey,
	}, nil
}
