package config

type Config struct {
	Http              Http              `mapstructure:"http"`
	Database          Database          `mapstructure:"database"`
	AddressValidation AddressValidation `mapstructure:"address_validation"`
	GoogleSheets      GoogleSheets      `mapstructure:"google_sheets"`
	Mandrill          Mandrill          `mapstructure:"mandrill"`
}

type Http struct {
	Public Api `mapstructure:"public"`
	Admin  Api `mapstructure:"admin"`
}

type Api struct {
	Address string `mapstructure:"address"`
}

type Database struct {
	ConnectionString string `mapstructure:"connection_string"`
}

type AddressValidation struct {
	ApiKey string `mapstructure:"api_key"`
}

type GoogleSheets struct {
	ServiceAccountKeyPath string `mapstructure:"service_account_key_path"`
}

type Mandrill struct {
	ApiKey string `mapstructure:"api_key"`
}
