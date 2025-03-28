# Voltaserve UI

Web app and extensible React component.

## Use as a React Component

Installation:

```shell
npm i @voltaserve/ui
```

Usage:

```tsx
import { Voltaserve } from '@voltaserve/ui'
import { createRoot } from 'react-dom/client'

createRoot(document.getElementById('root') as HTMLElement).render(
  <Voltaserve extensions={/*...*/} />
)
```

Build:

```shell
npm run build:rollup
```

## Use as a Web App

Install dependencies:

```shell
npm i
```

Run for development:

```shell
npm run dev
```

Build for production:

```shell
npm run build
```

Run for production:

```shell
go run .
```

Lint TypeScript code:

```shell
npm run lint
```

Format TypeScript code:

```shell
npm run format
```

Format Go code:

```shell
gofumpt -w . && \
gofmt -s -w . && \
goimports -w . && \
golangci-lint run --fix
```

Lint Go code:

```shell
golangci-lint run
```

Build Docker Image:

```shell
docker build -t voltaserve/ui .
```
