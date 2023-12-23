package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// MongoDB configuration
var mongoURI = "mongodb+srv://abhiyanampally:aBHIRAMY@cluster0.erhdu89.mongodb.net/"
var dbName = "testdb"
var userCollectionName = "notebooks"

type Query struct {
	Question string `json:"question " bson:"question "`
	Response string `json:"response " bson:"response "`
}

type NewNote struct {
	Title string `json:"title "`
	Text  string `json:"text "`
	Data  Data   `json:"data "`
}

type Data struct {
	ID      primitive.ObjectID `json:"id " bson:"_id "`
	Info    string             `json:"info " bson:"info "`
	Queries []Query            `json:"queries " bson:"queries "`
}

type Note struct {
	ID    primitive.ObjectID `json:"id " bson:"_id "`
	Title string             `json:"title " bson:"title "`
	Text  string             `json:"text " bson:"text "`
	Data  Data               `json:"data " bson:"data "`
}

type User_Data struct {
	Username  string `json:"username" bson:"username"`
	Password  string `json:"password" bson:"password"`
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
}

type Notebook struct {
	ID         primitive.ObjectID `json:"id " bson:"_id "`
	User_Data  User_Data          `json:"user_data " bson:"user_data "`
	Title      string             `json:"title " bson:"title "`
	Notes      []Note             `json:"notes " bson:"notes "`
	Created    time.Time          `json:"created " bson:"created "`
	LastAccess time.Time          `json:"lastAccess " bson:"lastAccess "`
	Data       Data               `json:"data " bson:"data "`
}

var client *mongo.Client

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user User_Data
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Check if the user with the given email already exists
	collection := client.Database(dbName).Collection(userCollectionName)
	existingUser := User_Data{}
	err = collection.FindOne(context.Background(), bson.M{"username": user.Username}).Decode(&existingUser)
	if err == nil {
		respondWithError(w, http.StatusConflict, "Account already exists with the provided username")
		return
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	user.Password = string(hashedPassword)

	// Insert the user into the database
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error inserting user into database")
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "User created successfully"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Find the user by username
	collection := client.Database(dbName).Collection(userCollectionName)
	var user User_Data
	err = collection.FindOne(context.Background(), bson.M{"username": credentials.Username}).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Compare the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Respond with success and user details
	response := map[string]interface{}{
		"status": "success",
		"user": map[string]interface{}{
			"username":  user.Username,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			// Add other user details as needed
		},
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func getAllNotebooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var notebooks []Notebook
	collection := client.Database("testdb").Collection("notebooks")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var notebook Notebook
		if err := cursor.Decode(&notebook); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		notebooks = append(notebooks, notebook)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(notebooks)
}

//get all the notebook data

func getNotebookByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var notebook Notebook
	collection := client.Database("testdb").Collection("notebooks")

	err = collection.FindOneAndUpdate(
		context.Background(),
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"lastAccess": time.Now()}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&notebook)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Notebook not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(notebook)
}

//get notebook by id ,returns all fields

func getNotebookByTitle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	title := params["title"]

	var notebook Notebook
	collection := client.Database("testdb").Collection("notebooks")

	err := collection.FindOne(context.Background(), bson.M{"title": title}).Decode(&notebook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	result := struct {
		ID         primitive.ObjectID `json:"id " bson:"_id "`
		Title      string             `json:"title " bson:"title "`
		LastAccess time.Time          `json:"lastAccess " bson:"lastAccess "`
	}{ID: notebook.ID, Title: notebook.Title, LastAccess: notebook.LastAccess}

	json.NewEncoder(w).Encode(result)
}

//get notebook by Title ,returns ID,Title,Last Access

