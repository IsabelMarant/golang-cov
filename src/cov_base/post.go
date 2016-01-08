/*post.go - post cov info to cov.baidu.com*/
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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type RespBody struct {
	Result string `json:"result"`
	ErrMsg string `json:"err_msg"`
}

func PostFile(queryHeader map[string]string, rootUrl, paramName, filename string, isPostHtml bool) error {
	var err error
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	if isPostHtml {
		fileWriter, err := bodyWriter.CreateFormFile(paramName, filename)
		if err != nil {
			return err
		}

		fd, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer fd.Close()

		_, err = io.Copy(fileWriter, fd)
		if err != nil {
			return err
		}
	}

	for key, value := range queryHeader {
		_ = bodyWriter.WriteField(key, value)
	}

	contentType := bodyWriter.FormDataContentType()
	err = bodyWriter.Close()
	if err != nil {
		return err
	}

	resp, err := http.Post(rootUrl, contentType, bodyBuf)
	if err != nil {
		return fmt.Errorf("http post err: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("post err: %d %s", resp.StatusCode, resp.Status)
	}

	var respBody RespBody
	body, _ := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &respBody); err != nil {
		return fmt.Errorf("body json Unmarshal failed:%s", err.Error())
	}

	if respBody.Result == "SUCC" {
		return nil
	}
	return fmt.Errorf("post cov err: %s", respBody.ErrMsg)
}
