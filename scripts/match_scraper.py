import db
import requests
import time
import datetime


def get_match_data(pl_match_id: int) -> dict:
    for _ in range(5):
        try:
            data = requests.get(
                "https://footballapi.pulselive.com/football/stats/match/"
                + str(pl_match_id),
                headers={
                    "content-type": "application/x-www-form-urlencoded; charset=UTF-8",
                    "origin": "https://www.premierleague.com",
                    "referer": "https://www.premierleague.com",
                },
            )
            break
        except:
            print("fetch failed, waiting for 2 s")
            time.sleep(2)
    return data.json()


STATS_MAP = {
    "possession_percentage": "possession",
    "ontarget_scoring_att": "shots_on_target",
    "total_scoring_att": "shots",
    "touches": "touches",
    "total_pass": "passes",
    "total_tackle": "tackles",
    "total_clearance": "clearances",
    "corner_taken": "corners",
    "total_offside": "offsides",
    "fk_foul_lost": "fouls_conceded",
}


def extract_match_data(pl_match_id: int, db_client):
    api_data = get_match_data(pl_match_id)

    entity = api_data["entity"]

    gameweek = entity["gameweek"]

    week = gameweek["gameweek"]

    season = gameweek["compSeason"]["label"]

    kickoff = entity["kickoff"]
    seconds = kickoff["millis"] / 1000
    start_at = datetime.datetime.fromtimestamp(seconds, tz=datetime.timezone.utc)

    location = entity["ground"]["name"] + ", " + entity["ground"]["city"]

    home, away = entity["teams"]

    home_club_id = home["team"]["club"]["abbr"]
    away_club_id = away["team"]["club"]["abbr"]

    home_data_id = home["team"]["id"]
    away_data_id = away["team"]["id"]

    stats_data = api_data["data"]

    home_stats = stats_data[str(home_data_id)]["M"]
    away_stats = stats_data[str(away_data_id)]["M"]

    home_lineup = {"club_id": home_club_id}
    away_lineup = {"club_id": away_club_id}

    for stat in home_stats:
        if stat["name"] in STATS_MAP:
            key = STATS_MAP[stat["name"]]
            home_lineup[key] = (
                stat["value"] if key == "possession" else int(stat["value"])
            )

    for stat in away_stats:
        if stat["name"] in STATS_MAP:
            key = STATS_MAP[stat["name"]]
            away_lineup[key] = (
                stat["value"] if key == "possession" else int(stat["value"])
            )

    for key in STATS_MAP.values():
        home_lineup[key] = home_lineup[key] if key in home_lineup else 0
        away_lineup[key] = away_lineup[key] if key in away_lineup else 0

    print("HOME", "===========================", "AWAY", "====")
    for key in home_lineup:
        print(
            home_lineup[key],
            "\t",
            key.ljust(20),
            "\t",
            away_lineup[key],
        )
    print("==========================================")

    ans = input("Insert? (Y/N): ")
    if ans == "Y":
        res = db_client.table("lineup").insert([home_lineup, away_lineup]).execute()
        print(res)
    else:
        raise Exception("insert cancelled")
    print("==========================================")

    insert_match_data = {
        "week": week,
        "season": season,
        "start_at": start_at.isoformat(),
        "location": location,
        "is_finished": start_at < datetime.datetime.now(tz=datetime.timezone.utc),
        "home_lineup_id": res.data[0]["id"],
        "away_lineup_id": res.data[1]["id"],
    }
    for key in insert_match_data:
        print(key, ":", insert_match_data[key])
    print("==========================================")

    ans = input("Insert? (Y/N): ")
    if ans == "Y":
        res = db_client.table("match").insert(insert_match_data).execute()
        print(res)
    else:
        raise Exception("insert cancelled")
    print("==========================================")

    return (
        res.data[0]["id"],
        insert_match_data["home_lineup_id"],
        insert_match_data["away_lineup_id"],
        home_club_id,
        away_club_id,
    )


if __name__ == "__main__":
    client = db.supabase_connect()

    extract_match_data(115836, client)
