from fileinput import filename
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from FireBaseStorageService import getJobMarketData

app = FastAPI(
    title="DSSFastAPIServerVercel",
    description="DSSFastAPIServerVercel",
    version="10.13.2025",
)

# 👇 CORS setup here
origins = [
    "http://localhost:3000",
    "https://devskillsets.com"
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/health")
async def root():
    return {"status": "DSS server OK"}

@app.get("/api/getJobData/{fileName}")
async def getJobData(fileName):
    print("getJobData Hit: "+ fileName)
    return getJobMarketData(fileName)



