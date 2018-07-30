// Copyright [2018] [Rafał Korepta]
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:generate go run script/gen_markdown_docs.go
package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugFlag      = "debug"
	configPathFlag = "cfg_path"
	configFlag     = "config"
)

var (
	// Version will be populated with binary semver by the linker
	// during the build process.
	// See https://blog.cloudflare.com/setting-go-variables-at-compile-time/
	// and https://golang.org/cmd/link/ in section Flags `-X importpath.name=value`.
	Version string

	// Commit will be populated with correct git commit id by the linker
	// during the build process.
	// See https://blog.cloudflare.com/setting-go-variables-at-compile-time/
	// and https://golang.org/cmd/link/ in section Flags `-X importpath.name=value`.
	Commit string
)

//RootCmd is the root of all cobra command in this project
var RootCmd = &cobra.Command{
	Use:   "portal-backend",
	Short: "The Email microservice",
	Long: `To get started run the serve subcommand which will start a server

	portal-backend serve

After that you can test it with the client subcommand:

	portal-backend sendmail a b c

Or over HTTP 1.1 with curl:

	curl -X POST -k https://localhost:9091/v1alpha1/emial -d '{"message":"abc"}''
`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		zap.L().Fatal("Failed to execute root command", zap.Error(err))
	}
}

func init() {
	//First initialize logger
	newLogger, err := zap.NewProduction(
		zap.AddStacktrace(zap.ErrorLevel),
		zap.Fields(
			zap.Field{
				Key:    "commit",
				Type:   zapcore.StringType,
				String: Commit,
			},
			zap.Field{
				Key:    "version",
				Type:   zapcore.StringType,
				String: Version,
			},
		))
	if err != nil {
		log.Fatalf("Unable to create logger. Error: %v", err)
	}
	zap.ReplaceGlobals(newLogger)

	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolP(DebugFlag, "d", false, "turn on debug logging")
	RootCmd.PersistentFlags().String(configPathFlag, ".", "Relative path where config resides")
	RootCmd.PersistentFlags().String(configFlag, ".portal-backend", "config file (default is $HOME/.portal-backend.yaml)")
	if err := viper.BindPFlags(RootCmd.PersistentFlags()); err != nil {
		zap.L().Error("Can not bind persistent flags")
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigName(viper.GetString(configFlag)) // name of config file (without extension)
	viper.AddConfigPath(viper.GetString(configPathFlag))
	viper.AddConfigPath("$HOME")
	viper.SetEnvPrefix("SMACC")
	viper.AutomaticEnv()

	// Update global logger if debug flag is chosen
	cfg := zap.NewProductionConfig()
	if viper.GetBool("debug") {
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}
	newLogger, err := cfg.Build(zap.AddStacktrace(zap.ErrorLevel),
		zap.Fields(
			zap.Field{
				Key:    "commit",
				Type:   zapcore.StringType,
				String: Commit,
			},
			zap.Field{
				Key:    "version",
				Type:   zapcore.StringType,
				String: Version,
			},
		))
	if err != nil {
		log.Fatalf("Unable to create logger. Error: %v", err)
	}
	zap.ReplaceGlobals(newLogger)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		zap.S().Errorw("Failed to read from config file",
			"configFile", viper.ConfigFileUsed(),
			"error", err)
	}
}
