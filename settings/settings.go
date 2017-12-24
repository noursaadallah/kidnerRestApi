package settings

import (
	"os"

	"github.com/spf13/viper"
)

type FabricSettings struct {
	ChannelId     string
	ChannelConfig string

	ChaincodeId      string
	ChaincodeVersion string
	ChaincodePath    string

	SDKConfig string

	AdminName      string
	AdminPwd       string
	StateStorePath string

	OrdKeyDir   string
	OrdCertDir  string
	OrdUsername string

	OrgKeyDir   string
	OrgCertDir  string
	OrgUsername string
}

type WebSettings struct {
	Address string
	Port    string
}

func findConfigFile(configPath string, configFileName string) error {
	path := configPath + "/" + configFileName + ".toml"
	if _, err := os.Stat(path); err != nil {
		gopath := os.Getenv("GOPATH")
		configPath = gopath + "/src/github.com/noursaadallah/kidner/settings"
	}

	viper.SetConfigName(configFileName)
	viper.AddConfigPath(configPath)
	return nil
}

func GetFabricSettings() (FabricSettings, error) {
	var fs FabricSettings
	_ = findConfigFile(".", "setup")
	err := viper.ReadInConfig()
	if err != nil {
		return fs, err
	}

	// Channel config
	fs.ChannelId = viper.GetString("Channel.Id")
	fs.ChannelConfig = viper.GetString("Channel.ConfigFile")

	// Chaincode config
	fs.ChaincodeId = viper.GetString("Chaincode.Id")
	fs.ChaincodeVersion = viper.GetString("Chaincode.Version")
	fs.ChaincodePath = viper.GetString("Chaincode.Path")

	// SDK config
	fs.SDKConfig = viper.GetString("SDKConfig.Path")

	// Admin config
	fs.AdminName = viper.GetString("Admin.Name")
	fs.AdminPwd = viper.GetString("Admin.Pwd")
	fs.StateStorePath = viper.GetString("Admin.StateStorePath")

	// OrdererUser config
	fs.OrdKeyDir = viper.GetString("OrdererUser.KeyDir")
	fs.OrdCertDir = viper.GetString("OrdererUser.CertDir")
	fs.OrdUsername = viper.GetString("OrdererUser.Username")

	// OrgUser config
	fs.OrgKeyDir = viper.GetString("OrgUser.KeyDir")
	fs.OrgCertDir = viper.GetString("OrgUser.CertDir")
	fs.OrgUsername = viper.GetString("OrgUser.Username")
	return fs, nil
}

func GetWebSettings() (WebSettings, error) {
	var ws WebSettings
	_ = findConfigFile(".", "setup")
	err := viper.ReadInConfig()
	if err != nil {
		return ws, err
	}

	// Web config
	ws.Address = viper.GetString("WebServer.Address")
	ws.Port = viper.GetString("WebServer.Port")

	return ws, nil
}
