package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/richjyoung/echopen"
	v310 "github.com/richjyoung/echopen/openapi/v3.1.0"
)

var (
	todos = map[uuid.UUID]*Todo{
		uuid.Must(uuid.FromString("11c7810d-6627-497a-91e9-e3dc4812ce30")): {
			ID:          uuid.Must(uuid.FromString("11c7810d-6627-497a-91e9-e3dc4812ce30")),
			Description: "New todo",
			Created:     time.Now(),
			Updated:     time.Now(),
			Tags:        []string{"urgent"},
		},
		uuid.Must(uuid.FromString("65071c49-6942-4fca-a11f-7cacd59c6e6d")): {
			ID:          uuid.Must(uuid.FromString("65071c49-6942-4fca-a11f-7cacd59c6e6d")),
			Description: "Completed todo",
			Created:     time.Now(),
			Updated:     time.Now(),
			Tags:        []string{"urgent"},
			Completed:   true,
		},
	}
	mu = sync.Mutex{}
)

func main() {
	api := echopen.New(
		"Responses",
		"1.0.0",
		echopen.WithSpecDescription("Example to show a CRUD-like API, based on Todo lists."),
		echopen.WithSpecTag(&v310.Tag{Name: "todo", Description: "Todo endpoints"}),
	)

	api.Engine.Logger.SetLevel(log.DEBUG)

	api.Spec.GetComponents().AddJSONResponse("UnexpectedErrorResponse", "Unexpected error", api.ToSchemaRef(Error{}))
	api.Spec.GetComponents().AddJSONResponse("BadRequestResponse", "Bad request", api.ToSchemaRef(Error{}))
	api.Spec.GetComponents().AddJSONResponse("NotFoundResponse", "Not found", api.ToSchemaRef(Error{}))

	todos := api.Group("/todo", echopen.WithGroupTags("todo"))

	todos.GET(
		"",
		getTodos,
		echopen.WithDescription("Get all todos for the user"),
		echopen.WithQueryStruct(GetTodosQuery{}),
		echopen.WithResponseStruct(fmt.Sprint(http.StatusOK), "Successful response", []Todo{}),
		echopen.WithResponseRef("default", "UnexpectedErrorResponse"),
	)

	todos.POST(
		"",
		newTodo,
		echopen.WithDescription("Create a new Todo"),
		echopen.WithRequestBodyStruct(echo.MIMEApplicationJSON, "New Todo", NewTodo{}),
		echopen.WithResponseStruct(fmt.Sprint(http.StatusCreated), "Successful response", Todo{}),
		echopen.WithResponseRef(fmt.Sprint(http.StatusBadRequest), "BadRequestResponse"),
		echopen.WithResponseRef("default", "UnexpectedErrorResponse"),
	)

	todos.GET(
		"/:id",
		getTodoByID,
		echopen.WithDescription("Get todo by ID"),
		echopen.WithPathParameter("id", "Todo ID", uuid.Must(uuid.FromString("11c7810d-6627-497a-91e9-e3dc4812ce30"))),
		echopen.WithResponseStruct(fmt.Sprint(http.StatusOK), "Successful response", Todo{}),
		echopen.WithResponseRef("default", "UnexpectedErrorResponse"),
	)

	todos.PATCH(
		"/:id",
		updateTodo,
		echopen.WithDescription("Update a todo"),
		echopen.WithPathParameter("id", "Todo ID", uuid.Must(uuid.FromString("11c7810d-6627-497a-91e9-e3dc4812ce30"))),
		echopen.WithRequestBodyStruct(echo.MIMEApplicationJSON, "Updated Todo", UpdateTodo{}),
		echopen.WithResponseStruct(fmt.Sprint(http.StatusOK), "Successful response", Todo{}),
		echopen.WithResponseRef(fmt.Sprint(http.StatusBadRequest), "BadRequestResponse"),
		echopen.WithResponseRef("default", "UnexpectedErrorResponse"),
	)

	todos.DELETE(
		"/:id",
		deleteTodo,
		echopen.WithDescription("Delete a todo"),
		echopen.WithPathParameter("id", "Todo ID", uuid.Must(uuid.FromString("11c7810d-6627-497a-91e9-e3dc4812ce30"))),
		echopen.WithResponseDescription(fmt.Sprint(http.StatusNoContent), "Successful response"),
		echopen.WithResponseRef(fmt.Sprint(http.StatusBadRequest), "BadRequestResponse"),
		echopen.WithResponseRef("default", "UnexpectedErrorResponse"),
	)

	// Serve the generated schema
	api.ServeYAMLSpec("/openapi.yml")
	api.ServeRapidoc("/", "/openapi.yml")

	// Write the full generated spec
	api.WriteYAMLSpec("openapi_out.yml")

	// Start the server
	api.Start("localhost:3000")
}

func getTodos(c echo.Context) error {
	qry := c.Get("query").(*GetTodosQuery)
	res := []*Todo{}

	mu.Lock()
	defer mu.Unlock()
	for _, v := range todos {
		if qry.Search != "" {
			if !strings.Contains(v.Description, qry.Search) {
				continue
			}
		}
		if !qry.ShowCompleted && v.Completed {
			continue
		}
		res = append(res, v)
	}

	return c.JSON(http.StatusOK, res)
}

func newTodo(c echo.Context) error {
	body := c.Get("body").(*NewTodo)

	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	if body.Description == "" {
		return c.JSON(http.StatusBadRequest, Error{Message: http.StatusText(http.StatusBadRequest)})
	}

	mu.Lock()
	defer mu.Unlock()
	todos[id] = &Todo{
		ID:          id,
		Description: body.Description,
		Tags:        body.Tags,
		Created:     time.Now(),
		Updated:     time.Now(),
	}

	return c.JSON(http.StatusCreated, todos[id])
}

func getTodoByID(c echo.Context) error {
	id := c.Get("path.id").(uuid.UUID)

	mu.Lock()
	defer mu.Unlock()
	if res, exists := todos[id]; exists {
		return c.JSON(http.StatusOK, res)
	}

	return c.JSON(http.StatusNotFound, Error{Message: http.StatusText(http.StatusNotFound)})
}

func updateTodo(c echo.Context) error {
	id := c.Get("path.id").(uuid.UUID)
	body := c.Get("body").(*UpdateTodo)

	mu.Lock()
	defer mu.Unlock()
	if res, exists := todos[id]; exists {
		res.Completed = body.Completed
		if body.Description != "" {
			res.Description = body.Description
		}
		if len(body.Tags) != 0 {
			res.Tags = body.Tags
		}
		res.Updated = time.Now()
		return c.JSON(http.StatusOK, res)
	}

	return c.JSON(http.StatusNotFound, Error{Message: http.StatusText(http.StatusNotFound)})
}

func deleteTodo(c echo.Context) error {
	id := c.Get("path.id").(uuid.UUID)

	mu.Lock()
	defer mu.Unlock()
	if _, exists := todos[id]; exists {
		delete(todos, id)
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusNotFound, Error{Message: http.StatusText(http.StatusNotFound)})
}
