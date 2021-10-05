package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"study_websocket_go/conf"
	"study_websocket_go/ws"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	banner = "\n" +
		"    ___    _____   _   _    ___   __   __         __      __ ___     ___     ___     ___     ___    _  __    ___    _____             ___     ___   \n" +
		"   / __|  |_   _| | | | |  |   \\  \\ \\ / /    o O O\\ \\    / /| __|   | _ )   / __|   / _ \\   / __|  | |/ /   | __|  |_   _|    o O O  / __|   / _ \\  \n" +
		"   \\__ \\    | |   | |_| |  | |) |  \\ V /    o      \\ \\/\\/ / | _|    | _ \\   \\__ \\  | (_) | | (__   | ' <    | _|     | |     o      | (_ |  | (_) | \n" +
		"   |___/   _|_|_   \\___/   |___/   _|_|_   TS__[O]  \\_/\\_/  |___|   |___/   |___/   \\___/   \\___|  |_|\\_\\   |___|   _|_|_   TS__[O]  \\___|   \\___/  \n" +
		" _|\"\"\"\"\"|_|\"\"\"\"\"|_|\"\"\"\"\"|_|\"\"\"\"\"|_| \"\"\" | {======|_|\"\"\"\"\"|_|\"\"\"\"\"|_|\"\"\"\"\"|_|\"\"\"\"\"|_|\"\"\"\"\"|_|\"\"\"\"\"|_|\"\"\"\"\"|_|\"\"\"\"\"|_|\"\"\"\"\"| {======|_|\"\"\"\"\"|_|\"\"\"\"\"| \n" +
		" \"`-0-0-'\"`-0-0-'\"`-0-0-'\"`-0-0-'\"`-0-0-'./o--000'\"`-0-0-'\"`-0-0-'\"`-0-0-'\"`-0-0-'\"`-0-0-'\"`-0-0-'\"`-0-0-'\"`-0-0-'\"`-0-0-'./o--000'\"`-0-0-'\"`-0-0-' \n\n" +
		" => Starting listen %s\n"
)

func main() {
	studyWS := conf.StudyWS
	e := echoInit(studyWS)
	signal := sigInit(e)

	if err := initHandler(studyWS, e, signal); err != nil {
		e.Logger.Error("InitHandler Error")
		os.Exit(1)
	}

	startServer(studyWS, e)
}

func echoInit(studyWS *conf.ViperConfig) (e *echo.Echo) {
	e = echo.New()

	// Middleware
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.POST, echo.GET, echo.PUT, echo.DELETE},
	}))

	e.HideBanner = true

	return e
}

func sigInit(e *echo.Echo) chan os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		os.Interrupt,
	)
	go func() {
		sig := <-quit
		e.Logger.Error("Got signal", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
		signal.Stop(quit)
		close(quit)
	}()
	return quit
}

func startServer(studyWS *conf.ViperConfig, e *echo.Echo) {
	// Start Server
	apiServer := fmt.Sprintf("0.0.0.0:%d", studyWS.GetInt("port"))
	e.Logger.Debugf("Starting server, Listen[%s]", apiServer)

	fmt.Printf(banner, apiServer)
	if err := e.Start(apiServer); err != nil {
		e.Logger.Fatal(err)
	}
}

func initHandler(studyGoroutine *conf.ViperConfig, e *echo.Echo, signal <-chan os.Signal) error {
	// ws
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/ws", ws.Hello)

	return nil
}
