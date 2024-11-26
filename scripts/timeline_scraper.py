import db
import requests
from bs4 import BeautifulSoup

EVENT_TYPE_MAP = {
    "Goal": "GOAL",
    "label.penalty.scored": "GOAL",
    "Yellow Card": "YELLOW",
    "Second Yellow Card (Red Card)": "RED",
    "Red Card": "RED",
    "Substitution": "SUB",
    "Own Goal": "OWN_GOAL",
}


def scrape_timeline(
    page_source: str,
    nm_event_class="event home",
    nm_ot_class="event home event--added-time",
) -> list:
    soup = BeautifulSoup(page_source, "html.parser")

    timeline_div = soup.find(
        "div", class_="eventLine timeLineEventsContainer is-completed"
    )

    events = []

    # normal events (no extra time)
    nm_events = timeline_div.find_all("div", class_=nm_event_class)
    for event_div in nm_events:
        event_data = {}

        event_data["minutes"] = event_div.find(
            "div", class_="event__minute"
        ).text.strip()

        home_events = event_div.find(
            "ul", class_="event__icons event__icons--home"
        ).find_all("li")
        event_data["home"] = extract_events(home_events)

        away_events = event_div.find(
            "ul", class_="event__icons event__icons--away"
        ).find_all("li")
        event_data["away"] = extract_events(away_events)

        event_data["after_half"] = not check_after_half(event_div)
        events.append(event_data)

    # overtime events
    ot_events = timeline_div.find_all("div", class_=nm_ot_class)
    for ot_div in ot_events:
        event_data = {}

        try:
            minutes, extra = ot_div.find("time", class_="min").text.strip().split(" ")
        except:
            minutes, extra = (
                ot_div.find("div", class_="event__minute")
                .text.strip()
                .replace(" ", "")
                .replace("\n", "")
                .split("+")
            )

        event_data["minutes"] = minutes
        event_data["extra"] = extra.replace("'", "").replace("+", "")

        home_events = ot_div.find(
            "ul", class_="event__icons event__icons--home"
        ).find_all("li")
        event_data["home"] = extract_events(home_events)

        away_events = ot_div.find(
            "ul", class_="event__icons event__icons--away"
        ).find_all("li")
        event_data["away"] = extract_events(away_events)

        event_data["after_half"] = not check_after_half(event_div)
        events.append(event_data)

    return events


def extract_events(event_tags: list) -> list:
    events = []
    for event_tag in event_tags:
        event_data = {}

        info_header = event_tag.find("header", class_="eventInfoHeader")
        if info_header is None:
            continue

        event_type = EVENT_TYPE_MAP[
            (info_header.find("span", class_="visuallyHidden")).text.strip()
        ]
        event_data["event"] = event_type

        player_info = event_tag.find("div", class_="eventPlayerInfo")

        if (
            event_type == "GOAL"
            or event_type == "YELLOW"
            or event_type == "RED"
            or event_type == "OWN_GOAL"
        ):
            scorer = player_info.find("a", class_="name").text.split(".")[1].strip()
            event_data["player1"] = scorer

            assist = player_info.find("div", class_="assist")
            if assist is not None:
                event_data["player2"] = assist.text.strip().replace("Ast. ", "")
        if event_type == "SUB":
            sub_off = event_tag.find("div", class_="eventInfoContent")
            off_player_info = sub_off.find("div", class_="eventPlayerInfo")
            event_data["player1"] = (
                off_player_info.find("a", class_="name")
                .text.split("\n")[0]
                .split(".")[1]
                .strip()
            )

            sub_on = event_tag.find("div", class_="eventInfoContent subOn")
            on_player_info = sub_on.find("div", class_="eventPlayerInfo")
            event_data["player2"] = (
                on_player_info.find("a", class_="name")
                .text.split("\n")[0]
                .split(".")[1]
                .strip()
            )

        events.append(event_data)

    return events


def check_after_half(tag):
    while tag.next:
        tag = tag.next
        try:
            if " ".join(tag.attrs["class"]) == "event ht":
                return True
        except:
            pass
    return False


def db_transform(events, match_id, db_client):
    match = db.find_match_by_id(match_id, db_client)
    rows = []
    for tl_event in events:
        extra = int(tl_event["extra"]) if "extra" in tl_event else None
        after_half = tl_event["after_half"]
        if int(tl_event["minutes"]) == 45 and extra is not None:
            after_half = False

        for home_event in tl_event["home"]:
            rows.append(
                {
                    "lineup_id": (
                        match["home_lineup_id"]
                        if home_event["event"] != "OWN_GOAL"
                        else match["away_lineup_id"]
                    ),
                    "minutes": int(tl_event["minutes"]),
                    "extra": extra,
                    "event": home_event["event"],
                    "player_id1": db.find_player_id_by_fullname(
                        home_event["player1"], db_client
                    ),
                    "player_id2": (
                        db.find_player_id_by_fullname(home_event["player2"], db_client)
                        if "player2" in home_event
                        else None
                    ),
                    "after_half": after_half,
                }
            )

        for away_event in tl_event["away"]:
            rows.append(
                {
                    "lineup_id": (
                        match["away_lineup_id"]
                        if away_event["event"] != "OWN_GOAL"
                        else match["home_lineup_id"]
                    ),
                    "minutes": int(tl_event["minutes"]),
                    "extra": extra,
                    "event": away_event["event"],
                    "player_id1": db.find_player_id_by_fullname(
                        away_event["player1"], db_client
                    ),
                    "player_id2": (
                        db.find_player_id_by_fullname(away_event["player2"], db_client)
                        if "player2" in away_event
                        else None
                    ),
                    "after_half": after_half,
                }
            )
    return rows


if __name__ == "__main__":
    MATCH_URL = "https://www.premierleague.com/match/115899"
    MATCH_ID = 13

    timeline_events = scrape_timeline(requests.get(MATCH_URL).text)

    client = db.supabase_connect()

    insert_data = db_transform(timeline_events, MATCH_ID, client)
    for row in insert_data:
        print(row)

    ans = input("Insert? (Y/N): ")
    if ans == "Y":
        res = client.table("lineup_event").insert(insert_data).execute()
        print(res)
