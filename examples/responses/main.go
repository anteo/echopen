package main

import (
	"fmt"
	"net/http"

	"github.com/anteo/echopen"
	v310 "github.com/anteo/echopen/openapi/v3.1.0"
	"github.com/labstack/echo/v4"
)

func main() {
	api := echopen.New(
		"Responses",
		"1.0.0",
		echopen.WithSpecDescription("Example to show reused responses, loosely based on the Petstore examples."),
		echopen.WithSpecTag(&v310.Tag{Name: "pets", Description: "Pets endpoints"}),
	)

	api.Spec.GetComponents().AddJSONResponse("UnexpectedErrorResponse", "unexpected error", api.ToSchemaRef(Error{}))
	api.Spec.GetComponents().AddJSONResponse("NotFoundResponse", "not found", api.ToSchemaRef(Error{}))
	api.Spec.GetComponents().AddJSONResponse("PetResponse", "pet response", api.ToSchemaRef(Pet{}))

	api.GET(
		"/pets",
		findPets,
		echopen.WithOperationID("findPets"),
		echopen.WithDescription("Finds all pets the user has access to"),
		echopen.WithTags("pets"),
		echopen.WithQueryStruct(FindPetsQuery{}),
		echopen.WithResponseStruct(fmt.Sprint(http.StatusOK), "pet response", []Pet{}),
		echopen.WithResponseRef("default", "UnexpectedErrorResponse"),
	)

	api.POST(
		"/pets",
		addPet,
		echopen.WithOperationID("addPet"),
		echopen.WithDescription("Creates a new pet in the store. Duplicates are allowed"),
		echopen.WithTags("pets"),
		echopen.WithRequestBodyStruct(echo.MIMEApplicationJSON, "Pet to add to the store", NewPet{}),
		echopen.WithResponseRef(fmt.Sprint(http.StatusOK), "PetResponse"),
		echopen.WithResponseRef("default", "UnexpectedErrorResponse"),
	)

	api.GET(
		"/pets/:id",
		findPetByID,
		echopen.WithOperationID("findPetByID"),
		echopen.WithDescription("Returns a user based on a single ID, if the user does not have access to the pet"),
		echopen.WithTags("pets"),
		// echopen.WithPathParameter(&echopen.PathParameterConfig{
		// 	Name:        "id",
		// 	Description: "ID of pet to fetch",
		// 	Schema: &v310.Schema{
		// 		Type:   v310.IntegerSchemaType,
		// 		Format: "int64",
		// 	},
		// }),
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
		// echopen.WithPathParameter(&echopen.PathParameterConfig{
		// 	Name:        "id",
		// 	Description: "ID of pet to delete",
		// 	Schema: &v310.Schema{
		// 		Type:   v310.IntegerSchemaType,
		// 		Format: "int64",
		// 	},
		// }),
		echopen.WithResponseDescription(fmt.Sprint(http.StatusNoContent), "pet deleted"),
		echopen.WithResponseRef(fmt.Sprint(http.StatusNotFound), "NotFoundResponse"),
		echopen.WithResponseRef("default", "UnexpectedErrorResponse"),
	)

	// Serve the generated schema
	api.ServeYAMLSpec("/openapi.yml")
	api.ServeSwaggerUI("/", "/openapi.yml", "5.10.3")

	// Write the full generated spec
	api.WriteYAMLSpec("openapi_out.yml")

	// Start the server
	api.Start("localhost:3000")
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
