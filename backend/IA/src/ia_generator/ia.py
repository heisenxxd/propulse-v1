import os
from dotenv import load_dotenv
from pydantic import SecretStr
from langchain_openai import ChatOpenAI
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.output_parsers import PydanticOutputParser
from model.proposta import Proposta as PropostaModel
from langchain_core.output_parsers import StrOutputParser
from .ai_models import PropostaConteudo
from pathlib import Path

load_dotenv()

api_key = os.getenv("OPENROUTER_API_KEY")
base_url = os.getenv("OPENROUTER_API_BASE")

if not api_key:
    raise ValueError("OPENROUTER_API_KEY não encontrada no .env")
if not base_url:
    raise ValueError("OPENROUTER_API_BASE não encontrada no .env")

llm = ChatOpenAI(
    model="minimax/minimax-m2:free",
    temperature=0.7,
    api_key=SecretStr(api_key),
    base_url=base_url,
    default_headers={
        "HTTP-Referer": "http://localhost:5000",
        "X-Title": "GeradorDePropostasIA"
    }
)

script_path = Path(__file__).resolve()
PROJECT_ROOT = script_path.parent.parent.parent
TEMPLATE_PATH = PROJECT_ROOT / "templates" / "base.html"
try:
    with open(TEMPLATE_PATH, "r", encoding="utf-8") as f:
        exemplo_html = f.read()
    print(f"Sucesso: Template '{TEMPLATE_PATH}' carregado.")
except FileNotFoundError:
    print(f"Arquivo de template não encontrado em: {TEMPLATE_PATH}")
    print("IA vai gerar o design do 0")
    exemplo_html = "Nenhum exemplo fornecido. Crie um design HTML profissional do zero."

parser = PydanticOutputParser(pydantic_object=PropostaConteudo)

prompt_template = """
Você é um assistente de IA especialista em duas coisas:
1. Redator de Propostas Comerciais (Copywriter)
2. Desenvolvedor Web Front-End (HTML/CSS/JS)

Sua tarefa é gerar uma proposta comercial completa, persuasiva e visualmente atraente 
em formato HTML, com base nas informações fornecidas e no exemplo de proposta fornecido

**Informações Base da Proposta:**
* **Empresa do Cliente:** {nome_empresa}
* **Nome do Contato:** {nome_cliente}
* **Título da Proposta:** {titulo}
* **Prompt/Instruções do Usuário:** {prompt}
* **Cores Sugeridas:** {cores}
* **Logo da Empresa:** {logo}
* **Logo do Cliente:** {logo_cliente}

### EXEMPLO DE PROPOSTA (Use como sua base de estilo e estrutura)
Este é um um exemplo de alta qualidade fornecido:
{exemplo_html}

**Requisitos da Resposta:**
1.  **HTML Completo:** Sua resposta deve ser um documento HTML completo, começando com `<!DOCTYPE html>` e terminando com `</html>`.
2.  **CSS Inline ou em Bloco:** TODO o CSS deve estar dentro do arquivo, seja em atributos `style="..."` ou em um bloco `<style>...</style>` no `<head>`. Não use links externos para CSS.
3.  **Design Moderno:** Use um design limpo, profissional e moderno (ex: flexbox, padding, fontes legíveis).
4.  **Conteúdo Persuasivo:** Use as informações base para gerar o conteúdo de todas as seções necessárias (Introdução, O Desafio do Cliente, Nossa Solução, Escopo, Próximos Passos).
5.  **Use as Cores:** Se as cores forem fornecidas, tente incorporá-las no design (ex: em títulos, botões).
6.  **REGRA ESTRITA:** Responda APENAS com o código HTML. Não inclua NENHUM texto, preâmbulo (como "Aqui está seu HTML...") ou explicação antes de `<!DOCTYPE html>` ou depois de `</html>`.

**Início da Resposta HTML:**
<!DOCTYPE html>
"""

prompt = ChatPromptTemplate.from_template(
    prompt_template,
    partial_variables={"format_instructions": parser.get_format_instructions()}
)

chain = prompt | llm | StrOutputParser()

async def gerar_html_proposta(proposta: PropostaModel) -> str:
    input_data = {
        "nome_empresa": proposta.nome_empresa,
        "nome_cliente": proposta.nome_cliente or "Não informado",
        "titulo": proposta.titulo,
        "prompt": proposta.prompt,
        "cores": ", ".join(proposta.cores) if proposta.cores else "Cores padrão (azul e cinza)",
        "logo": proposta.logo,
        "logo_cliente": proposta.logo_cliente,
        "exemplo_html": exemplo_html
    }
    try:
        resultado_html = await chain.ainvoke(input_data)
        return resultado_html
    except Exception as e:
        print(f"Erro ao gerar proposta: {e}")
        raise e