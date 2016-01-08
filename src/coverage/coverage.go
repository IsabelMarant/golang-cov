/*coverage.go*/
/*
modification history
--------------------
2015/11/17, by Xiaoye Jiang, create
*/
/*
DESCRIPTION
*/
package coverage

import (
	"fmt"
	"strconv"
	"strings"
)

import (
	"cov_base"
	"cov_conf"
)

func PostHeader(svnInfo map[string]string, codeCov *CodeCov, config *cov_conf.ConfCover) map[string]string {
	queryHeader := make(map[string]string, 0)
	for k, v := range svnInfo {
		queryHeader[k] = v
	}
	if codeCov.LineVaild > 0 {
		queryHeader["total_line"] = strconv.FormatInt(codeCov.LineVaild, 10)
		queryHeader["cover_line"] = strconv.FormatInt(codeCov.LineCovered, 10)
	}
	if codeCov.FuncVaild > 0 {
		queryHeader["total_function"] = strconv.FormatInt(codeCov.FuncVaild, 10)
		queryHeader["cover_function"] = strconv.FormatInt(codeCov.FuncCovered, 10)
	}

	queryHeader["module_id"] = strconv.FormatInt(config.Post.ModuleId, 10)
	queryHeader["method"] = config.Post.Method

	fmt.Println(queryHeader)
	return queryHeader
}

func GetCovInfo(covFile string, skip string) (*CodeCov, error) {
	var err error
	var codeCov CodeCov

	skipModules := strings.Split(skip, ",")

	codeCov, err = GetCodeCov(covFile, skipModules)
	if err != nil {
		return &codeCov, fmt.Errorf("GetCodeCov failed:%s", err.Error())
	}
	return &codeCov, nil
}

func Coverage(config *cov_conf.ConfCover) error {
	isPostHtml := false
	if config.Cover.HtmlPath != "" {
		err := HtmlOutput(config.Cover.CoverPath, config.Cover.HtmlPath)
		//if error, not exit
		if err != nil {
			isPostHtml = false
			fmt.Println("generate html err:", err.Error())
		} else {
			isPostHtml = true
			fmt.Println("generate html success: ", config.Cover.HtmlPath)
		}
	}

	svnInfo, err := cov_base.GetSvnInfo(config.Svn.SvnPath)
	if err != nil {
		fmt.Println("GetSvnInfo failed:", err.Error())
		return err
	}
	fmt.Println(svnInfo)

	codeCov, err := GetCovInfo(config.Cover.CoverPath, config.Cover.SkipModules)
	if err != nil {
		fmt.Println("GetCovInfo failed:", err.Error())
		return err
	}
	fmt.Println(*codeCov)

	postHeader := PostHeader(svnInfo, codeCov, config)
	if err = cov_base.PostFile(postHeader, config.Post.RootUrl, "report", config.Cover.HtmlPath, isPostHtml); err != nil {
		fmt.Println("post coverage failed:", err.Error())
		return err
	}
	fmt.Println("post coverage success")
	return nil
}
