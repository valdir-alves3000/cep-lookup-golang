# API de Consulta de CEP com Go

Este projeto é uma API de consulta de CEP em Go, utilizando as APIs [BrasilAPI](https://brasilapi.com.br/) e [ViaCEP](https://viacep.com.br/), com deploy feito na [Vercel](https://vercel.com/).

## 📑 Estrutura do Projeto

```
/api
  └── cep.go              # Código da API para consulta de CEP
/templates
  └── docs.html           # Documentação da API
/vercel.json              # Configuração do deploy na Vercel
```

## 🚀 Configuração para Deploy na Vercel

Este projeto utiliza um arquivo de configuração `vercel.json` para definir como a API e os arquivos estáticos serão servidos na Vercel.

### Arquivo `vercel.json`

```json
{
  "version": 2,
  "builds": [
    {
      "src": "api/*.go",
      "use": "@vercel/go"
    },
    {
      "src": "templates/**/*.html",
      "use": "@vercel/static"
    }
  ],
  "routes": [
    {
      "src": "/api/docs",
      "dest": "/templates/docs.html"
    },
    {
      "src": "/api/cep",
      "dest": "/api/cep.go"
    },
    {
      "src": "/(.*)",
      "status": 308,
      "headers": { "Location": "/api/docs" }
    }
  ]
}
```

## 📜 Endpoints

### 1. **Documentação da API**
   - **URL**: `/api/docs`
   - Exibe a documentação da API.

### 2. **Consulta de Endereço pelo CEP**
   - **URL**: `/api/cep?cep=<seu_cep>`
   - **Método**: `GET`
   - **Descrição**: Consulta o endereço com base no CEP informado.
   - **Parâmetro**:
     - `cep`: CEP para consulta, no formato `00000-000` ou `00000000`.
   - **Exemplo de Resposta**:
     ```json
     {
       "address": {
         "cep": "01001-000",
         "logradouro": "Praça da Sé",
         "bairro": "Sé",
         "localidade": "São Paulo",
         "uf": "SP"
       },
       "api": "BrasilAPI",
     }
     ```

## 📄 Código da API

O código da API está localizado em `api/cep.go`. A função `Handler` faz o tratamento do CEP e realiza a busca simultânea nas APIs BrasilAPI e ViaCEP.

### Estrutura dos Dados

- A API responde com os dados de endereço usando a estrutura `Address`:
  ```go
  type Address struct {
      CEP        string `json:"cep"`
      Logradouro string `json:"logradouro"`
      Bairro     string `json:"bairro"`
      Localidade string `json:"localidade"`
      UF         string `json:"uf"`
  }
  ```

### Função `Handler`

A função `Handler` realiza a seguinte lógica:
1. Configura o CORS.
2. Recebe o CEP da query string.
3. Inicia duas goroutines que consultam o CEP nas APIs BrasilAPI e ViaCEP.
4. Retorna o primeiro endereço válido encontrado ou exibe uma mensagem de erro caso ocorra um timeout.

### Exemplo de Função para Consultar o Endereço

```go
func handlerAddress(ctx context.Context, url string) (*Address, string, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, "", err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, "", err
    }
    defer resp.Body.Close()

    var address Address
    if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
        return nil, "", err
    }

    return &address, url, nil
}
```

## 🛠️ Como Fazer o Deploy na Vercel

1. Certifique-se de que você possui o [Vercel CLI](https://vercel.com/download) instalado.
2. Faça login na Vercel:
   ```bash
   vercel login
   ```
3. No diretório do projeto, execute:
   ```bash
   vercel deploy
   ```
4. Siga as instruções para configurar o projeto e escolher o ambiente de deploy.

## ⚙️ Tecnologias Utilizadas

- [Go](https://golang.org/) - Linguagem de programação utilizada para a API.
- [Vercel](https://vercel.com/) - Plataforma de deploy para a API e arquivos estáticos.
- [BrasilAPI](https://brasilapi.com.br/) e [ViaCEP](https://viacep.com.br/) - APIs de consulta de CEP.

