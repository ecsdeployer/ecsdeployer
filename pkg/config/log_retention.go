package config

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	logTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/invopop/jsonschema"
)

var (
	ErrInvalidLogRetention = errors.New("invalid log retention")
)

var logRetentionRegex *regexp.Regexp

type LogRetention struct {
	days int32
}

func init() {
	logRetentionRegex = regexp.MustCompile("^(forever|[1-9][0-9]*)$")
}

func (obj *LogRetention) Forever() bool {
	return obj.days == -1
}

func (obj *LogRetention) Days() int32 {
	return obj.days
}

func (obj *LogRetention) ToAwsInt32() *int32 {
	return &obj.days
}

func (obj *LogRetention) EqualsLogGroup(group logTypes.LogGroup) bool {
	if group.RetentionInDays == nil {
		return obj.Forever()
	}

	return obj.days == *group.RetentionInDays
}

func ParseLogRetention[T int32 | int64 | int | string](value T) (LogRetention, error) {

	strVal, isString := any(value).(string)

	if isString {
		if !logRetentionRegex.MatchString(strVal) {
			return LogRetention{}, fmt.Errorf("%w: Log retention must be number of days or 'forever'", ErrInvalidLogRetention)
		}

		if strVal == "forever" {
			return LogRetention{days: -1}, nil
		}

		days, err := strconv.ParseInt(strVal, 10, 32)
		if err != nil {
			// return LogRetention{}, err
			return LogRetention{}, fmt.Errorf("%w: %s", ErrInvalidLogRetention, err.Error())
		}

		return ParseLogRetention(days)
	}

	var intVal int32

	switch v := any(value).(type) {
	case int64:
		intVal = int32(v)
	case int32:
		intVal = v
	case int:
		intVal = int32(v)
	default:
		return LogRetention{}, fmt.Errorf("%w: somehow got a nonstandard integer to the log retention parser", ErrInvalidLogRetention)
	}

	if intVal <= 0 {
		return LogRetention{}, fmt.Errorf("%w: Log retention must be more than 1 day or 'forever'", ErrInvalidLogRetention)
	}

	return LogRetention{days: intVal}, nil
}

func (a *LogRetention) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	obj, err := ParseLogRetention(str)
	if err != nil {
		return NewValidationError(err)
	}

	*a = obj

	return nil
}

func (obj LogRetention) MarshalJSON() ([]byte, error) {
	if obj.Forever() {
		return []byte(`"forever"`), nil
	}
	return []byte(fmt.Sprintf("%d", obj.days)), nil
}

func (LogRetention) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:        "string",
				Const:       "forever",
				Description: "Retain logs forever",
			},
			{
				Type:        "string",
				Pattern:     "^[1-9][0-9]*$",
				Description: "The number of days to retain logs",
			},
			{
				Type:        "integer",
				Minimum:     1,
				Description: "The number of days to retain logs",
			},
		},
		Description: "How long to retain logs in CloudWatch logs",
	}
}
