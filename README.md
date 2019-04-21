# Stamped
Swagger to Postman.json

# Installation
```
brew tap hahnlee/brew
brew install stamped
```

# Usage
```
stamped --config config.json --url YOUR_URL --out SAVE_FILE_NAME
```
OR
```
stamped --config config.json --file SWAGGER_JSON.json --out SAVE_FILE_NAME
```

## Config file
```json
{
  "swagger": {
    "version": "2.0"
  },
  "postman": {
    "version": "2.1",
    "host": "HOST-URL"
  }
}
```

# License
[MIT](./LICENSE)
