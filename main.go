package main

import (
  "fmt"
  "strings"
)

// main function to run the game loop.
func main() {
  game := initializeGame()
  game.Render()

  // Game loop
  for {
    input := readInput()
    tokens := strings.Fields(input)

    if len(tokens) == 0 {
      fmt.Println("Please enter a command.")
      continue
    }

    switch tokens[0] {
    case "q":
      fmt.Println("Thanks for playing!")
      return
    case "d":
      game.DrawCards()
      game.Render()
    case "mv":
      if len(tokens) != 3 {
        fmt.Println("Invalid move command. Usage: move SRC DEST")
        continue
      }
      err := game.MoveCard(tokens[1], tokens[2])
      if err != nil {
        fmt.Println("Error:", err)
      }
      game.Render()
    default:
      fmt.Println("Unknown command. Available commands: d, mv , q")
  }
  }
}
