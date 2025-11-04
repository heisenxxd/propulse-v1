from pydantic import BaseModel, Field
from typing import List

class PropostaConteudo(BaseModel):
    introducao: str = Field(description="Um parágrafo de introdução caloroso e focado no cliente.")
    desafio_cliente: str = Field(description="Um resumo do desafio ou da dor que o cliente enfrenta.")
    nossa_solucao: str = Field(description="Descrição de como sua solução resolve o desafio.")
    escopo_servicos: List[str] = Field(description="Uma lista de entregáveis ou itens de serviço claros.")
    proximos_passos: str = Field(description="Chamada para ação clara, ex: 'Agendar uma reunião'.")