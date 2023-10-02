# alp-backup-system

Esse projeto foi desenvolvido para a disciplina de Aspectos de Linguagem de Programação, do Curso de BCC da UTFPR Campo Mourão, e tem como intuito criar um sistema simplificado para backup incremental de arquivos utilizando recursos da Linguagem Go.

## Clonando o Projeto
```zsh
# Via HTTPS
git clone https://github.com/afmireski/alp-backup-system.git
# ou via SSH
git clone git@github.com:afmireski/alp-backup-system.git
```

## Rodando o projeto
```zsh
# Vá até a pasta do projeto
cd alp-backup-system/

# Faça download dos módulos necessários
go mod tidy

# Execute main.go
go run ./cmd/server/main.go
```
