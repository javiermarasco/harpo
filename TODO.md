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
  - [x] Create IAM user
  - [ ] Assign the correct permissions to IAM User to handle secrets
    - [ ] secretsmanager:Name
    - [ ] secretsmanager:Description
    - [ ] secretsmanager:KmsKeyId
    - [ ] aws:RequestTag/${TagKey}
    - [ ] aws:ResourceTag/${TagKey}
    - [ ] aws:TagKeys
    - [ ] secretsmanager:ResourceTag/tag-key
    - [ ] secretsmanager:AddReplicaRegions
    - [ ] secretsmanager:ForceOverwriteReplicaSecret
    - [ ] resourcetypes: Secret*
    - [ ] secretsmanager:TagResource
    - [ ] secretsmanager:UntagResource
  - [ ] Learn how to create and delete secrets using GO SDK
  - [ ] Add support for tagging of secrets
    - Should I use tagresource? (https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_TagResource.html)
    - Should I use untagresource? (https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_UntagResource.html)
  - Information on what is returned as a secret from the SDK: (https://github.com/aws/aws-sdk-go-v2/blob/main/service/secretsmanager/types/types.go)
  - Repo with the SDK implementation (https://github.com/aws/aws-sdk-go-v2/tree/main/service/secretsmanager)

Add management for cloud providers:
 - Ideally the CLI should also be able to configure the needed "stuff" in the cloud provider to be used
   - In Azure create the SPN, assign it to the keyvault, grant permissions and retrieve the SPN clientid, secret, tenant (Also create a keyvault???)
   - In AWS Create the IAM, policy, assign the policy and retrieve the IAM id and secret (Also create the secrets manager?)

Agregar GCP vault
  - [ ] TODO