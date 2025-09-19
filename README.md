# Didactic Goggles

## Architecture

```plaintext
    ./cmd/cli/main.go               # initialization & entrypoint
    └─▶ ./internal/config/*         # configuration library
    └─▶ ./internal/cli-command/*    # parses args, initializes Use Case
        └─▶ ./internal/usecase/*    # (sub)program logic
            └─▶ ./internal/db/*     # persistence
            └─▶ ./internal/domain/* # business rules
```
