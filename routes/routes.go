package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"roomino/api"
	"roomino/middleware"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	store := cookie.NewStore([]byte("something-very-secret"))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Use(sessions.Sessions("mysession", store))
	r.Use(middleware.Cors())

	v1 := r.Group("")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "success"})
		})

		v1.POST("register", api.UserRegisterHandler())
		v1.POST("login", api.UserLoginHandler())

		authed := v1.Group("", middleware.JWT())
		{
			authed.GET("profile", api.GetUserProfileHandler())
			authed.POST("profile/unitinfo", api.UnitInfoHandler())
			authed.POST("profile/petupdate", api.UpdatePetHandler())
			authed.GET("profile/petupdate", api.GetPetHandler())
			authed.POST("profile/petregister", api.CreatePetHandler())
			authed.GET("profile/interests", api.GetInterestsHandler())
			authed.POST("profile/interests", api.CreateInterestsHandler())
			authed.GET("profile/complexunitinfo", api.GetComplexUnitinfoHandler())
			authed.POST("profile/searchinterests", api.SearchInterestswithcondHandler())
		}
	}

	return r
}
