FROM node:18-alpine

WORKDIR /app

COPY package*.json ./

RUN npm ci

COPY . .

RUN mkdir -p node_modules/.cache && chown -R node:node /app

USER node

CMD ["npm", "start"]
