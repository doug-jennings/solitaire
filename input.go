package main

import (
  "bufio"
  "fmt"
  "os"
  "strings"
)

// readInput reads a line of input from the user.
func readInput() string {
  reader := bufio.NewReader(os.Stdin)
  fmt.Print("\nEnter command: ")
  input, _ := reader.ReadString('\n')
  return strings.TrimSpace(input)
}
