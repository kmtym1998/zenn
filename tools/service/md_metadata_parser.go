package service

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	separator                     = "---"
	ARTICLE_TYPE_TECH ArticleType = "tech"
	ARTICLE_TYPE_IDEA ArticleType = "idea"
)

type ArticleType string

type MDMetadata struct {
	Title     string      `yaml:"title"`
	Emoji     string      `yaml:"emoji"`
	Type      ArticleType `yaml:"type"`
	Topics    []string    `yaml:"topics"`
	Published bool        `yaml:"published"`
}

func ParseMDMetadata(path string) (*MDMetadata, error) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v", os.Args[0], err)
		return nil, err
	}
	defer f.Close()

	var frontMatter []string
	sepCount := 0

	// １つ目のセパレータから２つ目のセパレータまでを読み込む
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == separator {
			if sepCount > 0 {
				break
			}
			sepCount += 1
			continue
		}
		frontMatter = append(frontMatter, line)
	}

	// パース
	buf := []byte(strings.Join(frontMatter, "\n"))
	metadata := MDMetadata{}
	err = yaml.Unmarshal(buf, &metadata)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v", os.Args[0], err)
		return nil, err
	}

	return &metadata, nil
}
