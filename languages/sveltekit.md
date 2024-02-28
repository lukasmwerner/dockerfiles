# SvelteKit Application Framework

This assumes that you already have the node adapter installed through this:
[Node servers • Docs • SvelteKit](https://kit.svelte.dev/docs/adapter-node)

```dockerfile
FROM node:18

WORKDIR /src
# Deps Stage
COPY package.json .
RUN npm i

# Build Stage
COPY . .
RUN npm run build

ENV HOST 0.0.0.0
ENV PORT 3000
# This is the default port sveltekit servers in prod listen to
EXPOSE 3000

# CHANGEME this is your domain and proto
ENV ORIGIN http://localhost:3000

# Run stage
CMD ["node", "build"]
```
