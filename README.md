# winboard
golang package to interact with the live windows clipboard (only tested on windows 11)

# usage
install:
```bash
go get github.com/horobimasu/winboard
```

funcs:
```go
winboard.ClearClipboardHistory() (bool, error) // this will not clear pinned items
winboard.GetClipboardHistory() ([]string, error)
winboard.GetPinnedItems() ([]string, error)
```
