package catalogbot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type telegramClient struct {
	token      string
	apiBase    string
	fileBase   string
	httpClient *http.Client
}

func newTelegramClient(token string) *telegramClient {
	return &telegramClient{
		token:      token,
		apiBase:    "https://api.telegram.org/bot" + token,
		fileBase:   "https://api.telegram.org/file/bot" + token,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type update struct {
	UpdateID      int            `json:"update_id"`
	Message       *message       `json:"message"`
	CallbackQuery *callbackQuery `json:"callback_query"`
}

type message struct {
	MessageID int           `json:"message_id"`
	From      *telegramUser `json:"from"`
	Chat      chat          `json:"chat"`
	Text      string        `json:"text"`
	Photo     []photoSize   `json:"photo"`
}

type telegramUser struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type chat struct {
	ID int64 `json:"id"`
}

type photoSize struct {
	FileID   string `json:"file_id"`
	FileSize int    `json:"file_size"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

type fileInfo struct {
	FileID   string `json:"file_id"`
	FilePath string `json:"file_path"`
}

type callbackQuery struct {
	ID      string        `json:"id"`
	From    *telegramUser `json:"from"`
	Message *message      `json:"message"`
	Data    string        `json:"data"`
}

type inlineKeyboardMarkup struct {
	InlineKeyboard [][]inlineKeyboardButton `json:"inline_keyboard"`
}

type inlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type replyKeyboardMarkup struct {
	Keyboard       [][]keyboardButton `json:"keyboard"`
	ResizeKeyboard bool               `json:"resize_keyboard"`
}

type keyboardButton struct {
	Text string `json:"text"`
}

func (c *telegramClient) getUpdates(ctx context.Context, offset int) ([]update, error) {
	values := url.Values{}
	values.Set("timeout", "25")
	if offset > 0 {
		values.Set("offset", strconv.Itoa(offset))
	}

	var response struct {
		OK     bool     `json:"ok"`
		Result []update `json:"result"`
	}
	if err := c.get(ctx, "getUpdates", values, &response); err != nil {
		return nil, err
	}
	if !response.OK {
		return nil, fmt.Errorf("telegram getUpdates returned ok=false")
	}
	return response.Result, nil
}

func (c *telegramClient) sendMessage(ctx context.Context, chatID int64, text string) error {
	payload := map[string]any{
		"chat_id":                  chatID,
		"text":                     text,
		"disable_web_page_preview": true,
	}
	var response struct {
		OK bool `json:"ok"`
	}
	if err := c.post(ctx, "sendMessage", payload, &response); err != nil {
		return err
	}
	if !response.OK {
		return fmt.Errorf("telegram sendMessage returned ok=false")
	}
	return nil
}

func (c *telegramClient) sendMessageMarkup(ctx context.Context, chatID int64, text string, replyMarkup any) error {
	payload := map[string]any{
		"chat_id":                  chatID,
		"text":                     text,
		"disable_web_page_preview": true,
		"reply_markup":             replyMarkup,
	}
	var response struct {
		OK bool `json:"ok"`
	}
	if err := c.post(ctx, "sendMessage", payload, &response); err != nil {
		return err
	}
	if !response.OK {
		return fmt.Errorf("telegram sendMessage returned ok=false")
	}
	return nil
}

func (c *telegramClient) answerCallbackQuery(ctx context.Context, callbackQueryID string) error {
	payload := map[string]any{"callback_query_id": callbackQueryID}
	var response struct {
		OK bool `json:"ok"`
	}
	if err := c.post(ctx, "answerCallbackQuery", payload, &response); err != nil {
		return err
	}
	if !response.OK {
		return fmt.Errorf("telegram answerCallbackQuery returned ok=false")
	}
	return nil
}

func (c *telegramClient) getFile(ctx context.Context, fileID string) (fileInfo, error) {
	values := url.Values{}
	values.Set("file_id", fileID)

	var response struct {
		OK     bool     `json:"ok"`
		Result fileInfo `json:"result"`
	}
	if err := c.get(ctx, "getFile", values, &response); err != nil {
		return fileInfo{}, err
	}
	if !response.OK {
		return fileInfo{}, fmt.Errorf("telegram getFile returned ok=false")
	}
	return response.Result, nil
}

func (c *telegramClient) downloadFile(ctx context.Context, filePath string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.fileBase+"/"+filePath, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("telegram file download returned status %d", res.StatusCode)
	}
	return io.ReadAll(res.Body)
}

func (c *telegramClient) get(ctx context.Context, method string, values url.Values, destination any) error {
	endpoint := c.apiBase + "/" + method + "?" + values.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("telegram %s returned status %d", method, res.StatusCode)
	}
	return json.NewDecoder(res.Body).Decode(destination)
}

func (c *telegramClient) post(ctx context.Context, method string, payload any, destination any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiBase+"/"+method, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("telegram %s returned status %d", method, res.StatusCode)
	}
	return json.NewDecoder(res.Body).Decode(destination)
}
