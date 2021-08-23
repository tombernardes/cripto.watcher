package configuration

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/adshao/go-binance/v2"
	log "github.com/sirupsen/logrus"
)

var Config Configuration

type Configuration struct {
	BinanceApiKey    string   `json:"binance-api-key,omitempty"`
	BinanceSecretKey string   `json:"binance-secret-key,omitempty"`
	LogLevel         string   `json:"log-level,omitempty"`
	ChartPeriod      string   `json:"chart-period,omitempty"`
	Tickers          []string `json:"tickers,omitempty"`
	Client           *binance.Client
	OpenOrders       int64 `json:"open-orders,omitempty"`
}

func init() {
	configureLog()
	updateConfiguration()
	Config.Client = binance.NewClient(Config.BinanceApiKey, Config.BinanceSecretKey)
	go asyncUpdate()
}

func asyncUpdate() {
	for {
		updateConfiguration()
		time.Sleep(time.Duration(10) * time.Second)
	}
}

func configureLog() {
	log.SetFormatter(&log.JSONFormatter{})
	Config.LogLevel = "info"
}

func getConfig() []byte {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	file, err := ioutil.ReadFile(basepath + "/configuration.json")
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Fatal("Error getting configuration file")
	}
	return []byte(file)
}

func updateConfiguration() {
	raw := getConfig()
	_ = json.Unmarshal(raw, &Config)

	level, err := log.ParseLevel(Config.LogLevel)
	if err != nil {
		log.WithFields(log.Fields{
			"prefix": "configuration.updateConfiguration",
		}).Error("Error on setting log level. \n" + err.Error())
	} else {
		if level != log.GetLevel() {
			log.WithFields(log.Fields{
				"prefix": "configuration.updateConfiguration",
			}).Info("Setting log level to " + strings.ToUpper(Config.LogLevel))
		}
		log.SetLevel(level)
	}
	time.Sleep(time.Duration(3) * time.Second)

}
