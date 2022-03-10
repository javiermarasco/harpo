Implementación

- [x] Crear una keyvault
- [x] Crear un SPN para manejar los secrets del keyvault
- [x] Dar permisos de secrets al SPN
- [x] Crear un packet o modulo de GO para la solución
- [x] Credenciales para acceder a keyvault como environment variables
  - [x] ClientID
  - [x] Secret
  - [x] TenantID
  - [x] VaultName
- [x] Inputs
  - [x] Path: "/infra/dev/cluster"
  - [x] SecretName: "Port"
  - [x] Value: "80"
- [ ] Tareas:
  - [x] Split del path por "/" (separador)
  - [x] Crear tags en base al path
    - [x] A: infra
    - [x] B: dev
    - [x] C: cluster
    - [x] SecretName: port
    - [x] Valor en el secret --> Value(80)
  - [x] Nombre del secret es un hash compuesto por: (Esto nos va a permitir mantener multiples versiones del mismo secret)
    - [x] Todos los tags + SecretName
  - [x] Write de secrets
    - [x] Usando argumentos y flags desde la linea de comandos
  - [x] Read de secrets
    - [x] Usando argumentos y flags desde la linea de comandos
  - [ ] Manejo de Keys
  - [ ] Manejo de Certs
  
Agregar AWS vault
  - [ ] TODO
Agregar GCP vault
  - [ ] TODO