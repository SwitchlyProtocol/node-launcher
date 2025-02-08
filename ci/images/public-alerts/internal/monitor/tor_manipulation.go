package monitor

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"public-alerts/internal/common"
	"public-alerts/internal/config"
	"public-alerts/internal/notify"

	"github.com/rs/zerolog/log"
	openapi "gitlab.com/thorchain/thornode/openapi/gen"
)

type FetchPoolsFunc func() ([]openapi.Pool, error)
type FetchRunePriceFunc func() (float64, error)

type TorManipulationMonitor struct {
	client    common.ThornodeDataFetcher
	lastAlert time.Time
}

func (tmm *TorManipulationMonitor) Name() string {
	return "SolvencyMonitor"
}

func NewTorManipulationMonitor() *TorManipulationMonitor {
	client, _ := common.NewThornodeClient()
	return &TorManipulationMonitor{
		client:    client,
		lastAlert: time.Time{}, // allow alerting immediately
	}
}

func fetchRunePrice() (float64, error) {
	url := "https://api.coingecko.com/api/v3/simple/price?ids=thorchain&vs_currencies=usd"
	response, err := http.Get(url)
	if err != nil {
		return 0.0, fmt.Errorf("get_rune_price error: %v", err)
	}
	defer response.Body.Close()
	var data map[string]map[string]float64
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return 0.0, fmt.Errorf("error decoding JSON: %v", err)
	}
	return data["thorchain"]["usd"], nil
}

func (tm *TorManipulationMonitor) checkTorManipulation(fetchPools FetchPoolsFunc, fetchRunePrice FetchRunePriceFunc) ([]string, error) {

	pools, err := fetchPools()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pools: %v", err)
	}

	runePrice, err := fetchRunePrice()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rune price: %v", err)
	}

	var alerts []string
	for _, pool := range pools {
		if pool.Asset == "THOR.TOR" {
			balanceRune, err := strconv.ParseFloat(pool.BalanceRune, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse BalanceRune: %v", err)
			}
			balanceAsset, err := strconv.ParseFloat(pool.BalanceAsset, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse BalanceAsset: %v", err)
			}
			torPriceInRune := balanceRune / balanceAsset
			torPriceInUsd := torPriceInRune * runePrice
			runePriceInTor := 1 / torPriceInRune

			torDelta := 100 * ((torPriceInUsd - 1) / 1)
			runeDelta := 100 * ((runePriceInTor - runePrice) / runePrice)
			// debug
			log.Info().Msgf("RUNE external Price: %.3f,  TOR Price in USD: $%.3f, tor delta: %.3f, rune delta: %.3f ", runePrice, torPriceInUsd, torDelta, runeDelta)
			if math.Abs(torDelta) > float64(config.Get().TorManipulationMonitor.TorPriceDeltaThreshold) {
				alert := fmt.Sprintf(
					"**🚨 ALERT: TOR PRICE 🚨**\n"+
						"> TOR Price in USD: $%.3f\n"+
						"> TOR Price Delta: %.3f%%\n"+
						"> ---------------------------\n"+
						"> Rune Price in TOR: %.3f\n"+
						"> Rune Price in USD: $%.3f\n"+
						"> Rune Price Delta: %.3f%%\n"+
						"> Check TOR Anchor pool depths and median value.",
					torPriceInUsd, torDelta, runePriceInTor, runePrice, runeDelta,
				)
				alerts = append(alerts, alert)
			}
		}
	}

	return alerts, nil
}

func (tm *TorManipulationMonitor) Check() ([]notify.Alert, error) {
	log.Info().Msg("Checking Tor Manipulation...")
	if time.Since(tm.lastAlert) < time.Duration(config.Get().TorManipulationMonitor.AlertCooldownSeconds)*time.Second {
		return nil, nil
	}

	alerts, err := tm.checkTorManipulation(tm.client.GetDerivedPools, fetchRunePrice)
	if err != nil {
		return nil, err
	}

	var notifyAlerts []notify.Alert
	for _, msg := range alerts {
		notifyAlerts = append(notifyAlerts, notify.Alert{Webhooks: config.Get().Webhooks.Security, Message: msg})
	}

	if len(notifyAlerts) > 0 {
		tm.lastAlert = time.Now()
	}

	return notifyAlerts, nil
}
