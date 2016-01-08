/*svn.go - Get svn info*/
/*
modification history
--------------------
2015/11/15, by Xiaoye Jiang, create
*/
/*
DESCRIPTION
*/
package cov_base

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func GetSvnInfo(svnPath string) (map[string]string, error) {
	prefixSvn := map[string]string{
		"author":      "Last Changed Author:",
		"revision":    "Last Changed Rev:",
		"commit_time": "Last Changed Date:",
	}
	svnInfo := make(map[string]string, len(prefixSvn))
	out, err := exec.Command("svn", "info", svnPath).Output()
	if err != nil {
		return svnInfo, err
	}
	strList := strings.Split(string(out), "\n")
	for _, line := range strList {
		for key, value := range prefixSvn {
			if strings.HasPrefix(line, value) {
				svnInfo[key] = strings.TrimSpace(strings.TrimPrefix(line, value))
			}
		}
	}
	if value, ok := svnInfo["commit_time"]; ok {
		date := value[:strings.IndexByte(value, '(')-1]
		timeForm := "2006-01-02 15:04:05 +0800 MST"
		tmp, err := time.Parse(timeForm, fmt.Sprintf("%s CST", date))
		if err != nil {
			return svnInfo, err
		}
		svnInfo["commit_time"] = fmt.Sprintf("%d", tmp.Unix())
	}
	return svnInfo, nil
}
