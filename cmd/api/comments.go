package main

import (
  //"encoding/json"
  "fmt"
  "net/http"
  "errors"
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

func (a *application)displayCommentHandler(w http.ResponseWriter, r *http.Request) {
// Get the id from the URL /v1/comments/:id so that we
// can use it to query teh comments table. We will 
// implement the readIDParam() function later
id, err := a.readIDParam(r)
if err != nil {
    a.notFoundResponse(w, r)
    return 
}
// Call Get() to retrieve the comment with the specified id
comment, err := a.commentModel.Get(id)
if err != nil {
    switch {
        case errors.Is(err, data.ErrRecordNotFound):
           a.notFoundResponse(w, r)
        default:
           a.serverErrorResponse(w, r, err)
    }
    return 
}
// display the comment
data := envelope {
  "comment": comment,
}
err = a.writeJSON(w, http.StatusOK, data, nil)
if err != nil {
a.serverErrorResponse(w, r, err)
return 
}

}

func (a *application)updateCommentHandler(w http.ResponseWriter,r *http.Request) {
  // Get the id from the URL
  id, err := a.readIDParam(r)
  if err != nil {
    a.notFoundResponse(w, r)
    return 
  }
  // Call Get() to retrieve the comment with the specified id
comment, err := a.commentModel.Get(id)
if err != nil {
    switch {
        case errors.Is(err, data.ErrRecordNotFound):
           a.notFoundResponse(w, r)
        default:
           a.serverErrorResponse(w, r, err)
    }
    return 
}
// Use our temporary incomingData struct to hold the data
// Note: I have changed the types to pointer to differentiate
// between the client leaving a field empty intentionally
// and the field not needing to be updated
var incomingData struct {
  Content  *string  `json:"content"`
  Author   *string  `json:"author"`
} 
// perform the decoding
err = a.readJSON(w, r, &incomingData)
if err != nil {
    a.badRequestResponse(w, r, err)
    return
}
// We need to now check the fields to see which ones need updating
// if incomingData.Content is nil, no update was provided
if incomingData.Content != nil {
    comment.Content = *incomingData.Content
}
// if incomingData.Author is nil, no update was provided
if incomingData.Author != nil {
  comment.Author = *incomingData.Author
}

// Before we write the updates to the DB let's validate
v := validator.New()
data.ValidateComment(v, comment)
if !v.IsEmpty() {
   a.failedValidationResponse(w, r, v.Errors)  
   return
}
// perform the update
err = a.commentModel.Update(comment)
if err != nil {
   a.serverErrorResponse(w, r, err)
   return 
}
data := envelope {
            "comment": comment,
      }
err = a.writeJSON(w, http.StatusOK, data, nil)
if err != nil {
   a.serverErrorResponse(w, r, err)
   return 
}
}

func (a *application)deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
  id, err := a.readIDParam(r)
   if err != nil {
       a.notFoundResponse(w, r)
       return 
   }
   err = a.commentModel.Delete(id)

   if err != nil {
       switch {
           case errors.Is(err, data.ErrRecordNotFound):
              a.notFoundResponse(w, r)
           default:
              a.serverErrorResponse(w, r, err)
       }
       return 
   }
   // display the comment
   data := envelope {
    "message": "comment successfully deleted",
}
err = a.writeJSON(w, http.StatusOK, data, nil)
if err != nil {
a.serverErrorResponse(w, r, err)
}

}

func (a *application)listCommentsHandler(w http.ResponseWriter,r *http.Request) {
  comments, err := a.commentModel.GetAll()
if err != nil {
    a.serverErrorResponse(w, r, err)
    return
  }

data := envelope {
    "comments": comments,
   }
err = a.writeJSON(w, http.StatusOK, data, nil)
if err != nil {
    a.serverErrorResponse(w, r, err)
  }
}
