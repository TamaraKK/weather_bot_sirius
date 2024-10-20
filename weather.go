package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/briandowns/openweathermap"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var userTimers = make(map[int64]*time.Timer)

func isValidCity(city string) (bool, error) {
	city = strings.ReplaceAll(city, " ", "+")
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, os.Getenv("API_TOKEN"))
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return false, err
	}

	if data["cod"] == "404" {
		return false, nil
	}
	return true, nil
}

func sendWeather(bot *tgbotapi.BotAPI, chatID int64, frequency string) {
	city, err := getCityByChatID(chatID)
	if err != nil {
		log.Println("Error getting city from database:", err)
		return
	}

	if timer, exists := userTimers[chatID]; exists {
		timer.Stop()
	}

	owm, _ := openweathermap.NewCurrent("c", "ru", os.Getenv("API_TOKEN"))
	err = owm.CurrentByName(city)
	if err != nil {
		log.Printf("Ошибка при получении погоды для города '%s': %s\n", city, err.Error())
		return
	}

	message := fmt.Sprintf("🌤 Прогноз погоды для города %s:\n\n"+
		"📝 Описание: %s\n"+
		"🌡 Температура: %.1f°C\n"+
		"🌬 Ощущается как: %.1f°C\n",
		city, owm.Weather[0].Description, owm.Main.Temp, owm.Main.FeelsLike)

	msg := tgbotapi.NewMessage(chatID, message)
	_, err = bot.Send(msg)
	if err != nil {
		log.Println("Error sending message:", err)
	}

	var duration time.Duration
	switch frequency {
	case "1_minute":
		duration = 1 * time.Minute
	case "1_hour":
		duration = 1 * time.Hour
	case "6_hours":
		duration = 6 * time.Hour
	}

	timer := time.AfterFunc(duration, func() {
		sendWeather(bot, chatID, frequency)
	})
	userTimers[chatID] = timer
}