func patchUpdateNoteData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	notebookID := params["notebookID"]
	noteID := params["noteID"]

	notebookObjectID, err := primitive.ObjectIDFromHex(notebookID)
	if err != nil {
		http.Error(w, "Invalid Notebook ID format", http.StatusBadRequest)
		return
	}

	noteObjectID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		http.Error(w, "Invalid Note ID format", http.StatusBadRequest)
		return
	}

	var data Data
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection := client.Database("testdb").Collection("notebooks")
	update := bson.M{"$set": bson.M{"notes.$.data": data, "lastAccess": time.Now()}}
	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": notebookObjectID, "notes._id": noteObjectID}, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
		http.Error(w, "Note not found or no changes applied", http.StatusNotFound)
		return
	}

	// Retrieve the updated note data
	var updatedData Data
	err = collection.FindOne(context.Background(), bson.M{"_id": notebookObjectID, "notes._id": noteObjectID}).Decode(&updatedData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedData)
}

// changes note data
func postNotebook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var notebookReq struct {
		Title string `json:"title "`
		Notes []Note `json:"notes "`
	}

	err := json.NewDecoder(r.Body).Decode(&notebookReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	notebook := Notebook{
		Title:      notebookReq.Title,
		Notes:      notebookReq.Notes,
		Created:    time.Now(),
		LastAccess: time.Now(),
		Data:       Data{}, // You can initialize Data as needed
	}

	collection := client.Database("testdb").Collection("notebooks")
	_, err = collection.InsertOne(context.Background(), notebook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(notebook)
}

// creates a new notebook
func getLastAccessDate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	notebookID := params["notebookID"]

	objectID, err := primitive.ObjectIDFromHex(notebookID)
	if err != nil {
		http.Error(w, "Invalid Notebook ID format", http.StatusBadRequest)
		return
	}

	var notebook Notebook
	collection := client.Database("testdb").Collection("notebooks")

	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&notebook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := struct {
		LastAccess time.Time `json:"lastAccess " bson:"lastAccess "`
	}{LastAccess: notebook.LastAccess}

	json.NewEncoder(w).Encode(response)
}

func updateLastAccess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	collection := client.Database("testdb").Collection("notebooks")
	update := bson.M{"$set": bson.M{"lastAccess": time.Now()}}
	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
		http.Error(w, "Notebook not found or no changes applied", http.StatusNotFound)
		return
	}

	// Retrieve the updated last access time
	var updatedNotebook Notebook
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&updatedNotebook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		LastAccess time.Time `json:"lastAccess " bson:"lastAccess "`
	}{LastAccess: updatedNotebook.LastAccess}

	json.NewEncoder(w).Encode(response)
}

//updates the last access time of a notebook

func patchUpdateNotebookData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var dataPatch struct {
		Title      *string    `json:"title "`
		Notes      []Note     `json:"notes "`
		Created    *time.Time `json:"created "`
		LastAccess *time.Time `json:"lastAccess "`
		Data       struct {
			Info    string  `json:"info "`
			Queries []Query `json:"queries "`
			Notes   []Note  `json:"notes "`
		} `json:"data "`
	}

	err = json.NewDecoder(r.Body).Decode(&dataPatch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updateFields := bson.M{}
	if dataPatch.Title != nil {
		updateFields["title"] = *dataPatch.Title
	}
	if len(dataPatch.Notes) > 0 {
		updateFields["notes"] = dataPatch.Notes
	}
	if dataPatch.Created != nil {
		updateFields["created"] = *dataPatch.Created
	}
	if dataPatch.LastAccess != nil {
		updateFields["lastAccess"] = *dataPatch.LastAccess
	}
	if dataPatch.Data.Info != "" || len(dataPatch.Data.Queries) > 0 || len(dataPatch.Data.Notes) > 0 {
		updateFields["data"] = dataPatch.Data
	}

	collection := client.Database("testdb").Collection("notebooks")
	update := bson.M{"$set": updateFields}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve the updated data
	var updatedNotebook Notebook
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&updatedNotebook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedNotebook)
}

//updates the data of a notebook

