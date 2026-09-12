package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/iqox/nieve/internal/rank"
)

var RM *rank.Manager

type RankCommand struct {
	Target []cmd.Target `cmd:"jugador"`
	RankID int          `cmd:"rango"`
}

func (c RankCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if s, ok := src.(*player.Player); ok {
		if !RM.PlayerRank(s.XUID()).HasPermission("nieve.permission.setrank") {
			o.Errorf("No tienes permiso para usar este comando.")
			return
		}
	}

	if len(c.Target) == 0 {
		o.Errorf("Debes especificar un jugador.")
		return
	}

	for _, t := range c.Target {
		p, ok := t.(*player.Player)
		if !ok {
			continue // el target no es un jugador (podría ser otra entidad)
		}
		if err := RM.SetPlayerRank(p.XUID(), c.RankID); err != nil {
			o.Errorf("Error: %v", err)
			continue
		}
		o.Printf("Rango de %s actualizado a %d.", p.Name(), c.RankID)
	}
}

func RegisterRankCommand(rm *rank.Manager) {
	RM = rm
	cmd.Register(cmd.New("rank", "Asigna un rango a un jugador.", []string{"rango"}, RankCommand{}))
}
