package note

import (
	noteHandler "febre/internals/handlers/note"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/websocket/v2"
)

var jwtmidleware = jwtware.New(jwtware.Config{
	SigningKey: []byte("jwtKey"),
})

func SetupNoteRoutes(router fiber.Router) {
	note := router.Group("/note")
	// Create a Note
	note.Post("/", noteHandler.CreateNotes)
	// Read all Notes
	note.Get("/", jwtmidleware, noteHandler.GetNotes)
	// // Read one Note
	note.Get("/:noteId", noteHandler.GetNote)
	// // Update one Note
	note.Put("/:noteId", noteHandler.UpdateNote)
	// // Delete one Note
	note.Delete("/:noteId", noteHandler.DeleteNote)
	//restricted Routes
	// Upgrade for WebSocket
	note.Use("/ws", noteHandler.WebSocketUpgrade)
	// Create WebSocket Endpoint
	note.Get("/ws/:id", websocket.New(noteHandler.WebSocket))
	////JWT Registeration
	note.Post("/register", noteHandler.Registration)
	//Login with JWT
	note.Post("/login", noteHandler.Login)

}
