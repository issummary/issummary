package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/issummary/issummary/api"
	"github.com/issummary/issummary/issummary"
	"github.com/joho/godotenv"
	"github.com/mpppk/gitany"
	"github.com/mpppk/gitany/etc"
	_ "github.com/mpppk/gitany/gitlab"

	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var classPrefixKey = "class-prefix"
var spPrefixKey = "sp-prefix"
var gitServiceTypeKey = "git-service"
var baseURLKey = "base-url"
var tokenKey = "token"
var portKey = "port"
var gidKey = "gid"

var RootCmd = &cobra.Command{
	Use:   "issummary",
	Short: "issue summary viewer",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		config, err := generateIssummaryConfigFromViper()
		if err != nil {
			panic(err)
		}

		fmt.Printf("%#v\n", config)

		protocolAndHost := strings.Split(config.GitServiceBaseURL, "://")
		protocol := protocolAndHost[0]
		host := protocolAndHost[1]

		serviceConfig := &etc.ServiceConfig{ // FIXME
			Host:     host,
			Type:     config.GitServiceType,
			Token:    config.Token,
			Protocol: protocol,
		}

		gitanyClient, err := gitany.GetClient(context.Background(), serviceConfig) // FIXME
		if err != nil {
			panic(err)
		}

		client := issummary.New(gitanyClient)

		statikFS, err := fs.New()
		if err != nil {
			log.Fatal(err)
		}

		http.Handle("/", http.FileServer(statikFS))

		http.HandleFunc("/api/works", api.GetWorksJsonHandleFunc(ctx, client, config))
		http.HandleFunc("/api/milestones", api.GetMilestonesJsonHandleFunc(ctx, client, config))
		err = http.ListenAndServe(fmt.Sprintf(":%v", config.Port), nil)
		if err != nil {
			panic(err)
		}
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	err := godotenv.Load()
	if err == nil {
		log.Println(".env file found")
	}

	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.issummary.yaml)")

	RootCmd.PersistentFlags().String(tokenKey, "", "git repository service token")
	viper.BindPFlag(tokenKey, RootCmd.PersistentFlags().Lookup(tokenKey))

	RootCmd.PersistentFlags().String(gitServiceTypeKey, "github", "git service type")
	viper.BindPFlag(gitServiceTypeKey, RootCmd.PersistentFlags().Lookup(gitServiceTypeKey))

	RootCmd.PersistentFlags().Int(portKey, 8080, "Listen port")
	viper.BindPFlag(portKey, RootCmd.PersistentFlags().Lookup(portKey))

	RootCmd.PersistentFlags().String(gidKey, "", "Group ID list")
	viper.BindPFlag(gidKey, RootCmd.PersistentFlags().Lookup(gidKey))

	RootCmd.PersistentFlags().String(spPrefixKey, "S", "prefix of Story Point label")
	viper.BindPFlag(spPrefixKey, RootCmd.PersistentFlags().Lookup(spPrefixKey))

	RootCmd.PersistentFlags().String(classPrefixKey, "C:", "prefix of class label")
	viper.BindPFlag(classPrefixKey, RootCmd.PersistentFlags().Lookup(classPrefixKey))

	RootCmd.PersistentFlags().String(baseURLKey, "https://github.com", "base URL of git service")
	viper.BindPFlag(baseURLKey, RootCmd.PersistentFlags().Lookup(baseURLKey))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".issummary")      // name of config file (without extension)
	viper.AddConfigPath(os.Getenv("HOME")) // adding home directory as first search path

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	viper.SetEnvPrefix("issummary") // will be uppercased automatically
	viper.AutomaticEnv()            // read in environment variables that match

	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
}

func generateIssummaryConfigFromViper() (*issummary.Config, error) {
	gidStr := viper.GetString(gidKey)
	gids := strings.Split(gidStr, ",")

	if len(gids) == 0 {
		return nil, errors.New("gid is empty")
	}

	return &issummary.Config{
		Port:              viper.GetInt(portKey),
		Token:             viper.GetString(tokenKey),
		GitServiceBaseURL: viper.GetString(baseURLKey),
		GitServiceType:    viper.GetString(gitServiceTypeKey),
		SPLabelPrefix:     viper.GetString(spPrefixKey),
		ClassLabelPrefix:  viper.GetString(classPrefixKey),
		GIDs:              gids,
	}, nil
}
