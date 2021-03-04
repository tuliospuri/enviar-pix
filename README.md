# (POC) Envio de Pix via API com Golang, Siri e Gerencianet

![](https://img.shields.io/badge/Pix-Gerencianet-orange) ![](https://img.shields.io/github/go-mod/go-version/tuliospuri/enviar-pix?filename=src%2Fapp%2Fgo.mod) ![](https://img.shields.io/badge/App-Heroku-blueviolet) ![](https://img.shields.io/badge/Apple-Siri-blue)

## Vídeo demonstração

https://youtu.be/MXgAtg6XZ3A

![](https://img.youtube.com/vi/MXgAtg6XZ3A/0.jpg)

## Como executar

1. `./build.sh` para gerar a imagem Docker;
2. Criar o `.env` (instruções abaixo);
3. `./run.sh air` compilar e executar o código Golang;

## Variáveis de ambiente

Descrição das variáveis de ambiente do arquivo `.env`.

- `PORT`: porta do servidor HTTP;
- `CERT`: certificado em uma linha;
- `KEY`: chave do certificado em uma linha;
- `API`: URL da API Pix;
- `ID`: client_id da API Pix;
- `SECRET`: client_secret da API Pix;

### Conversão do certificado

Neste [Gist](https://gist.github.com/tuliospuri/bc4abcaee428ba456c0bae5752f7532a) existem instruções para converter um `.p12` em dois arquivos `.pem`. E depois para converter seus valores para uma única linha para ser usado como variáveis de ambiente, exemplo: Docker, Heroku, etc.
