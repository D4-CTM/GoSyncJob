FROM oven/bun:1

WORKDIR /frontend/app

COPY bun.lock package.json ./

RUN bun install

COPY . .

RUN bun run build

EXPOSE 4173

CMD ["bun", "run", "preview"]
