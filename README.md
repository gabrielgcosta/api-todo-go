# Go Todo API 🚀

Uma API REST simples e eficiente para gerenciamento de tarefas (Todo List), desenvolvida em **Go** (Golang) com banco de dados **PostgreSQL** e empacotada com **Docker**.

Este projeto foi construído utilizando a biblioteca padrão do Go (`net/http`), aproveitando os novos recursos de roteamento nativos do Go 1.22+.

---

## 🛠️ Tecnologias Utilizadas

- **Go 1.26** (utilizando apenas a biblioteca padrão `net/http` para roteamento e servidor)
- **PostgreSQL 15** (banco de dados relacional)
- **lib/pq** (driver nativo de PostgreSQL para Go)
- **Docker & Docker Compose** (para conteinerização e fácil execução local)

---

## 🚀 Como Executar o Projeto

Você não precisa ter o Go ou PostgreSQL instalados na sua máquina local, apenas o **Docker** e o **Docker Compose**.

### Passos para rodar:

1. Clone este repositório:
   ```bash
   git clone https://github.com/gabrielgcosta/api-todo-go.git
   cd api
   ```

2. Inicie os containers com Docker Compose:
   ```bash
   docker compose up --build -d
   ```

A API estará rodando em `http://localhost:8080` e criará automaticamente a tabela de tarefas no banco de dados na primeira inicialização.

---

## 📌 Rotas da API

### 1. Criar uma Tarefa
* **Rota:** `POST /tasks`
* **JSON de Entrada:**
  ```json
  {
    "title": "Estudar estruturas de dados em Go"
  }
  ```
* **Resposta (201 Created):**
  ```json
  {
    "id": 1,
    "title": "Estudar estruturas de dados em Go",
    "finished": false,
    "created_at": "2026-06-08T22:37:58Z"
  }
  ```

### 2. Listar todas as Tarefas
* **Rota:** `GET /tasks`
* **Resposta (200 OK):**
  ```json
  [
    {
      "id": 1,
      "title": "Estudar estruturas de dados em Go",
      "finished": false,
      "created_at": "2026-06-08T22:37:58Z"
    }
  ]
  ```

### 3. Atualizar uma Tarefa
* **Rota:** `PUT /tasks/{id}`
* **JSON de Entrada:**
  ```json
  {
    "title": "Estudar estruturas de dados em Go (Concluído)",
    "finished": true
  }
  ```
* **Resposta (200 OK):** `Task updated successfully`

### 4. Deletar uma Tarefa
* **Rota:** `DELETE /tasks/{id}`
* **Resposta (204 No Content)**

---

## 🏗️ Arquitetura & Recursos Recentes

A API adota práticas de arquitetura modular, focada em testabilidade, segurança e observabilidade:
- **Middleware de Log**: Registra automaticamente no stdout detalhes de cada requisição (método, caminho, IP, status code HTTP, duração e causa raiz detalhada de eventuais falhas).
- **Tratamento de Erros Centralizado (`apierror`)**: Retorna erros padronizados em JSON (`{"error": "message"}`). Utiliza `errors.As` para extrair detalhes estruturados e ocultar erros sensíveis de banco de dados do cliente final, registrando o erro raiz detalhado apenas no log interno.
- **Worker de Eventos Assíncronos (`worker`)**: Um processador em segundo plano rodando em uma goroutine. Os handlers HTTP despacham eventos de CRUD para canais em memória, liberando a resposta HTTP para o cliente imediatamente.

---

## 🧪 Como Executar os Testes

O projeto possui uma suíte completa de testes unitários que cobre a lógica do worker, tratamento de erros e handlers HTTP (com mocks para evitar a dependência com o Postgres).

Para rodar os testes e verificar a cobertura do código localmente:
```powershell
go test -cover ./...
```

---

## 📁 Estrutura do Projeto

* `main.go`: Ponto de inicialização do banco, do worker assíncrono, roteador nativo e injeção do middleware de logs.
* `apierror/`: Pacote para formatação JSON de erros e compatibilidade com erros envelopados.
* `database/`: Pacote para inicialização e migração automática do esquema do banco de dados.
* `middleware/`: Contém os middlewares da aplicação (como o Logger).
* `task/`: Pacote de domínio contendo a definição da entidade, a interface e implementação de persistência, e os handlers HTTP.
* `worker/`: Pacote responsável pela goroutine e canal de processamento assíncrono de eventos de tarefas.
* `Dockerfile`: Configura o build em multi-estágio da aplicação para gerar um container Alpine leve.
* `docker-compose.yml`: Orquestra o container da aplicação e do banco PostgreSQL.
