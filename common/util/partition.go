package util

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tel4vn/fins-microservices/model"
)

func GetJoinPartTemplate(partTemplateList []model.PartitionTemplate) string {
	partitionList := make([]string, 0, len(partTemplateList))
	for _, item := range partTemplateList {
		index := item.IndexTemplate
		key := item.KeyTemplate
		body := item.BodyTemplate
		var indexKey string
		if key != "" {
			if body != "" {
				indexKey = fmt.Sprintf("%v", index) + ":" + key + ":1"
			} else {
				indexKey = fmt.Sprintf("%v", index) + ":" + key + ":0"
			}
		} else {
			indexKey = fmt.Sprintf("%v", index)
		}
		partitionList = append(partitionList, indexKey)
	}
	return strings.Join(partitionList, ";")
}

func GetPartitionContentTemplate(templateContent string) ([]model.PartitionTemplate, bool) {
	countLeftDelimiter := 0
	wrongFormat := false

	RightDelimiter := strings.Split(templateContent, "}}")
	lengthRightDelimiter := len(RightDelimiter)
	partTemplateList := make([]model.PartitionTemplate, 0, lengthRightDelimiter)

	for index, elem := range RightDelimiter {
		leftDelimiter := strings.Split(elem, "{{")
		if len(leftDelimiter) == 1 {
			bodyTemplate := leftDelimiter[0]
			if bodyTemplate != "" {
				partTemplate := model.PartitionTemplate{
					BodyTemplate:  bodyTemplate,
					KeyTemplate:   "",
					IndexTemplate: index + 1,
				}
				partTemplateList = append(partTemplateList, partTemplate)
			}
		} else if len(leftDelimiter) == 2 {
			countLeftDelimiter++
			keyTemplate := leftDelimiter[1]
			checkFormat := CheckPatternKeyTemplate(keyTemplate)
			if !checkFormat {
				wrongFormat = true
				break
			} else {
				bodyTemplate := leftDelimiter[0]
				partTemplate := model.PartitionTemplate{
					BodyTemplate:  bodyTemplate,
					KeyTemplate:   keyTemplate,
					IndexTemplate: index + 1,
				}
				partTemplateList = append(partTemplateList, partTemplate)
			}
		} else {
			wrongFormat = true
			break
		}
	}

	// the number of '{{' should be equal to '}}'. Otherwise wrong format
	if lengthRightDelimiter != countLeftDelimiter+1 {
		wrongFormat = true
	}
	if wrongFormat {
		return nil, true
	} else {
		return partTemplateList, false
	}
}

func CheckPatternKeyTemplate(templateKey string) bool {
	r, _ := regexp.Compile("^[a-zA-Z][a-zA-Z0-9_/-:]+$")
	return r.MatchString(templateKey)
}

func CheckTemplate(templateContent string, isContent bool) (templateContentNew string, keys []string, ok bool) {
	countLeftDelimiter := 0
	wrongFormat := false

	RightDelimiter := strings.Split(templateContent, "}}")
	lengthRightDelimiter := len(RightDelimiter)

	for _, elem := range RightDelimiter {
		leftDelimiter := strings.Split(elem, "{{")
		if len(leftDelimiter) == 1 {
			bodyTemplate := leftDelimiter[0]
			if bodyTemplate != "" {
				keys = append(keys, "")
				templateContentNew += bodyTemplate
			}
		} else if len(leftDelimiter) == 2 {
			countLeftDelimiter++
			keyTemplate := leftDelimiter[1]
			keys = append(keys, keyTemplate)
			if !isContent {
				checkFormat := CheckPatternKeyTemplate(keyTemplate)
				if !checkFormat {
					wrongFormat = true
					break
				} else {
					bodyTemplate := leftDelimiter[0]
					templateContentNew += bodyTemplate
				}
			} else {
				bodyTemplate := leftDelimiter[0]
				templateContentNew += bodyTemplate
			}
		} else {
			wrongFormat = true
			break
		}
	}

	// the number of '{{' should be equal to '}}'. Otherwise wrong format
	if lengthRightDelimiter != countLeftDelimiter+1 {
		ok = true
	}
	if wrongFormat {
		return "", keys, true
	} else {
		return templateContentNew, keys, false
	}
}
