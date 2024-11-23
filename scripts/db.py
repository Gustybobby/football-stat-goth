import supabase
import os
from dotenv import load_dotenv


def supabase_connect() -> supabase.Client:
    load_dotenv()
    return supabase.create_client(os.getenv("SUPABASE_URL"), os.getenv("SUPABASE_KEY"))


def get_bucket_url() -> str:
    load_dotenv()
    bucket_url = os.getenv("SUPABASE_BUCKET_URL")
    if bucket_url is None:
        raise Exception("invalid SUPABASE_BUCKET_URL")
    return bucket_url


def find_club_id(club_name: str, client: supabase.Client) -> str:
    res = client.table("club").select("id").filter("name", "eq", club_name).execute()
    return res.data[0]["id"]


def find_match_by_id(match_id: int, client: supabase.Client):
    res = client.table("match").select("*").filter("id", "eq", match_id).execute()
    return res.data[0]


def find_player_id_by_fullname(fullname: str, client: supabase.Client) -> int:
    words = fullname.split(" ")
    if len(words) == 1:
        fn = ""
        ln = words[0]
    else:
        fn = words[0]
        ln = " ".join(words[1:])
    res = (
        client.table("player")
        .select("id")
        .filter("firstname", "eq", fn)
        .filter("lastname", "eq", ln)
    ).execute()
    try:
        return res.data[0]["id"]
    except:
        raise Exception(fullname + " is invalid")
