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
