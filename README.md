
# 📁 folder-structure-generator

Gere automaticamente **imagens e arquivos `.txt` com a estrutura de diretórios** do seu projeto a partir de um arquivo `structure.yaml`.  
Ideal para documentação técnica, onboarding de times, ou ilustrações em README e apresentações.

---

## ✅ Exemplos de saída

```
src/
├─ api/
│  └─ cupons/
│     ├─ controller.js
│     └─ routes.js
├─ integrations/
│  ├─ redis.js
│  └─ cupons/
│     ├─ constants.js
│     └─ cuponsService.js
└─ middleware/
   └─ cache.js
```

Também é gerada automaticamente a imagem `.png` com tema personalizado.

---

## 📦 Instalação

```bash
# 1. Inicializa o projeto com um nome qualquer
go mod init treegen

# 2. Instala os pacotes usados pelo seu script
go get gopkg.in/yaml.v3
go get golang.org/x/image/...

# 3. Agora sim, execute
go run main.go -theme=dracula
```

---

## 🚀 Como usar

1. Crie um arquivo `structure.yaml` na raiz do projeto com o seguinte formato:

```yaml
src:
  api:
    cupons:
      __files__:
        - controller.js
        - routes.js
  integrations:
    __files__:
      - redis.js
    cupons:
      __files__:
        - constants.js
        - cuponsService.js
  middleware:
    __files__:
      - cache.js
```

> **Observações**  
> - Use `__files__:` para listar arquivos em uma pasta que também tem subpastas.  
> - Ou use `arquivo.js:` com valor `null` se preferir estilo mais compacto.

2. Execute:

```bash
go run main.go -theme=dracula
```

---

## 🎨 Temas disponíveis

- `dark` (padrão)
- `light`
- `dracula`

Você pode especificar o tema assim:

```bash
go run main.go -theme=light
```

---

## 📂 Output

Os arquivos são gerados automaticamente na pasta `/output`:

```
output/
├─ folder_structure.txt   # Representação em árvore
└─ folder_structure.png   # Imagem colorida estilo terminal
```

---

## 🛠 Estrutura

- `main.go` – código principal do gerador
- `structure.yaml` – arquivo de entrada com a estrutura da árvore
- `output/` – onde os arquivos gerados são salvos

---

## 📋 Exemplo de uso em projetos

- Documentação técnica (em PDF, Notion, Confluence)
- READMEs de repositórios
- Scripts de onboarding
- Apresentações de arquitetura

---

## 🧪 Testado com

- Go 1.21+
- VS Code e terminal padrão
- UTF-8 compatível (sem dependências nativas)

---

## 📄 2025 – Desenvolvido por Rodrigo Godoi 🧙🏼‍♂️

