# feign-go

📦 Biblioteca declarativa de chamadas HTTP estilo OpenFeign, escrita em Go.

## Recursos

- Métodos: GET, POST, PUT, DELETE
- Suporte a Token via `TokenProvider` (Authorization: Bearer)
- Extensível: logging, retry, fallback etc.

## Exemplo

```go
client := feign.NewClientWithToken("http://localhost:8080", tokenProvider)

var user User
client.Get("/users/1", &user)
