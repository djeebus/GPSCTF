package db

var playersTable = `
CREATE TABLE players (
	"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	"gameId" INTEGER,
	"name" TEXT,
	"latitude" TEXT DEFAULT NULL,
	"longitude" TEXT DEFAULT NULL,
	"lastUpdate" DATETIME DEFAULT NULL,

	FOREIGN KEY(gameId) REFERENCES games(id)
)`

type Player struct {
	PlayerID int64  `json:"playerId"`
	GameID   int64  `json:"gameId"`
	Name     string `json:"name"`
}

func CreatePlayer(gameId int64, name string) (*Player, error) {
	query := `INSERT INTO players (gameId, name) VALUES (?, ?)`
	statement, err := sqlDb.Prepare(query)
	if err != nil {
		return nil, err
	}

	result, err := statement.Exec(gameId, name)
	if err != nil {
		return nil, err
	}

	lastRowId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	player := &Player{
		PlayerID: lastRowId,
		GameID:   gameId,
		Name:     name,
	}

	err = addOnePlayer(gameId)
	if err != nil {
		return nil, err
	}

	return player, err
}

func GetPlayer(playerId int64) (*Player, error) {
	query := `SELECT gameId, name FROM players WHERE id = ?`
	stmt, err := sqlDb.Prepare(query)
	if err != nil {
		return nil, err
	}

	row, err := stmt.Query(playerId)
	if err != nil {
		return nil, err
	}

	var gameId int64
	var name string

	for row.Next() {
		err = row.Scan(&gameId, &name)
		if err != nil {
			return nil, err
		}

		player := Player{
			PlayerID: playerId,
			GameID: gameId,
			Name: name,
		}
		return &player, nil
	}

	return nil, &PlayerNotFoundError{PlayerID: playerId}
}


func GetGamePlayers(gameId int64) ([]*Player, error) {
	query := `SELECT id, name FROM players WHERE gameId = ?`
	stmt, err := sqlDb.Prepare(query)
	if err != nil {
		return nil, err
	}

	row, err := stmt.Query(gameId)
	if err != nil {
		return nil, err
	}

	players := make([]*Player, 0)
	var id int64
	var name string

	for row.Next() {
		err = row.Scan(&id, &name)
		if err != nil {
			return nil, err
		}

		player := &Player{
			PlayerID: id,
			GameID:   gameId,
			Name:     name,
		}

		players = append(players, player)
	}

	return players, nil
}
