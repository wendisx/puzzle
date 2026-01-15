package config

type SwagConfig struct {
	URLs                 []string     `yaml:"urls"`
	DocExpansion         string       `yaml:"docExpansion"`
	DomID                string       `yaml:"domID"`
	InstanceName         string       `yaml:"instanceName"`
	DeepLinking          bool         `yaml:"deepLinking"`
	PersistAuthorization bool         `yaml:"persistAuthorization"`
	SyntaxHighlight      bool         `yaml:"syntaxHightlight"`
	OAuth                *OAuthConfig `yaml:"oauth"`
}

type OAuthConfig struct {
	ClientId string `yaml:"chienId"`
	Realm    string `yaml:"realm"`
	AppName  string `yaml:"appName"`
}
