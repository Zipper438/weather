package api

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

// Мок HTTP-клиента
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	return m.DoFunc(req)
}

// Тест на успешный ответ
func TestGetWeather_Success(t *testing.T) {
	// Подготовка мока
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"current": {"temp_c": 15.5}}`
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	// Инициализация клиента с моком
	weatherClient := &WeatherClient{
		Client:  mockClient,
		BaseURL: "http://fake-api",
		APIKey:  "test-key",
	}

	// Вызов тестируемой функции
	temp, err := weatherClient.GetWeather("Moscow")
	if err != nil {
		t.Fatalf("Ожидалась успешная обработка, но получена ошибка: %v", err)
	}

	// Проверка результата
	expected := 15.5
	if temp != expected {
		t.Errorf("Ожидалась температура %.1f, получено %.1f", expected, temp)
	}
}

// Тест на ошибку HTTP 404
func TestGetWeather_HTTPError(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(strings.NewReader(`Not Found`)),
			}, nil
		},
	}

	weatherClient := &WeatherClient{Client: mockClient, APIKey: "test-key"}
	_, err := weatherClient.GetWeather("InvalidCity")

	if err == nil {
		t.Fatal("Ожидалась ошибка, но её нет")
	}

	expectedError := "HTTP 404: Not Found"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Ожидалась ошибка '%s', получено: '%v'", expectedError, err)
	}
}

// Тест на некорректный JSON
func TestGetWeather_InvalidJSON(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{invalid-json}`)),
			}, nil
		},
	}

	weatherClient := &WeatherClient{Client: mockClient, APIKey: "test-key"}
	_, err := weatherClient.GetWeather("Moscow")

	if err == nil || !strings.Contains(err.Error(), "парсинга JSON") {
		t.Errorf("Ожидалась ошибка парсинга JSON, получено: %v", err)
	}
}
