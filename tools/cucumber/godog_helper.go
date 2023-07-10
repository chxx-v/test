package cucumber

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/cucumber/godog"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

type gherkinFinder struct {
	features []godog.Feature
}

func (g *gherkinFinder) Visit(node ast.Node, entering bool) ast.WalkStatus {
	if c, ok := node.(*ast.CodeBlock); ok {
		if bytes.Equal(c.Info, []byte(string("gherkin"))) && entering {
			f := godog.Feature{Name: fmt.Sprintf("%d", len(g.features)+1), Contents: c.Literal}
			g.features = append(g.features, f)
		}
		return ast.SkipChildren
	}
	return ast.GoToNext
}

func ExtractGherkinSnippetFromMarkdown(files []string) ([]godog.Feature, error) {
	features := []godog.Feature{}
	for i, fileName := range files {
		f, err := ioutil.ReadFile(fileName)
		if err != nil {
			return nil, err
		}

		f = markdown.NormalizeNewlines(f)
		exts := parser.CommonExtensions // parser.OrderedListStart | parser.NoEmptyLineBeforeBlock
		p := parser.NewWithExtensions(exts)
		root := markdown.Parse(f, p)

		fn := &gherkinFinder{}
		ast.Walk(root, fn)
		// all default feature name will be "1"
		// assign different feature name for the features map
		fn.features[0].Name = strconv.Itoa(i)
		features = append(features, fn.features...)
	}
	return features, nil
}

func MDFinder(mdPath string) ([]string, error) {
	// recursively find all files end with .feature.md
	mdRegEx, err := regexp.Compile("(?i)^*\\.feature\\.md$")
	if err != nil {
		return nil, err
	}

	mdFiles := []string{}
	err = filepath.Walk(mdPath, func(path string, info os.FileInfo, err error) error {
		if err == nil && mdRegEx.MatchString(info.Name()) {
			mdFiles = append(mdFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return mdFiles, nil
}
