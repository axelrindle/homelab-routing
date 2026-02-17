package app

import "go.uber.org/zap"

func (a *App) sync() {
	a.Logger.Debug("building configuration")

	src, err := a.fetchTraefikRouters()
	if err != nil {
		a.Logger.Error("failed to fetch Traefik routers", zap.Error(err))
		return
	}

	a.configuration = a.buildTraefikConfiguration(src)
	a.ready = true

	a.Logger.Debug("routers loaded", zap.Int("count", len(a.configuration.HTTP.Routers)))
}
