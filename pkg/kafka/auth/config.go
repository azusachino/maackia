package auth

import (
	"fmt"
	"strings"

	"github.com/AzusaChino/maackia/pkg/config/tlscfg"
	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	none      = "none"
	kerberos  = "kerberos"
	tls       = "tls"
	plaintext = "plaintext"
)

var authTypes = []string{
	none,
	kerberos,
	tls,
	plaintext,
}

// AuthenticationConfig describes the configuration properties needed authenticate with kafka cluster
type AuthenticationConfig struct {
	Authentication string          `mapstructure:"type"`
	Kerberos       KerberosConfig  `mapstructure:"kerberos"`
	TLS            tlscfg.Options  `mapstructure:"tls"`
	PlainText      PlainTextConfig `mapstructure:"plaintext"`
}

// configure authentication into sarama config
func (config *AuthenticationConfig) SetConfiguration(saramaConfig *sarama.Config, logger *zap.Logger) error {
	authentication := strings.ToLower(config.Authentication)
	if strings.Trim(authentication, " ") == "" {
		authentication = none
	}
	if config.Authentication == tls || config.TLS.Enabled {
		err := setTLSConfiguration(&config.TLS, saramaConfig, logger)
		if err != nil {
			return err
		}
	}
	switch authentication {
	case none:
		return nil
	case tls:
		return nil
	case kerberos:
		setKerberosConfiguration(&config.Kerberos, saramaConfig)
		return nil
	case plaintext:
		err := setPlainTextConfiguration(&config.PlainText, saramaConfig)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("Unknown/Unsupported authentication method %s to kafka cluster", config.Authentication)
	}
}

// InitFromViper loads authentication configuration from viper flags.
func (config *AuthenticationConfig) InitFromViper(configPrefix string, v *viper.Viper) {
	config.Authentication = v.GetString(configPrefix + suffixAuthentication)
	config.Kerberos.ServiceName = v.GetString(configPrefix + kerberosPrefix + suffixKerberosServiceName)
	config.Kerberos.Realm = v.GetString(configPrefix + kerberosPrefix + suffixKerberosRealm)
	config.Kerberos.UseKeyTab = v.GetBool(configPrefix + kerberosPrefix + suffixKerberosUseKeyTab)
	config.Kerberos.Username = v.GetString(configPrefix + kerberosPrefix + suffixKerberosUsername)
	config.Kerberos.Password = v.GetString(configPrefix + kerberosPrefix + suffixKerberosPassword)
	config.Kerberos.ConfigPath = v.GetString(configPrefix + kerberosPrefix + suffixKerberosConfig)
	config.Kerberos.KeyTabPath = v.GetString(configPrefix + kerberosPrefix + suffixKerberosKeyTab)

	var tlsClientConfig = tlscfg.ClientFlagsConfig{
		Prefix: configPrefix,
	}

	config.TLS = tlsClientConfig.InitFromViper(v)
	if config.Authentication == tls {
		config.TLS.Enabled = true
	}

	config.PlainText.Username = v.GetString(configPrefix + plainTextPrefix + suffixPlainTextUsername)
	config.PlainText.Password = v.GetString(configPrefix + plainTextPrefix + suffixPlainTextPassword)
	config.PlainText.Mechanism = v.GetString(configPrefix + plainTextPrefix + suffixPlainTextMechanism)
}
