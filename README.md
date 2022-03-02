Objetivo:

- Crear un CLI para azure Keyvault que permita manejar secrets como si estuvieran en carpetas
- Utilizar GO como lenguaje de programación
- Usar un SPN para el acceso a los secrets de keyvault
- Inicialmente solo manejo de secrets 
  - A futuro mismo funcionamiento pero para keys y certs
- Agregar soporte para AWS
- Agregar soporte para GCP

# Pueden seguir el progreso de este proyecto en mis streams en twitch o Youtube (Para mas información pueden revisar aquí -> https://linktr.ee/javi__codes)


# You can follow the progress of this project on my live streams on Twitch or Youtube (for more information check here -> https://linktr.ee/javi__codes )

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
  - [ ] Read de secrets
    - [ ] Usando argumentos y flags desde la linea de comandos
  - [ ] Manejo de Keys
  - [ ] Manejo de Certs

Agregar AWS vault
  - [ ] TODO
Agregar GCP vault
  - [ ] TODO
