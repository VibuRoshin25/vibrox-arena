# vibrox-arena

`vibrox-arena` is the server-authoritative game and bot service for the Vibrox
Systems Lab. Its first experiment is an unbeatable tic-tac-toe bot using
minimax. The response exposes search metrics so visitors can inspect the
algorithm rather than only play against it.

## Interfaces

- gRPC on `:8100`
- HTTP health check at `GET :8054/healthz`

See [`proto/arena.proto`](proto/arena.proto) for the game contract.

## Development

```bash
go test ./...
go run ./cmd/server
```

The service targets Go `1.26.6` and will remain pinned until the planned Go
`1.27` migration.

The service is stateless: each request carries the board before the human move.
Arena validates the board and move, applies the move, calculates the bot reply,
and returns the resulting board with minimax metrics. Durable match state can be
introduced later only if a concrete experiment requires it.
