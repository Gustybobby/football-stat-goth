import db
import requests
from bs4 import BeautifulSoup
from datetime import datetime, timezone


POS_ENUM = {
    "Forward": "FWD",
    "Midfielder": "MFD",
    "Defender": "DEF",
    "Goalkeeper": "GK",
}


def scrape_player(page_source: str, client) -> tuple[dict, dict]:
    soup = BeautifulSoup(page_source, "html.parser")

    data = {}
    club_data = {"season": "2024/25"}

    firstname_div = soup.find("div", class_="player-header__name-first")
    if firstname_div is not None:
        data["firstname"] = (
            soup.find("div", class_="player-header__name-first").get_text().strip()
        )
    else:
        data["firstname"] = ""

    data["lastname"] = (
        soup.find("div", class_="player-header__name-last").get_text().strip()
    )

    try:
        club_data["no"] = soup.find(
            "div",
            class_="player-header__player-number player-header__player-number--large",
        ).get_text()
    except:
        no = input("Missing Player No, Please input: ")
        club_data["no"] = no

    data["nationality"] = soup.find(
        "span", class_="player-overview__player-country"
    ).get_text()

    data["position"] = POS_ENUM[
        (
            soup.find("div", class_="player-overview__label", string="Position")
            .find_parent("div")
            .find_all("div")[1]
            .get_text()
        )
    ]

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


PLAYER_PROFILE_URLS = []


if __name__ == "__main__":
    client = db.supabase_connect()

    for profile_url in PLAYER_PROFILE_URLS:
        print("reading", profile_url)
        player, club_player = scrape_player(requests.get(profile_url).text, client)

        print("===============PLAYER===============")
        for key in player:
            print(key, ":", player[key])
        print("=============CLUB_PLAYER============")
        for key in club_player:
            print(key, ":", club_player[key])
        print("====================================")

        ans = input("Insert? (Y/N): ")

        if ans == "Y":
            res = client.table("player").insert(player).execute()
            print(res)
            club_player["player_id"] = res.data[0]["id"]
            res = client.table("club_player").insert(club_player).execute()
            print(res)
