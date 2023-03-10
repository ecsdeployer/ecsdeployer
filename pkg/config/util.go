package config

import (
	"errors"
)

func ExtractCommonTaskAttrs(obj any) (*CommonTaskAttrs, error) {

	if thing, ok := obj.(IsTaskStruct); ok {
		cta := thing.GetCommonTaskAttrs()
		return &cta, nil
	}

	// OLD VERSION
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
