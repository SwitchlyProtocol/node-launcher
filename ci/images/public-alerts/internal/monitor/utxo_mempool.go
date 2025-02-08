package monitor

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"public-alerts/internal/config"
	"public-alerts/internal/notify"
	"time"

	"github.com/rs/zerolog/log"
)

type MempoolInfo struct {
	Result struct {
		Bytes int `json:"bytes"`
	} `json:"result"`
}

type MempoolMonitor struct {
	Observations map[string][]int
	LastAlert    map[string]time.Time
}

func (mm *MempoolMonitor) Name() string {
	return "MempoolMonitor"
}

func NewMempoolMonitor() *MempoolMonitor {

	return &MempoolMonitor{
		Observations: make(map[string][]int),
		LastAlert:    make(map[string]time.Time), // allow alerting immediately
	}
}

func SafeMean(numbers []float64) (float64, error) {
	if len(numbers) == 0 {
		return 0, errors.New("cannot calculate the mean of an empty slice")
	}

	sum := 0.0
	for _, number := range numbers {
		sum += number
	}

	mean := sum / float64(len(numbers))
	return mean, nil
}

func (mm *MempoolMonitor) Check() ([]notify.Alert, error) {
	log.Info().Msg("Checking mempool...")
	cfg := config.Get()
	return mm.CheckMempool(mm.fetchMempoolSize, &cfg)
}

func (mm *MempoolMonitor) CheckMempool(fetchMempoolSize func(daemonUrl string) (int, error), cfg *config.Config) ([]notify.Alert, error) {

	var alerts []notify.Alert

	for _, chain := range cfg.UtxoMempoolMonitor.Chains {
		log.Info().Msgf("Checking mempool for chain: %s", chain)
		daemonUrl := cfg.UtxoMempoolMonitor.ChainDaemonURLs[chain]

		mempoolSize, err := fetchMempoolSize(daemonUrl)
		log.Info().Msgf("Mempool size for chain: %s is %d", chain, mempoolSize)
		if err != nil {
			return nil, err
		}

		mm.Observations[chain] = append(mm.Observations[chain], mempoolSize)
		if len(mm.Observations[chain]) > cfg.UtxoMempoolMonitor.AlertWindow {
			mm.Observations[chain] = mm.Observations[chain][1:]
		}

		// get mean of last n Observations
		var floatSlice []float64
		for _, val := range mm.Observations[chain] {
			floatSlice = append(floatSlice, float64(val))
		}
		meanObs, err := SafeMean(floatSlice)
		if err != nil {
			log.Error().Msgf("Error calculating mean: %v for chain: %s", err, chain)
			continue
		}

		if len(mm.Observations[chain]) >= cfg.UtxoMempoolMonitor.AlertObservations &&
			float64(mempoolSize) > float64(cfg.UtxoMempoolMonitor.AlertSizeThreshold[chain]) &&
			float64(mempoolSize) > float64(cfg.UtxoMempoolMonitor.AlertFactor)*meanObs &&
			time.Since(mm.LastAlert[chain]) > time.Hour {

			msg := fmt.Sprintf("%s mempool size %.2fx over %d minute mean (%.2fMb): %.2fMb",
				chain,
				float64(mempoolSize)/meanObs,
				cfg.UtxoMempoolMonitor.AlertWindow,
				meanObs/1e6,
				float64(mempoolSize)/1e6,
			)
			log.Info().Msg(msg)
			alerts = append(alerts, notify.Alert{Webhooks: cfg.Webhooks.Info, Message: msg})
			mm.LastAlert[chain] = time.Now()
		}
	}

	return alerts, nil
}

func (mm *MempoolMonitor) fetchMempoolSize(daemonUrl string) (int, error) {

	// Prepare the request body
	reqBody := map[string]string{"method": "getmempoolinfo"}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return 0, fmt.Errorf("error marshalling JSON: %w", err)
	}

	// Create new request with a body
	req, err := http.NewRequest("POST", daemonUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Logging response for debugging
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %w", err)
	}

	// Decode the JSON response into the struct
	var mempoolInfo MempoolInfo
	if err := json.Unmarshal(bodyBytes, &mempoolInfo); err != nil {
		return 0, fmt.Errorf("error decoding JSON: %w", err)
	}

	return mempoolInfo.Result.Bytes, nil
}
