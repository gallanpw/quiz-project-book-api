package routes

import (
	"quiz-project-book-api/controllers"
	"quiz-project-book-api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// Endpoint untuk users (login)
	users := api.Group("/users")
	{
		users.POST("/login", controllers.LoginUser)
		users.POST("/register", controllers.RegisterUser)
	}

	// Grup endpoint yang dilindungi JWT
	protected := api.Group("/")
	protected.Use(middlewares.JWTAuthMiddleware())
	{
		// Endpoint untuk logout
		protected.POST("/users/logout", controllers.LogoutUser)

		// Endpoint untuk Categories
		categories := protected.Group("/categories")
		{
			categories.GET("", controllers.GetAllCategories)
			categories.POST("", controllers.CreateCategory)
			categories.GET("/:id", controllers.GetCategoryByID)
			categories.PUT("/:id", controllers.UpdateCategory)
			categories.DELETE("/:id", controllers.DeleteCategory)
			categories.GET("/:id/books", controllers.GetBooksByCategory)
		}

		// Endpoint untuk Books
		books := protected.Group("/books")
		{
			books.GET("", controllers.GetAllBooks)
			books.POST("", controllers.CreateBook)
			books.GET("/:id", controllers.GetBookByID)
			books.DELETE("/:id", controllers.DeleteBook)
			books.PUT("/:id", controllers.UpdateBook)
		}
	}
}
