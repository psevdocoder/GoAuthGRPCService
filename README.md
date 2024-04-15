# GoAuthGRPCService
Сервис аутентификации и авторизации использующий gRPC. Генерирует пару токенов Access-Refresh (где Access генерируется по стандарту JWT) регистрирует пользователей и выдает пользовательскую роль в системе.

<details>
  <summary>Хендлеры сервиса</summary>

1. `Register`:
   - Request: `RegisterRequest(username: string, password: string)`
   - Response: `RegisterResponse(user_id: uint32)`

2. `Login`:
   - Request: `LoginRequest(username: string, password: string, app_id: uint32)`
   - Response: `LoginResponse(access_token: string, refresh_token: string)`

3. `GetRole`:
   - Request: `GetRoleRequest(user_id: uint32)`
   - Response: `GetRoleResponse(role: uint32)`
</details>
