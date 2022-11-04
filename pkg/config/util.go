package config

import (
	"errors"
)

func ExtractCommonTaskAttrs(obj interface{}) (*CommonTaskAttrs, error) {
	switch v := obj.(type) {
	case *Service:
		return &v.CommonTaskAttrs, nil

	case *ConsoleTask:
		return &v.CommonTaskAttrs, nil

	case *PreDeployTask:
		return &v.CommonTaskAttrs, nil

	case *CronJob:
		return &v.CommonTaskAttrs, nil

	case Service:
		return &v.CommonTaskAttrs, nil

	case ConsoleTask:
		return &v.CommonTaskAttrs, nil

	case PreDeployTask:
		return &v.CommonTaskAttrs, nil

	case CronJob:
		return &v.CommonTaskAttrs, nil

	case CommonTaskAttrs:
		return &v, nil

	case *CommonTaskAttrs:
		return v, nil

	default:
		return nil, errors.New("provided struct does not embed CommonTaskAttrs")
	}

}
