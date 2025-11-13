# Environment variables should be set in .env file or exported before running make commands
# See .env.example for required variables

# Auto-load .env file if it exists (using shell)
ifneq (,$(wildcard ./.env))
    $(shell export $$(grep -v '^#' .env | xargs))
endif

GOSTRIPE_PORT?=4001
API_PORT?=4002

## build: builds all binaries
build: clean build_front build_back
	@printf "All binaries built!\n"

## clean: cleans all binaries and runs go clean
clean:
	@echo "Cleaning..."
	@- rm -f dist/*
	@go clean
	@echo "Cleaned!"

## build_front: builds the front end
build_front:
	@echo "Building front end..."
	@go build -o dist/gostripe ./cmd/web
	@echo "Front end built!"

## build_back: builds the back end
build_back:
	@echo "Building back end..."
	@go build -o dist/gostripe_api ./cmd/api
	@echo "Back end built!"

## start: starts front and back end
start: stop start_front start_back
	
## start_front: starts the front end
start_front: build_front
	@echo "Starting the front end..."
	@if [ -f .env ]; then \
		export $$(grep -v '^#' .env | xargs); \
	fi; \
	if [ -z "$$STRIPE_KEY" ] || [ -z "$$STRIPE_SECRET" ] || [ -z "$$DSN" ]; then \
		echo "Error: STRIPE_KEY, STRIPE_SECRET, and DSN environment variables must be set"; \
		echo "Create a .env file or export them before running make start"; \
		exit 1; \
	fi; \
	env STRIPE_KEY=$$STRIPE_KEY STRIPE_SECRET=$$STRIPE_SECRET DSN=$$DSN ./dist/gostripe -port=${GOSTRIPE_PORT} &
	@sleep 1
	@echo "Front end running on port ${GOSTRIPE_PORT}!"

## start_back: starts the back end
start_back: build_back
	@echo "Starting the back end..."
	@if [ -f .env ]; then \
		export $$(grep -v '^#' .env | xargs); \
	fi; \
	if [ -z "$$STRIPE_KEY" ] || [ -z "$$STRIPE_SECRET" ] || [ -z "$$DSN" ]; then \
		echo "Error: STRIPE_KEY, STRIPE_SECRET, and DSN environment variables must be set"; \
		echo "Create a .env file or export them before running make start"; \
		exit 1; \
	fi; \
	env STRIPE_KEY=$$STRIPE_KEY STRIPE_SECRET=$$STRIPE_SECRET DSN=$$DSN ./dist/gostripe_api -port=${API_PORT} &
	@sleep 1
	@echo "Back end running on port ${API_PORT}!"

## stop: stops the front and back end
stop: stop_front stop_back
	@echo "All applications stopped"

## stop_front: stops the front end
stop_front:
	@echo "Stopping the front end..."
	@-pkill -9 -f "gostripe" 2>/dev/null || true
	@-if lsof -ti :${GOSTRIPE_PORT} >/dev/null 2>&1; then \
		lsof -ti :${GOSTRIPE_PORT} | xargs kill -9 2>/dev/null || true; \
	fi
	@sleep 1
	@echo "Stopped front end"

## stop_back: stops the back end
stop_back:
	@echo "Stopping the back end..."
	@-pkill -9 -f "gostripe_api" 2>/dev/null || true
	@-if lsof -ti :${API_PORT} >/dev/null 2>&1; then \
		lsof -ti :${API_PORT} | xargs kill -9 2>/dev/null || true; \
	fi
	@sleep 1
	@echo "Stopped back end"

