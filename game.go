package main

import (
  "fmt"
  "math/rand"
  "strconv"
  "strings"
  "time"
)

// GameState represents the current state of the game.
type GameState struct {
  Stock       []Card       // The stock pile (draw pile)
  Waste       []Card       // The waste pile (cards drawn from the stock)
  Foundations [4][]Card    // Four foundation piles, one for each suit
  Tableaus    [7][]Card    // Seven tableau piles
}

// createDeck generates a standard 52-card deck.
func createDeck() []Card {
  var deck []Card
  for _, suit := range Suits {
    for _, rank := range Ranks {
      deck = append(deck, Card{Suit: suit, Rank: rank, FaceUp: false})
    }
  }
  return deck
}

// shuffleDeck shuffles the cards in the deck.
func shuffleDeck(deck []Card) {
  source := rand.NewSource(time.Now().UnixNano())
  rand.New(source)
  rand.Shuffle(len(deck), func(i, j int) {
    deck[i], deck[j] = deck[j], deck[i]
  })
}

// dealToTableau deals cards to the tableau piles according to Solitaire rules.
func dealToTableau(deck []Card) ([]Card, [7][]Card) {
  var tableaus [7][]Card
  index := 0
  for i := 0; i < 7; i++ {
    for j := 0; j <= i; j++ {
      card := deck[index]
      index++
      // The last card in the pile is face-up
      if j == i {
        card.FaceUp = true
      }
      tableaus[i] = append(tableaus[i], card)
    }
  }
  // Return the remaining deck and the tableau piles
  return deck[index:], tableaus
}

// initializeGame sets up a new game and returns the initial game state.
func initializeGame() GameState {
  deck := createDeck()
  shuffleDeck(deck)
  remainingDeck, tableaus := dealToTableau(deck)

  game := GameState{
    Stock:    remainingDeck,
    Waste:    []Card{},
    Tableaus: tableaus,
  }

  // Initialize empty foundations
  for i := range game.Foundations {
    game.Foundations[i] = []Card{}
  }

  return game
}

// DrawCards handles drawing three cards from the stock to the waste.
func (g *GameState) DrawCards() {
  if len(g.Stock) == 0 {
    // If the stock is empty, recycle the waste into the stock without reversing
    g.Stock = g.Waste
    g.Waste = []Card{}
    return
  }

  // Draw up to three cards
  drawCount := 3
  if len(g.Stock) < 3 {
    drawCount = len(g.Stock)
  }

  // Move cards from stock to waste
  for i := 0; i < drawCount; i++ {
    card := g.Stock[0]
    g.Stock = g.Stock[1:]
    card.FaceUp = true
    g.Waste = append(g.Waste, card)
  }
}

// MoveCard handles moving a card from one pile to another.
func (g *GameState) MoveCard(source string, destination string) error {
  // Parse source
  srcPile, err := g.getPile(source)
  if err != nil {
    return err
  }

  // Parse destination
  destPile, err := g.getPile(destination)
  if err != nil {
    return err
  }

  // Get the cards to move
  cardsToMove, err := g.getCardsToMove(source, srcPile)
  if err != nil {
    return err
  }

  // Attempt to move decreasing number of cards until valid move is found
  for i := 0; i < len(cardsToMove); i++ {
    subCards := cardsToMove[i:]
    if g.isValidMove(subCards, destPile, destination) {
      // Perform the move
      *destPile = append(*destPile, subCards...)
      *srcPile = (*srcPile)[:len(*srcPile)-len(subCards)]

      // Flip the next card in the source tableau if needed
      if strings.HasPrefix(source, "T") && len(*srcPile) > 0 {
        lastCard := &(*srcPile)[len(*srcPile)-1]
        if !lastCard.FaceUp {
          lastCard.FaceUp = true
        }
      }
      return nil
    }
  }

  return fmt.Errorf("invalid move")
}

// getPile returns a pointer to the pile specified by the identifier.
func (g *GameState) getPile(identifier string) (*[]Card, error) {
  switch {
  case identifier == "W":
    return &g.Waste, nil
  case strings.HasPrefix(identifier, "T"):
    index, err := strconv.Atoi(identifier[1:])
    if err != nil || index < 1 || index > 7 {
      return nil, fmt.Errorf("invalid tableau identifier")
    }
    return &g.Tableaus[index-1], nil
  case strings.HasPrefix(identifier, "F"):
    index, err := strconv.Atoi(identifier[1:])
    if err != nil || index < 1 || index > 4 {
      return nil, fmt.Errorf("invalid foundation identifier")
    }
    return &g.Foundations[index-1], nil
  default:
    return nil, fmt.Errorf("unknown pile identifier")
}
}

// getCardsToMove returns the cards to move from the source pile.
func (g *GameState) getCardsToMove(source string, pile *[]Card) ([]Card, error) {
  if len(*pile) == 0 {
    return nil, fmt.Errorf("source pile is empty")
  }

  // For waste or foundations, only the top card can be moved
  if source == "W" || strings.HasPrefix(source, "F") {
    topCard := (*pile)[len(*pile)-1]
    return []Card{topCard}, nil
  }

  // For tableaus, multiple face-up cards can be moved
  if strings.HasPrefix(source, "T") {
    // Find the index where face-up cards start
    var i int
    for i = len(*pile) - 1; i >= 0; i-- {
      if !(*pile)[i].FaceUp {
        break
      }
    }
    // Return all face-up cards starting from faceUpStart
    cardsToMove := (*pile)[i + 1:]
    return cardsToMove, nil
  }

  return nil, fmt.Errorf("invalid source pile")
}

// isValidMove checks whether moving the cards to the destination pile is valid.
func (g *GameState) isValidMove(cards []Card, destPile *[]Card, destIdentifier string) bool {
  if len(cards) == 0 {
    return false
  }
  movingCard := cards[0]

  if strings.HasPrefix(destIdentifier, "F") {
    // Moving to a foundation
    if len(cards) != 1 {
      return false // Can only move one card to a foundation at a time
    }
    // Get the foundation index
    index, err := strconv.Atoi(destIdentifier[1:])
    if err != nil || index < 1 || index > 4 {
      return false
    }
    destSuit := Suits[index - 1]
    if movingCard.Suit != destSuit {
      return false // The card's suit does not match the foundation's suit
    }
    if len(*destPile) == 0 {
      return movingCard.Rank == Ace
    }
    topDestCard := (*destPile)[len(*destPile)-1]
    return rankValue(movingCard.Rank) == rankValue(topDestCard.Rank)+1

  } else if strings.HasPrefix(destIdentifier, "T") {
    // Moving to a tableau
    if len(*destPile) == 0 {
      return movingCard.Rank == King
    }
    topDestCard := (*destPile)[len(*destPile)-1]
    return isOppositeColor(movingCard, topDestCard) && rankValue(movingCard.Rank)+1 == rankValue(topDestCard.Rank)
  } else {
    // Can't move to waste or stock
    return false
  }
}

