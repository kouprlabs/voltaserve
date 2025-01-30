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
  <Voltaserve extensions={...} />
)
```

Build:

```shell
bun run build:rollup
```

## Use as a Web App

Install dependencies:

```shell
bun i
```

Run for development:

```shell
bun run dev
```

Build for production:

```shell
bun run build
```

Run for production:

```shell
go run .
```

Lint TypeScript code:

```shell
bun run lint
```

Format TypeScript code:

```shell
bun run format
```

Lint Go code:

```shell
golangci-lint run
```

Build Docker Image:

```shell
docker build -t voltaserve/ui .
```
