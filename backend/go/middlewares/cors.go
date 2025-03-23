package middlewares

import (

    "github.com/gin-contrib/cors"

    "github.com/gin-gonic/gin"   
    
)

func CorsMiddleware() gin.HandlerFunc {

    config := cors.DefaultConfig()

    config.AllowOrigins = []string{"*"}

    return cors.New(config) 
    
}