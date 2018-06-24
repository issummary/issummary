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
	"github.com/issummary/issummary/gitlab"
	"github.com/joho/godotenv"
	"github.com/mpppk/gitany"
	"github.com/mpppk/gitany/etc"
	gitanygitlab "github.com/mpppk/gitany/gitlab"

	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "issummary",
	Short: "issue summary viewer",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		config, err := generateIssummaryConfig()
		if err != nil {
			panic(err)
		}

		fmt.Printf("%#v\n", config)

		gitany.RegisterClientGenerator(&gitanygitlab.ClientGenerator{}) // FIXME
		serviceConfig := &etc.ServiceConfig{                            // FIXME
			Host:     "gitlab.com",
			Type:     "gitlab",
			Token:    config.Token,
			Protocol: "https",
		}

		gitanyClient, err := gitany.GetClient(context.Background(), serviceConfig) // FIXME
		if err != nil {
			panic(err)
		}

		client := gitlab.New(config.Token, gitanyClient)

		if config.GitServiceBaseURL != "" {
			client.SetBaseURL(config.GitServiceBaseURL)
		}

		worksBodyFunc := func(body []byte) (interface{}, error) {
			workManager := gitlab.NewWorkManager()
			for _, gid := range config.GIDs {
				works, err := client.ListGroupWorks(ctx, gid, "LC", "S")

				if err != nil {
					return nil, err
				}

				workManager.AddWorks(works)
			}

			if err := workManager.ResolveDependencies(); err != nil {
				return nil, err
			}
			sortedWorks, err := workManager.GetSortedWorks()

			if err != nil {
				return nil, err
			}
			return sortedWorks, nil
		}

		milestonesBodyFunc := func(body []byte) (interface{}, error) {
			var allMilestones []*gitlab.Milestone
			for _, gid := range config.GIDs {
				milestones, err := client.ListGroupMilestones(gid)

				if err != nil {
					panic(err)
				}

				allMilestones = append(allMilestones, milestones...)
			}

			return allMilestones, nil
		}

		statikFS, err := fs.New()
		if err != nil {
			log.Fatal(err)
		}

		http.Handle("/", http.FileServer(statikFS))

		http.HandleFunc("/api/works", api.CreateJsonHandleFunc(worksBodyFunc))
		http.HandleFunc("/api/milestones", api.CreateJsonHandleFunc(milestonesBodyFunc))
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
	RootCmd.PersistentFlags().String("token", "", "git repository service token")
	viper.BindPFlag("token", RootCmd.PersistentFlags().Lookup("token"))
	RootCmd.PersistentFlags().Int("port", 8080, "Listen port")
	viper.BindPFlag("port", RootCmd.PersistentFlags().Lookup("port"))
	RootCmd.PersistentFlags().String("gid", "", "Group ID list")
	viper.BindPFlag("gid", RootCmd.PersistentFlags().Lookup("gid"))
	fmt.Println("base-url")
	fmt.Println(viper.GetString("base-url"))
	RootCmd.PersistentFlags().String("base-url", viper.GetString("base-url"), "GitLab base URL")
	viper.BindPFlag("base-url", RootCmd.PersistentFlags().Lookup("base-url"))
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

	viper.Set("base-url", viper.GetString("base_url"))
}

type Config struct {
	Port              int
	Token             string
	GitServiceBaseURL string
	GIDs              []string
}

func generateIssummaryConfig() (*Config, error) {
	gidStr := viper.GetString("gid")
	gids := strings.Split(gidStr, ",")

	if len(gids) == 0 {
		return nil, errors.New("gid is empty")
	}

	return &Config{
		Port:              viper.GetInt("port"),
		Token:             viper.GetString("token"),
		GitServiceBaseURL: viper.GetString("base-url"),
		GIDs:              gids,
	}, nil
}
