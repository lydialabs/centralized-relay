package bitcoin

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	MongoURI    string
	MongoDB     string
	BitcoinRPC  string
	AdminWallet string
	StartBlock  int
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")

	// Load the .env file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Error reading config file, %s", err)
	}

	startBlock, err := strconv.Atoi(viper.GetString("BITCOIN_START_BLOCK"))
	if err != nil {
		return nil, fmt.Errorf("Invalid BITCOIN_START_BLOCK value")
	}

	config := &Config{
		MongoURI:    viper.GetString("DATABASE_MONGO_URI"),
		MongoDB:     viper.GetString("DATABASE_MONGO_DATABASE"),
		BitcoinRPC:  viper.GetString("BITCOIN_RPC_1"),
		AdminWallet: viper.GetString("BITCOIN_ADMIN_WALLET"),
		StartBlock:  startBlock,
	}

	return config, nil
}
