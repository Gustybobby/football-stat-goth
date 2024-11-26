import db
import requests
import time
import random
from bs4 import BeautifulSoup
from datetime import datetime, timezone


POS_ENUM = {
    "Forward": "FWD",
    "Midfielder": "MFD",
    "Defender": "DEF",
    "Goalkeeper": "GK",
}


def scrape_player(page_source: str, client) -> tuple[dict | None, dict | None]:
    soup = BeautifulSoup(page_source, "html.parser")

    data = {}
    club_data = {"season": "2024/25"}

    firstname_div = soup.find("div", class_="player-header__name-first")
    if firstname_div is not None:
        data["firstname"] = firstname_div.get_text().strip()
    else:
        data["firstname"] = ""

    lastname_div = soup.find("div", class_="player-header__name-last")
    if lastname_div is not None:
        data["lastname"] = lastname_div.get_text().strip()
    else:
        data["lastname"] = input("Missing Lastname, Please input: ")

    try:
        data["height"] = (
            soup.find("div", class_="player-overview__label", string="Height")
            .find_parent("div")
            .find_all("div")[1]
            .get_text()
            .replace("cm", "")
        )
    except:
        height = input("Missing Height(cm), Please input: ")
        data["height"] = height

    try:
        day, month, year = (
            soup.find("div", class_="player-overview__label", string="Date of Birth")
            .find_parent("div")
            .find_all("div")[1]
            .get_text()
            .strip()
            .split(" ")[0]
        ).split("/")

        data["dob"] = datetime(
            int(year), int(month), int(day), tzinfo=timezone.utc
        ).isoformat(timespec="milliseconds")
    except:
        data["dob"] = input("Missing DOB(ISO), Please input: ")

    if db.player_exists(data, client):
        random_sleep = random.randint(0, 1)
        print("PLAYER EXISTS, sleeping for", random_sleep, "s")
        time.sleep(random_sleep)
        return None, None

    try:
        club_data["no"] = soup.find(
            "div",
            class_="player-header__player-number player-header__player-number--large",
        ).get_text()
    except:
        no = input("Missing Player No, Please input: ")
        club_data["no"] = no

    try:
        data["nationality"] = soup.find(
            "span", class_="player-overview__player-country"
        ).get_text()
    except:
        data["nationality"] = input("Missing Nationality, Please input: ")

    try:
        data["position"] = POS_ENUM[
            (
                soup.find("div", class_="player-overview__label", string="Position")
                .find_parent("div")
                .find_all("div")[1]
                .get_text()
            )
        ]
    except:
        data["position"] = input("Missing Position, Please input: ")

    try:
        club_div = soup.find("div", class_="player-overview__label", string="Club")
        if club_div is None:
            club_name = "not_found"
            raise "club not found"
        else:
            club_name = (
                club_div.find_parent("div").find_all("div")[1].get_text().strip()
            )
        club_data["club_id"] = db.find_club_id(club_name, client)
    except:
        print("found club name:", "'" + club_name + "'")
        club_id = input("Missing Club ID, Please input (or NULL): ")
        club_data["club_id"] = club_id

    data["image"] = (
        db.get_bucket_url()
        + "/"
        + replace_special(
            (
                data["lastname"]
                if data["firstname"] == ""
                else "_".join(
                    [
                        data["firstname"].replace(" ", "_"),
                        data["lastname"].replace(" ", "_"),
                    ]
                )
            ).lower()
        )
        + ".webp"
    )

    return data, club_data


def replace_special(string: str):
    return (
        string.replace("é", "e")
        .replace("ë", "e")
        .replace("í", "i")
        .replace("î", "i")
        .replace("ï", "i")
        .replace("ø", "o")
        .replace("ö", "o")
        .replace("ã", "a")
        .replace("á", "a")
        .replace("ñ", "n")
        .replace("ú", "u")
        .replace("ü", "u")
    )


def insert_players(player_urls: list, db_client) -> None:
    for profile_url in player_urls:
        print("Reading", profile_url)

        trial = 0
        while trial <= 5:
            try:
                source = requests.get(profile_url).text
                break
            except:
                fetch_sleep = 2
                print("fetch failed, retrying in", fetch_sleep, "s")
                time.sleep(fetch_sleep)
                trial += 5

        player, club_player = scrape_player(source, db_client)

        if player is None:
            print("====================================")
            continue

        print("===============PLAYER===============")
        for key in player:
            print(key, ":", player[key])
        print("=============CLUB_PLAYER============")
        for key in club_player:
            print(key, ":", club_player[key])
        print("====================================")

        ans = input("Insert? (Y/N): ")
        if ans == "Y":
            res = db_client.table("player").insert(player).execute()
            print(res)
            club_player["player_id"] = res.data[0]["id"]
            res = db_client.table("club_player").insert(club_player).execute()
            print(res)

        print("====================================")


if __name__ == "__main__":
    PLAYER_PROFILE_URLS = []
    client = db.supabase_connect()
    insert_players(PLAYER_PROFILE_URLS, client)
