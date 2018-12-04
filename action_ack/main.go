package main

import (
	"os"

	docopt "github.com/docopt/docopt-go"
	"github.com/kovetskiy/lorg"
	"github.com/kovetskiy/toml"
	karma "github.com/reconquest/karma-go"
)

var (
	logger           *lorg.Log
	version          = "[manual build]"
	definedMessenger = "mattermost"
	usage            = "chattixd " + version + ` for ` + definedMessenger + `

Usage:
  chattixd [--config <path>]

Options:
    -c --config <path>  Path to config file 
                         [default: /etc/chattix/chattixd.conf]
                                                                               
`
)

func main() {
	destiny := karma.Describe(
		"method", "main",
	).Describe(
		"version", version,
	)

	logger = lorg.NewLog()
	conf := &config{}

	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		logger.Fatal(destiny.Format(err, "can't parse args"))
	}

	configFile := args["--config"].(string)

	if _, err := os.Stat(configFile); !os.IsNotExist(err) {
		_, err = toml.DecodeFile(configFile, conf)
		if err != nil {
			logger.Fatal(
				destiny.Format(
					err,
					"can't read config file %s",
					configFile,
				),
			)
		}
	} else {
		logger.Infof(
			"config file %s doesn't exist, use default configuration",
			configFile,
		)

		_, err = toml.Decode(defaultConfiguration, conf)
		if err != nil {
			logger.Fatal(
				destiny.Describe(
					"default configuration", defaultConfiguration,
				).Describe(
					"error", err,
				).Reason(
					"can't parse default configuration",
				),
			)
		}
	}

	parseEnvironmentVariables(conf)

	actionService := newActionACKService(
		conf,
		logger,
		definedMessenger,
	)

	actionService.run()
}
