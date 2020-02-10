# gorocrate
RO-Crate Library in Go


```go
package main

import (
  "encoding/json"
  "fmt"
  "os"
)

func main() {
  roCrate := NewROCrate()
  b, err := json.Marshal(roCrate)
  if err != nil {
    fmt.Println("error:", err)
  }
  os.Stdout.Write(b)
}
```