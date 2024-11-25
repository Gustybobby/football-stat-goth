import db
import requests
from bs4 import BeautifulSoup


def scrape_lineups(page_source: str):
    lineups = []

    soup = BeautifulSoup(page_source, "html.parser")

    player_lis = soup.find_all("li", class_="player")

    for player in player_lis:
        player = player.find("a")
        lineups.append(
            "https://www.premierleague.com"
            + "/".join(player.attrs["href"].split("/")[:3])
        )

    return lineups


def scrape_pitch_position(
    page_source: str,
    home_lineup_id: int,
    away_lineup_id: int,
    home_club_id: str,
    away_club_id: str,
    db_client,
):
    soup = BeautifulSoup(page_source, "html.parser")

    home_lineup = soup.find(
        "div", class_="teamList mcLineUpContainter homeLineup"
    ).find("div", class_="matchLineupTeamContainer")

    away_lineup = soup.find(
        "div", class_="teamList mcLineUpContainter awayLineup"
    ).find("div", class_="matchLineupTeamContainer")

    home_pos_map, home_subs = scrape_lineup_position(home_lineup, True)
    away_pos_map, away_subs = scrape_lineup_position(away_lineup, False)

    home_pitch = soup.find("div", class_="team home pitchPositonsContainer")
    rows = home_pitch.find_all("div", class_="row")

    lineup_players = []

    for i, row in enumerate(rows):
        pins = row.find_all("div")
        for j, pin in enumerate(pins):
            lineup_players.append(
                {
                    "lineup_id": home_lineup_id,
                    "player_id": db.find_player_id_by_club_no(
                        int(pin.text), home_club_id, "2024/25", db_client
                    ),
                    "position_no": i * 10 + j,
                    "position": home_pos_map[pin.text],
                }
            )

    for i, sub in enumerate(home_subs):
        lineup_players.append(
            {
                "lineup_id": home_lineup_id,
                "player_id": db.find_player_id_by_club_no(
                    sub, home_club_id, "2024/25", db_client
                ),
                "position_no": 100 + i,
                "position": "SUB",
            }
        )

    away_pitch = soup.find("div", class_="team away pitchPositonsContainer")
    rows = away_pitch.find_all("div", class_="row")

    for i, row in enumerate(rows):
        pins = row.find_all("div")
        for j, pin in enumerate(pins):
            lineup_players.append(
                {
                    "lineup_id": away_lineup_id,
                    "player_id": db.find_player_id_by_club_no(
                        int(pin.text), away_club_id, "2024/25", db_client
                    ),
                    "position_no": i * 10 + j,
                    "position": away_pos_map[pin.text],
                }
            )

    for i, sub in enumerate(away_subs):
        lineup_players.append(
            {
                "lineup_id": away_lineup_id,
                "player_id": db.find_player_id_by_club_no(
                    sub, away_club_id, "2024/25", db_client
                ),
                "position_no": 100 + i,
                "position": "SUB",
            }
        )

    res = db_client.table("lineup_player").insert(lineup_players).execute()
    print("inserted lineup players into their positions", res)
    print("====================================================")


def scrape_lineup_position(container, is_home: bool):
    con_class = (
        "startingLineUpContainer squadList home"
        if is_home
        else "startingLineUpContainer squadList"
    )

    POS_ENUM = {
        "Forward": "FWD",
        "Midfielder": "MFD",
        "Defender": "DEF",
        "Goalkeeper": "GK",
    }

    no_pos_map = {}
    count = 0

    SUB = {"GK": [], "DEF": [], "MFD": [], "FWD": []}

    for ul in container.find_all("ul", class_=con_class):
        for li in ul.find_all("li"):
            number = li.find("div", class_="number").text
            position = li.find("span", class_="position").find("span").text
            if count < 11:
                no_pos_map[number] = POS_ENUM[position]
            else:
                SUB[POS_ENUM[position]].append(int(number))

            count += 1
    subs = (
        sorted(SUB["GK"]) + sorted(SUB["DEF"]) + sorted(SUB["MFD"]) + sorted(SUB["FWD"])
    )
    return (no_pos_map, subs)
