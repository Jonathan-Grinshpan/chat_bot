package main

import (
	"fmt"

	chatBot "go_chatbot/chatbot"
	"log"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "chatbot.html")
}
func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/chatbot", chatBot.ChatbotHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
}

func main() {
	fmt.Println("Starting server on :3000...")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":3000", nil))

}
