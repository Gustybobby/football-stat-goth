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


def scrape_player(page_source: str, client) -> dict:
    soup = BeautifulSoup(page_source, "html.parser")

    data = {}

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
        data["no"] = soup.find(
            "div",
            class_="player-header__player-number player-header__player-number--large",
        ).get_text()
    except:
        no = input("Missing Player No, Please input: ")
        data["no"] = no

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

    data["height"] = (
        soup.find("div", class_="player-overview__label", string="Height")
        .find_parent("div")
        .find_all("div")[1]
        .get_text()
        .replace("cm", "")
    )

    try:
        club_div = soup.find("div", class_="player-overview__label", string="Club")
        if club_div is None:
            club_name = "not_found"
            raise "club not found"
        else:
            club_name = (
                club_div.find_parent("div").find_all("div")[1].get_text().strip()
            )
        data["club_id"] = db.find_club_id(club_name, client)
    except:
        print("found club name:", "'" + club_name + "'")
        club_id = input("Missing Club ID, Please input: ")
        data["club_id"] = club_id

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

    return data


def replace_special(string: str):
    return (
        string.replace("é", "e")
        .replace("ë", "e")
        .replace("í", "i")
        .replace("ø", "o")
        .replace("ö", "o")
        .replace("ã", "a")
        .replace("ñ", "n")
        .replace("ú", "u")
        .replace("ü", "u")
    )


PLAYER_PROFILE_URLS = []


if __name__ == "__main__":
    client = db.supabase_connect()

    for profile_url in PLAYER_PROFILE_URLS:
        player = scrape_player(requests.get(profile_url).text, client)

        print("===============PLAYER===============")
        for key in player:
            print(key, ":", player[key])
        print("====================================")

        ans = input("Insert? (Y/N): ")

        if ans == "Y":
            res = client.table("player").insert(player).execute()
            print(res)
