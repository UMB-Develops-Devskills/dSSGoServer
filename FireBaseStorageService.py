import firebase_admin
from firebase_admin import credentials, storage
import json
from dotenv import load_dotenv
import os

load_dotenv()

firebase_json = os.getenv("FIREBASE_CONFIG")

if not firebase_json:
    raise ValueError("Missing FIREBASE_CONFIG environment variable")

firebase_config = json.loads(firebase_json)
cred = credentials.Certificate(firebase_config)


firebase_admin.initialize_app(cred, {
       'storageBucket': 'devskillsets-2a367.firebasestorage.app'
   })


def getJobMarketData(fileName):
    bucket = storage.bucket()  
    blob = bucket.blob(fileName)

    if not blob.exists():
        raise RuntimeError(f"File not found: {fileName}")

    
    file_content = blob.download_as_text()
    
    data = json.loads(file_content)
    print("Sending data")
    return data