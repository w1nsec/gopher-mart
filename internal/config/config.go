package config

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"gopher-mart/internal/domain"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DBurl               string
	SrvAddr             string
	RemoteServiceAddr   string
	Secret              string
	LogLevel            string
	CookieName          string
	CookieHoursLifeTime time.Duration
	WorkersCount        uint
	RetryTimer          time.Duration
	RetryAttempts       uint
}

func LoadConfig() *Config {
	var (
		envSrvAddr = os.Getenv("RUN_ADDRESS")
		envDBurl   = os.Getenv("DATABASE_URI")
		envRSAddr  = os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	)

	config := &Config{
		DBurl:             envDBurl,
		SrvAddr:           envSrvAddr,
		RemoteServiceAddr: envRSAddr,
	}

	var flagSrvAddr, flagDBurl, flagRSAddr string
	dbUsage := `URL address for DB postgres
	Ex. user:password@127.0.0.1:5432`
	flag.StringVar(&flagSrvAddr, "a", "127.0.0.1:8000", "address for http server")
	flag.StringVar(&flagDBurl, "d", "127.0.0.1:5432", dbUsage)
	flag.StringVar(&flagRSAddr, "r", "127.0.0.1:8001", "address for remote accrual system")
	flag.Parse()

	if config.SrvAddr == "" {
		config.SrvAddr = flagSrvAddr
	}

	if config.DBurl == "" {
		config.DBurl = flagDBurl
	}

	if config.RemoteServiceAddr == "" {
		config.RemoteServiceAddr = flagRSAddr
	}

	err := LoadEnvfileConfig(config)
	if err != nil {
		fmt.Println("------------------------------------------------------------")
		log.Error().Err(err).Msg("set default values")

		// set default values
		config.Secret = domain.Secret
		config.LogLevel = domain.LogLevel

		config.CookieName = domain.CookieName
		config.CookieHoursLifeTime = domain.CookieHoursLifeTime
		config.WorkersCount = domain.WorkersCount
		config.RetryAttempts = domain.RetryAttempts
		config.RetryTimer = domain.RetryTimer
	}

	return config
}

func LoadEnvfileConfig(config *Config) error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	secret, exists := os.LookupEnv("SECRET")
	if exists {
		config.Secret = secret
	}

	logLvl, exists := os.LookupEnv("LOG_LEVEL")
	if exists {
		config.LogLevel = logLvl
	} else {
		config.LogLevel = domain.LogLevel
	}

	// tables
	orders, exists := os.LookupEnv("TableOrders")
	if exists {
		domain.TableOrders = orders
	}
	users, exists := os.LookupEnv("TableUsers")
	if exists {
		domain.TableUsers = users
	}
	balance, exists := os.LookupEnv("TableBalance")
	if exists {
		domain.TableBalance = balance
	}
	withdraws, exists := os.LookupEnv("TableWithdraws")
	if exists {
		domain.TableWithdraws = withdraws
	}

	val, exists := os.LookupEnv("WorkersCount")
	if exists {
		WorkersCount, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			WorkersCount = uint64(domain.WorkersCount)
		}
		config.WorkersCount = uint(WorkersCount)
		domain.WorkersCount = uint(WorkersCount)
	} else {
		config.WorkersCount = domain.WorkersCount
	}

	val, exists = os.LookupEnv("RetryTimer")
	if exists {
		RetryTimer, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			return err
		}
		config.RetryTimer = time.Duration(RetryTimer) * time.Second
		domain.RetryTimer = config.RetryTimer
	} else {
		config.RetryTimer = domain.RetryTimer
	}

	val, exists = os.LookupEnv("RetryAttempts")
	if exists {
		RetryAttempts, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			return err
		}
		config.RetryAttempts = uint(RetryAttempts)
		domain.RetryAttempts = uint(RetryAttempts)
	} else {
		config.RetryAttempts = domain.RetryAttempts
	}

	CookieName, exists := os.LookupEnv("CookieName")
	if exists {
		config.CookieName = CookieName
		domain.CookieName = CookieName
	} else {
		config.CookieName = domain.CookieName
	}

	CookieHoursLifeTime, exists := os.LookupEnv("CookieHoursLifeTime")
	if exists {
		val, err := strconv.ParseUint(CookieHoursLifeTime, 10, 64)
		if err != nil {
			return err
		}

		config.CookieHoursLifeTime = time.Duration(val) * time.Hour
	} else {
		config.CookieHoursLifeTime = domain.CookieHoursLifeTime
	}
	return nil
}
