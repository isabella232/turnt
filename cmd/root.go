// Copyright © 2017 Dave Greene davepgreene@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"errors"
	"net/url"
	"encoding/json"
	"time"
	"strconv"

	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"github.com/rapid7/turnt/utils"
	"github.com/rapid7/turnt/lib"
	"strings"
)

type configStruct struct {
	file string
	method string
	payload string
	digest string
	headers map[string]string
	identity string
	secret string
}

var conf configStruct
var requestUrl *url.URL

// RootCmd represents the base command when called without any subcommands
var TurntCmd = &cobra.Command{
	Use:   "turnt [url]",
	Short: "A helper app to make requests against Turnstile",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("URL is a required argument")
		}

		u, err := url.Parse(args[0])

		if err != nil {
			return err
		}

		requestUrl = u

		return nil
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Validate flags before we execute the main function
		if conf.identity == "" {
			return errors.New("Identity is required")
		}

		if conf.secret == "" {
			return errors.New("Secret is required")
		}

		if !utils.AlgorithmIsSupported(conf.digest) {
			errStr := fmt.Sprintf("Turnstile currently supports the following encryption algorithms: %s. You specified %s.",
				utils.GetSupportedAlgorithms(),
				conf.digest)
			return errors.New(errStr)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		date := time.Now().Unix()
		conf.headers["date"] = strconv.Itoa(int(date))

		uri := requestUrl.Path
		host := requestUrl.Host

		algorithm := utils.GetAlgorithmType(conf.digest)
		digest := lib.GenerateDigest(algorithm, conf.payload)
		conf.headers["digest"] = digest

		method := strings.ToUpper(conf.method)

		log.WithFields(log.Fields{"method": method}).Info("Using method")
		log.WithFields(log.Fields{"uri": uri}).Info("Using URI")
		log.WithFields(log.Fields{"host": host}).Info("Using host")
		log.WithFields(log.Fields{"date": date}).Info("Using date")
		log.WithFields(log.Fields{"identity": conf.identity}).Info("Using identity")
		log.WithFields(log.Fields{"digest": digest}).Info("Using digest")

		signature := lib.GenerateSignature(algorithm, conf.identity, conf.secret, digest, method, uri, host, date)
		log.WithFields(log.Fields{"signature": signature}).Info("Using signature")

		authorization := lib.GenerateAuthorization(algorithm, conf.identity, signature)
		log.WithFields(log.Fields{"authorization": authorization}).Info("Using authorization")
		conf.headers["authorization"] = authorization

		err, out := lib.GenerateRequest(method, requestUrl.String(), conf.payload, conf.headers)
		if err != nil {
			return err
		}

		log.Infof("Response: \n%s", out.String())

		return nil
	},
}

func Execute() {
	if err := TurntCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	log.SetLevel(log.InfoLevel)
	formatter := prefixed.TextFormatter{
		FullTimestamp: true,
		TimestampFormat:time.RFC822,
	}
	log.SetFormatter(&formatter)

	setOptionalFlags()
	setRequiredFlags()
}

// Set optional flags
func setOptionalFlags() {
	TurntCmd.Flags().StringVarP(&conf.method, "method", "X", "GET", "HTTP request method")

	TurntCmd.Flags().StringVarP(&conf.payload, "payload", "d", "{}", "HTTP request payload")
	TurntCmd.Flags().StringVar(&conf.digest, "digest", "SHA256", "Digest signing scheme")

	headerStr := TurntCmd.Flags().StringP("header", "H", "{}", "HTTP request headers")

	// Get a head start and unmarshal provided headers
	err := json.Unmarshal([]byte(*headerStr), &conf.headers)
	if err != nil {
		panic(err)
	}
}

// Set required flags
func setRequiredFlags() {
	TurntCmd.Flags().StringVarP(&conf.identity, "identity", "u", "", "Identity key for the request")
	TurntCmd.Flags().StringVarP(&conf.secret, "secret", "p", "", "Secret key for the request")
}
