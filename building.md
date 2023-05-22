## Building this plugin

### Frontend

1. Install dependencies

```sh
npm install
```

2. Build plugin in development mode

```sh
npm run dev
```

3. Build plugin in production mode

```sh
npm run build
```

### Backend

1. Install dependencies

```sh
go mod tidy
```

2. Build backend plugin binaries for Linux, Windows and Darwin:

```sh
mage -v
```

3. List all available Mage targets for additional commands:

```sh
mage -l
```
