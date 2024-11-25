import requests
import db
import lineup_scraper
import player_scraper
import timeline_scraper

if __name__ == "__main__":
    MATCH_URL = "https://www.premierleague.com/match/115886"
    MATCH_ID = 60

    page_source = requests.get(MATCH_URL).text

    client = db.supabase_connect()

    players = lineup_scraper.scrape_lineups(page_source)

    player_scraper.insert_players(players, client)

    timeline_events = timeline_scraper.scrape_timeline(
        page_source,
        nm_event_class="event away",
        nm_ot_class="event away event--added-time",
    )

    insert_data = timeline_scraper.db_transform(timeline_events, MATCH_ID, client)
    for row in insert_data:
        print(row)

    ans = input("Insert? (Y/N): ")
    if ans == "Y":
        res = client.table("lineup_event").insert(insert_data).execute()
        print(res)
