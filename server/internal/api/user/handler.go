package user

// The handler.go file is responsible for mapping the HTTP routes to their respective handler
// functions in the user service. In other words, it "registers" the services our API provides by
// linking them to specific URL paths and HTTP methods.
//
// For instance, it associates the HTTP POST method at the "/users" path with the HandleCreateUser
// function of the API struct, meaning when a POST request is sent to "/users",
// the HandleCreateUser function will handle that request.
//
// The RegisterHandlers method is typically called during the setup stage of the server,
// and it's how the server knows what function to execute for each API endpoint.

import "github.com/go-chi/chi/v5"

func (api *API) RegisterHandlers(r chi.Router) {
	r.Post("/users", api.HandleCreateUser)
}