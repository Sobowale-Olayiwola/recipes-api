// Recipes API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/PacktPublishing/Building-Distributed-Applications-in-Gin.
//
// Schemes: http
// Host: localhost:8080
// BasePath: /
// Version: 1.0.0
// Contact: Olayiwola Sobowale <layitheinfotechguru@gmail.com> https://github.com/Sobowale-Olayiwola
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"recipes-api/handlers"

	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var recipesHandler *handlers.RecipesHandler
var authHandler *handlers.AuthHandler

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collectionRecipes := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	status := redisClient.Ping(ctx)
	log.Println(status)

	recipesHandler = handlers.NewRecipesHandler(ctx, collectionRecipes, redisClient)

	collectionUsers := client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")
	authHandler = handlers.NewAuthHandler(ctx, collectionUsers)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// tokenValue := c.GetHeader("Authorization")
		// claims := &handlers.Claims{}
		// tkn, err := jwt.ParseWithClaims(tokenValue, claims, func(t *jwt.Token) (interface{}, error) {
		// 	return []byte(os.Getenv("JWT_SECRET")), nil
		// })
		// if err != nil {
		// 	c.AbortWithStatus(http.StatusUnauthorized)
		// }

		// if tkn == nil || !tkn.Valid {
		// 	c.AbortWithStatus(http.StatusUnauthorized)
		// }
		session := sessions.Default(c)
		sessionToken := session.Get("token")
		if sessionToken == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Not logged",
			})
			c.Abort()
		}
		c.Next()
	}
}

func main() {
	router := gin.Default()

	store, _ := redisStore.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	router.Use(sessions.Sessions("recipes_api", store))
	authorized := router.Group("/")
	authorized.Use(AuthMiddleware())
	authorized.POST("/recipes", recipesHandler.NewRecipeHandler)
	authorized.PUT("/recipes/:id",
		recipesHandler.UpdateRecipeHandler)
	authorized.DELETE("/recipes/:id",
		recipesHandler.DeleteRecipeHandler)
	authorized.GET("/recipes/:id",
		recipesHandler.GetOneRecipeHandler)
	// router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/refresh", authHandler.RefreshHandler)
	router.POST("/signout", authHandler.SignOutHandler)
	//router.GET("/recipes/search", SearchRecipesHandler)
	// router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	// router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	// router.GET("/recipes/:id", recipesHandler.GetOneRecipeHandler)
	router.Run()
}
