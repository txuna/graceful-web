push:
	@docker build --rm --no-cache -t txuna/graceful:1.0.8 -f Dockerfile .
	@docker push txuna/graceful:1.0.8

deploy:
	@helm upgrade --install graceful install/helm/graceful -n graceful --create-namespace

remove:
	@helm uninstall graceful -n graceful