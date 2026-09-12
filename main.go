package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/iqox/nieve/internal/command"
	"github.com/iqox/nieve/internal/handler"
	"github.com/iqox/nieve/internal/rank"
	"github.com/pelletier/go-toml"
	slogrus "github.com/samber/slog-logrus/v2"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger
var Server *server.Server

func main() {
	Logger = logrus.New()
	Logger.Formatter = &logrus.TextFormatter{ForceColors: true}
	Logger.Level = logrus.DebugLevel

	log := slog.New(slogrus.Option{Logger: Logger}.NewLogrusHandler())

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	conf, err := readConfig(log)
	if err != nil {
		Logger.Fatalln(err)
	}

	Server = conf.New()
	Server.CloseOnProgramEnd()

	Server.World().SetTime(18_000)
	Server.World().StopTime()

	Server.World().StopWeatherCycle()
	Server.World().StopRaining()
	Server.World().StopThundering()

	Server.Listen()
	Logger.Infoln("Server is now listening for connections.")

	store, err := rank.NewJSONStorage("player_ranks.json")
	if err != nil {
		Logger.Fatalln(err)
	}

	RankManager := rank.NewManager(0, store)

	RankManager.RegisterRank(&rank.Rank{
		ID: 0, Display: "§7{name}", Chat: "§8[§7Jugador§8]§r §7{name} §8» §f{msg}", Priority: 0,
	})

	RankManager.RegisterRank(&rank.Rank{
		ID: 1, Display: "§c{name}", Chat: "§8[§cAdministrador§8]§r §c{name} §8» §f{msg}", Priority: 255,
		Permissions: map[string]bool{"nieve.permission.fly": true, "nieve.permission.vanish": true},
	})

	command.RegisterRankCommand(RankManager)

	for p := range Server.Accept() {
		xuid := p.XUID()
		if err := RankManager.LoadPlayer(xuid); err != nil {
			Logger.Errorln(err)
		}

		p.Handle(handler.New(p, xuid, RankManager))

		r := RankManager.PlayerRank(xuid)
		nametag := strings.ReplaceAll(r.Display, "{name}", p.Name())
		p.SetNameTag(nametag)

		Logger.Infof("%s connected with rank id: %d\n", p.Name(), RankManager.PlayerRank(xuid).ID)
	}

	Logger.Warnln("Accept loop exited — server is shutting down.")
}

func readConfig(log *slog.Logger) (server.Config, error) {
	c := server.DefaultConfig()
	var zero server.Config
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return zero, fmt.Errorf("encode default config: %v", err)
		}
		if err := os.WriteFile("config.toml", data, 0644); err != nil {
			return zero, fmt.Errorf("create default config: %v", err)
		}
		return c.Config(log)
	}
	data, err := os.ReadFile("config.toml")
	if err != nil {
		return zero, fmt.Errorf("read config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return zero, fmt.Errorf("decode config: %v", err)
	}
	return c.Config(log)
}
