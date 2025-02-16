package script

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func ShowScriptSelector(dir string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var scripts []string
	for _, f := range files {
		// アンダースコアで始まるファイルを除外
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".star") && !strings.HasPrefix(f.Name(), "_") {
			// パスの区切り文字を'/'に統一
			path := filepath.ToSlash(filepath.Join(dir, f.Name()))
			scripts = append(scripts, path)
		}
	}

	// スクリプトが見つからない場合はデフォルトを使用
	if len(scripts) == 0 {
		return "main.star", nil
	}

	// スクリプトの選択
	var selected string
	prompt := &survey.Select{
		Message: "デバッグスクリプトを選択してください:",
		Options: scripts,
	}
	err = survey.AskOne(prompt, &selected)
	if err != nil {
		return "", fmt.Errorf("script selection failed: %v", err)
	}

	return selected, nil
}
