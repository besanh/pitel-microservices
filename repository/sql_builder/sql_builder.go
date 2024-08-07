package sql_builder

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/tel4vn/fins-microservices/common/util"
)

// Define structs to parse JSON
type QueryCondition struct {
	Field    string           `json:"field,omitempty"`
	Operator string           `json:"operator,omitempty"`
	Value    any              `json:"value,omitempty"`
	And      []QueryCondition `json:"and,omitempty"`
	Or       []QueryCondition `json:"or,omitempty"`
}

func BuildQueryByJson(val any) (queryStr string, err error) {
	var query QueryCondition
	if err = util.ParseAnyToAny(val, &query); err != nil {
		return
	}
	return BuildQuery(query), nil
}

// Recursive function to build conditions string
func BuildQuery(condition QueryCondition) string {
	var conditions []string

	if condition.Field != "" && condition.Operator != "" {
		value := formatValue(condition.Value)
		conditionStr := fmt.Sprintf("%s %s %s", condition.Field, condition.Operator, value)
		if condition.Operator == "IN" {
			conditionStr = fmt.Sprintf("%s %s %s", condition.Field, condition.Operator, value)
		}
		conditions = append(conditions, conditionStr)
	}

	if len(condition.And) > 0 {
		var andQueryConditions []string
		for _, cond := range condition.And {
			andQueryConditions = append(andQueryConditions, BuildQuery(cond))
		}
		conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(andQueryConditions, " AND ")))
	}

	if len(condition.Or) > 0 {
		var orQueryConditions []string
		for _, cond := range condition.Or {
			orQueryConditions = append(orQueryConditions, BuildQuery(cond))
		}
		conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(orQueryConditions, " OR ")))
	}

	return strings.Join(conditions, " AND ")
}

// Function to format value for SQL query
func formatValue(value any) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("'%s'", v)
	case time.Time:
		return fmt.Sprintf("'%s'", v.Format(time.RFC3339))
	default:
		reflectValue := reflect.ValueOf(v)
		if reflectValue.Kind() == reflect.Slice {
			lenn := reflectValue.Len()
			if lenn == 0 {
				return "''"
			}
			var values []string
			for i := 0; i < lenn; i++ {
				values = append(values, formatValue(reflectValue.Index(i).Interface()))
			}
			return fmt.Sprintf("(%s)", strings.Join(values, ","))
		}
		return fmt.Sprintf("%v", v)
	}
}

func EqualQuery(field string, value any) QueryCondition {
	return QueryCondition{
		Field:    field,
		Operator: "=",
		Value:    value,
	}
}

func NotEqualQuery(field string, value any) QueryCondition {
	return QueryCondition{
		Field:    field,
		Operator: "!=",
		Value:    value,
	}
}

func GreaterThanQuery(field string, value any) QueryCondition {
	return QueryCondition{
		Field:    field,
		Operator: ">",
		Value:    value,
	}
}

func GreaterThanOrEqualQuery(field string, value any) QueryCondition {
	return QueryCondition{
		Field:    field,
		Operator: ">=",
		Value:    value,
	}
}

func LessThanQuery(field string, value any) QueryCondition {
	return QueryCondition{
		Field:    field,
		Operator: "<",
		Value:    value,
	}
}

func LessThanOrEqualQuery(field string, value any) QueryCondition {
	return QueryCondition{
		Field:    field,
		Operator: "<=",
		Value:    value,
	}
}

func LikeQuery(field string, value any) QueryCondition {
	return QueryCondition{
		Field:    field,
		Operator: "LIKE",
		Value:    value,
	}
}

func InQuery(field string, value any) QueryCondition {
	return QueryCondition{
		Field:    field,
		Operator: "IN",
		Value:    value,
	}
}
