package rank

type Storage interface {
	PlayerRank(xuid string) (int, error)
	SavePlayerRank(xuid string, rankID int) error
}
