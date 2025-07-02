# 📄 Documentação Extensiva do Projeto: MyHouseTask

**Última Atualização:** 02 de julho de 2025

## 1. 📘 Visão Geral e Objetivo

O **MyHouseTask** é uma aplicação multiplataforma concebida para simplificar a gestão de tarefas domésticas em ambientes familiares ou partilhados. O objetivo central é criar um ecossistema digital onde as responsabilidades podem ser distribuídas de forma justa e transparente, promovendo a colaboração e a organização.

O sistema deve ser:

* **Reativo:** As atualizações devem ser refletidas em tempo real em todos os dispositivos conectados.
* **Eficiente:** A arquitetura deve ser otimizada para um baixo consumo de recursos, especialmente na aplicação móvel.
* **Acessível:** Com uma aplicação móvel nativa como cliente principal e um painel web para administração.

## 2. 🏗️ Arquitetura Detalhada

A arquitetura do projeto foi desenhada para ser modular, escalável e de fácil manutenção, seguindo uma clara separação de responsabilidades.

### Fluxo Geral

* **Utilizadores Finais (Famílias):** Interagem principalmente através da aplicação móvel **Flutter**. A comunicação com o backend é feita via **gRPC** para performance e **WebSocket** para atualizações em tempo real.
* **Administradores:** Interagem através de um painel web renderizado no servidor com **Go + Templ**. A comunicação é feita internamente no backend (HTTP/renderização direta) e também utiliza **WebSocket** para dados ao vivo.

### Componentes Principais

#### A. Backend (Go)

O cérebro da aplicação, escrito em Go para alta performance e concorrência.

* **Camada de Entrada (API):**
    * **Servidor gRPC:** A principal porta de entrada para a aplicação móvel Flutter. Oferece uma comunicação binária, rápida e fortemente tipada, ideal para mobile. Expõe serviços como `UserService`, `TaskService`, etc.
    * **Servidor HTTP (com `chi`):** Serve o painel administrativo web. Os *handlers* nesta camada são responsáveis por receber pedidos HTTP, chamar a camada de serviço e renderizar os componentes **Templ**.
    * **Servidor WebSocket:** Mantém conexões persistentes com todos os clientes (Flutter e Web) para "empurrar" atualizações em tempo real, como a mudança de status de uma tarefa.
* **Camada de Serviços (`internal/core/service`):**
    * Contém a lógica de negócio pura, agnóstica em relação à forma como foi chamada (gRPC ou HTTP).
    * Exemplos: `DashboardService` (para agregar dados para o painel), `UserService` (para gerir utilizadores), `TaskAssignmentService` (para a lógica de sorteio de tarefas).
    * Esta camada é a única que interage com a camada de acesso a dados.
* **Camada de Acesso a Dados (DAL - `internal/data/db`):**
    * A interação com a base de dados é totalmente gerida pelo **`sqlc`**.
    * As queries SQL são escritas em ficheiros `.sql`. O `sqlc` lê estas queries e gera automaticamente o código Go (structs e métodos) para executar essas operações de forma segura e tipada. Isto elimina a necessidade de escrever SQL "boilerplate" no código Go.

#### B. Base de Dados (MySQL)

Um banco de dados relacional robusto para armazenar de forma persistente todos os dados da aplicação: utilizadores, famílias, tarefas, rotinas, etc.

#### C. Clientes (Frontend)

* **Flutter (App Móvel):** A experiência principal para o utilizador final. Consome os serviços gRPC e subscreve as atualizações via WebSocket.
* **Templ + htmx (Painel Web):** Uma interface web leve e rápida para administração, renderizada no servidor.

## 3. 🛠️ Tecnologias e Justificativas

