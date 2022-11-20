package note

import (
	"febre/database"
	"febre/internals/model"
	redis "febre/pkg/service"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"log"
	"time"
)

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
	"user3": "password3",
}

func GetNotes(c *fiber.Ctx) error {
	//db := database.DB
	var notes []model.Note
	redis.SetToRedis("bbb")
	fmt.Println(redis.GetToRedis())
	// find all notes in the database
	//db.Find(&notes)

	// If no note is present return an error
	if len(notes) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No notes present***", "data": nil})
	}

	// Else return notes
	return c.JSON(fiber.Map{"status": "success", "message": "Notes Found", "data": notes})
}
func CreateNotes(c *fiber.Ctx) error {
	db := database.DB
	note := new(model.Note)

	// Store the body in the note and return error if encountered
	err := c.BodyParser(note)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	// Add a uuid to the note
	note.ID = uuid.New()
	// Create the Note and return error if encountered
	err = db.Create(&note).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create note", "data": err})
	}

	// Return the created note
	return c.JSON(fiber.Map{"status": "success", "message": "Created Note", "data": note})
}
func GetNote(c *fiber.Ctx) error {
	db := database.DB
	var note model.Note

	// Read the param noteId
	id := c.Params("noteId")

	// Find the note with the given Id
	db.Find(&note, "id = ?", id)

	// If no such note present return an error
	if note.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No note present", "data": nil})
	}

	// Return the note with the Id
	return c.JSON(fiber.Map{"status": "success", "message": "Notes Found", "data": note})
}
func UpdateNote(c *fiber.Ctx) error {
	type updateNote struct {
		Title    string `json:"title"`
		SubTitle string `json:"sub_title"`
		Text     string `json:"Text"`
	}
	db := database.DB
	var note model.Note

	// Read the param noteId
	id := c.Params("noteId")

	// Find the note with the given Id
	db.Find(&note, "id = ?", id)

	// If no such note present return an error
	if note.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No note present", "data": nil})
	}

	// Store the body containing the updated data and return error if encountered
	var updateNoteData updateNote
	err := c.BodyParser(&updateNoteData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	// Edit the note
	note.Title = updateNoteData.Title
	note.SubTitle = updateNoteData.SubTitle
	note.Text = updateNoteData.Text

	// Save the Changes
	db.Save(&note)

	// Return the updated note
	return c.JSON(fiber.Map{"status": "success", "message": "Notes Found", "data": note})
}
func DeleteNote(c *fiber.Ctx) error {
	db := database.DB
	var note model.Note

	// Read the param noteId
	id := c.Params("noteId")

	// Find the note with the given Id
	db.Find(&note, "id = ?", id)

	// If no such note present return an error
	if note.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No note present", "data": nil})
	}

	// Delete the note and return error if encountered
	err := db.Delete(&note, "id = ?", id).Error

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Failed to delete note", "data": nil})
	}

	// Return success message
	return c.JSON(fiber.Map{"status": "success", "message": "Deleted Note"})
}
func WebSocketUpgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}
func WebSocket(c *websocket.Conn) {
	log.Println(c.Locals("allowed"))  // true
	log.Println(c.Params("id"))       // 123
	log.Println(c.Query("v"))         // 1.0
	log.Println(c.Cookies("session")) // ""

	// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
	var (
		mt  int
		msg []byte
		err error
	)
	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", msg)

		if err = c.WriteMessage(mt, msg); err != nil {
			log.Println("write:", err)
			break
		}
	}

}
func Registration(c *fiber.Ctx) error {
	credentials := new(model.CredentialsJsonLess)
	err := c.BodyParser(credentials)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "RPeview your input", "data": err})
	}
	expirationTime := time.Now().Add(time.Minute * 5)
	users[credentials.Username] = credentials.Password
	claims := &model.Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Subject:   "mahdi",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("jwtKey"))
	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "StatusInternalServerError", "data": err})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Created Note", "data": tokenString})
}

func Login(c *fiber.Ctx) error {
	//db := database.DB
	credentials := new(model.Credentials)
	err := c.BodyParser(credentials)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "RPeview your input", "data": err})
	}
	expectedPassword, ok := users[credentials.Username]

	if !ok || expectedPassword != credentials.Password {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "StatusUnauthorized", "data": err})
	} ////////////////
	expirationTime := time.Now().Add(time.Minute * 5)
	users[credentials.Username] = credentials.Password
	claims := &model.Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Subject:   "mahdi",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("jwtKey"))
	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "StatusInternalServerError", "data": err})
	}
	temp := redis.GetToRedis()
	fmt.Println(temp)
	return c.JSON(fiber.Map{"status": "success", "message": temp, "data": tokenString})

	//tokenStr := credentials.JwtToken
	//claims := &model.Claims{}
	//
	//tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
	//	return []byte("jwtKey"), nil
	//})
	//if err != nil {
	//	if err == jwt.ErrSignatureInvalid {
	//
	//		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "StatusUnauthorized", "data": err})
	//	}
	//	return c.Status(500).JSON(fiber.Map{"status": "error", "message": "StatusBadRequest", "data": err})
	//
	//}
	//if !tkn.Valid {
	//	return c.Status(500).JSON(fiber.Map{"status": "error", "message": "StatusUnauthorized", "data": err})
	//}
	////database.SetToRedis(client, "aaa")
	//return c.JSON(fiber.Map{"status": credentials.Username, "message": "Created Note", "data": "hello"})
	//////////////
	//expirationTime := time.Now().Add(time.Minute * 5)

	//claims := &model.Claims{
	//	Username: credentials.Username,
	//	StandardClaims: jwt.StandardClaims{
	//		ExpiresAt: expirationTime.Unix(),
	//		Subject:   "mahdi",
	//	},
	//}
	//
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//tokenString, err := token.SignedString([]byte("jwtKey"))
	//if err != nil {
	//	fmt.Println(err)
	//	return c.Status(500).JSON(fiber.Map{"status": "error", "message": "StatusInternalServerError", "data": err})
	//}
	//db.Find(&credentials1.Password, "SELECT password FROM public.\"Auth\" WHERE username=?", credentials.Username)
	//return c.JSON(fiber.Map{"status": "success", "message": "Created Note", "data": tokenString})
}
