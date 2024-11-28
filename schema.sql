CREATE TYPE "user_role" AS ENUM ('USER','ADMIN');
CREATE TYPE "player_position" AS ENUM ('GK','DEF','MFD','FWD','SUB');
CREATE TYPE "event_type" AS ENUM ('GOAL','OWN_GOAL','YELLOW','RED','SUB');
CREATE TYPE "fantasy_transaction_type" AS ENUM ('BUY','SELL');

CREATE TABLE "club" (
    id          CHAR(3) PRIMARY KEY,
    name        VARCHAR(64) NOT NULL,
    short_name  VARCHAR(32) NOT NULL,
    stadium     VARCHAR(64),
    logo        VARCHAR(255) NOT NULL,
    est         INT4 NOT NULL
);

CREATE TABLE "player" (
    id          SERIAL PRIMARY KEY,
    firstname   VARCHAR(64) NOT NULL,
    lastname    VARCHAR(64) NOT NULL,
    dob         TIMESTAMP NOT NULL,
    height      INT2 NOT NULL,
    nationality VARCHAR(64) NOT NULL,
    position    player_position NOT NULL,
    image       VARCHAR(255)
);

CREATE TABLE "club_player" (
    club_id     CHAR(3),
    player_id   INTEGER,
    season      VARCHAR(16),
    no          INT2 NOT NULL,

    CONSTRAINT pk_club_player           PRIMARY KEY (club_id, player_id, season),  
    CONSTRAINT fk_club_player_club      FOREIGN KEY (club_id) REFERENCES "club"(id),
    CONSTRAINT fk_club_player_player    FOREIGN KEY (player_id) REFERENCES "player"(id)
);

CREATE TABLE "lineup" (
    id                  SERIAL PRIMARY KEY,
    club_id             CHAR(3) NOT NULL,
    possession          NUMERIC(4,1) NOT NULL DEFAULT 0,
    shots_on_target     INT2 NOT NULL DEFAULT 0,
    shots               INT2 NOT NULL DEFAULT 0,
    touches             INT2 NOT NULL DEFAULT 0,
    passes              INT2 NOT NULL DEFAULT 0,
    tackles             INT2 NOT NULL DEFAULT 0,
    clearances          INT2 NOT NULL DEFAULT 0,
    corners             INT2 NOT NULL DEFAULT 0,
    offsides            INT2 NOT NULL DEFAULT 0,
    fouls_conceded      INT2 NOT NULL DEFAULT 0,

    CONSTRAINT fk_lineup_club FOREIGN KEY (club_id) REFERENCES "club"(id)
);

CREATE TABLE "lineup_player" (
    lineup_id       INTEGER NOT NULL,
    player_id       INTEGER NOT NULL,
    no              INT2 NOT NULL,
    position_no     INT2 NOT NULL,
    position        player_position NOT NULL,

    CONSTRAINT pk_lineup_player                 PRIMARY KEY (lineup_id, player_id),
    CONSTRAINT unique_lineup_id_position_no     UNIQUE (lineup_id, position_no),
    CONSTRAINT unique_lineup_id_no              UNIQUE (lineup_id, no),
    CONSTRAINT fk_lineup_player_lineup          FOREIGN KEY (lineup_id) REFERENCES "lineup"(id),
    CONSTRAINT fk_lineup_player_player          FOREIGN KEY (player_id) REFERENCES "player"(id)
);

CREATE INDEX "idx_lineup_player_player_id" ON "lineup_player" USING BTREE(player_id);

CREATE TABLE "lineup_event" (
    id          SERIAL PRIMARY KEY,
    lineup_id   INTEGER NOT NULL,
    player_id1  INTEGER,
    player_id2  INTEGER,
    event       event_type NOT NULL,
    minutes     INT2 NOT NULL,
    extra       INT2,
    after_half  BOOLEAN NOT NULL DEFAULT false,

    CONSTRAINT fk_lineup_player1 FOREIGN KEY (lineup_id,player_id1) REFERENCES "lineup_player"(lineup_id,player_id),
    CONSTRAINT fk_lineup_player2 FOREIGN KEY (lineup_id,player_id2) REFERENCES "lineup_player"(lineup_id,player_id)
);

CREATE INDEX "idx_lineup_event_lineup_id" ON "lineup_event" USING BTREE(lineup_id);

