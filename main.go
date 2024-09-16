package main

import (
    "bufio"
    "fmt"
    "math/rand"
    "os"
    "strconv"
    "strings"
    "time"
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

// Suits slice for iteration
var Suits = []Suit{Hearts, Diamonds, Clubs, Spades}

// Foundation suits mapping
var foundationSuits = []Suit{Hearts, Diamonds, Clubs, Spades}

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

// Ranks slice for iteration
var Ranks = []Rank{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}

// Card represents a playing card.
type Card struct {
    Suit   Suit
    Rank   Rank
    FaceUp bool
}

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
    rand.Seed(time.Now().UnixNano())
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

// renderCard returns a string representation of a card.
func renderCard(c Card) string {
    if c.FaceUp {
        return fmt.Sprintf("[%2s%s]", c.Rank, c.Suit)
    }
    return "[XX]"
}

// renderPile returns a string representation of a pile of cards.
func renderPile(pile []Card) string {
    s := ""
    for _, card := range pile {
        s += renderCard(card) + " "
    }
    return strings.TrimSpace(s)
}

// renderTableau returns a string representation of a tableau pile.
func renderTableau(pile []Card) string {
    s := ""
    for _, card := range pile {
        s += renderCard(card) + " "
    }
    return strings.TrimSpace(s)
}

// Render displays the current game state to the console.
func (g *GameState) Render() {
    fmt.Println("\nFoundations:")
    for i, foundation := range g.Foundations {
        suit := getFoundationSuit(i)
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
    for i, tableau := range g.Tableaus {
        fmt.Printf(" %d: %s\n", i+1, renderTableau(tableau))
    }
}

// readInput reads a line of input from the user.
func readInput() string {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("\nEnter command: ")
    input, _ := reader.ReadString('\n')
    return strings.TrimSpace(input)
}

// DrawCards handles drawing three cards from the stock to the waste.
func (g *GameState) DrawCards() {
    if len(g.Stock) == 0 {
        // If the stock is empty, recycle the waste into the stock without reversing
        g.Stock = g.Waste
        g.Waste = []Card{}
        fmt.Println("Recycled waste pile back into stock.")
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
    srcPile, _, err := g.getPile(source)
    if err != nil {
        return err
    }

    // Parse destination
    destPile, _, err := g.getPile(destination)
    if err != nil {
        return err
    }

    // Get the card(s) to move
    cardsToMove, err := g.getCardsToMove(source, srcPile)
    if err != nil {
        return err
    }

    // Validate the move
    if !g.isValidMove(cardsToMove, destPile, destination) {
        return fmt.Errorf("invalid move")
    }

    // Perform the move
    *destPile = append(*destPile, cardsToMove...)
    *srcPile = (*srcPile)[:len(*srcPile)-len(cardsToMove)]

    // Flip the next card in the source tableau if needed
    if strings.HasPrefix(source, "T") && len(*srcPile) > 0 {
        lastCard := &(*srcPile)[len(*srcPile)-1]
        if !lastCard.FaceUp {
            lastCard.FaceUp = true
        }
    }

    return nil
}

// getPile returns a pointer to the pile specified by the identifier.
func (g *GameState) getPile(identifier string) (*[]Card, int, error) {
    switch {
    case identifier == "W":
        return &g.Waste, -1, nil
    case strings.HasPrefix(identifier, "T"):
        index, err := strconv.Atoi(identifier[1:])
        if err != nil || index < 1 || index > 7 {
            return nil, -1, fmt.Errorf("invalid tableau identifier")
        }
        return &g.Tableaus[index-1], index - 1, nil
    case strings.HasPrefix(identifier, "F"):
        index, err := strconv.Atoi(identifier[1:])
        if err != nil || index < 1 || index > 4 {
            return nil, -1, fmt.Errorf("invalid foundation identifier")
        }
        return &g.Foundations[index-1], index - 1, nil
    default:
        return nil, -1, fmt.Errorf("unknown pile identifier")
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
        faceUpStart := i + 1
        if faceUpStart >= len(*pile) {
            return nil, fmt.Errorf("no face-up cards to move")
        }

        // Return all face-up cards starting from faceUpStart
        cardsToMove := (*pile)[faceUpStart:]
        return cardsToMove, nil
    }

    return nil, fmt.Errorf("invalid source pile")
}

// isValidMove checks whether moving the cards to the destination pile is valid.
func (g *GameState) isValidMove(cards []Card, destPile *[]Card, destIdentifier string) bool {
    if len(cards) == 0 {
        return false
    }
    // TODO: moving a pile of multiple cards should be possible.
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
        destSuit := getFoundationSuit(index - 1)
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

func isOppositeColor(a, b Card) bool {
    redSuits := map[Suit]bool{Hearts: true, Diamonds: true}
    return redSuits[a.Suit] != redSuits[b.Suit]
}

func getFoundationSuit(index int) Suit {
    return foundationSuits[index]
}

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
        case "exit", "quit":
            fmt.Println("Thanks for playing!")
            return
        case "draw":
            game.DrawCards()
            game.Render()
        case "move":
            if len(tokens) != 3 {
                fmt.Println("Invalid move command. Usage: move SOURCE DESTINATION")
                continue
            }
            err := game.MoveCard(tokens[1], tokens[2])
            if err != nil {
                fmt.Println("Error:", err)
            }
            game.Render()
        default:
            fmt.Println("Unknown command. Available commands: draw, move, exit")
        }
    }
}

