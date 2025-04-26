from pymongo import MongoClient
from langdetect import detect, DetectorFactory
import re

DetectorFactory.seed = 0

client = MongoClient("mongodb://localhost:27017")
db = client["booksV2"]
collection = db["works"]

cursor = collection.find({})

count = 0
exceptionCount = 0
docCount = 0

for doc in cursor:
    work = doc.get('work', {})

    title = work.get("title", "")
    summary = work.get('summary', "")
    ratings = work.get('ratings', "")



    if len(summary) > 0:
        summaryLang = detect(summary)
        nonEnglishFlag = summaryLang != 'en'
        ratingsFlag = ratings == -1
        boxedFlag = "box set" in title.lower() or "boxed set" in title.lower()
        libraryFlag = re.match(r"^Librarian's note", summary, re.I)

        if nonEnglishFlag or ratingsFlag or boxedFlag or libraryFlag:
            collection.delete_one({"_id": doc["_id"]})
            count += 1

    docCount += 1

    print("{} documents processed".format(docCount), end="\r")

print("{} non-English or invalid documents dropped".format(count))
client.close()