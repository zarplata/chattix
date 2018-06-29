package main

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/kovetskiy/lorg"
	"github.com/kovetskiy/toml"
	karma "github.com/reconquest/karma-go"
)

var (
	logger      *lorg.Log
	version     = "[manual build]"
	definedChat = "mattermost"
	usage       = "chattixd " + version + `

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

	_, err = toml.DecodeFile(args["--config"].(string), conf)
	if err != nil {
		logger.Fatal(
			destiny.Format(
				err,
				"can't read config file %s",
				args["--config"].(string),
			),
		)
	}

	actionService := newActionACKService(
		conf,
		logger,
		definedChat,
	)

	actionService.run()
}
