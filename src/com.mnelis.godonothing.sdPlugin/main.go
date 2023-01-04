// Needs to be the "main" package since StreamDeck will call our executable.
package main

import (
	"encoding/json"
	"flag"
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

		messageStr := string(message)

		log.Println("Got message from StreamDeck: ", messageStr)

		// Parse our JSON message into a map of strings.
		// https://stackoverflow.com/questions/28859941/whats-the-golang-equivalent-of-converting-any-json-to-standard-dict-in-python
		var eventData map[string]interface{}
		err = json.Unmarshal(message, &eventData)
		if err != nil {
			log.Println("Unable to parse event JSON: ", err)
			return
		}

		event := eventData["event"]
		// action := eventData["action"]
		context := eventData["context"]

		if event == "deviceDidConnect" {
			deviceId := eventData["device"]
			log.Println("Reigstered plugin with StreamDeck device:", deviceId)
		} else if event == "keyDown" {
			jsonPayload, ok := eventData["payload"].(map[string]interface{})
			if !ok {
				log.Println("Unable to read payload from JSON event data")
				return
			} else {
				log.Println("Got KeyDown event")
			}
			settings := jsonPayload["settings"]
			coordinates := jsonPayload["coordinates"]
			userDesiredState := jsonPayload["userDesiredState"]
			onKeyDown(context, settings, coordinates, userDesiredState)
		} else if event == "keyUp" {
			jsonPayload, ok := eventData["payload"].(map[string]interface{})
			if !ok {
				log.Println("Unable to read payload from JSON event data")
				return
			} else {
				log.Println("Got KeyUp event")
			}
			settings := jsonPayload["settings"]
			coordinates := jsonPayload["coordinates"]
			userDesiredState := jsonPayload["userDesiredState"]
			onKeyUp(context, settings, coordinates, userDesiredState)
		} else if event == "willAppear" {
			jsonPayload, ok := eventData["payload"].(map[string]interface{})
			if !ok {
				log.Println("Unable to read payload from JSON event data")
				return
			} else {
				log.Println("Got WillAppear event")
			}
			settings := jsonPayload["settings"]
			coordinates := jsonPayload["coordinates"]
			onWillAppear(context, settings, coordinates)
		}
	}
}

func main() {
	logFile := setupLogging()
	defer logFile.Close()

	log.Println("Starting Golang DoNothing StreamDeck Plugin")
	var stBuilder strings.Builder
	for _, arg := range os.Args[1:] {
		stBuilder.WriteString(fmt.Sprintf("%s ", arg))
	}
	log.Println("Command line arguments: ", stBuilder.String())

	var inPort string
	var inPluginUUID string
	var inRegisterEvent string
	var inInfo string
	flag.StringVar(&inPort, "port", "", "-port <port number>")
	flag.StringVar(&inPluginUUID, "pluginUUID", "", "-pluginUUID <UUID of the plugin>")
	flag.StringVar(&inRegisterEvent, "registerEvent", "", "-registerEvent <name of the event to register your plugin>")
	flag.StringVar(&inInfo, "info", "", "-info <info about the StreamDeck environment>")
	flag.Parse() // after declaring flags we need to call it

	log.Println("Port:", inPort)
	log.Println("Plugin UUID:", inPluginUUID)
	log.Println("Register Event:", inRegisterEvent)
	log.Println("Info:", inInfo)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	conn, _, connectErr := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://127.0.0.1:%s", inPort), nil)
	if connectErr != nil {
		log.Fatal("Unable to connect to StreamDeck:", connectErr)
	}
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
	registerErr := conn.WriteMessage(websocket.TextMessage, []byte(connect_message))
	if registerErr != nil {
		log.Fatal("Unable to send Register message to StreamDeck:", registerErr)
	}

	log.Println("Plugin registered")

	// Keep running until we are interrupted
	<-interrupt

	log.Println("Interrupt signal caught, cleaning up.")

	// Cleanly close the connection by sending a close message and then
	// waiting (with timeout) for the server to close the connection.
	disconnectErr := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if disconnectErr != nil {
		log.Println("Unable to cleanly disconnect from StreamDeck:", disconnectErr)
		return
	}
}