func patchUpdateNoteTitle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	notebookID := params["notebookID"]
	noteID := params["noteID"]

	notebookObjectID, err := primitive.ObjectIDFromHex(notebookID)
	if err != nil {
		http.Error(w, "Invalid Notebook ID format", http.StatusBadRequest)
		return
	}

	noteObjectID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		http.Error(w, "Invalid Note ID format", http.StatusBadRequest)
		return
	}

	var noteTitle struct {
		Title string `json:"title "`
	}

	err = json.NewDecoder(r.Body).Decode(&noteTitle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection := client.Database("testdb").Collection("notebooks")
	update := bson.M{
		"$set": bson.M{
			"notes.$.title":      noteTitle.Title,
			"notes.$.lastAccess": time.Now(),
		},
	}
	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": notebookObjectID, "notes._id": noteObjectID}, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
		http.Error(w, "Note not found or no changes applied", http.StatusNotFound)
		return
	}

	// Retrieve the updated note data
	var updatedNote Note
	err = collection.FindOne(context.Background(), bson.M{"_id": notebookObjectID, "notes._id": noteObjectID}).Decode(&updatedNote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedNote)
}

//updates the note title and return the changed value

func removeNotebookByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	collection := client.Database("testdb").Collection("notebooks")
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func postNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	notebookID := params["notebookID"]

	objectID, err := primitive.ObjectIDFromHex(notebookID)
	if err != nil {
		http.Error(w, "Invalid Notebook ID format", http.StatusBadRequest)
		return
	}

	var newNote NewNote
	err = json.NewDecoder(r.Body).Decode(&newNote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a Note from the NewNote structure
	note := Note{
		Title: newNote.Title,
		Text:  newNote.Text,
		Data:  newNote.Data,
	}

	collection := client.Database("testdb").Collection("notebooks")
	update := bson.M{"$push": bson.M{"notes": note}, "$set": bson.M{"lastAccess": time.Now()}}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(note)
}
func removeNoteByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	notebookID := params["notebookID"]
	noteID := params["noteID"]

	notebookObjectID, err := primitive.ObjectIDFromHex(notebookID)
	if err != nil {
		http.Error(w, "Invalid Notebook ID format", http.StatusBadRequest)
		return
	}

	noteObjectID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		http.Error(w, "Invalid Note ID format", http.StatusBadRequest)
		return
	}

	collection := client.Database("testdb").Collection("notebooks")
	update := bson.M{"$pull": bson.M{"notes": bson.M{"_id": noteObjectID}}, "$set": bson.M{"lastAccess": time.Now()}}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": notebookObjectID}, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getDataByNotebookID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	notebookID := params["notebookID"]

	objectID, err := primitive.ObjectIDFromHex(notebookID)
	if err != nil {
		http.Error(w, "Invalid Notebook ID format", http.StatusBadRequest)
		return
	}

	var notebook Notebook
	collection := client.Database("testdb").Collection("notebooks")
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&notebook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var allData []Data
	for _, note := range notebook.Notes {
		allData = append(allData, note.Data)
	}

	json.NewEncoder(w).Encode(allData)
}

func getDataByNoteAndID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	notebookID := params["notebookID"]
	noteID := params["noteID"]

	notebookObjectID, err := primitive.ObjectIDFromHex(notebookID)
	if err != nil {
		http.Error(w, "Invalid Notebook ID format", http.StatusBadRequest)
		return
	}

	noteObjectID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		http.Error(w, "Invalid Note ID format", http.StatusBadRequest)
		return
	}

	var notebook Notebook
	collection := client.Database("testdb").Collection("notebooks")
	err = collection.FindOne(context.Background(), bson.M{"_id": notebookObjectID}).Decode(&notebook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	for _, note := range notebook.Notes {
		if note.ID == noteObjectID {
			json.NewEncoder(w).Encode(note.Data)
			return
		}
	}

	http.Error(w, "Note not found", http.StatusNotFound)
}

func postData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	notebookID := params["notebookID"]
	noteID := params["noteID"]

	notebookObjectID, err := primitive.ObjectIDFromHex(notebookID)
	if err != nil {
		http.Error(w, "Invalid Notebook ID format", http.StatusBadRequest)
		return
	}

	noteObjectID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		http.Error(w, "Invalid Note ID format", http.StatusBadRequest)
		return
	}

	var data Data
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection := client.Database("testdb").Collection("notebooks")
	update := bson.M{"$set": bson.M{"notes.$.data": data, "lastAccess": time.Now()}}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": notebookObjectID, "notes._id": noteObjectID}, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func removeAllNotebooks(w http.ResponseWriter, r *http.Request) {
	collection := client.Database("testdb").Collection("notebooks")
	_, err := collection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// getAllNotesByNotebookID gets all notes of a notebook by ID
