from PIL import Image
import os
from dotenv import load_dotenv

load_dotenv()

DIR = os.getenv("PNG_DIR")
SAVE_DIR = "./downloads"

for file in os.listdir(DIR):
    image = Image.open(DIR + "/" + file)
    image.save(
        SAVE_DIR + "/" + file.split(".")[0] + ".webp", "webp", optimize=True, quality=80
    )
