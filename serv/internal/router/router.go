package router

import (
	"bsquared.network/b2-message-channel-serv/internal/boot"
	"bsquared.network/b2-message-channel-serv/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func InitRoutes(initVal *boot.Initialization, cfg config.AppConfig) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "UPDATE", "PATCH"},
		AllowHeaders:     []string{"Origin, X-Requested-With, APP-ID, APP-SECRET, Content-Type, Accept, Authorization"},
		ExposeHeaders:    []string{"Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))
	//r.Use(sentrygin.New(sentrygin.Options{
	//	Repanic: true,
	//}))
	//
	//r.Use(sentrygin.New(sentrygin.Options{
	//	Repanic: true,
	//}))
	//r.Use(gin.CustomRecovery(middlewares.ErrorHandler))
	//if err := sentry.Init(sentry.ClientOptions{
	//	Dsn:           cfg.Sentry.Url,
	//	EnableTracing: true,
	//	// Set TracesSampleRate to 1.0 to capture 100%
	//	// of transactions for performance monitoring.
	//	// We recommend adjusting this value in production,
	//	TracesSampleRate: cfg.Sentry.SampleRate,
	//	Debug:            cfg.Sentry.Env != "prod",
	//	Environment:      cfg.Sentry.Env,
	//	Release:          cfg.Sentry.Release,
	//}); err != nil {
	//	log.Error("Sentry initialization failed: %v", err)
	//}
	//r.Use(middlewares.JwtMiddleWare(cfg.Api.Secret))
	//r.Use(middlewares.FixTokenMiddleWare(cfg.Api.Secret))
	//BybitApi := r.Group("/api/bybit")
	//{
	//	BybitApi.GET("/verify-by-wallet", initVal.BybitInit.BybitCtrl.VerifyByWallet)
	//}
	//StakeApi := r.Group("/stake")
	//{
	//	StakeApi.GET("/rewords", initVal.StakeInit.StakeCtrl.Rewords)
	//	StakeApi.GET("/histories", initVal.StakeInit.StakeCtrl.Histories)
	//	StakeApi.GET("/statistical", initVal.StakeInit.StakeCtrl.Statistical)
	//}
	MessageApi := r.Group("/message")
	{
		MessageApi.GET("/records", initVal.MessageInit.MessageCtrl.Records)
	}
	return r
}
