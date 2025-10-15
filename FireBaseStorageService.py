import firebase_admin
from firebase_admin import credentials, storage
import json
from dotenv import load_dotenv
import os

load_dotenv()

cred = credentials.Certificate(os.environ["GOOGLE_APPLICATION_CREDENTIALS"])
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
    return data