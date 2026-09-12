package handler

import (
	"fmt"
	"strings"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/iqox/nieve/internal/rank"
)

type PlayerHandler struct {
	player.NopHandler
	p    *player.Player
	xuid string
	rm   *rank.Manager
}

func New(p *player.Player, xuid string, rm *rank.Manager) *PlayerHandler {
	return &PlayerHandler{p: p, xuid: xuid, rm: rm}
}

func (h *PlayerHandler) HandleChat(ctx *player.Context, message *string) {
	ctx.Cancel()

	r := h.rm.PlayerRank(h.xuid)
	formatted := strings.ReplaceAll(r.Chat, "{name}", h.p.Name())
	formatted = strings.ReplaceAll(formatted, "{msg}", *message)

	fmt.Fprintln(chat.Global, formatted)
}
