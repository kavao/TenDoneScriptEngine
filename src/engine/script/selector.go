package script

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
)

func ShowScriptSelector(debugDir string) (string, error) {
	// デバッグディレクトリ内の.starファイルを検索
	var scripts []string
	err := filepath.Walk(debugDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".star" {
			// パスの区切り文字を'/'に統一
			path = filepath.ToSlash(path)
			scripts = append(scripts, path)
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to list debug scripts: %v", err)
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
