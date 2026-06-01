from fastapi import FastAPI
from pydantic import BaseModel
from sentence_transformers import SentenceTransformer
from routes.embeddings import router as embedding_router

app = FastAPI()
app.include_router(embedding_router)

@app.get("/health")
def health():
    return {"status": "ok"}