package main

import (
  //"encoding/json"
  "fmt"
  "net/http"
  // import the data package which contains the definition for Comment
  "github.com/amari03/qod/internal/data"
  "github.com/amari03/qod/internal/validator"

)

func (a *application)createCommentHandler(w http.ResponseWriter,r *http.Request) { 
// create a struct to hold a comment
// we use struct tags[``] to make the names display in lowercase
var incomingData struct {
	Content  string  `json:"content"`
	Author   string  `json:"author"`
} 

// perform the decoding
err := a.readJSON(w, r, &incomingData)
if err != nil {
	a.badRequestResponse(w, r, err)
	return
}

// Copy into a domain model (what your notes ask for)
comment := &data.Comment{
  Content: incomingData.Content,
  Author:  incomingData.Author,
}
// Initialize a Validator instance
v := validator.New()
// Do the validation
data.ValidateComment(v, comment)
if !v.IsEmpty() {
    a.failedValidationResponse(w, r, v.Errors)  // implemented later
    return
}

// Add the comment to the database table
err = a.commentModel.Insert(comment)
if err != nil {
    a.serverErrorResponse(w, r, err)
    return
}

fmt.Fprintf(w, "%+v\n", incomingData)      // delete this
// Set a Location header. The path to the newly created comment
headers := make(http.Header)
headers.Set("Location", fmt.Sprintf("/v1/comments/%d", comment.ID))

// Send a JSON response with 201 (new resource created) status code
data := envelope{
  "comment": comment,
}
err = a.writeJSON(w, http.StatusCreated, data, headers)
if err != nil {
a.serverErrorResponse(w, r, err)
return
}


}
