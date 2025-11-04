from fastapi import FastAPI, HTTPException, BackgroundTasks
from fastapi.responses import FileResponse
import os
import traceback

from model.proposta import Proposta as PropostaModel
from src.ia_generator.ia import gerar_html_proposta
from src.ia_generator.pdf_generator import converter_html_para_pdf

app = FastAPI()

@app.post("/gerarproposta/pdf_dynamic")
async def criar_proposta_pdf_dinamica(
    proposta: PropostaModel,
    background_tasks: BackgroundTasks
):
    
    caminho_pdf = None
    try:
        html_gerado = await gerar_html_proposta(proposta)

        caminho_pdf = await converter_html_para_pdf(html_gerado)
        
        background_tasks.add_task(os.remove, caminho_pdf)

        return FileResponse(
            path=caminho_pdf,
            media_type='application/pdf',
            filename=f"proposta_{proposta.nome_empresa}.pdf"
        )
        
    except Exception as e:
        traceback.print_exc() 
        if caminho_pdf and os.path.exists(caminho_pdf):
            os.remove(caminho_pdf)
            
        raise HTTPException(status_code=500, detail=f"Erro ao gerar proposta: {str(e)}")