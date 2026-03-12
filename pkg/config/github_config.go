package config

type (
	Repo struct {
		Name string `yaml:"name" json:"name"`
		Ref  string `yaml:"ref" json:"ref"`
	}
	GithubConfig struct {
		ApiHost     string `yaml:"apiHost" json:"apiHost"`
		RawHost     string `yaml:"rawHost" json:"rawHost"`
		AccessToken string `yaml:"accessToken" json:"accessToken"`
		UserName    string `yaml:"userName" json:"userName"`
		UserEmail   string `yaml:"userEmail" json:"userEmail"`
		Repos       []Repo `yaml:"repos" json:"repos"`
		ActiveRepo  int    `yaml:"activeRepo" json:"activeRepo"`
	}
)

func initGithubConfig() GithubConfig {
	return GithubConfig{
		ApiHost:     "https://api.github.com/",
		AccessToken: "",
		UserName:    "git",
		UserEmail:   "git@github.com",
	}
}
