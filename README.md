# Didactic Goggles

## Structure

```plaintext
    ./cmd/cli/main.go               # initialization & entrypoint
    └─▶ ./internal/config/*         # configuration library
    └─▶ ./internal/cli-command/*    # parses args, initializes Use Case
        └─▶ ./internal/usecase/*    # (sub)program logic
            └─▶ ./internal/db/*     # persistence
            └─▶ ./internal/domain/* # business rules
```

### Naming Conventions

| Type             | Convention                                                        |
| ---------------- | ----------------------------------------------------------------- |
| Directories      | kebab-case                                                        |
| Individual Files | snake_case                                                        |
| Package Names    | lowercase, no separators (e.g. packagename ) - prefer single word |
| Test Files       | testdata/* (see https://pkg.go.dev/cmd/go#hdr-Test_packages)      |
