CREATE DATABASE "gomoku";
\connect "gomoku"

CREATE TABLE "players" (
  "player_id" SERIAL PRIMARY KEY,
  "nickname" VARCHAR(20) NOT NULL UNIQUE,
  "password" VARCHAR(20) NOT NULL,
  "score" INT NOT NULL DEFAULT 0
);

CREATE TYPE "game_type" AS ENUM ('pvp', 'pva');

CREATE TABLE "games" (
  "game_id" SERIAL PRIMARY KEY,
  "type" "game_type" NOT NULL,
  "board" JSONB NOT NULL,
  
  "current_player_id" INT NOT NULL REFERENCES "players" ("player_id"),
  "winner_player_id" INT NULL REFERENCES "players" ("player_id"),

  "first_player_id" INT NOT NULL REFERENCES "players" ("player_id"),
  "second_player_id" INT NULL REFERENCES "players" ("player_id"),
  
  "last_activity" timestamp NOT NULL
);

CREATE INDEX "IDX_players_player_id" ON "players" USING BTREE ("player_id");
CREATE INDEX "IDX_players_nickname" ON "players" USING BTREE ("nickname");

CREATE INDEX "IDX_games_game_id" ON "games" USING BTREE ("game_id");
CREATE INDEX "IDX_games_current_player_id" ON "games" USING BTREE ("current_player_id");
CREATE INDEX "IDX_games_winner_player_id" ON "games" USING BTREE ("winner_player_id");
