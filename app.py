from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from FireBaseStorageService import getJobMarketData
from services.request_logger import log_request, get_request_stats

app = FastAPI(
    title="DSSFastAPIServerVercel",
    description="DSSFastAPIServerVercel",
    version="10.13.2025",
)

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
async def getJobData(fileName: str, request: Request):
    log_request(
        path=str(request.url.path),
        method=request.method,
        params={"fileName": fileName}
    )
    print(f"getJobData Hit: {fileName}")
    return getJobMarketData(fileName)


@app.get("/api/stats")
async def get_stats():
    """Expose the collected request statistics."""
    return get_request_stats()