| Tecnologia | Justificativa                                                                                                                              |
| :--- | :----------------------------------------------------------------------------------------------------------------------------------------- |
| **Go** | Excelente para backend de alta performance, concorrência nativa (goroutines) e compilação para um único binário, simplificando o deploy. |
| **MySQL** | Banco de dados relacional maduro, fiável e amplamente utilizado.                                                                           |
| **sqlc** | Ferramenta "database-first" que gera código Go idiomático e seguro a partir de SQL, acelerando o desenvolvimento da camada de dados.          |
| **gRPC** | Protocolo de comunicação de alta performance da Google, ideal para a comunicação entre o backend e a app móvel Flutter devido à sua eficiência. |
| **Templ** | Motor de templates para Go que permite criar interfaces web renderizadas no servidor de forma segura e componentizada, ideal para o painel admin. |
| **WebSocket** | Protocolo de comunicação bidirecional que permite ao servidor enviar dados para os clientes em tempo real sem que eles precisem de pedir. |
| **Flutter** | Framework da Google para construir aplicações móveis nativas para iOS e Android a partir de uma única base de código.                   |

## 4. ✅ Progresso Atual e Artefatos Gerados

* **Base de Dados:** O schema inicial foi modelado e os ficheiros de migration (`.up.sql`, `.down.sql`) foram criados para as tabelas principais (users, families, tasks, etc.).
* **`sqlc`:** O ficheiro `sqlc.yaml` está configurado e o comando `sqlc generate` foi executado com sucesso, gerando os modelos e métodos Go na pasta `internal/data/db/gen`.
* **Servidores:** O esqueleto dos servidores gRPC e HTTP (`chi`) foi inicializado no `cmd/server/main.go`.
* **Painel Administrativo:**
    * O design visual completo do dashboard foi criado como um ficheiro HTML estático com TailwindCSS.
    * Este design foi decomposto em componentes **Templ** reutilizáveis, seguindo a estrutura `layout/`, `pages/` e `components/` em `internal/web/view/`.
    * Foi implementada a lógica para abrir um modal de "Criar Família" usando JavaScript nativo.

## 5. 🔜 Plano de Ação Detalhado

1.  **Implementar a Camada de Serviços (Go):**
    * **Tarefa 5.1:** Criar e implementar o `DashboardService` para buscar os dados agregados necessários para o painel (contadores, listas recentes), utilizando os métodos gerados pelo `sqlc`.
    * **Tarefa 5.2:** Criar e implementar o `UserService` com as funções básicas de CRUD (Create, Read, Update, Delete) e autenticação.
    * **Tarefa 5.3:** Criar os outros serviços de negócio conforme necessário (`FamilyService`, `TaskService`).
2.  **Conectar os Handlers aos Serviços:**
    * **Tarefa 6.1 (Painel Web):** Modificar o `DashboardHandler` (`internal/web/handler/`) para chamar o `DashboardService`, obter os dados reais e passá-los para o componente `DashboardPage.templ`.
    * **Tarefa 6.2 (gRPC):** Implementar os métodos do servidor gRPC (ex: `CreateUser` no `user_server.go`) para que eles chamem os métodos correspondentes no `UserService`.
3.  **Implementar Autenticação (JWT):**
    * **Tarefa 7.1:** No `UserService`, criar um método `Authenticate(email, password)` que valide as credenciais e gere um token JWT.
    * **Tarefa 7.2:** Criar um endpoint de login (tanto gRPC como HTTP) que chame este serviço.
    * **Tarefa 7.3:** Implementar middlewares (gRPC Interceptor e Chi Middleware) para validar o token JWT em rotas protegidas.
4.  **Desenvolvimento do Frontend Flutter:**
    * **Tarefa 8.1:** Configurar o projeto Flutter para usar os ficheiros `.proto` e gerar o código cliente gRPC.
    * **Tarefa 8.2:** Construir as telas de Login/Registo, consumindo os endpoints de autenticação gRPC.
    * **Tarefa 8.3:** Desenvolver a tela principal que lista as tarefas do utilizador.

## 6. 🚀 Como Executar o Projeto (Ambiente de Desenvolvimento)

1.  **Base de Dados:** Inicie uma instância do MySQL (preferencialmente via Docker, usando o comando já documentado).
2.  **Dependências Go:** Execute `go mod tidy` para garantir que todas as dependências estão baixadas.
3.  **Geração de Código:** Execute `sqlc generate` se houver alterações nas queries SQL.
4.  **Executar o Backend:** Execute `go run ./cmd/server/main.go` para iniciar os servidores gRPC e HTTP.
```