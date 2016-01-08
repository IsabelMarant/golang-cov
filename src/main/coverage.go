/*main.go*/
/*
modification history
--------------------
2015/11/15, by Xiaoye Jiang, create
*/
/*
DESCRIPTION
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

import (
	"cov_conf"
	"coverage"
)

var (
	help     *bool   = flag.Bool("h", false, "to show help")
	confRoot *string = flag.String("c", "../conf/", "path of config file")
	showVer  *bool   = flag.Bool("v", false, "to show version of coverage")
)

func main() {
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		return
	}
	if *showVer {
		fmt.Println("coverage: version 1.0.0.1\n")
		return
	}

	confPath := path.Join(*confRoot, "coverage.conf")
	config, err := cov_conf.CovConfLoad(confPath, *confRoot)
	if err != nil {
		fmt.Println("main(): err in CovConfLoad(): ", err.Error())
		os.Exit(1)
	}

	err = coverage.Coverage(config)
	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
