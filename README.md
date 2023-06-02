# Static analysis package

### Using custom golangci plugin
- Documentation reference https://golangci-lint.run/contributing/new-linters/#how-to-add-a-private-linter-to-golangci-lint
- To add a custom private plugin the golangci must be build with `CGO_ENABLED=1` flag
- It is required to have own build golangci binary, because this way we ensure the same used dependencies for building
the our own custom plugin and the golangci. Without it we would not be able to plug-in our custom plugin
- Once we have a ready binary we can built `cmd/plugin` `CGO_ENABLED=1 go build -buildmode=plugin -o wraperrchecker.so cmd/plugin/plugin.go`
- Then we can add this custom built in plugin in the runtime using golangci config file

### Wrap error checker
- Checks if we wrap our errors correctly through our codebase

```
    // invalid
    if err != nil {
        fmt.Errorf("something went wrong")
    }
    
    // valid
    if err != nil {
        fmt.Errorf("something went wrong: %w", err)
    }
```