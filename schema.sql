CREATE TYPE "user_role" AS ENUM ('USER','ADMIN');
CREATE TYPE "player_position" AS ENUM ('GK','DEF','MFD','FWD','SUB');

CREATE TABLE "club" (
    id        CHAR(3) PRIMARY KEY,
    name      VARCHAR(64) NOT NULL,
    stadium   VARCHAR(64),
    logo      VARCHAR(255) NOT NULL,
    est       INT4 NOT NULL
);

CREATE TABLE "player" (
    id          SERIAL PRIMARY KEY,
    club_id     CHAR(3),
    no          INT2 NOT NULL,
    firstname   VARCHAR(64) NOT NULL,
    lastname    VARCHAR(64) NOT NULL,
    dob         TIMESTAMP NOT NULL,
    height      INT2 NOT NULL,
    nationality VARCHAR(64) NOT NULL,
    position    player_position NOT NULL,
    image       VARCHAR(255),

    CONSTRAINT unique_club_id_no UNIQUE (club_id, no),
    CONSTRAINT fk_player_club FOREIGN KEY (club_id) REFERENCES "club"(id)
);

CREATE TABLE "lineup" (
    id                  SERIAL PRIMARY KEY,
    club_id             CHAR(3) NOT NULL,
    goals               INT2 NOT NULL DEFAULT 0,
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
    lineup_id INTEGER NOT NULL,
    player_id INTEGER NOT NULL,
    position_no INT2 NOT NULL,
    goals INT2 NOT NULL DEFAULT 0,
    yellow_cards INT2 NOT NULL DEFAULT 0,
    red_cards INT2 NOT NULL DEFAULT 0,

    CONSTRAINT pk_lineup_player PRIMARY KEY (lineup_id, player_id),
    CONSTRAINT fk_lineup_player_lineup FOREIGN KEY (lineup_id) REFERENCES "lineup"(id),
    CONSTRAINT fk_lineup_player_player FOREIGN KEY (player_id) REFERENCES "player"(id)
);

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
    token VARCHAR(64) PRIMARY KEY,
    username VARCHAR(128) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_session_user FOREIGN KEY (username) REFERENCES "user"(username)
);