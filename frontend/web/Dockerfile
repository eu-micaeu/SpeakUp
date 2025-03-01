FROM node:18-alpine

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

# Garantindo permissões para o diretório node_modules
# Adicionando permissões para o node_modules
RUN chown -R node:node /app

USER node

CMD ["npm", "start"]
