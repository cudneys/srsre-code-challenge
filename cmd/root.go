/*
Copyright © 2021 GPC SRE

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"

	"github.com/cudneys/srsre-code-challenge/password"
	pw "github.com/cudneys/srsre-code-challenge/password"
	tools "github.com/cudneys/srsre-code-challenge/tools"

	"github.com/Depado/ginprom"
	docs "github.com/cudneys/srsre-code-challenge/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Password Generator API
// @version 1.0

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

var cfgFile string
var host string
var port string

// Doc returned for successes
type PasswordResponse struct {
	Password string `json:"password"`
	Length   int    `json:"length"`
	Duration int64  `json:"generated_in_ns"`
	Digits   int    `json:"digits"`
	Symbols  int    `json:"symbols"`
	AR       bool   `json:"allow_repeat"`
	Sha256   string `json:"sha_256_sum"`
	Sha512   string `json:"sha_512_sum"`
}

// Doc returned for errors
type PasswordError struct {
	Error   string `json:"error"`
	Length  int    `json:"length"`
	Digits  int    `json:"digits"`
	Symbols int    `json:"symbols"`
	AR      bool   `json:"allow_repeat"`
}

type ValidationResponse struct {
	Password  string `json:"password"`
	Validates bool   `json:"validates"`
	Error     string `json:"error"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "srsre-code-challenge",
	Short: "Sr SRE Code Challenge",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("PORT: %s\n", port)
		fmt.Printf("HOST: %s\n", host)
		runServer(host, port)
	},
}

// @BasePath /api/v1
// PasswordGenerator godoc
// @Summary Validates Passwords
// @Schemes
// @Description Validates Passwoed
// @Tags Password
// @Accept json
// @Produce json
// @Param password query string true "The password to validate"
// @Success 200 {object} ValidationResponse
// @Failure 499 {object} ValidationResponse
// @Router /validate [get]
// Handles the GET /validate request
func getValidate(c *gin.Context) {
	password := c.Query("password")
	v, err := pw.Validate(password)
	if err != nil {
		c.JSON(499, ValidationResponse{password, false, err.Error()})
		return
	}
	c.JSON(200, ValidationResponse{password, v, ""})
}

// @BasePath /api/v1
// PasswordGenerator godoc
// @Summary Generates Passwords
// @Schemes
// @Description Generate Passwoed
// @Tags Password
// @Accept json
// @Produce json
// @Param length query int true "Length of the password"
// @Param digits query int false "Number of digits (Default: Length/4)"
// @Param symbols query int false "Number of symbols (Default: Length/4)"
// @Param allow_repeat query bool false "Allow repeated chars (Default: true)"
// @Success 200 {object} PasswordResponse
// @Failure 400 {object} PasswordError
// @Router /generate [get]
// Handles the GET /generate request
func getPassword(c *gin.Context) {
	start := time.Now()

	var allowRepeat bool

	length, err := strconv.Atoi(c.DefaultQuery("length", "64"))
	if err != nil {
		c.JSON(400, PasswordError{fmt.Sprintf("Invalid Length Parameter: %s", c.DefaultQuery("length", "64")), 0, 0, 0, true})
		return
	}

	digits, _ := strconv.Atoi(c.DefaultQuery("digits", "0"))
	symbols, _ := strconv.Atoi(c.DefaultQuery("symbols", "0"))
	repeat := c.DefaultQuery("allow_repeat", "true")

	if repeat == "false" {
		allowRepeat = false
	} else {
		allowRepeat = true
	}

	digits = password.GetDefaultValue(digits, length)
	symbols = password.GetDefaultValue(symbols, length)

	if length < 24 {
		c.JSON(400, PasswordError{"Acceptable passwords must be at least 24 characters in length", length, digits, symbols, allowRepeat})
		return
	}

	password, err := pw.Generate(length, digits, symbols, allowRepeat)

	if err != nil {
		c.JSON(400, PasswordError{err.Error(), length, digits, symbols, allowRepeat})
		return
	}

	sha256Sum, _ := tools.GetSum(password, "256")
	sha512Sum, _ := tools.GetSum(password, "512")

	duration := time.Since(start)

	c.JSON(200, PasswordResponse{password, length, duration.Nanoseconds(), digits, symbols, allowRepeat, sha256Sum, sha512Sum})
}

// Build and launch the HTTP API server
func runServer(host, port string) {
	bindAddr := strings.Join([]string{host, port}, ":")

	router := gin.Default()

	p := ginprom.New(
		ginprom.Engine(router),
		ginprom.Namespace("password"),
		ginprom.Subsystem("api"),
		ginprom.Path("/metrics"),
	)
	router.Use(p.Instrument())

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := router.Group("/api/v1")
	{
		v1.GET("/generate", getPassword)
		v1.GET("/validate", getValidate)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.Run(bindAddr)
}

// Execute the command
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.srsre-code-challenge.yaml)")
	rootCmd.Flags().StringVarP(&host, "host", "H", tools.GetEnvValue("BIND_HOST", "0.0.0.0"), "The host to bind to (Env Var: BIND_HOST)")
	rootCmd.Flags().StringVarP(&port, "port", "p", tools.GetEnvValue("BIND_PORT", "2112"), "The port to bind to (Env Var: BIND_PORT)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

// Initalize config
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigName(".srsre-code-challenge")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
