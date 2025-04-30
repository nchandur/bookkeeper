from pymongo import MongoClient
from sentence_transformers import SentenceTransformer
from time import perf_counter
import re

client = MongoClient("mongodb://localhost:27017")
db = client["booksV2"]
collection = db["works"]

model = SentenceTransformer('all-mpnet-base-v2')
# model = SentenceTransformer('all-MiniLM-L6-v2')
# model = SentenceTransformer('BAAI/bge-large-zh-v1.5')


cursor = collection.find({})

count = 0

start = perf_counter()

for doc in cursor:
    work = doc.get('work', {})

    title = work.get('title', "")
    genres = work.get('genres', [])
    summary = work.get('summary', "")
    summary = re.sub(r'[^a-zA-Z0-9]', ' ', summary.lower())
    
    if isinstance(genres, list):
        text = "The book belongs to these genres: {}. Here is the summary of the book: {}".format(", ".join(genres), summary)
    else:
        text = "Genres unavailable for this book. Here is the summary of the book: {}".format(summary)

    vector = model.encode(text)

    collection.update_one(
        {"_id": doc["_id"]},
        {"$set": {"work.embedding": vector.tolist()}},
        )

    count += 1
    print("{} documents processed".format(count), end="\r")

print("{} documents embedded in {} seconds".format(count, round(perf_counter() - start, 2)))
client.close()