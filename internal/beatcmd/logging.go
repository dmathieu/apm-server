// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package beatcmd

import (
	"strings"

	"github.com/elastic/elastic-agent-libs/logp"
)

var (
	logVerbose        bool
	logStderr         bool
	logDebugSelectors []string
	logEnvironment    logpEnvironmentVar
	logOptions        []logp.Option
)

func configureLogging(cfg *Config) error {
	logpConfig := logp.DefaultConfig(logEnvironment.env)
	logpConfig.Beat = "apm-server"
	if cfg.Logging != nil {
		if err := cfg.Logging.Unpack(&logpConfig); err != nil {
			return err
		}
	}

	// Apply command line flags to the logging configuration.
	if logpConfig.Level > logp.InfoLevel && logVerbose {
		logpConfig.Level = logp.InfoLevel
	}
	if len(logDebugSelectors) > 0 {
		for _, selectors := range logDebugSelectors {
			logpConfig.Selectors = append(logpConfig.Selectors, strings.Split(selectors, ",")...)
		}
		logpConfig.Level = logp.DebugLevel
	}
	if logStderr {
		logpConfig.ToStderr = true
	}
	for _, opt := range logOptions {
		opt(&logpConfig)
	}

	return logp.Configure(logpConfig)
}
