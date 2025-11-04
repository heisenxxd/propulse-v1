from pydantic import BaseModel, Field
from typing import Optional, List
from datetime import datetime
from uuid import UUID

class Proposta(BaseModel):
    id: UUID
    titulo: str
    nome_empresa: str = Field(alias="nomeEmpresa")
    nome_cliente: str = Field(alias="nomeCliente")
    prompt: str
    cores: List[str]
    logo: Optional[str] = None
    logo_cliente: Optional[str] = Field(default=None, alias="logoCliente")
    status: str
    arquivo_final: Optional[str] = Field(default=None, alias="arquivoFinal")
    data_criacao: datetime = Field(alias="dataCriacao")
    last_update: datetime = Field(alias="lastUpdate")
    class Config:
        populate_by_name = True