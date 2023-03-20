package config

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/BurntSushi/toml"
)

// NewConfig returns a new Config
func NewConfig(cfgFile string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(cfgFile, &config); err != nil {
		return nil, fmt.Errorf("error decoding config: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("error validating config: %w", err)
	}
	return &config, nil
}

type Config struct {
	Credentials Credentials `toml:"credentials"`
	Location    string      `toml:"location"`
}

func (c *Config) GetCredentials() (azcore.TokenCredential, error) {
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("error validating config: %w", err)
	}

	creds, err := c.Credentials.Auth()
	if err != nil {
		return nil, fmt.Errorf("failed to get authentication token: %w", err)
	}

	return creds, nil
}

func (c *Config) Validate() error {
	if _, err := c.Credentials.Auth(); err != nil {
		return fmt.Errorf("failed to validate credentials: %w", err)
	}

	return nil
}

type Credentials struct {
	Name        string `toml:"name"`
	Description string `toml:"description"`

	TenantID       string `toml:"tenant_id"`
	ClientID       string `toml:"client_id"`
	SubscriptionID string `toml:"subscription_id"`
	ClientSecret   string `toml:"client_secret"`
	// ClientOptions is the azure identity client options that will be used to authenticate
	// againsts an azure cloud. This is a heavy handed approach for now, defining the entire
	// ClientOptions here, but should allow users to use this provider with AzureStack or any
	// other azure cloud besides Azure proper (like Azure China, Germany, etc).
	ClientOptions azcore.ClientOptions `toml:"client_options"`
}

func (c Credentials) Validate() error {
	if c.TenantID == "" {
		return fmt.Errorf("missing tenant_id")
	}

	if c.ClientID == "" {
		return fmt.Errorf("missing client_id")
	}

	if c.SubscriptionID == "" {
		return fmt.Errorf("missing subscription_id")
	}

	if c.ClientSecret == "" {
		return fmt.Errorf("missing subscription_id")
	}

	return nil
}

func (c Credentials) Auth() (azcore.TokenCredential, error) {
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("validating credentials: %w", err)
	}

	o := &azidentity.ClientSecretCredentialOptions{ClientOptions: c.ClientOptions}
	cred, err := azidentity.NewClientSecretCredential(c.TenantID, c.ClientID, c.ClientSecret, o)
	if err != nil {
		return nil, err
	}
	return cred, nil
}
