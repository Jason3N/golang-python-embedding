package models

import (
	"github.com/pgvector/pgvector-go"
)

type StoredEmbeddings struct {
	ID 		  int			  `json:"id"`		
	Content   string		  `json:"content"` 	
	Embedding pgvector.Vector `json:"embedding"`
}