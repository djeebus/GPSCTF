package db

var gamesTable = `
CREATE TABLE games (
	"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	"name" TEXT,
	"playerCount" INTEGER NOT NULL DEFAULT 0,
	"latitude" REAL NOT NULL,
	"longitude" REAL NOT NULL,
	"radius" INTEGER NOT NULL
)`

type Game struct {
	GameID      int64   `json:"gameId"`
	Name        string  `json:"name"`
	PlayerCount int     `json:"playerCount"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Radius      uint8   `json:"radius"`
}

func CreateGame(name string, latitude float64, longitude float64, radius uint8) (*Game, error) {
	query := `INSERT INTO games (name, latitude, longitude, radius) VALUES (?, ?, ?, ?)`
	statement, err := sqlDb.Prepare(query)
	if err != nil {
		return nil, err
	}

	result, err := statement.Exec(name, latitude, longitude, radius)
	if err != nil {
		return nil, err
	}

	lastRowId, err := result.LastInsertId()
	game := &Game{
		GameID:      lastRowId,
		Name:        name,
		PlayerCount: 0,
		Latitude:    latitude,
		Longitude:   longitude,
		Radius:      radius,
	}
	return game, err
}

func GetGame(gameID int64) (*Game, error) {
	query := `SELECT name, playerCount FROM games WHERE id = ?`
	stmt, err := sqlDb.Prepare(query)
	if err != nil {
		return nil, err
	}

	row, err := stmt.Query(gameID)
	if err != nil {
		return nil, err
	}

	var name string
	var count int

	for row.Next() {
		err = row.Scan(&name, &count)
		if err != nil {
			return nil, err
		}

		game := Game{
			GameID:      gameID,
			Name:        name,
			PlayerCount: count,
		}
		return &game, nil
	}

	return nil, &GameNotFoundError{GameID: gameID}
}

func addOnePlayer(gameId int64) error {
	query := `UPDATE games SET playerCount = playerCount + 1 WHERE id = ?`
	stmt, err := sqlDb.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(gameId)
	if err != nil {
		return err
	}

	return nil
}

func ListGames() (*[]Game, error) {
	query := `SELECT id, name, playerCount FROM games`
	stmt, err := sqlDb.Prepare(query)
	if err != nil {
		return nil, err
	}

	row, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	games := make([]Game, 0)
	var id int64
	var name string
	var count int

	for row.Next() {
		err = row.Scan(&id, &name, &count)
		if err != nil {
			return nil, err
		}

		game := Game{
			GameID:      id,
			Name:        name,
			PlayerCount: count,
		}
		games = append(games, game)
	}

	return &games, nil
}
