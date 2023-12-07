package main

import (
	"github.com/labstack/echo/v4"
	"github.com/richjyoung/echopen"
	v310 "github.com/richjyoung/echopen/openapi/v3.1.0"
)

const Description = `
This is a sample server Petstore server.  You can find out more about
Swagger at [http://swagger.io](http://swagger.io) or on [irc.freenode.net, #swagger](http://swagger.io/irc/).  For
this sample, you can use the api key "special-key" to test the authorization filters.`

func main() {
	// Create a new echOpen wrapper
	api := echopen.New(
		"Swagger Petstore",
		"1.0.0",
		echopen.WithSpecDescription(Description),
		echopen.WithSpecLicense(&v310.License{Name: "Apache 2.0", URL: "http://www.apache.org/licenses/LICENSE-2.0.html"}),
		echopen.WithSpecTermsOfService("http://swagger.io/terms/"),
		echopen.WithSpecContact(&v310.Contact{Email: "apiteam@swagger.io"}),
		echopen.WithSpecServer(&v310.Server{URL: "http://petstore.swagger.io/v2"}),
		echopen.WithSpecExternalDocs(&v310.ExternalDocs{
			Description: "Find out more about Swagger",
			URL:         "http://swagger.io",
		}),
		echopen.WithSpecTag(&v310.Tag{
			Name:        "pet",
			Description: "Everything about your Pets",
			ExternalDocs: &v310.ExternalDocs{
				Description: "Find out more",
				URL:         "http://swagger.io",
			},
		}),
		echopen.WithSpecTag(&v310.Tag{
			Name:        "store",
			Description: "Access to Petstore orders",
		}),
		echopen.WithSpecTag(&v310.Tag{
			Name:        "user",
			Description: "Operations about user",
			ExternalDocs: &v310.ExternalDocs{
				Description: "Find out more about our store",
				URL:         "http://swagger.io",
			},
		}),
	)

	api.Spec.GetComponents().AddSecurityScheme("petstore_auth", &v310.SecurityScheme{
		Type: v310.OAuth2SecuritySchemeType,
		Flows: &v310.OAuthFlows{
			Implicit: &v310.OAuthFlow{
				AuthorizationURL: "http://petstore.swagger.io/oauth/dialog",
				Scopes: map[string]string{
					"write:pets": "modify pets in your account",
					"read:pets":  "read your pets",
				},
			},
		},
	})

	api.Spec.GetComponents().AddSecurityScheme("api_key", &v310.SecurityScheme{
		Type: v310.APIKeySecuritySchemeType,
		Name: "api_key",
		In:   "header",
	})

	api.Spec.GetComponents().AddRequestBody("Pet", &v310.RequestBody{
		Content: map[string]*v310.MediaTypeObject{
			echo.MIMEApplicationJSON: {
				Schema: api.ToSchemaRef(Pet{}),
			},
		},
		Required: true,
	})

	api.Spec.GetComponents().AddRequestBody("UserArray", &v310.RequestBody{
		Content: map[string]*v310.MediaTypeObject{
			echo.MIMEApplicationJSON: {
				Schema: api.ToSchemaRef([]User{}),
			},
		},
		Required: true,
	})

	petGroup := api.Group("/pet", echopen.WithGroupTags("pet"))

	petGroup.POST(
		"",
		noop,
		echopen.WithOperationID("addPet"),
		echopen.WithSummary("Add a new pet to the store"),
		echopen.WithSecurityRequirement("petstore_auth", []string{"write:pets", "read:pets"}),
		echopen.WithRequestBodyRef("Pet"),
		echopen.WithResponseDescription("405", "Invalid input"),
	)

	petGroup.PUT(
		"",
		noop,
		echopen.WithOperationID("updatePet"),
		echopen.WithSummary("Update an existing pet"),
		echopen.WithSecurityRequirement("petstore_auth", []string{"write:pets", "read:pets"}),
		echopen.WithRequestBodyRef("Pet"),
		echopen.WithResponseDescription("400", "Invalid ID supplied"),
		echopen.WithResponseDescription("404", "Pet not found"),
		echopen.WithResponseDescription("405", "Validation exception"),
	)

	petGroup.GET(
		"/findByStatus",
		noop,
		echopen.WithOperationID("findPetsByStatus"),
		echopen.WithSummary("Finds Pets by status"),
		echopen.WithDescription("Multiple status values can be provided with comma separated strings"),
		echopen.WithSecurityRequirement("petstore_auth", []string{"write:pets", "read:pets"}),
		echopen.WithQueryStruct(FindPetsByStatusQuery{}),
		echopen.WithResponseStruct("200", "successful operation", []Pet{}),
		echopen.WithResponseDescription("400", "Invalid status value"),
	)

	petGroup.GET(
		"/findByTags",
		noop,
		echopen.WithOperationID("findPetsByTags"),
		echopen.WithSummary("Finds Pets by tags"),
		echopen.WithDescription("Muliple tags can be provided with comma separated strings. Use tag1, tag2, tag3 for testing."),
		echopen.WithSecurityRequirement("petstore_auth", []string{"write:pets", "read:pets"}),
		echopen.WithQueryStruct(FindPetsByTagsQuery{}),
		echopen.WithResponseStruct("200", "successful operation", []Pet{}),
		echopen.WithResponseDescription("400", "Invalid tag value"),
		echopen.WithDeprecated(),
	)

	petGroup.GET(
		"/:petId",
		noop,
		echopen.WithOperationID("getPetById"),
		echopen.WithSummary("Find pet by ID"),
		echopen.WithDescription("Returns a single pet"),
		echopen.WithSecurityRequirement("api_key", []string{}),
		echopen.WithPathParameter("petId", "ID of pet to return", int64(1234)),
		echopen.WithResponseStruct("200", "successful operation", Pet{}),
		echopen.WithResponseDescription("400", "Invalid ID supplied"),
		echopen.WithResponseDescription("404", "Pet not found"),
		echopen.WithResponseDescription("default", "successful response"),
	)

	petGroup.POST(
		"/:petId",
		noop,
		echopen.WithOperationID("updatePetWithForm"),
		echopen.WithPathParameter("petId", "ID of pet to return", int64(1234)),
		echopen.WithSecurityRequirement("petstore_auth", []string{"write:pets", "read:pets"}),
		echopen.WithRequestBodyStruct(echo.MIMEApplicationJSON, "", UpdatePet{}),
		echopen.WithResponseDescription("405", "Invalid input"),
	)

	petGroup.DELETE(
		"/:petId",
		noop,
		echopen.WithOperationID("deletePet"),
		echopen.WithSummary("Deletes a pet"),
		echopen.WithPathParameter("petId", "Pet id to delete", int64(1234)),
		echopen.WithHeaderParameter("api_key", "", ""),
		echopen.WithSecurityRequirement("petstore_auth", []string{"write:pets", "read:pets"}),
		echopen.WithResponseDescription("400", "Invalid ID supplied"),
		echopen.WithResponseDescription("404", "Pet not found"),
	)

	petGroup.POST(
		"/:petId/uploadImage",
		noop,
		echopen.WithOperationID("uploadFile"),
		echopen.WithSummary("uploads an image"),
		echopen.WithPathParameter("petId", "ID of pet to update", int64(1234)),
		echopen.WithSecurityRequirement("petstore_auth", []string{"write:pets", "read:pets"}),
		echopen.WithRequestBodySchema("application/octet-stream", &v310.Schema{
			Type:   "string",
			Format: "binary",
		}),
		echopen.WithResponseStruct("200", "successful operation", ApiResponse{}),
	)

	store := api.Group("/store", echopen.WithGroupTags("store"))

	store.GET(
		"/inventory",
		noop,
		echopen.WithOperationID("getInventory"),
		echopen.WithSummary("Returns pet inventories by status"),
		echopen.WithDescription("Returns a map of status codes to quantities"),
		echopen.WithSecurityRequirement("api_key", []string{}),
		echopen.WithResponseType("200", "successful operation", map[string]int32{}),
	)

	store.POST(
		"/order",
		noop,
		echopen.WithOperationID("placeOrder"),
		echopen.WithSummary("Place an order for a pet"),
		echopen.WithRequestBodyStruct(echo.MIMEApplicationJSON, "order placed for purchasing the pet", Order{}),
		echopen.WithResponseStruct("200", "successful operation", Order{}),
		echopen.WithResponseDescription("400", "Invalid Order"),
	)

	store.GET(
		"/order/:orderId",
		noop,
		echopen.WithOperationID("getOrderById"),
		echopen.WithSummary("Find purchase order by ID"),
		echopen.WithDescription("For valid response try integer IDs with value >= 1 and <= 10. Other values will generated exceptions"),
		echopen.WithPathParameter("orderId", "ID of order that needs to be fetched", int64(0)),
		echopen.WithResponseStruct("200", "successful operation", Order{}),
		echopen.WithResponseDescription("400", "Invalid ID supplied"),
		echopen.WithResponseDescription("404", "Order not found"),
	)

	store.DELETE(
		"/order/:orderId",
		noop,
		echopen.WithOperationID("deleteOrder"),
		echopen.WithSummary("Delete purchase order by ID"),
		echopen.WithDescription("For valid response try integer IDs with positive integer value. Negative or non-integer values will generate API errors"),
		echopen.WithPathParameter("orderId", "ID of the order that needs to be deleted", int64(0)),
		echopen.WithResponseDescription("400", "Invalid ID supplied"),
		echopen.WithResponseDescription("404", "Order not found"),
	)

	user := api.Group("/user", echopen.WithGroupTags("user"))

	user.POST(
		"",
		noop,
		echopen.WithOperationID("createUser"),
		echopen.WithSummary("Create user"),
		echopen.WithDescription("This can only be done by the logged in user."),
		echopen.WithRequestBodyStruct(echo.MIMEApplicationJSON, "Created user object", User{}),
		echopen.WithResponseDescription("default", "successful operation"),
	)

	user.POST(
		"/createWithArray",
		noop,
		echopen.WithOperationID("createUsersWithArrayInput"),
		echopen.WithSummary("Creates list of users with given input array"),
		echopen.WithRequestBodyRef("UserArray"),
		echopen.WithResponseDescription("default", "successful operation"),
	)

	user.POST(
		"/createWithList",
		noop,
		echopen.WithOperationID("createUsersWithListInput"),
		echopen.WithSummary("Creates list of users with given input array"),
		echopen.WithRequestBodyRef("UserArray"),
		echopen.WithResponseDescription("default", "successful operation"),
	)

	user.GET(
		"/login",
		noop,
		echopen.WithOperationID("loginUser"),
		echopen.WithSummary("Logs user into the system"),
		echopen.WithQueryStruct(LoginUserQuery{}),
		echopen.WithResponseType("200", "successful operation", "token"),
		echopen.WithResponseDescription("400", "Invalid username/password supplied"),
	)

	user.GET(
		"/logout",
		noop,
		echopen.WithOperationID("logoutUser"),
		echopen.WithSummary("Logs out current logged in user session"),
		echopen.WithResponseDescription("default", "successful operation"),
	)

	user.GET(
		"/:username",
		noop,
		echopen.WithOperationID("getUserByName"),
		echopen.WithSummary("Get user by user name"),
		echopen.WithPathParameter("username", "The name that needs to be fetched. Use user1 for testing. ", "username"),
		echopen.WithResponseStruct("200", "successful operation", User{}),
		echopen.WithResponseDescription("400", "Invalid username supplied"),
		echopen.WithResponseDescription("404", "User not found"),
	)

	user.PUT(
		"/:username",
		noop,
		echopen.WithOperationID("updateUser"),
		echopen.WithSummary("This can only be done by the logged in user."),
		echopen.WithPathParameter("username", "name that need to be updated", "username"),
		echopen.WithRequestBodyStruct(echo.MIMEApplicationJSON, "Updated user object", User{}),
		echopen.WithResponseDescription("400", "Invalid user supplied"),
		echopen.WithResponseDescription("404", "User not found"),
	)

	user.DELETE(
		"/:username",
		noop,
		echopen.WithOperationID("deleteUser"),
		echopen.WithSummary("This can only be done by the logged in user."),
		echopen.WithPathParameter("username", "The name that needs to be deleted", "username"),
		echopen.WithResponseDescription("400", "Invalid username supplied"),
		echopen.WithResponseDescription("404", "User not found"),
	)

	// Serve the generated schema
	api.ServeYAMLSpec("/openapi.yml")
	api.ServeSwaggerUI("/", "/openapi.yml", "5.10.3")

	// Write the full generated spec
	api.WriteYAMLSpec("openapi_out.yml")

	// Start the server
	api.Start("localhost:3000")
}

func noop(c echo.Context) error {

	return nil
}
