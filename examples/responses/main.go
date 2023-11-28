package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/richjyoung/echopen"
	v310 "github.com/richjyoung/echopen/openapi/v3.1.0"
)

func main() {
	api := echopen.New(
		"Responses",
		"1.0.0",
		echopen.WithSchemaDescription("Example to show reused responses, loosely based on the Petstore examples."),
		echopen.WithSchemaTag(&v310.Tag{Name: "pets", Description: "Pets endpoints"}),
	)

	api.Schema.GetComponents().AddResponse("UnexpectedErrorResponse", &v310.Response{
		Description: echopen.PtrTo("unexpected error"),
		Content: map[string]*v310.MediaTypeObject{
			"application/json": {
				Schema: api.ToSchemaRef(Error{}),
			},
		},
	})

	api.Schema.GetComponents().AddResponse("NotFoundResponse", &v310.Response{
		Description: echopen.PtrTo("not found"),
		Content: map[string]*v310.MediaTypeObject{
			"application/json": {
				Schema: api.ToSchemaRef(Error{}),
			},
		},
	})

	api.Schema.GetComponents().AddResponse("PetResponse", &v310.Response{
		Description: echopen.PtrTo("pet response"),
		Content: map[string]*v310.MediaTypeObject{
			"application/json": {
				Schema: api.ToSchemaRef(Pet{}),
			},
		},
	})

	api.GET(
		"/pets",
		findPets,
		echopen.WithOperationID("findPets"),
		echopen.WithDescription("Finds all pets the user has access to"),
		echopen.WithTags("pets"),
		echopen.WithQueryStruct(FindPetsQuery{}),
		echopen.WithResponseBody(fmt.Sprint(http.StatusOK), "pet response", []Pet{}),
		echopen.WithResponseRef("default", "UnexpectedErrorResponse"),
	)

	api.POST(
		"/pets",
		addPet,
		echopen.WithOperationID("addPet"),
		echopen.WithDescription("Creates a new pet in the store. Duplicates are allowed"),
		echopen.WithTags("pets"),
		echopen.WithRequestBody("Pet to add to the store", NewPet{}),
		echopen.WithResponseRef(fmt.Sprint(http.StatusOK), "PetResponse"),
		echopen.WithResponseRef("default", "UnexpectedErrorResponse"),
	)

	api.GET(
		"/pets/:id",
		findPetByID,
		echopen.WithOperationID("findPetByID"),
		echopen.WithDescription("Returns a user based on a single ID, if the user does not have access to the pet"),
		echopen.WithTags("pets"),
		echopen.WithPathParameter(&echopen.PathParameter{
			Name:        "id",
			Description: "ID of pet to fetch",
			Schema: &v310.Schema{
				Type:   v310.IntegerSchemaType,
				Format: "int64",
			},
		}),
		echopen.WithResponseRef(fmt.Sprint(http.StatusOK), "PetResponse"),
		echopen.WithResponseRef(fmt.Sprint(http.StatusNotFound), "NotFoundResponse"),
		echopen.WithResponseRef("default", "UnexpectedErrorResponse"),
	)

	api.DELETE(
		"/pets/:id",
		deletePet,
		echopen.WithOperationID("deletePet"),
		echopen.WithDescription("deletes a single pet based on the ID supplied"),
		echopen.WithTags("pets"),
		echopen.WithPathParameter(&echopen.PathParameter{
			Name:        "id",
			Description: "ID of pet to delete",
			Schema: &v310.Schema{
				Type:   v310.IntegerSchemaType,
				Format: "int64",
			},
		}),
		echopen.WithResponse(fmt.Sprint(http.StatusNoContent), "pet deleted"),
		echopen.WithResponseRef(fmt.Sprint(http.StatusNotFound), "NotFoundResponse"),
		echopen.WithResponseRef("default", "UnexpectedErrorResponse"),
	)

	// Serve the generated schema
	api.ServeYAMLSchema("/openapi.yml")
	api.ServeUI("/", "/openapi.yml", "5.10.3")

	// Start the server
	api.Start("localhost:3030")
}

func findPets(c echo.Context) error {
	return nil
}

func addPet(c echo.Context) error {
	return nil
}

func findPetByID(c echo.Context) error {
	return nil
}

func deletePet(c echo.Context) error {
	return nil
}
