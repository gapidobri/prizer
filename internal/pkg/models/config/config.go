package config

type Config struct {
	Database          Database          `mapstructure:"database"`
	AddressValidation AddressValidation `mapstructure:"address_validation"`
	Mailchimp         Mailchimp         `mapstructure:"mailchimp"`
	Mandrill          Mandrill          `mapstructure:"mandrill"`
}

type Database struct {
	ConnectionString string `mapstructure:"connection_string"`
}

type AddressValidation struct {
	ApiKey string `mapstructure:"api_key"`
}

type Mailchimp struct {
	ApiKey string `mapstructure:"api_key"`
}

type Mandrill struct {
	User   string `mapstructure:"user"`
	ApiKey string `mapstructure:"api_key"`
}
