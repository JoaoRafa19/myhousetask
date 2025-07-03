# Documentação do Projeto: go-grpc-teste

## 1. Visão Geral do Projeto

O projeto **go-grpc-teste** é uma aplicação backend desenvolvida em Go, que utiliza gRPC para comunicação entre serviços, SQLC para geração de código de acesso a banco de dados, e segue uma arquitetura modularizada. O objetivo principal é gerenciar entidades como categorias, famílias, usuários, tarefas e eventos de calendário, com suporte a autenticação, permissões e administração via interface web.

### Principais Tecnologias Utilizadas

- **Go (Golang):** Linguagem principal do backend.
- **gRPC:** Comunicação eficiente entre serviços.
- **SQLC:** Geração de código Go a partir de queries SQL.
- **Docker/Docker Compose:** Containerização e orquestração de serviços.
- **Tailwind CSS:** Estilização da interface web.
- **Templ (templ):** Geração de templates HTML.
- **PostgreSQL:** Banco de dados relacional.
- **Migrações SQL:** Controle de versão do banco de dados.

---

## 2. Estrutura do Projeto

```
go-grpc-teste/
│
├── category/                # Lógica de domínio para categorias
│   └── category.go
│
├── cmd/                     # Pontos de entrada da aplicação
│   └── server/
│       ├── grpc/main.go     # Inicialização do servidor gRPC
│       └── main/main.go     # Inicialização do servidor principal
│
├── db/
│   ├── gen/                 # Código gerado pelo SQLC
│   ├── migrations/          # Scripts de migração do banco de dados
│   ├── query/               # Queries SQL utilizadas pelo SQLC
│   └── sqlc.yaml            # Configuração do SQLC
│
├── internal/
│   ├── core/services/       # Serviços de domínio (ex: dashboard)
│   └── web/                 # Lógica e componentes da interface web
│       ├── handlers/        # Handlers HTTP
│       └── view/            # Templates e componentes visuais
│
├── migrator/                # Lógica para rodar migrações
├── proto/                   # Definições de Protobuf (gRPC)
├── web/static/              # Arquivos estáticos (CSS)
├── Dockerfile               # Containerização da aplicação
├── docker-compose.yml       # Orquestração de containers
├── go.mod / go.sum          # Gerenciamento de dependências Go
└── readme.md                # Documentação inicial
```

---

## 3. Fluxo de Desenvolvimento

### 3.1. Criação e Evolução do Banco de Dados

- As migrações SQL são criadas em `db/migrations/` para versionar o banco.
- Cada entidade (ex: categoria, família, usuário) possui scripts `.up.sql` (criação) e `.down.sql` (remoção).
- Para aplicar migrações, utilize ferramentas como [golang-migrate](https://github.com/golang-migrate/migrate) ou comandos customizados no projeto.

### 3.2. Geração de Código com SQLC

- Queries SQL são escritas em `db/query/query.sql`.
- O SQLC gera código Go para acesso ao banco em `db/gen/`.
- Para gerar o código, execute:
  ```sh
  sqlc generate
  ```
- O arquivo `db/sqlc.yaml` define as configurações de geração.

### 3.3. Definição de APIs com Protobuf/gRPC

- As APIs são definidas em arquivos `.proto` dentro de `proto/`.
- Para gerar o código Go dos serviços gRPC:
  ```sh
  protoc --go_out=. --go-grpc_out=. proto/*.proto
  ```
- O servidor gRPC é inicializado em `cmd/server/grpc/main.go`.

### 3.4. Backend Web e Handlers

- Handlers HTTP estão em `internal/web/handlers/`.
- Os templates HTML são gerados com Templ em `internal/web/view/`.
- Componentes visuais reutilizáveis ficam em `internal/web/view/components/`.

### 3.5. Frontend e Estilização

- O frontend utiliza Tailwind CSS, configurado em `tailwind.config.js`.
- Os arquivos CSS estão em `web/static/css/`.

### 3.6. Containerização

- O `Dockerfile` define a imagem da aplicação.
- O `docker-compose.yml` orquestra múltiplos serviços (ex: app, banco de dados).

---

## 4. Estado Atual do Projeto

- **Banco de dados:** Diversas tabelas já criadas (categorias, famílias, usuários, tarefas, eventos, convites).
- **Migrações:** Scripts de migração versionados e organizados.
- **SQLC:** Queries e código gerado para acesso ao banco.
- **gRPC:** Definições de proto para categorias e usuários.
- **Backend Web:** Handlers e templates para dashboard e administração.
- **Containerização:** Dockerfile e docker-compose prontos, mas ainda não versionados no git.
- **Pendências:** Alguns arquivos modificados não estão com commit, e arquivos Docker ainda não estão versionados.

---

## 5. Próximos Passos Sugeridos

### 5.1. Organização e Versionamento

- Adicionar e commitar arquivos Docker (`Dockerfile`, `docker-compose.yml`).
- Garantir que todas as alterações em arquivos Go e SQL estejam versionadas.

### 5.2. Testes

- Implementar testes unitários e de integração para handlers, serviços e queries.
- Garantir cobertura mínima para as principais funcionalidades.

### 5.3. Documentação

- Expandir o `readme.md` com instruções de setup, build, testes e deploy.
- Documentar endpoints gRPC e exemplos de uso.

### 5.4. Segurança e Autenticação

- Implementar autenticação JWT ou similar para APIs.
- Adicionar controle de permissões para rotas administrativas.

### 5.5. Integração Contínua

- Configurar pipeline CI/CD (ex: GitHub Actions) para build, testes e deploy automático.

### 5.6. Melhorias de UX/UI

- Refinar componentes visuais e responsividade.
- Adicionar feedbacks visuais para ações do usuário.

### 5.7. Monitoramento e Observabilidade

- Adicionar logs estruturados.
- Configurar métricas e alertas básicos.

---

## 6. Processos de Geração de Código

### 6.1. Geração de Código SQLC

1. Escreva queries SQL em `db/query/query.sql`.
2. Execute `sqlc generate` para gerar código Go.
3. Utilize os métodos gerados em `db/gen/querier.go` nos serviços.

### 6.2. Geração de Código gRPC

1. Defina serviços e mensagens em arquivos `.proto` em `proto/`.
2. Execute:
   ```sh
   protoc --go_out=. --go-grpc_out=. proto/*.proto
   ```
3. Implemente os serviços gerados no backend.

### 6.3. Migrações de Banco

1. Crie scripts `.up.sql` e `.down.sql` em `db/migrations/`.
2. Aplique migrações com ferramenta apropriada:
   ```sh
   migrate -path db/migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up
   ```

### 6.4. Build e Deploy com Docker

1. Build da imagem:
   ```sh
   docker build -t go-grpc-teste .
   ```
2. Subida dos serviços:
   ```sh
   docker-compose up -d
   ```

---

## 7. Boas Práticas

- Sempre versionar código gerado e scripts de migração.
- Manter as definições de proto e queries SQL sincronizadas com o código.
- Escrever testes para novas funcionalidades.
- Documentar endpoints e fluxos críticos.
- Utilizar variáveis de ambiente para configurações sensíveis.

---

## 8. Referências Úteis

- [Documentação do Go](https://golang.org/doc/)
- [gRPC em Go](https://grpc.io/docs/languages/go/)
- [SQLC](https://docs.sqlc.dev/)
- [Docker](https://docs.docker.com/)
- [Tailwind CSS](https://tailwindcss.com/docs/installation)
- [Templ](https://templ.guide/)

