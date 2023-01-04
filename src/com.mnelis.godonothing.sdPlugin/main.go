// Needs to be the "main" package since StreamDeck will call our executable.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/gorilla/websocket"
)

func setupLogging() *os.File {
	logFile, err := os.CreateTemp("", "streamdeck-godonothing-log-")
	if err != nil {
		log.Fatalf("error creating temp file: %v", err)
	}

	log.SetOutput(logFile)

	return logFile
}

func onKeyDown(context interface{}, settings interface{}, coordinates interface{}, userDesiredState interface{}) {
	log.Println("KeyDown called")
}

func onKeyUp(context interface{}, settings interface{}, coordinates interface{}, userDesiredState interface{}) {
	log.Println("KeyUp called")
}

func onWillAppear(context interface{}, settings interface{}, coordinates interface{}) {
	log.Println("WillAppear called")
}

func handleMessagesForever(conn *websocket.Conn) {
	for {
		log.Println("Waiting for messages from StreamDeck")

		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Unable to read message from web socket: ", err)
			return
		}

		log.Println("Got message from StreamDeck: ", message)

		// Parse our JSON message into a map of strings.
		// https://stackoverflow.com/questions/28859941/whats-the-golang-equivalent-of-converting-any-json-to-standard-dict-in-python
		var eventDataRoot map[string]interface{}
		err = json.Unmarshal(message, &eventDataRoot)
		if err != nil {
			log.Println("Unable to parse event JSON: ", err)
			return
		}
		eventData, ok := eventDataRoot["data"].(map[string]interface{})
		if !ok {
			log.Println("Unable to parse event JSON root node")
			return
		}

		event := eventData["event"]
		// action := eventData["action"]
		context := eventData["context"]

		if event == "keyDown" {
			jsonPayload := eventData["payload"].(map[string]interface{})
			if !ok {
				log.Println("Unable to read payload from JSON event data")
				return
			}
			settings := jsonPayload["settings"]
			coordinates := jsonPayload["coordinates"]
			userDesiredState := jsonPayload["userDesiredState"]
			onKeyDown(context, settings, coordinates, userDesiredState)
		} else if event == "keyUp" {
			jsonPayload := eventData["payload"].(map[string]interface{})
			if !ok {
				log.Println("Unable to read payload from JSON event data")
				return
			}
			settings := jsonPayload["settings"]
			coordinates := jsonPayload["coordinates"]
			userDesiredState := jsonPayload["userDesiredState"]
			onKeyUp(context, settings, coordinates, userDesiredState)
		} else if event == "willAppear" {
			jsonPayload := eventData["payload"].(map[string]interface{})
			if !ok {
				log.Println("Unable to read payload from JSON event data")
				return
			}
			settings := jsonPayload["settings"]
			coordinates := jsonPayload["coordinates"]
			onWillAppear(context, settings, coordinates)
		}
	}
}

func main() {
	inPort := os.Args[1]
	inPluginUUID := os.Args[2]
	inRegisterEvent := os.Args[3]
	// inInfo := os.Args[4]

	logFile := setupLogging()
	defer logFile.Close()

	log.Println("Starting Golang DoNothing StreamDeck Plugin")
	var stBuilder strings.Builder
	for _, arg := range os.Args[1:] {
		stBuilder.WriteString(fmt.Sprintf("%s, ", arg))
	}
	log.Println("Command line arguments: ", stBuilder.String())

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	conn, _, _ := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://127.0.0.1:%s", inPort), nil)
	defer conn.Close()

	log.Println("Connected to StreamDeck server")

	go handleMessagesForever(conn)

	log.Println("Handlers setup")

	log.Println("Registering")

	// Send register message to StreamDeck "server"
	connect_message_template := `{
		"event": "%s",
		"uuid": "%s"
	}`
	connect_message := fmt.Sprintf(connect_message_template, inRegisterEvent, inPluginUUID)
	conn.WriteMessage(websocket.TextMessage, []byte(connect_message))

	log.Println("Plugin registered")

	// Keep running until we are interrupted
	<-interrupt

	log.Println("Interrupt signal caught, cleaning up.")

	// Cleanly close the connection by sending a close message and then
	// waiting (with timeout) for the server to close the connection.
	err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		return
	}
}
