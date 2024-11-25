CREATE TYPE "user_role" AS ENUM ('USER','ADMIN');
CREATE TYPE "player_position" AS ENUM ('GK','DEF','MFD','FWD','SUB');
CREATE TYPE "event_type" AS ENUM ('GOAL','OWN_GOAL','YELLOW','RED','SUB');

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
    position_no     INT2 NOT NULL,
    position        player_position NOT NULL,

    CONSTRAINT pk_lineup_player                 PRIMARY KEY (lineup_id, player_id),
    CONSTRAINT unique_lineup_id_position_no     UNIQUE (lineup_id, position_no),
    CONSTRAINT fk_lineup_player_lineup          FOREIGN KEY (lineup_id) REFERENCES "lineup"(id),
    CONSTRAINT fk_lineup_player_player          FOREIGN KEY (player_id) REFERENCES "player"(id)
);

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
    cost        INTEGER NOT NULL,
    points      INTEGER,
    rating      INTEGER,
    
    CONSTRAINT fk_fantasy_player_club      FOREIGN KEY (club_id) REFERENCES "club"(id),
    CONSTRAINT fk_fantasy_player_player    FOREIGN KEY (player_id) REFERENCES "player"(id)
);