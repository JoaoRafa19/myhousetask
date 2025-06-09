# 🧹 Projeto: Myhousetask – Gerenciador de Tarefas Domésticas Compartilhadas

## 📘 Visão Geral

O **Myhousetask** é um aplicativo colaborativo para gerenciamento de tarefas domésticas entre membros de uma família. O objetivo é promover organização e divisão justa das tarefas, de forma clara, acessível e interativa, com sorteio de tarefas rotineiras e possibilidade de redistribuição entre os membros.

---

## 🧾 Funcionalidades principais

* **👨‍👩‍👧‍👦 Sistema de famílias**
    * Criação de famílias por um usuário.
    * Entrada via código ou link de convite.
    * Gerenciamento colaborativo de tarefas.

* **📋 Tarefas e rotinas**
    * Criação de tarefas únicas e rotineiras.
    * Sorteio diário de tarefas entre os membros.
    * Redistribuição de tarefa (somente com aceite do novo membro).
    * Status de tarefas: `pendente`, `em andamento`, `concluída`.

* **📅 Calendário**
    * Visualização de tarefas por dia.
    * Edição colaborativa de tarefas não rotineiras.
    * Execução automática de rotinas diárias.

* **🔔 Notificações em tempo real**
    * Atualizações de status e responsabilidades via **WebSocket**.

---

## 🛠️ Tecnologias Utilizadas

| Tecnologia  | Uso principal                              |
| :---------- | :----------------------------------------- |
| **Go** | Backend (linguagem principal)              |
| **sqlc** | Geração de código a partir de SQL          |
| **MySQL** | Banco de dados relacional                  |
| **WebSocket** | Comunicação em tempo real                  |
| **gRPC** | Comunicação entre microserviços e Flutter  |
| **Templ** | Frontend SSR opcional em Go                |
| **React** | Frontend web alternativo (uso leve)        |
| **Flutter** | Frontend mobile principal                  |

---

## 🏗️ Arquitetura do Sistema

***Projeto***

![img](/docs/arc_proto.png)

A arquitetura do **Myhousetask** é desenhada para ser moderna, reativa e eficiente, separando claramente as responsabilidades de cada parte do sistema, conforme ilustrado no diagrama visual (`docs/arc_proto.png`).

### 1. Clientes (Frontend)

* **Flutter (App Mobile):** Cliente principal. Prioriza **gRPC** para comunicação rápida e se conecta via **WebSocket** para atualizações em tempo real.
* **Go + Templ (App Web):** Alternativa web com renderização no servidor (SSR). Interage diretamente com as Regras de Negócio e usa **WebSocket** para sincronização.

### 2. Backend (Servidor Go)

* **Camada de Entrada (API):** Expõe **API gRPC**, **API REST (`chi`)**, e um **Servidor WebSocket**.
* **Camada de Serviços (Regras de Negócio):** Contém a lógica principal do sistema (`internal/services/`), de forma agnóstica à API.
* **Camada de Acesso a Dados (DAL):** Usa **`sqlc`** para gerar uma interface Go tipada e segura a partir de queries SQL.

### 3. Banco de Dados

* **MySQL:** Armazena todas as informações de forma persistente.

---

## 📂 Estrutura Inicial do Projeto

```
myapp/
├── db/
│   ├── migrations/        # Scripts SQL (schema do banco)
│   ├── query/             # Queries SQL com nomes para o sqlc
│   └── gen/               # Código gerado por sqlc (structs + funções)
├── internal/
│   ├── api/               # Handlers HTTP (REST e gRPC)
│   ├── ws/                # Lógica de WebSocket
│   └── services/          # Regras de negócio
├── cmd/
│   └── server/            # main.go (servidor principal)
├── sqlc.yaml              # Configuração do sqlc
└── go.mod
```

---

## 📜 Progresso atual

* ✅ Schema SQL modelado.
* ✅ `sqlc.yaml` configurado e `sqlc generate` funcionando.
* ✅ Servidor Go inicializado com `chi`.

---

## 🧠 Decisões de Design e Convenções

* **Nomenclatura no Banco de Dados:** Tabelas são `minúsculas` (ex: `users`) e colunas são `MAIÚSCULAS` (ex: `NAME`). Essa convenção melhora a legibilidade das queries SQL e diferencia claramente as estruturas do banco no código Go gerado.
* **Tratamento de Nulos com `sqlc`:** A opção `emit_pointers_for_null_types: false` foi usada para que campos `NULL` no banco sejam gerados como tipos zero em Go (ex: string vazia `""`) em vez de ponteiros (como `*sql.NullString`). Isso simplifica a manipulação de dados na camada de serviço, evitando verificações constantes de `nil`.
* **Chaves Primárias:** Utiliza-se `CHAR(36)` para UUIDs como chaves primárias para garantir identificadores únicos e não sequenciais, o que é mais seguro e escalável.

---

## 📈 Próximos passos detalhados

1.  **Criar rotas REST para CRUD de Usuários**
    * `POST /users` - Criar novo usuário.
    * `GET /users/{id}` - Obter dados de um usuário.
    * `PUT /users/{id}` - Atualizar dados de um usuário.
    * `DELETE /users/{id}` - Desativar um usuário.
2.  **Definir padrão de autenticação com JWT**
    * Criar rota `POST /login` que retorna um token JWT.
    * Implementar um middleware em `chi` para validar o token em rotas protegidas.
3.  **Implementar sorteio automático de tarefas**
    * Criar um worker/goroutine que executa uma vez ao dia.
    * O worker deve buscar as tarefas rotineiras, sortear os responsáveis entre os membros ativos da família e salvar os resultados.
4.  **Criar canal WebSocket para atualização de status**
    * Definir um hub de WebSocket que gerencia as conexões.
    * Quando uma tarefa for atualizada (ex: via API REST), o serviço correspondente deve notificar o hub, que transmitirá a mensagem para os clientes relevantes.
5.  **Iniciar o frontend em Flutter (ou Templ)**
    * Começar a construir as telas de login, registro e a tela principal de tarefas.

---

## Dependencias

- Projeto
  
  ```
  github.com/go-sql-driver/mysql v1.9.2
  github.com/golang-migrate/migrate/v4 v4.18.3
  google.golang.org/grpc v1.65.0
  google.golang.org/protobuf v1.36.6
  ```

- protocgen

  ```
  google.golang.org/protobuf/cmd/protoc-gen-go@latest
  ```

- sqlc

  ```
  go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
  ```

- gRPC

  ```
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  ```

- templ
  ```
  go install github.com/a-h/templ/cmd/templ@latest
  ```



## 🚀 Comandos Úteis

### Gerar arquivo Go a partir de arquivo .proto

```shell
protoc --go_out=. --go_opt=paths=source_relative \
         --go-grpc_out=. --go-grpc_opt=paths=source_relative \
         proto/category.proto
```

### Criar ou executar o banco MySQL no terminal via Docker

```shell
docker run -d \
  --name mysql-myhousetask \
  -p 3308:3306 \
  -e MYSQL_DATABASE=myhousetask \
  -e MYSQL_USER=user \
  -e MYSQL_PASSWORD=root \
  -e MYSQL_ROOT_PASSWORD=root \
  mysql:latest
```

### Gerar arquivos do sqlc
```
sqlc -f db/sqlc.yaml generate
```



```
templ generate
```