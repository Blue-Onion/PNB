.PHONY: run

run:
	@bash -c '\
		set -m; \
		(cd backend && exec go run ./cmd -serve) & backend_pid=$$!; \
		(cd frontend && exec bun run dev) & frontend_pid=$$!; \
		trap "kill -- -$$backend_pid -$$frontend_pid 2>/dev/null || true; wait" EXIT; \
		trap "exit 130" INT; \
		trap "exit 143" TERM; \
		while kill -0 "$$backend_pid" "$$frontend_pid" 2>/dev/null; do sleep 1; done; \
		if kill -0 "$$backend_pid" 2>/dev/null; then \
			wait "$$frontend_pid"; \
		else \
			wait "$$backend_pid"; \
		fi'
