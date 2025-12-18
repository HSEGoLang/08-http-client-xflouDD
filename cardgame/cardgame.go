//go:build !solution

package cardgame

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const defaultBaseURL = "https://deckofcardsapi.com/api/deck"

// Client представляет клиента для работы с API карточной колоды
type Client struct {
	baseURL string
	client  *http.Client
	output  io.Writer
}

// NewClient создаёт новый клиент с настройками по умолчанию
func NewClient() *Client {
	return &Client{
		baseURL: defaultBaseURL,
		client:  http.DefaultClient,
		output:  nil, // nil означает вывод в stdout через fmt
	}
}

// PlayGame запускает игру "угадай карту до дамы"
// userGuess - количество карт, которое, по мнению пользователя, нужно снять
// Возвращает true, если пользователь угадал, false если нет
// В этом коде есть ОШИБКИ! Найди и исправь их.
func (c *Client) PlayGame(userGuess int) (bool, error) {
	// Создаём и перетасовываем колоду
	resp, err := c.client.Get(c.baseURL + "/new/shuffle/?deck_count=1")
	if err != nil {
		return false, fmt.Errorf("failed to create deck: %w", err)
	}
	defer resp.Body.Close()

	var deckResp DeckResponse
	if err := json.NewDecoder(resp.Body).Decode(&deckResp); err != nil {
		return false, fmt.Errorf("failed to decode deck response: %w", err)
	}

	deckID := deckResp.DeckID
	realCount := 0

	// Вытягиваем карты, пока не найдём даму
	for {
		// ОШИБКА 1: запрашиваем неправильное количество карт
		drawResp, err := c.client.Get(fmt.Sprintf("%s/%s/draw/?count=1", c.baseURL, deckID))
		if err != nil {
			return false, fmt.Errorf("failed to draw card: %w", err)
		}
		defer drawResp.Body.Close()

		var draw DrawResponse
		if err := json.NewDecoder(drawResp.Body).Decode(&draw); err != nil {
			return false, fmt.Errorf("failed to decode draw response: %w", err)
		}

		// ОШИБКА 2: неправильно работаем с массивом cards
		card := draw.Cards[0]
		realCount++

		// ОШИБКА 3: неправильный доступ к полям карты
		c.printf("%s of %s\n", card.Value, card.Suit)

		if card.Value == "QUEEN" {
			break
		}
	}

	// Проверяем результат
	if realCount == userGuess {
		c.printf("Вы угадали!\n")
		return true, nil
	} else {
		c.printf("Вы проиграли! Правильный ответ: %d\n", realCount)
		return false, nil
	}
}

func (c *Client) printf(format string, args ...interface{}) {
	if c.output != nil {
		fmt.Fprintf(c.output, format, args...)
	} else {
		fmt.Printf(format, args...)
	}
}

// PlayGame - вспомогательная функция для обратной совместимости
func PlayGame(userGuess int) (bool, error) {
	return NewClient().PlayGame(userGuess)
}