func getAllNotesByNotebookID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	notebookID := params["notebookID"]

	objectID, err := primitive.ObjectIDFromHex(notebookID)
	if err != nil {
		http.Error(w, "Invalid Notebook ID format", http.StatusBadRequest)
		return
	}

	var notebook Notebook
	collection := client.Database("testdb").Collection("notebooks")
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&notebook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(notebook.Notes)
}
func removeAllData(w http.ResponseWriter, r *http.Request) {
	collection := client.Database("testdb").Collection("notebooks")
	update := bson.M{"$set": bson.M{"notes.$[].data": Data{}}}
	_, err := collection.UpdateMany(context.Background(), bson.M{}, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func removeDataByNotebookID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	notebookID := params["notebookID"]

	notebookObjectID, err := primitive.ObjectIDFromHex(notebookID)
	if err != nil {
		http.Error(w, "Invalid Notebook ID format", http.StatusBadRequest)
		return
	}

	collection := client.Database("testdb").Collection("notebooks")
	update := bson.M{"$set": bson.M{"notes.$[].data": Data{}}}
	_, err = collection.UpdateMany(context.Background(), bson.M{"_id": notebookObjectID}, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func removeDataByNoteAndID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	notebookID := params["notebookID"]
	noteID := params["noteID"]

	notebookObjectID, err := primitive.ObjectIDFromHex(notebookID)
	if err != nil {
		http.Error(w, "Invalid Notebook ID format", http.StatusBadRequest)
		return
	}

	noteObjectID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		http.Error(w, "Invalid Note ID format", http.StatusBadRequest)
		return
	}

	collection := client.Database("testdb").Collection("notebooks")
	update := bson.M{"$set": bson.M{"notes.$.data": Data{}}}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": notebookObjectID, "notes._id": noteObjectID}, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getNoteByTitle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	notebookID := params["notebookID"]
	noteTitle := params["noteTitle"]

	notebookObjectID, err := primitive.ObjectIDFromHex(notebookID)
	if err != nil {
		http.Error(w, "Invalid Notebook ID format", http.StatusBadRequest)
		return
	}

	var noteData struct {
		ID    primitive.ObjectID `json:"id " bson:"_id "`
		Title string             `json:"title " bson:"title "`
		Text  string             `json:"text " bson:"text "`
		Data  Data               `json:"data " bson:"data "`
	}

	collection := client.Database("testdb").Collection("notebooks")

	// Define the projection
	projection := bson.M{
		"_id":         1,
		"notes.$":     1, // Include the entire notes array
		"notes.text":  1, // Include the text field from notes
		"notes.title": 1, // Include the title field from notes
		"notes.data":  1, // Include the data field from notes
		"notes.id":    1, // Include the id field from notes
	}

	err = collection.FindOne(
		context.Background(),
		bson.M{"_id": notebookObjectID, "notes.title": noteTitle},
		options.FindOne().SetProjection(projection),
	).Decode(&noteData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(noteData)
}

// Fetch Note by ID to get Id ,title ,text
func patchUpdateNoteText(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	notebookID := params["notebookID"]
	noteTitle := params["noteTitle"]

	notebookObjectID, err := primitive.ObjectIDFromHex(notebookID)
	if err != nil {
		http.Error(w, "Invalid Notebook ID format", http.StatusBadRequest)
		return
	}

	var newText struct {
		Text string `json:"text "`
	}

	err = json.NewDecoder(r.Body).Decode(&newText)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection := client.Database("testdb").Collection("notebooks")
	update := bson.M{
		"$set": bson.M{
			"notes.$.text":       newText.Text,
			"notes.$.lastAccess": time.Now(),
		},
	}

	result, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": notebookObjectID, "notes.title": noteTitle},
		update,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
		http.Error(w, "Note not found or no changes applied", http.StatusNotFound)
		return
	}

	// Retrieve the updated note data
	var updatedNote Note
	err = collection.FindOne(
		context.Background(),
		bson.M{"_id": notebookObjectID, "notes.title": noteTitle},
	).Decode(&updatedNote)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedNote)
}

