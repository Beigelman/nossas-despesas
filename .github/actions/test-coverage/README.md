# Test Coverage Action

Este action executa testes Go com cobertura de código, valida se a cobertura está acima de um limite mínimo e opcionalmente comenta no PR com o relatório.

## Funcionalidades

- ✅ Executa testes Go com cobertura
- ✅ Gera relatório detalhado em Markdown
- ✅ Classifica cobertura como Excelente (≥80%), Boa (≥60%) ou Baixa (<60%)
- ✅ **Comenta automaticamente no PR** com o relatório (configurável)
- ✅ **Valida limite mínimo** após comentar no PR
- ✅ Funciona em qualquer evento (PR, push, etc.)

## Ordem de Execução

O action executa os steps na seguinte ordem:

1. **Install tparse** - Instala ferramenta para parsing de testes
2. **Run tests with coverage** - Executa os testes Go
3. **Parse coverage and generate report** - Processa cobertura e gera relatório
4. **Comment on PR** - Comenta no PR (se habilitado e for um PR)
5. **Validate coverage threshold** - Valida se cobertura está acima do limite

**Importante**: A validação do limite mínimo acontece **após** o comentário no PR, garantindo que o relatório seja sempre postado.

## Inputs

| Input                | Descrição                           | Obrigatório | Padrão                |
| -------------------- | ----------------------------------- | ----------- | --------------------- |
| `coverage-threshold` | Limite mínimo de cobertura (0-100)  | Não         | `40`                  |
| `test-packages`      | Pacotes para testar                 | Não         | `./internal/...`      |
| `comment-on-pr`      | Comentar no PR com o relatório      | Não         | `true`                |
| `github-token`       | Token do GitHub para comentar no PR | Não         | `${{ github.token }}` |

## Outputs

| Output                | Descrição                                 |
| --------------------- | ----------------------------------------- |
| `coverage-percentage` | Percentual de cobertura obtido            |
| `coverage-status`     | Status da cobertura (Excelente/Boa/Baixa) |
| `coverage-emoji`      | Emoji representando o status da cobertura |

## Exemplo de Uso

### Uso básico (com comentário automático no PR):
```yaml
- name: Run tests with coverage
  id: test-coverage
  uses: ./.github/actions/test-coverage
```

### Uso com configurações personalizadas:
```yaml
- name: Run tests with coverage
  id: test-coverage
  uses: ./.github/actions/test-coverage
  with:
    coverage-threshold: '70'
    working-directory: .
    test-packages: './internal/...'
    comment-on-pr: 'true'
```

### Uso sem comentário no PR:
```yaml
- name: Run tests with coverage
  id: test-coverage
  uses: ./.github/actions/test-coverage
  with:
    comment-on-pr: 'false'
```

## Comportamento

- **Cobertura ≥ 80%**: 🟢 Excelente
- **Cobertura ≥ 60%**: 🟡 Boa  
- **Cobertura < 60%**: 🔴 Baixa

### Validação de Limite

Se a cobertura estiver abaixo do `coverage-threshold`, o action falha na **última etapa**, mas:
- ✅ O relatório é gerado normalmente
- ✅ O comentário no PR é postado (se habilitado)
- ❌ O pipeline falha na validação final

### Comentário no PR

O comentário no PR só é feito quando:
- `comment-on-pr` está definido como `true` (padrão)
- O evento é um `pull_request`
- O `github-token` tem permissões adequadas

## Dependências

- Go 1.23.8+
- `tparse` (instalado automaticamente)
- `bc` (para cálculos matemáticos)
- `gh` (GitHub CLI, disponível no runner do GitHub Actions) 