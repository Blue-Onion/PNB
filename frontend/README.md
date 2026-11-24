# txShield frontend

From the repository root, run `make run` and open http://localhost:3000.
Install the frontend dependencies with `bun install` in this directory first.

Copy `.env.example` to `.env` and set `NEXT_PUBLIC_API_BASE_URL` to the
browser-accessible Go backend URL. The example uses http://localhost:8080.
Restart the development server or rebuild the production bundle after changing it.

When the public URL is unset, the frontend uses the same-origin Next.js rewrite.
Its destination is configured by `BACKEND_URL` at build time, as used by Docker.

The service sends `GET /api/check/{hexcode}`. A successful response has
`Success: true` and a report in `Data`. Validation failures return HTTP 400
with `Success: false` and an error string in `Data`.
