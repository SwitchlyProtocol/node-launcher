package config

import (
	"fmt"

	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

////////////////////////////////////////////////////////////////////////////////
// Monitor Configuration
////////////////////////////////////////////////////////////////////////////////

type MonitorConfig interface {
	Validate() error // Validates the monitor-specific configuration
}

// ///////////////////////
// ChainLagMonitorConfig
// ///////////////////////
type ChainLagMonitorConfig struct {
	MaxChainLag map[string]int
}

func (c ChainLagMonitorConfig) Validate() error {
	// Iterate over the map and check if any values are 0
	for k, v := range c.MaxChainLag {
		if v == 0 {
			// Use fmt.Errorf to format the error message with the chain key
			return fmt.Errorf("ChainLag Monitor testValue cannot be 0 for chain %s", k)
		}
	}

	return nil
}

// NewChainLagMonitorConfig creates a new ChainLagMonitorConfig with default settings.
func NewChainLagMonitorConfig() ChainLagMonitorConfig {
	return ChainLagMonitorConfig{
		MaxChainLag: map[string]int{
			"BCH":  3,
			"BTC":  3,
			"BNB":  1800,
			"DOGE": 30,
			"ETH":  70,
			"LTC":  6,
			"GAIA": 175,
			"AVAX": 900,
		},
	}
}

/////////////////////////
// SolvencyMonitorConfig
/////////////////////////

type SolvencyMonitorConfig struct {
	AlertWindowThreshold  int
	AlertPercentThreshold float64
	AlertUSDThreshold     float64
	AlertCooldownSeconds  int
}

func (s SolvencyMonitorConfig) Validate() error {
	// TODO: Implement validation

	return nil
}

// NewSolvencyMonitorConfig creates a new SolvencyMonitorConfig with default settings.
func NewSolvencyMonitorConfig() SolvencyMonitorConfig {
	return SolvencyMonitorConfig{
		AlertWindowThreshold:  60,
		AlertPercentThreshold: 0.02,
		AlertUSDThreshold:     5000,
		AlertCooldownSeconds:  60 * 60 * 12, // 12 hours
	}
}

// ///////////////////////
// StuckOutboundMonitorConfig
// ///////////////////////
type StuckOutboundMonitorConfig struct {
	BlockAgeThreshold int
}

func (sobm StuckOutboundMonitorConfig) Validate() error {
	// TODO(Orion): add validation
	return nil
}

func NewStuckOutboundMonitorConfig() StuckOutboundMonitorConfig {
	return StuckOutboundMonitorConfig{
		BlockAgeThreshold: 7200, // ~12 hours
	}
}

// ///////////////////////
// ChainUpdateMonitorConfig
// ///////////////////////
type DaemonConfig struct {
	Name      string
	Github    string
	LatestTag string
}

type ChainUpdateMonitorConfig struct {
	Daemons map[string]DaemonConfig
	DataDir string
}

func (c ChainUpdateMonitorConfig) Validate() error {
	// TODO: Implement validation

	return nil
}

func NewChainUpdateMonitorConfig() ChainUpdateMonitorConfig {

	// Define the daemons map
	daemons := map[string]DaemonConfig{
		"binance-smart": {"binance-smart", "bnb-chain/bsc", ""},
		"bitcoin":       {"bitcoin", "bitcoin/bitcoin", ""},
		"bitcoin-cash":  {"bitcoin-cash", "bitcoin-cash-node/bitcoin-cash-node", ""},
		"dogecoin":      {"dogecoin", "dogecoin/dogecoin", ""},
		"ethereum":      {"ethereum", "ethereum/go-ethereum", ""},
		"gaia":          {"gaia", "cosmos/gaia", ""},
		"litecoin":      {"litecoin", "litecoin-project/litecoin", ""},
		"avalanche":     {"avalanche", "ava-labs/avalanchego", ""},
		"prysm":         {"prysm", "prysmaticlabs/prysm", ""},
	}

	return ChainUpdateMonitorConfig{DataDir: viper.GetString("dataDir"), Daemons: daemons}
}

/////////////////////////
// SecurityUpdatesConfig
/////////////////////////

type SecurityUpdatesMonitorConfig struct {
	Repos []string
}

func (sup SecurityUpdatesMonitorConfig) Validate() error {
	// TODO: Implement validation

	return nil
}

func NewSecurityUpdatesMonitorConfig() SecurityUpdatesMonitorConfig {

	return SecurityUpdatesMonitorConfig{Repos: []string{"bnb-chain/tss-lib"}}
}

/////////////////////////
// TorManipulationMonitorConfig
/////////////////////////

type TorManipulationMonitorConfig struct {
	AlertCooldownSeconds   int // period of time before alert will fire agian, prevent spam
	TorPriceDeltaThreshold int // Percentage,value of 5 = 5%
}

func (tmm TorManipulationMonitorConfig) Validate() error {
	// TODO: Implement validation

	return nil
}

func NewTorManipulationMonitorConfig() TorManipulationMonitorConfig {

	return TorManipulationMonitorConfig{
		AlertCooldownSeconds:   60 * 60 * 2, // 2 hours
		TorPriceDeltaThreshold: 5,           // 5%
	}
}

/////////////////////////
// UtxoMempoolMonitorConfig
/////////////////////////

type UtxoMempoolMonitorConfig struct {
	ChainDaemonURLs    map[string]string `yaml:"chain_daemon_urls"`
	Chains             []string          `yaml:"chains"`
	AlertSizeThreshold map[string]int    `yaml:"alert_size_threshold"`
	AlertFactor        float64           `yaml:"alert_factor"`
	AlertWindow        int               `yaml:"alert_window"`
	AlertObservations  int               `yaml:"alert_observations"`
}

func (c UtxoMempoolMonitorConfig) Validate() error {
	if len(c.Chains) == 0 {
		return fmt.Errorf("at least one chain must be specified")
	}

	if c.AlertFactor <= 0 {
		return fmt.Errorf("alert_factor must be greater than 0")
	}

	if c.AlertWindow <= 0 {
		return fmt.Errorf("alert_window must be greater than 0")
	}

	if c.AlertObservations <= 0 {
		return fmt.Errorf("alert_observations must be greater than 0")
	}

	for _, chain := range c.Chains {
		if _, ok := c.AlertSizeThreshold[chain]; !ok {
			return fmt.Errorf("alert_size_threshold for chain %s is not specified", chain)
		}
	}

	return nil
}

func NewUtxoMempoolMonitorConfig() UtxoMempoolMonitorConfig {
	return UtxoMempoolMonitorConfig{
		ChainDaemonURLs: map[string]string{
			"bitcoin":      viper.GetString("BITCOIN_DAEMON_URL"),
			"bitcoin-cash": viper.GetString("BITCOIN_CASH_DAEMON_URL"),
			"dogecoin":     viper.GetString("DOGECOIN_DAEMON_URL"),
			"litecoin":     viper.GetString("LITECOIN_DAEMON_URL"),
		},
		Chains: []string{
			"bitcoin",
			"bitcoin-cash",
			"dogecoin",
			"litecoin",
		},
		AlertSizeThreshold: map[string]int{
			"bitcoin":      5000000,
			"bitcoin-cash": 32000000,
			"dogecoin":     3000000,
			"litecoin":     10000000,
		},
		AlertFactor:       5,
		AlertWindow:       60,
		AlertObservations: 5,
	}
}

////////////////////////////////////////////////////////////////////////////////
// Configuration
////////////////////////////////////////////////////////////////////////////////

type Webhooks struct {
	Slack     string `mapstructure:"slack"`
	Discord   string `mapstructure:"discord"`
	PagerDuty string `mapstructure:"pagerduty"`
}

type Config struct {
	Endpoints struct {
		ThornodeAPI   string `mapstructure:"thornode_api"`
		ThornodeRPC   string `mapstructure:"thornode_rpc"`
		NineRealmsAPI string `mapstructure:"ninerealms_api"`
		MidgardAPI    string `mapstructure:"midgard_api"`
		ExplorerURL   string `mapstructure:"explorer_url"`
	} `mapstructure:"endpoints"`
	Webhooks struct {
		Activity Webhooks `mapstructure:"activity"`
		Info     Webhooks `mapstructure:"info"`
		Updates  Webhooks `mapstructure:"updates"`
		Security Webhooks `mapstructure:"security"`
		Errors   Webhooks `mapstructure:"errors"`
	} `mapstructure:"webhooks"`
	// each monitor can have its own configuration params
	ChainLagMonitor        ChainLagMonitorConfig
	SolvencyMonitor        SolvencyMonitorConfig
	StuckOutboundMonitor   StuckOutboundMonitorConfig
	ChainUpdateMonitor     ChainUpdateMonitorConfig
	SecurityUpdatesMonitor SecurityUpdatesMonitorConfig
	TorManipulationMonitor TorManipulationMonitorConfig
	UtxoMempoolMonitor     UtxoMempoolMonitorConfig
}

// //////////////////////////////////////////////////////////////////////////////
// Init
// //////////////////////////////////////////////////////////////////////////////
var config Config

func init() {

	assert := func(err error) {
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to bind environment variable")
		}
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	// data dir
	viper.SetDefault("dataDir", "/data")
	assert(viper.BindEnv("dataDir", "DATA_DIR"))

	// Initialize ChainLagMonitor with hardcoded values
	config.ChainLagMonitor = NewChainLagMonitorConfig()
	config.SolvencyMonitor = NewSolvencyMonitorConfig()
	config.StuckOutboundMonitor = NewStuckOutboundMonitorConfig()
	config.ChainUpdateMonitor = NewChainUpdateMonitorConfig()
	config.SecurityUpdatesMonitor = NewSecurityUpdatesMonitorConfig()
	config.TorManipulationMonitor = NewTorManipulationMonitorConfig()
	config.UtxoMempoolMonitor = NewUtxoMempoolMonitorConfig()
	// TODO - validate each monitor configuration

	// endpoints
	assert(viper.BindEnv("endpoints.thornode_api", "ENDPOINTS_THORNODE_API"))
	assert(viper.BindEnv("endpoints.thornode_rpc", "ENDPOINTS_THORNODE_RPC"))
	assert(viper.BindEnv("endpoints.ninerealms_api", "ENDPOINTS_NINEREALMS_API"))
	assert(viper.BindEnv("endpoints.midgard_api", "ENDPOINTS_MIDGARD_API"))
	assert(viper.BindEnv("endpoints.explorer_url", "ENDPOINTS_EXPLORER_URL"))
	// webhooks - activity
	assert(viper.BindEnv("webhooks.activity.slack", "WEBHOOKS_ACTIVITY_SLACK"))
	assert(viper.BindEnv("webhooks.activity.discord", "WEBHOOKS_ACTIVITY_DISCORD"))
	// webhooks - security
	assert(viper.BindEnv("webhooks.security.slack", "WEBHOOKS_ACTIVITY_SLACK"))
	assert(viper.BindEnv("webhooks.security.discord", "WEBHOOKS_ACTIVITY_DISCORD"))
	// webhooks - errors
	assert(viper.BindEnv("webhooks.errors.slack", "WEBHOOKS_ERRORS_SLACK"))

	// Unmarshal the configuration into the config struct
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal().Err(err).Msg("Unable to unmarshal config")
	}
}

func Get() Config {
	return config
}
