package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"time"

	B "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/common/urltest"
	"github.com/sagernet/sing-box/experimental/libbox"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/service"
	"github.com/sagernet/sing/service/filemanager"
	"github.com/sagernet/sing/service/pause"
)

func main() {
	options := option.Options{}
	content, err := os.ReadFile("./hconfigs/full.json")
	json.Unmarshal(content, &options)
	fmt.Println(string(content))
	libbox.Setup("./hconfigs/", "./hconfigs/", "./hconfigs/", true)
	ctx, cancel := context.WithCancel(context.Background())
	ctx = filemanager.WithDefault(ctx, "./hconfigs/", "./hconfigs/", os.Getuid(), os.Getgid())
	urlTestHistoryStorage := urltest.NewHistoryStorage()
	ctx = service.ContextWithPtr(ctx, urlTestHistoryStorage)

	instance, err := B.New(B.Options{
		Context: ctx,
		Options: options,
	})
	if err != nil {
		cancel()
		return
	}
	commandServer = libbox.NewCommandServer(&CommandServerHandler{}, 100)
	commandServer.Start()
	service := libbox.NewBoxService(
		ctx,
		cancel,
		instance,
		service.FromContext[pause.Manager](ctx),
		urlTestHistoryStorage,
	)
	// instance.Start()
	service.Start()

	commandServer.SetService(&service)
	<-time.Tick(1 * time.Second)
	fmt.Println("command group update")

	groupInfoOnlyClient := libbox.NewCommandClient(
		&CommandClientHandler{
			logger: log.NewNOPFactory().NewLogger("[GroupInfoOnly Command Client]"),
		},
		&libbox.CommandClientOptions{
			Command:        libbox.CommandGroupInfoOnly,
			StatusInterval: 3000000000, //300ms debounce
		},
	)
	groupInfoOnlyClient.Connect()
	groupClient := libbox.NewCommandClient(
		&CommandClientHandler{
			logger: log.NewNOPFactory().NewLogger("[GroupInfoOnly Command Client]"),
		},
		&libbox.CommandClientOptions{
			Command:        libbox.CommandGroup,
			StatusInterval: 3000000000, //300ms debounce
		},
	)
	groupClient.Connect()
	for i := 0; i < 4; i++ {
		<-time.Tick(1000 * time.Millisecond)
		fmt.Println("selecting auto")
		libbox.NewStandaloneCommandClient().SelectOutbound("Select", "Auto")
		<-time.Tick(1000 * time.Millisecond)
		fmt.Println("selecting outbound")
		libbox.NewStandaloneCommandClient().SelectOutbound("Select", "test")
	}
	fmt.Println("===========Finished many bounce")
	<-time.Tick(20000 * time.Millisecond)
	fmt.Println("===========selecting final auto")
	libbox.NewStandaloneCommandClient().SelectOutbound("Select", "Auto")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, os.Kill)
	<-sigCh

}
