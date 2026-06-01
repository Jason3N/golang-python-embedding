from fastapi import APIRouter
from models.embeddings import EmbedRequest, EmbedResponse
from sentence_transformers import SentenceTransformer

model = SentenceTransformer("all-MiniLM-L6-v2")
router = APIRouter()

@router.post("/embed", response_model=EmbedResponse)
def embed(req: EmbedRequest):
    embedding = model.encode(req.text).tolist()
    return EmbedResponse(embedding=embedding)
