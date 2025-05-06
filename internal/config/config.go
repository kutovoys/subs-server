package config

import (
	"fmt"

	"github.com/alecthomas/kong"
)

var CLIConfig Config

func Parse(version string) {
	ctx := kong.Parse(&CLIConfig,
		kong.Name("subs-server"),
		kong.Description("Subs Server: A server for managing subscriptions"),
		kong.Vars{
			"version": version,
		},
	)
	_ = ctx
}

type Config struct {
	Source                string      `name:"source" short:"s" default:"filesystem" enum:"filesystem" help:"Source type for files (filesystem)" env:"SOURCE"`
	Location              string      `name:"location" short:"l" required:"true" help:"Path to files (directory for filesystem source)" env:"LOCATION"`
	Port                  int         `name:"port" short:"p" default:"2115" help:"Port to listen on" env:"PORT"`
	Host                  string      `name:"host" short:"h" default:"0.0.0.0" help:"Host to listen on" env:"HOST"`
	Debug                 bool        `name:"debug" short:"d" help:"Enable debug mode (verbose logging and endpoint listing)" env:"DEBUG"`
	ProfileTitle          string      `name:"profile-title" help:"Profile title (will be base64 encoded)" default:"Subs-Server" env:"PROFILE_TITLE"`
	ProfileUpdateInterval string      `name:"profile-update-interval" help:"Profile update interval in hours" default:"12" env:"PROFILE_UPDATE_INTERVAL"`
	ProfileWebPageURL     string      `name:"profile-web-page-url" help:"Profile web page URL" default:"https://github.com/kutovoys/subs-server" env:"PROFILE_WEB_PAGE_URL"`
	SupportURL            string      `name:"support-url" help:"Support URL" default:"https://github.com/kutovoys/subs-server" env:"SUPPORT_URL"`
	Version               VersionFlag `name:"version" help:"Print version information and quit"`
}

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println("Subs Server: A server for managing subscriptions")
	fmt.Printf("Version:\t %s\n", vars["version"])
	fmt.Printf("GitHub: https://github.com/kutovoys/subs-server\n")
	app.Exit(0)
	return nil
}
