# txShield frontend

From the repository root, run `make run` and open http://localhost:3000.
Install the frontend dependencies with `bun install` in this directory first if needed.

The browser posts calldata to `/api/check-code`. Next.js forwards the request to
the Go backend at `http://localhost:8080`. Set `BACKEND_URL` in the frontend
process environment if the Go service uses a different origin.

The API returns both `result` and `error`. Successful checks return a report
and an empty error string. Failed checks return a null result and an error message.
