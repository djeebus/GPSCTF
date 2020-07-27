GPS CTF Backend
===============

This is the initial commit of the API and game processor. Both are colocated in a single command. 

Architecture
------------
- ["api"](#api-module) module contains all the code that manages web traffic
- ["app"](#app-module) module contains the logic of processing a game
- "cmd" module contains the cli
- "db" module contains database handling code 

Creating a new game via the API creates a GameProcessor and starts a goroutine that waits for players. Creating a player registers it with the Game

api module
----------

- renderers.go: helpers for rendering json errors and responses
- router.go: ties routes to their handlers 
- schema.go: support code for using jsonschema to validate api requests
- status.go: view to test that the server is up
- games.go, players.go: implements the actual apis
- player_client.go: websocket handling code
- router_test.go: this is the full functional test that plays through a game

app module
----------

- worker.go: main background processor. one is created for the lifetime of the app; it manages the game loop and coordinates connecting clients to games
- [game_processor.go](#gameprocessor-class): the game loop. this manages the game. one is created per game
- util.go: implements the math for processing the players. VERY likely it is all wrong, but it wasn't the goal of this app.

GameProcessor class
-------------------

The game processor exists in two phases. 
1) Waiting for players. After the game is created, the processor waits for players to connect. 
2) In Progress. When the /games/{gameId}/:start API is called, the game processor is changed to inProgress=true, and starts processing player location updates. Once a player gets within `minimumFeetToWin`, it sends out a "winner" message and ends the game.

TODO:
----
- Use a scalable database (postgres? mysql? something weirder?)
- Handle crash recovery. If the app closes and reopens, the games in progress are lost. Will need to push game and player updates back to the db, and put unreceived messages into a message queue (rabbit mq?) for pick up later.
- Verify the math in app/util.go. HIGHLY likely it's wrong.
- Support alternative game types.
- Support authentication
- Support ping/pong on the websockets
- Make websocket processing code more resilient. Current implementation has no timeouts for stale clients