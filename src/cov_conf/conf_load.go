/*conf_load.go - load cov_conf for coverage*/
/*
modification history
--------------------
2015/11/15, by Xiaoye Jiang, create
*/
/*
DESCRIPTION
*/
package cov_conf

import (
	"fmt"
)

import (
	"code.google.com/p/gcfg"
)

type ConfCover struct {
	Cover ConfCov
	Post  ConfPost
	Svn   ConfSvn
}

type ConfCov struct {
	CoverPath   string
	SkipModules string
	HtmlPath    string
}

type ConfPost struct {
	ModuleId int64
	Method   string
	RootUrl  string
}

type ConfSvn struct {
	SvnPath string
}

func CovConfLoad(filePath string, confRoot string) (*ConfCover, error) {
	var cfg ConfCover

	//read config from file
	err := gcfg.ReadFileInto(&cfg, filePath)
	if err != nil {
		return nil, err
	}

	err = cfg.Check(confRoot)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

//TODO: check cfg.Cover.SkipModules
//TODO: check cfg.Cover.htmlpath
func (cfg *ConfCover) Check(confRoot string) error {
	if cfg.Cover.CoverPath == "" {
		return fmt.Errorf("Cover filePath must not be empty")
	}

	if cfg.Post.ModuleId <= 0 {
		return fmt.Errorf("ModuleId[%d] is too small", cfg.Post.ModuleId)
	}
	if cfg.Post.Method == "" {
		return fmt.Errorf("Method must not be empty")
	}
	if cfg.Post.RootUrl == "" {
		return fmt.Errorf("RootUrl must not be empty")
	}

	if cfg.Svn.SvnPath == "" {
		return fmt.Errorf("SvnPath must not be empty")
	}

	return nil
}
