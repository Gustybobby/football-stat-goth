from bs4 import BeautifulSoup


def read_html(file_dir: str) -> str:
    with open(file_dir, "r") as file:
        html = file.read().replace("\n", "")
        file.close()
        return html


def scrape_lineups(page_source: str):
    lineups = []

    soup = BeautifulSoup(page_source, "html.parser")

    player_lis = soup.find_all("li", class_="player")

    for player in player_lis:
        player = player.find("a")
        lineups.append("https://" + "/".join(player.attrs["href"][2:].split("/")[:3]))

    return lineups


if __name__ == "__main__":
    players = scrape_lineups(read_html("pages/115889.html"))
    print(players)
