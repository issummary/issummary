package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mpppk/issummary/api"
	"github.com/mpppk/issummary/gitlab"
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
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		client := gitlab.New(os.Getenv("GITLAB_TOKEN"))

		if os.Getenv("GITLAB_BASEURL") != "" {
			client.SetBaseURL(os.Getenv("GITLAB_BASEURL"))
		}

		gidList := strings.Split(os.Getenv("GITLAB_PID"), ",")

		worksBodyFunc := func(body []byte) (interface{}, error) {
			workManager := gitlab.NewWorkManager()
			for _, gid := range gidList {
				works, err := client.ListGroupWorks(gid, "LC", "S")

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
			for _, gid := range gidList {
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
		err = http.ListenAndServe(":8080", nil)
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
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".issummary")      // name of config file (without extension)
	viper.AddConfigPath(os.Getenv("HOME")) // adding home directory as first search path
	viper.AutomaticEnv()                   // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
