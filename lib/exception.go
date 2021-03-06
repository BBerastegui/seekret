package lib

import (
	"io/ioutil"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v2"
)

type exceptionYaml struct {
	Rule    *string
	Object  *string
	Line    *int
	Content *string
}

type Exception struct {
	Rule    *string
	Object  *regexp.Regexp
	Line    *int
	Content *regexp.Regexp
}

func (s *Seekret) AddException(exception Exception) {
	s.exceptionList = append(s.exceptionList, exception)
}

func (s *Seekret) LoadExceptionsFromFile(file string) error {
	var exceptionYamlList []exceptionYaml

	if file == "" {
		return nil
	}

	filename, _ := filepath.Abs(file)
	yamlData, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlData, &exceptionYamlList)
	if err != nil {
		return err
	}

	for _, v := range exceptionYamlList {
		exception := Exception{
			Rule: v.Rule,
			Line: v.Line,
		}
		if v.Object != nil {
			exception.Object = regexp.MustCompile("(?i)" + *v.Object)
		}
		if v.Content != nil {
			exception.Content = regexp.MustCompile("(?i)" + *v.Content)
		}
		s.AddException(exception)
	}

	return nil
}

func exceptionCheck(exceptionList []Exception, secret Secret) bool {
	for _, e := range exceptionList {
		match := true

		if match && e.Rule != nil && *e.Rule != secret.Rule.Name {
			match = false
		}
		if match && e.Line != nil && *e.Line != secret.Nline {
			match = false
		}
		if match && e.Object != nil && !(*e.Object).MatchString(secret.Object.Name) {
			match = false
		}
		if match && e.Content != nil && !(*e.Content).MatchString(secret.Line) {
			match = false
		}

		if match {
			return true
		}
	}

	return false
}
