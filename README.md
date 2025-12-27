<div align="center">
  <img src="./frontend/web/public/logo.png" width="75" alt="SpeakUp Logo">
</div>

<br>

<div align="center">
   
  [![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/) [![MongoDB](https://img.shields.io/badge/MongoDB-47A248?style=for-the-badge&logo=mongodb&logoColor=white)](https://www.mongodb.com/) [![Svelte](https://img.shields.io/badge/Svelte-FF3E00?style=for-the-badge&logo=svelte&logoColor=white)](https://svelte.dev/) [![React Native](https://img.shields.io/badge/React_Native-61DAFB?style=for-the-badge&logo=react&logoColor=black)](https://reactnative.dev/) [![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/) [![Expo](https://img.shields.io/badge/Expo-000020?style=for-the-badge&logo=expo&logoColor=white)](https://expo.dev/)
  
</div>

<br>

<p align="justify">
Uma plataforma inovadora que utiliza inteligência artificial para auxiliar no aprendizado de novos idiomas de maneira eficiente e personalizada. 
Com foco na melhoria contínua, a IA analisa seus textos em tempo real, identificando erros de digitação e gramaticais.
A cada erro encontrado, a ferramenta destaca as correções necessárias, oferecendo explicações detalhadas sobre o que foi corrigido e o motivo.
Dessa forma, você não apenas corrige suas falhas, mas também aprende com cada interação, desenvolvendo suas habilidades linguísticas de maneira prática e envolvente.
</p>

> **📱 NOVIDADE: App Mobile disponível!** → [Ver documentação mobile](./mobile/START_HERE.md)

---

## 🚀 Tecnologias utilizadas

### ⚙️ GoLang
<p align="justify">
GoLang, ou Go, é uma linguagem de programação desenvolvida pela Google, projetada para ser eficiente, com ótimo desempenho em ambientes de alta concorrência. Suas principais vantagens incluem sintaxe simples, gerenciamento de memória automático e alta velocidade de compilação, tornando-a ideal para a criação de sistemas distribuídos, servidores web e soluções de grande escala.
</p>

### 🧡 SvelteKit & ⚛️ React Native
<p align="justify">
SvelteKit é um framework web moderno construído sobre o Svelte, oferecendo alta performance e gerando menos código boilerplate. Para o desenvolvimento mobile, utilizamos React Native, permitindo criar aplicativos nativos para iOS e Android com a flexibilidade do ecossistema JavaScript.
</p>

### 🐳 Docker
<p align="justify">
Docker é uma plataforma que permite empacotar e isolar aplicações em contêineres. Isso garante que o software funcione de maneira consistente em diferentes ambientes, desde o desenvolvimento até a produção. Docker simplifica a gestão de dependências e facilita o deployment de aplicativos, otimizando a utilização de recursos e a portabilidade entre plataformas.
</p>

### 🍃 MongoDB
<p align="justify">
MongoDB é um banco de dados NoSQL orientado a documentos que armazena dados em formato JSON-like. Essa estrutura permite flexibilidade no esquema dos dados, o que facilita a escalabilidade e a integração com aplicações modernas. Ideal para grandes volumes de dados e operações em tempo real, MongoDB é amplamente utilizado em aplicações web, IoT e Big Data.
</p>

---

## ✨ Principais funcionalidades

- ✅ **Correção em Tempo Real:** A IA analisa os textos à medida que você digita, identificando e corrigindo erros gramaticais, ortográficos e de digitação instantaneamente.
- 📘 **Explicações Detalhadas:** Cada correção vem acompanhada de uma explicação sobre o motivo do erro e a regra gramatical ou ortográfica aplicada, facilitando o aprendizado.
- 🧠 **Feedback Personalizado:** A plataforma se adapta ao seu nível de proficiência e oferece feedback personalizado para melhorar suas habilidades linguísticas, focando nas áreas em que você apresenta mais dificuldades.
- 📈 **Histórico de Aprendizado:** Todas as correções feitas ao longo do tempo ficam armazenadas em um histórico pessoal, permitindo que você acompanhe seu progresso e revise suas melhorias.

---

## 🛠️ Como rodar?

### 🌐 Web (Frontend + Backend)

<p align="justify">
Para rodar o projeto web, certifique-se de ter o Docker e o Docker Compose instalados em sua máquina. Siga os passos abaixo:
</p>

```bash
# 1. Clone o repositório do projeto
git clone https://github.com/eu-micaeu/speakup.git

# 2. Navegue até o diretório raiz do projeto
cd speakup

# 3. Inicie os serviços com Docker Compose
docker-compose up
```

<p align="justify"> Após a inicialização, você poderá acessar a aplicação no navegador em: <strong>http://localhost:3000</strong> </p>

### 📱 Mobile (React Native + Expo)

<p align="justify">
Para rodar o aplicativo mobile, siga os passos detalhados no guia completo:
</p>

**[📖 Ver guia completo de instalação mobile →](./mobile/SETUP.md)**

```bash
# 1. Entre na pasta mobile
cd mobile

# 2. Instale as dependências
npm install

# 3. Configure o .env com o IP do backend
cp .env.example .env

# 4. Inicie o app
npm start
```

<p align="justify">
Escaneie o QR Code com o app <strong>Expo Go</strong> no seu celular para testar o app!
</p>

---

## 👥 GitHub dos Desenvolvedores

<div>
  <a href="https://github.com/eu-micaeu" target="_blank">
    <img src="https://avatars.githubusercontent.com/u/69124656?v=4" alt="Micael" width="100" style="border-radius: 50%;">
  </a>
  &nbsp;&nbsp;&nbsp;&nbsp;
  <a href="https://github.com/ryanmiura" target="_blank">
    <img src="https://avatars.githubusercontent.com/u/57397826?v=4" alt="Ryan" width="100" style="border-radius: 50%;">
  </a>
</div>
