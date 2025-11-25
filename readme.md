# 🚀 Propulse - V1

![Tecnologia](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Tecnologia](https://img.shields.io/badge/Python-3776AB?style=for-the-badge&logo=python&logoColor=white)
![Tecnologia](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![Tecnologia](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![Tecnologia](https://img.shields.io/badge/LangChain-8A2BE2?style=for-the-badge)

Propulse é uma aplicação focada em automação e inteligência artificial, especializada na geração de propostas personalizadas.

## Sobre o Projeto

Este é um sistema de microsserviços para geração de propostas comerciais em PDF. O projeto combina um backend de alta performance em **Go (Gin)** com um serviço de IA em **Python (FastAPI + LangChain)**.

O usuário envia os dados da proposta via API, a IA gera um HTMl e, em seguida, o **Playwright** renderiza o PDF final em um navegador headless, garantindo 100% de fidelidade de design (incluindo CSS complexo como o do Tailwind).

## ✨ Features Principais

- **API principal (Go):** Gerencia propostas e o status dos jobs de geração.
    
- **Microsserviço de IA (Python):** Utiliza **LangChain** para processar prompts e gerar o conteúdo das propostas.
    
- **Ambiente Local:** 100% "containerizado" com **Docker** para fácil execução e replicação do ambiente de desenvolvimento.

- **Geração de PDF com HTML/CSS:** Usa Playwright para renderizar HTML complexo (com CSS moderno como Tailwind) com perfeita fidelidade.
  
- **Engenharia de Prompt:** A IA gera o HTML do zero baseado em templates e modelos pré-aprovados.

    

## 🛠️ Tech Stack (Tecnologias Utilizadas)

|Categoria|Tecnologia|
|---|---|
|**Backend**|Go (Gin-Gonic), Python (FastAPI)|
|**Inteligência Artificial**|LangChain|
|**Banco de Dados**|PostgreSQL|
|**DevOps**|Docker, Docker Compose|
|**Conceitos**|Microsserviços, System Design, APIs REST|

## 🏗️ Arquitetura

O sistema é dividido nos seguintes serviços "containerizados" (gerenciados pelo `docker-compose.yml`):

- **`backend` (Go):** O serviço principal que orquestra as regras de negócio e expõe a API.
    
- **`ia-service` (Python):** Um serviço dedicado para as operações de IA, consumido pelo `backend`.
    
- **`db` (PostgreSQL):** O banco de dados relacional para persistência dos dados.
    

## 📖 Principais Endpoints da API

Todas as rotas são prefixadas com `/proposta/`.

|Método|Rota|Descrição|
|---|---|---|
|`POST`|`/`|Cria uma nova proposta.|
|`GET`|`/`|Lista todas as propostas existentes.|
|`GET`|`/:id`|Busca uma proposta específica pelo seu ID.|
|`PATCH`|`/:id`|Atualiza o status ou título da proposta pelo ID.|
|`DELETE`|`/:id`|Deleta uma proposta pelo seu ID.|
|`POST`|`/:id/regerar`|Dispara um novo job para regerar o conteúdo de uma proposta existente.|

## 🚀 Como Executar (Ambiente de Desenvolvimento Local)

O projeto é totalmente "containerizado", facilitando a configuração do ambiente.

1. **Clone o repositório**
    
2. **Configure as variáveis de ambiente:**
    
    - Existe um arquivo de exemplo em `backend/.env.example`.
        
    - Copie este arquivo para `backend/.env`
        
        ```
        cp backend/.env.example backend/.env
        ```
        
    - Edite o novo arquivo `backend/.env` e preencha com suas chaves de API e credenciais do banco de dados.

3. **Escolha o modelo de IA:**

   - Copie o ID do modelo escolhido em: https://openrouter.ai/models?q=free, entre em ./backend/IA/src/ia_generator/ia.py e na linha 23 ```model="x-ai/grok-4.1-fast:free",``` coloque o ID do modelo.
        
5. **Suba os containers:**
    
    ```
    docker-compose up -d --build
    ```
    Ao subir o conteiner será criado o banco de dados com a tabela!
    
6. Acesse o serviço de backend na porta mapeada (ex: `http://localhost:8080`).

# 🕹️ Testando os Endpoints (API)

Use o Postman, Insomnia ou cURL para testar a API.

### Criar Nova Proposta

**Request:** `POST http://localhost:8080/proposta`

**Body (JSON):**

```json
{
  "titulo": "Proposta de Teste (Docker)",
  "nomeEmpresa": "Empresa Cliente Teste",
  "nomeCliente": "Teste",
  "prompt": "Focar nos benefícios de usar Docker para a equipe de desenvolvimento, mencionando a padronização dos ambientes.",
  "cores": ["#007bff", "#343a40"],
  "logo": "https://exemplo.com/logo-empresa.png",
  "logoCliente": "https://exemplo.com/logo-cliente.png",
  "status": "rascunho",
}
```

**Sucesso (Resposta):**
A API retornará um JSON com a proposta criada, incluindo o `id` e o `arquivoFinal` (ex: `uploads/propostas/proposta_...pdf`). Verifique a pasta `./uploads` no seu computador\!
