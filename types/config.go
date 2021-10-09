package types

type Config struct {
	Addr         string `default:":3000"`
	SessionName  string
	GithubAppId  string
	ClientID     string
	AllowSignup  string
	ClientSecret string
	TokenExp     int `default:"600"`
}
