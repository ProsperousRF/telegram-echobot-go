package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main(){
    err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
    // https://api.telegram.org/bot<token>/METHOD_NAME
    botToken := os.Getenv("BOT_API")

    botApi := "https://api.telegram.org/bot"
    botUrl := botApi + botToken
    offset := 0


    for ;; {
        updates, err := getUpdates(botUrl, offset)
        if err != nil{
            log.Println("Some error found", err.Error())
        }

        for _, update := range updates{
            err = respond(botUrl, update)
            if err != nil{
                log.Println("Some error found", err.Error())
            }
            offset = update.UpdateId + 1
        }
        fmt.Println(updates)
    }
}

func getUpdates(botUrl string, offset int) ([]Update, error) {
    resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(
        offset))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil{
        return nil, err
    }
    var restResponse RestResponse
    err = json.Unmarshal(body, &restResponse)
    if err != nil{
        return nil, err
    }
    return restResponse.Result, nil
}

func respond(botUrl string, update Update) (error){
	var botMessage BotMessage
    botMessage.ChatId = update.Message.Chat.ChatId
    botMessage.Text = update.Message.Text
    buf, err := json.Marshal(botMessage)
    if err != nil{
        return err
    }
    _, err = http.Post(botUrl + "/sendMessage", "application/json",
        bytes.NewBuffer(buf))
    if err != nil{
        return err
    }
    return nil
}