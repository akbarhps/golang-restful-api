package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go-api/app"
	"go-api/middleware"
	"go-api/model/comment"
	"go-api/model/like"
	"go-api/model/post"
	"go-api/model/resource"
	"go-api/model/user"
)

func main() {
	app.Init()
	validate := validator.New()

	router := gin.Default()
	router.Use(middleware.JWTValidator())
	router.Use(gin.CustomRecovery(middleware.PanicHandler))

	// repositories
	userRepository := user.NewRepository()
	postRepository := post.NewRepository()
	likeRepository := like.NewRepository()
	commentRepository := comment.NewRepository()
	resourceRepository := resource.NewRepository()

	// services
	userService := user.NewService(validate, userRepository)
	postService := post.NewService(validate, postRepository, resourceRepository, likeRepository, commentRepository)
	likeService := like.NewService(validate, likeRepository)
	commentService := comment.NewService(validate, commentRepository)

	// controllers
	userController := user.NewController(userService)
	postController := post.NewController(postService)
	likeController := like.NewController(likeService)
	commentController := comment.NewController(commentService)

	// Routes
	user.InitRoutes(router, userController)
	post.InitRoutes(router, postController)
	like.InitRoutes(router, likeController)
	comment.InitRoutes(router, commentController)

	err := router.Run(":3000")
	if err != nil {
		panic(err)
	}
}
