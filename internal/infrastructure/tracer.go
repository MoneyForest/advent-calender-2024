package infrastructure

import "context"

type TracerProviderWrapper interface {
	Shutdown(context.Context) error
}

func InitTracer(env string) (TracerProviderWrapper, error) {
	switch env {
	case "dev":
		return InitDatadog() // または InitDatadogWithDDOTel()
	default:
		return InitJaeger()
	}
}
