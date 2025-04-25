from pymongo import MongoClient

client = MongoClient("mongodb://localhost:27017/")
db = client["booksV2"]
collection = db["works"]

pipeline = [
    {"$group": {
        "_id": "$work.title",
        "ids": {"$addToSet": "$_id"},
        "count": {"$sum": 1}
    }},
    {"$match": {
        "count": {"$gt": 1}
    }}
]

duplicates = collection.aggregate(pipeline)

ids_to_delete = []

for doc in duplicates:
    ids = doc["ids"]
    ids_to_delete.extend(ids[1:])

if ids_to_delete:
    result = collection.delete_many({"_id": {"$in": ids_to_delete}})

print("{} duplicate documents dropped".format(result.deleted_count))
client.close()