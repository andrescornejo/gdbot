package gdbot

import (
	"fmt"
	"os"
	"strconv"
)

// 	// GuildID = snowflake.GetEnv("GUILD_ID")
// 	GuildID = snowflake.ID(456983540853374986)

func LoadConfig() (*Config, error) {
	token := os.Getenv("TOKEN")
	if token == "" {
		return nil, fmt.Errorf("TOKEN environment variable is required")
	}

	nodeName := os.Getenv("NODE_NAME")
	if nodeName == "" {
		return nil, fmt.Errorf("NODE_NAME environment variable is required")
	}

	nodeAddress := os.Getenv("NODE_ADDRESS")
	if nodeAddress == "" {
		return nil, fmt.Errorf("NODE_ADDRESS environment variable is required")
	}

	nodePassword := os.Getenv("NODE_PASSWORD")
	if nodePassword == "" {
		return nil, fmt.Errorf("NODE_PASSWORD environment variable is required")
	}

	nodeSecureStr := os.Getenv("NODE_SECURE")
	nodeSecure, err := strconv.ParseBool(nodeSecureStr)
	if err != nil {
		return nil, fmt.Errorf("NODE_SECURE must be 'true' or 'false', got: %s", nodeSecureStr)
	}

	return &Config{
		Token:        token,
		NodeName:     nodeName,
		NodeAddress:  nodeAddress,
		NodePassword: nodePassword,
		NodeSecure:   nodeSecure,
	}, nil
}

type Config struct {
	Token        string
	NodeName     string
	NodeAddress  string
	NodePassword string
	NodeSecure   bool
}