CREATE TABLE "match" (
    id              SERIAL PRIMARY KEY,
    home_lineup_id  INTEGER NOT NULL UNIQUE,
    away_lineup_id  INTEGER NOT NULL UNIQUE,
    season          VARCHAR(16) NOT NULL,
    week            INT2 NOT NULL,
    location        VARCHAR(64) NOT NULL,
    start_at        TIMESTAMP NOT NULL,
    is_finished     BOOLEAN NOT NULL DEFAULT false,

    CONSTRAINT fk_match_home_lineup FOREIGN KEY (home_lineup_id) REFERENCES "lineup"(id),
    CONSTRAINT fk_match_away_lineup FOREIGN KEY (away_lineup_id) REFERENCES "lineup"(id)
);

CREATE TABLE "user" (
    username        VARCHAR(128) PRIMARY KEY,
    password_hash   TEXT NOT NULL,
    firstname       VARCHAR(128) NOT NULL,
    lastname        VARCHAR(128) NOT NULL,
    role            user_role NOT NULL DEFAULT 'USER'
);

CREATE TABLE "session" (
    token       VARCHAR(64) PRIMARY KEY,
    username    VARCHAR(128) NOT NULL,
    expires_at  TIMESTAMP NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_session_user FOREIGN KEY (username) REFERENCES "user"(username)
);

CREATE TABLE "fantasy_player" (
    id          SERIAL PRIMARY KEY,
    club_id     CHAR(3),
    player_id   INTEGER,
    
    CONSTRAINT fk_fantasy_player_club      FOREIGN KEY (club_id) REFERENCES "club"(id),
    CONSTRAINT fk_fantasy_player_player    FOREIGN KEY (player_id) REFERENCES "player"(id)
);

CREATE TABLE "fantasy_team" (
    id          SERIAL PRIMARY KEY,
    username    VARCHAR(128) NOT NULL,
    season      VARCHAR(16) NOT NULL,
    budget      INT4 NOT NULL,

    CONSTRAINT unique_username_season   UNIQUE (username, season),
    CONSTRAINT fk_fantasy_team_user     FOREIGN KEY (username) REFERENCES "user"(username)
);

CREATE TABLE "fantasy_team_player" (
    fantasy_team_id     INTEGER,
    fantasy_player_id   INTEGER,

    CONSTRAINT pk_fantasy_team_player                   PRIMARY KEY (fantasy_team_id, fantasy_player_id),
    CONSTRAINT fk_fantasy_team_player_fantasy_team      FOREIGN KEY (fantasy_team_id) REFERENCES "fantasy_team"(id),
    CONSTRAINT fk_fantasy_team_player_fantasy_player    FOREIGN KEY (fantasy_player_id) REFERENCES "fantasy_player"(id)
);

CREATE TABLE "fantasy_transaction" (
    id                  SERIAL PRIMARY KEY,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    cost                INT4 NOT NULL,
    type                fantasy_transaction_type NOT NULL,
    fantasy_team_id     INTEGER NOT NULL,
    fantasy_player_id   INTEGER NOT NULL,

    CONSTRAINT fk_fantasy_team_player_fantasy_team      FOREIGN KEY (fantasy_team_id) REFERENCES "fantasy_team"(id),
    CONSTRAINT fk_fantasy_team_player_fantasy_player    FOREIGN KEY (fantasy_player_id) REFERENCES "fantasy_player"(id)
);

CREATE OR REPLACE FUNCTION fantasy_player_team_transaction()
RETURNS TRIGGER
AS $$
DECLARE
    current_budget INT;
BEGIN
    IF NEW.type = 'BUY' THEN
        SELECT
            "fantasy_team".budget
        FROM "fantasy_team"
        WHERE "fantasy_team".id = NEW.fantasy_team_id
        INTO current_budget;

        IF current_budget < NEW.cost THEN
            RAISE EXCEPTION 'budget is not sufficient';
        END IF;

        UPDATE "fantasy_team"
        SET budget = budget - NEW.cost
        WHERE "fantasy_team".id = NEW.fantasy_team_id;

        INSERT INTO "fantasy_team_player" (
            fantasy_team_id,
            fantasy_player_id
        ) VALUES (
            NEW.fantasy_team_id,
            NEW.fantasy_player_id
        );
    ELSEIF NEW.type = 'SELL' THEN
        UPDATE "fantasy_team"
        SET budget = budget + NEW.cost
        WHERE "fantasy_team".id = NEW.fantasy_team_id;

        DELETE FROM "fantasy_team_player"
        WHERE
            "fantasy_team_player".fantasy_team_id = NEW.fantasy_team_id AND
            "fantasy_team_player".fantasy_player_id = NEW.fantasy_player_id;
    ELSE
        RAISE EXCEPTION 'invalid transaction type';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER fantasy_player_team_transaction
BEFORE INSERT
ON fantasy_transaction
FOR EACH ROW
EXECUTE FUNCTION fantasy_player_team_transaction();