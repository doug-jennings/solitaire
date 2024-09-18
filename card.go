package main

import (
  "fmt"
)

// Suit represents the suit of a card.
type Suit string

// Rank represents the rank of a card.
type Rank string

// Suits
const (
  Hearts   Suit = "♥"
  Diamonds Suit = "♦"
  Clubs    Suit = "♣"
  Spades   Suit = "♠"
)

// Ranks
const (
  Ace   Rank = "A"
  Two   Rank = "2"
  Three Rank = "3"
  Four  Rank = "4"
  Five  Rank = "5"
  Six   Rank = "6"
  Seven Rank = "7"
  Eight Rank = "8"
  Nine  Rank = "9"
  Ten   Rank = "10"
  Jack  Rank = "J"
  Queen Rank = "Q"
  King  Rank = "K"
)

// Card represents a playing card.
type Card struct {
  Suit   Suit
  Rank   Rank
  FaceUp bool
}

// Suits slice for iteration
var Suits = []Suit{Hearts, Diamonds, Clubs, Spades}

// Ranks slice for iteration
var Ranks = []Rank{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}

// renderCard returns a string representation of a card.
func renderCard(c Card) string {
  if c.FaceUp {
    cardStr := fmt.Sprintf("[%2s%s]", c.Rank, c.Suit)
    // Check if the card is red
    if c.Suit == Hearts || c.Suit == Diamonds {
      // ANSI code for red text
      return "\033[31m" + cardStr + "\033[0m"
    }
    return cardStr
  }
  return "[░░░]"
}

// rankValue returns the numerical value of a rank.
func rankValue(rank Rank) int {
  ranks := map[Rank]int{
    Ace:   1,
    Two:   2,
    Three: 3,
    Four:  4,
    Five:  5,
    Six:   6,
    Seven: 7,
    Eight: 8,
    Nine:  9,
    Ten:   10,
    Jack:  11,
    Queen: 12,
    King:  13,
  }
  return ranks[rank]
}

// isOppositeColor checks if two cards have opposite colors.
func isOppositeColor(a, b Card) bool {
  redSuits := map[Suit]bool{Hearts: true, Diamonds: true}
  return redSuits[a.Suit] != redSuits[b.Suit]
}