//Change Note Text and return the changed text

func getAllDataByNoteID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	notebookID := params["notebookID"]
	noteID := params["noteID"]

	notebookObjectID, err := primitive.ObjectIDFromHex(notebookID)
	if err != nil {
		http.Error(w, "Invalid Notebook ID format", http.StatusBadRequest)
		return
	}

	noteObjectID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		http.Error(w, "Invalid Note ID format", http.StatusBadRequest)
		return
	}

	var notebook Notebook
	collection := client.Database("testdb").Collection("notebooks")

	err = collection.FindOne(
		context.Background(),
		bson.M{"_id": notebookObjectID, "notes._id": noteObjectID},
	).Decode(&notebook)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(notebook.Notes[0].Data)
}

//Fetch All Data of a Note by ID

func main() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://abhiyanampally:aBHIRAMY@cluster0.erhdu89.mongodb.net/"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.Background())

	router := mux.NewRouter()

	router.HandleFunc("/notebooks", getAllNotebooks).Methods("GET")
	router.HandleFunc("/notebooks/{id}", getNotebookByID).Methods("GET")
	router.HandleFunc("/notebooks/{title}", getNotebookByTitle).Methods("GET")

	router.HandleFunc("/notebooks", postNotebook).Methods("POST")
	router.HandleFunc("/notebooks/{id}", removeNotebookByID).Methods("DELETE")
	router.HandleFunc("/notebook/{notebookID}", postNote).Methods("POST")
	router.HandleFunc("/notebooks/{notebookID}/notes/{noteID}", removeNoteByID).Methods("DELETE")
	router.HandleFunc("/notebooks/{notebookID}/lastaccess", updateLastAccess).Methods("POST")
	router.HandleFunc("/notebooks/{notebookID}/data", getDataByNotebookID).Methods("GET")
	router.HandleFunc("/notebooks/{notebookID}/notes/{noteID}/data", getDataByNoteAndID).Methods("GET")
	router.HandleFunc("/notebooks/{notebookID}/notes/{noteID}/data", postData).Methods("POST")
	router.HandleFunc("/notebooks/{notebookID}/notes/{noteID}/data", removeDataByNoteAndID).Methods("DELETE")
	router.HandleFunc("/notebooks/{notebookID}/data", removeDataByNotebookID).Methods("DELETE")
	router.HandleFunc("/notebooks/removeAll", removeAllNotebooks).Methods("DELETE")
	router.HandleFunc("/notebooks/removeAllData", removeAllData).Methods("DELETE")
	router.HandleFunc("/notebooks/{id}/data", patchUpdateNotebookData).Methods("PATCH")
	router.HandleFunc("/notebooks/{notebookID}/notes/{noteID}/data", patchUpdateNoteData).Methods("PATCH")
	router.HandleFunc("/notebooks/{notebookID}/notes/{noteID}/title", patchUpdateNoteTitle).Methods("PATCH")
	router.HandleFunc("/notebooks/{notebookID}/lastaccessdate", getLastAccessDate).Methods("GET")
	router.HandleFunc("/notebooks/{notebookID}/notes/{noteTitle}", getNoteByTitle).Methods("GET")
	router.HandleFunc("/notebooks/{notebookID}/notes/{noteTitle}/text", patchUpdateNoteText).Methods("PATCH")
	router.HandleFunc("/notebooks/{notebookID}/notes/{noteID}/alldata", getAllDataByNoteID).Methods("GET")
	router.HandleFunc("/notebooks/{notebookID}/notes", getAllNotesByNotebookID).Methods("GET")

	router.HandleFunc("/signup", SignUpHandler).Methods("POST")
	router.HandleFunc("/login", LoginHandler).Methods("POST")

	// CORS setup
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Update with your allowed origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Wrap the router with the CORS handler
	handler := corsHandler.Handler(router)

	log.Fatal(http.ListenAndServe(":8000", handler))
}