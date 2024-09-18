package main

import (
  "fmt"
  "strings"
  "os"
  "os/exec"
)

// renderPile returns a string representation of a pile of cards.
func renderPile(pile []Card) string {
  s := ""
  for _, card := range pile {
    s += renderCard(card) + " "
  }
  return strings.TrimSpace(s)
}

// Render displays the current game state to the console.
func (g *GameState) Render() {
  cmd := exec.Command("cmd", "/c", "cls")
  cmd.Stdout = os.Stdout
  cmd.Run()
  fmt.Println("\nFoundations:")
  for i, foundation := range g.Foundations {
    suit := Suits[i]
    fmt.Printf(" %s: %s\n", suit, renderPile(foundation))
  }

  fmt.Println("\nStock:")
  if len(g.Stock) > 0 {
    fmt.Printf(" [%d cards]\n", len(g.Stock))
  } else {
    fmt.Println(" Empty")
  }

  fmt.Println("\nWaste:")
  if len(g.Waste) > 0 {
    // Display up to the top three cards of the waste pile
    start := len(g.Waste) - 3
    if start < 0 {
      start = 0
    }
    topCards := g.Waste[start:]
    fmt.Println(renderPile(topCards))
  } else {
    fmt.Println(" Empty")
  }

  fmt.Println("\nTableaus:")
  // Determine the maximum height of the tableaus
  maxHeight := 0
  for _, tableau := range g.Tableaus {
    if len(tableau) > maxHeight {
      maxHeight = len(tableau)
    }
  }

  // Build the tableau lines
  tableauLines := make([]string, maxHeight)
  for i := 0; i < maxHeight; i++ {
    line := ""
    for _, tableau := range g.Tableaus {
      if i < len(tableau) {
        line += fmt.Sprintf("%s ", renderCard(tableau[i]))
      } else {
        line += "      "
      }
    }
    tableauLines[i] = line
  }

  // Print the tableau labels
  labels := ""
  for i := range g.Tableaus {
    labels += fmt.Sprintf(" T%d   ", i+1)
  }
  fmt.Println(labels)

  // Print the tableau lines
  for _, line := range tableauLines {
    fmt.Println(line)
  }
}
