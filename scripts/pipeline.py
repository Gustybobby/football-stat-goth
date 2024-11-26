import requests
import db
import time
import match_scraper
import lineup_scraper
import player_scraper
import timeline_scraper

if __name__ == "__main__":
    PL_MATCH_ID = 115835

    client = db.supabase_connect()

    MATCH_URL = "https://www.premierleague.com/match/" + str(PL_MATCH_ID)
    MATCH_ID, HLID, ALID, HCID, ACID = match_scraper.extract_match_data(
        PL_MATCH_ID, client
    )

    for _ in range(5):
        try:
            page_source = requests.get(MATCH_URL).text
            break
        except:
            time.sleep(2)
            print("fetch failed, waiting for 2 s")

    players = lineup_scraper.scrape_lineups(page_source)

    player_scraper.insert_players(players, client)

    lineup_scraper.scrape_pitch_position(page_source, HLID, ALID, HCID, ACID, client)

    for keyword in ["home", "away"]:
        timeline_events = timeline_scraper.scrape_timeline(
            page_source,
            nm_event_class="event " + keyword,
            nm_ot_class="event " + keyword + " event--added-time",
        )

        insert_data = timeline_scraper.db_transform(timeline_events, MATCH_ID, client)
        for row in insert_data:
            print(row)

        ans = input("Insert? (Y/N/R): ")
        if ans == "Y":
            res = client.table("lineup_event").insert(insert_data).execute()
            print(res)
        elif ans == "R":
            continue

        break
