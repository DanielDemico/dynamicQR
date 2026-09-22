AGENTS.md — Spec-Driven Development

Este projeto utiliza specs como fonte de verdade para requisitos, comportamento esperado, implementação e validação.

O agente deve seguir as regras abaixo sempre que criar, alterar, implementar ou validar uma funcionalidade.

1. Organização obrigatória das Specs

Todas as specs devem estar dentro de uma pasta _SPECS.

A estrutura recomendada é:

_SPECS/
├── INDEX.md
├── <nome-da-spec-1>/
│   ├── SPEC.md
│   └── ...arquivos auxiliares, quando necessários
├── <nome-da-spec-2>/
│   ├── SPEC.md
│   └── ...
└── ...

Regras

Nunca criar uma spec fora de _SPECS.

Cada spec deve possuir seu próprio diretório dentro de _SPECS.

O arquivo principal de uma spec deve ser SPEC.md.

O nome do diretório deve ser curto, descritivo e estável. Evite nomes genéricos como new-spec ou test.

Arquivos auxiliares relacionados exclusivamente àquela spec devem permanecer dentro do diretório da própria spec.

Exemplo:

_SPECS/
└── autenticacao-jwt/
    ├── SPEC.md
    ├── examples.md
    └── fixtures.json

2. Toda Spec deve declarar seu estado

Toda SPEC.md deve informar explicitamente:

se foi aprovada;

se já foi implementada.

Utilize o seguinte cabeçalho mínimo:

# <Nome da Spec>

## Status

- **Aprovada:** não
- **Implementada:** não

Os valores devem ser atualizados durante o ciclo de vida da spec.

Estados válidos:

Aprovada: não | sim
Implementada: não | parcialmente | sim

parcialmente pode ser utilizado quando apenas parte da especificação foi implementada. Nesse caso, a própria spec deve deixar claro quais partes ainda faltam.

Regra importante

Uma spec não deve ser considerada aprovada apenas porque o agente conseguiu implementá-la.

Aprovação e implementação são estados diferentes:

Proposta → Aprovada → Em implementação → Implementada → Validada

A existência da funcionalidade no código não substitui a aprovação da spec.

3. Índice central de Specs

A pasta _SPECS deve possuir um arquivo:

_SPECS/INDEX.md

Esse arquivo funciona como o índice oficial de todas as specs do projeto.

Sempre que uma spec for criada, alterada de estado, implementada ou invalidada, o índice deve ser atualizado.

Template do índice

# Índice de Specs

| ID | Spec | Diretório | Aprovada | Implementada | E2E | Status | Observações |
|---|---|---|---|---|---|---|---|
| SPEC-001 | Nome da funcionalidade | `./nome-da-spec/` | não | não | pendente | proposta | Descrição curta |

## Convenções

- **Aprovada:** `não` ou `sim`
- **Implementada:** `não`, `parcialmente` ou `sim`
- **E2E:** `não aplicável`, `pendente`, `criado`, `passando` ou `falhando`
- **Status:** `proposta`, `aprovada`, `em implementação`, `implementada`, `validada`, `bloqueada` ou `cancelada`

Regras do índice

O agente deve:

criar uma entrada no INDEX.md ao criar uma nova spec;

atribuir um identificador único, como SPEC-001, SPEC-002, etc.;

manter o caminho da spec atualizado;

atualizar os estados conforme o trabalho evolui;

registrar observações importantes, limitações ou bloqueios;

não remover uma spec do índice apenas porque ela foi concluída, cancelada ou substituída.

O índice representa o histórico e inventário das specs, não apenas as que estão pendentes.

4. Template obrigatório de uma Spec

Novas specs devem seguir, no mínimo, esta estrutura:

# <Título da Spec>

## Metadata

- **ID:** SPEC-XXX
- **Aprovada:** não
- **Implementada:** não
- **Criada em:** YYYY-MM-DD
- **Atualizada em:** YYYY-MM-DD

## Objetivo

Descreva o problema que esta spec resolve e o resultado esperado.

## Contexto

Explique informações relevantes, restrições, decisões existentes e dependências.

## Requisitos

1. ...
2. ...
3. ...

## Comportamento esperado

Descreva o comportamento observável do sistema.

## Critérios de aceitação

- [ ] ...
- [ ] ...
- [ ] ...

## Casos de erro / bordas

- ...

## E2E

- **Script:** `...`
- **Comando:** `...`
- **Status:** pendente

## Notas de implementação

Detalhes técnicos que não são requisitos, mas ajudam a implementação.

## Pendências

- ...

O agente pode adicionar seções quando necessário, mas não deve remover as informações essenciais de status, comportamento esperado e critérios de aceitação.

5. Fluxo obrigatório ao criar ou implementar uma Spec

Ao criar uma Spec

Criar a pasta _SPECS/<nome-da-spec>/.

Criar _SPECS/<nome-da-spec>/SPEC.md usando o template.

Atribuir um ID único.

Marcar Aprovada: não e Implementada: não.

Atualizar _SPECS/INDEX.md.

Ao aprovar uma Spec

Alterar Aprovada para sim dentro da SPEC.md.

Atualizar o estado correspondente no INDEX.md.

Não alterar os requisitos aprovados silenciosamente.

Caso os requisitos mudem materialmente após a aprovação, registrar a mudança na própria spec e atualizar seu estado conforme a necessidade do projeto.

Ao implementar uma Spec

Ler a SPEC.md completa antes de alterar o código.

Verificar seus requisitos e critérios de aceitação.

Implementar o comportamento descrito.

Criar ou atualizar os testes automatizados necessários.

Criar o teste E2E descrito na seção de E2E.

Executar a validação.

Atualizar Implementada e o INDEX.md de acordo com o resultado real.

Não marcar como implementada apenas porque o código compila.

6. Testes E2E são parte da Spec

O sistema deve criar testes end-to-end (E2E) para validar o comportamento esperado sempre que a funcionalidade possuir um fluxo observável de ponta a ponta.

Esses testes devem ser executados por scripts reproduzíveis, e não depender de validação manual como única evidência.

Objetivo

O E2E deve responder, de maneira automatizada:

"Dado o estado inicial esperado, o sistema produz a resposta/comportamento definido na spec?"

Regras

Cada spec aplicável deve possuir pelo menos um cenário E2E relevante.

O E2E deve validar resultado observável, não detalhes internos desnecessários.

O teste deve usar entradas representativas e verificar as saídas esperadas.

O script deve retornar código de erro quando a validação falhar.

O comando para executar o teste deve estar documentado na SPEC.md.

O agente deve executar o script após implementar a funcionalidade.

Quando um E2E falhar, o agente deve tratar a falha como evidência de que a spec ainda não está validada.

Estrutura sugerida

scripts/
└── e2e/
    ├── spec-001-autenticacao.sh
    ├── spec-002-criacao-usuario.sh
    └── ...

Também é permitido usar Python, Node.js, PowerShell ou outra tecnologia adequada, desde que o script seja reproduzível e tenha uma interface clara de execução.

Exemplo:

./scripts/e2e/spec-001-autenticacao.sh

O script deve indicar claramente sucesso ou falha.

7. Relação entre Spec, implementação e E2E

Considere os três como partes de um mesmo ciclo:

SPEC.md
   ↓
Implementação
   ↓
Testes automatizados
   ↓
E2E por script
   ↓
Resultado observável
   ↓
Atualização da Spec + INDEX.md

Uma implementação que não passa pelos critérios de aceitação não deve ser registrada como concluída.

Um E2E que passa não significa automaticamente que todos os requisitos da spec foram atendidos; os critérios de aceitação continuam sendo a referência.

8. Quando solicitar uma nova Skill

O agente pode identificar situações nas quais a execução recorrente de uma tarefa exige conhecimento, ferramentas ou procedimentos que ainda não estão bem definidos.

Quando isso acontecer, o agente deve solicitar a criação de uma Skill em vez de improvisar um processo complexo repetidamente.

Uma nova Skill deve ser considerada quando houver, por exemplo:

uma tarefa recorrente com um procedimento específico;

necessidade de uma ferramenta ou integração especializada;

conhecimento de domínio que deva ser reutilizado em várias specs;

uma sequência de validações que precise ser padronizada;

um processo complexo que esteja sendo repetido e gerando inconsistência;

necessidade de encapsular comandos, padrões ou boas práticas do projeto.

Formato da solicitação

Quando identificar essa necessidade, registre na spec ou no índice algo como:

### Skill necessária

- **Nome sugerido:** <nome-da-skill>
- **Motivo:** <por que a skill é necessária>
- **Objetivo:** <o que ela deverá padronizar ou automatizar>
- **Escopo:** <quais tarefas deverão utilizá-la>

A solicitação de Skill não deve bloquear todo o trabalho automaticamente. O agente deve avaliar se existe uma solução simples e segura para continuar e, quando a especialização for realmente necessária, registrar explicitamente a necessidade.

9. Regras de consistência

O agente deve sempre preservar a consistência entre:

_SPECS/<spec>/SPEC.md
_SPECS/INDEX.md
código implementado
scripts de E2E

Se qualquer um desses elementos ficar desatualizado, o trabalho não está completo.

Alterações relevantes de comportamento devem primeiro ser refletidas na spec ou, quando a mudança surgir durante a implementação, registradas como alteração explícita da especificação.

O agente não deve inventar aprovação, implementação ou sucesso de testes. Os estados devem representar o estado real do projeto.

10. Checklist operacional do Agente

Antes de considerar uma tarefa concluída, verificar:

[ ] A spec está dentro de `_SPECS/`.
[ ] A spec possui `SPEC.md`.
[ ] A spec possui ID único.
[ ] A spec informa se foi aprovada.
[ ] A spec informa se foi implementada.
[ ] `_SPECS/INDEX.md` possui a entrada correspondente.
[ ] Os critérios de aceitação estão definidos.
[ ] Existe E2E quando a funcionalidade é testável ponta a ponta.
[ ] O E2E é executável por script.
[ ] O E2E foi executado.
[ ] O resultado real foi registrado.
[ ] Uma Skill foi solicitada quando a especialização recorrente se mostrou necessária.

10. instructions.md — Guia operacional do sistema

O projeto deve possuir um arquivo instructions.md na raiz do repositório. Esse arquivo é a referência operacional para qualquer pessoa ou agente que precise instalar, configurar, executar, testar, validar ou depurar o sistema.

O instructions.md deve conter tudo o que for necessário para rodar o sistema de forma reproduzível, evitando que informações essenciais fiquem apenas implícitas no código, na memória do desenvolvedor ou em mensagens de execução.

Conteúdo mínimo esperado

O agente deve criar e manter o instructions.md com, quando aplicável:

pré-requisitos de sistema;

versões de linguagens, runtimes, SDKs e ferramentas;

dependências do projeto;

instruções de instalação;

configuração de variáveis de ambiente;

arquivos .env necessários e quais valores precisam ser fornecidos;

serviços externos necessários;

bancos de dados e como inicializá-los;

migrações;

comandos para iniciar o ambiente de desenvolvimento;

comandos para executar o sistema em produção, quando aplicável;

portas utilizadas e serviços correspondentes;

comandos para executar testes automatizados;

comandos para executar os testes E2E;

comandos para build;

comandos úteis para logs, diagnóstico e troubleshooting;

estrutura relevante dos diretórios;

dependências entre serviços;

procedimentos especiais necessários para reproduzir o ambiente;

limitações ou configurações que possam causar falhas comuns.

Estrutura sugerida

# Instructions

## Pré-requisitos

## Instalação

## Configuração

## Variáveis de ambiente

## Dependências externas

## Banco de dados

## Como executar

## Como executar os testes

## Como executar os E2E

## Build / Deploy

## Serviços e portas

## Troubleshooting

## Estrutura do projeto

Regras

Criar instructions.md quando ele ainda não existir.

Atualizar instructions.md sempre que uma alteração mudar a forma de instalar, configurar, executar, testar ou depurar o sistema.

Não documentar apenas o caminho feliz: incluir pré-requisitos e configurações necessárias para uma execução limpa.

Nunca colocar segredos reais no instructions.md; documentar apenas quais variáveis são necessárias e como obtê-las/configurá-las de forma segura.

Os comandos documentados devem ser compatíveis com o estado atual do projeto. O agente deve evitar comandos fictícios ou desatualizados.

Sempre que possível, preferir comandos que possam ser copiados e executados diretamente.

Se houver múltiplos ambientes (por exemplo, desenvolvimento, testes e produção), documentar separadamente como cada um deve ser executado.

Caso uma spec introduza uma nova dependência operacional, atualização de ambiente, novo serviço, nova porta ou novo procedimento de execução, o instructions.md deve ser atualizado junto com a implementação.

Relação com Specs

O instructions.md não substitui as specs. Ele descreve como rodar e operar o sistema; as specs descrevem o comportamento e os requisitos que o sistema deve cumprir.

Quando uma alteração de uma spec exigir alguma mudança operacional, o trabalho só deve ser considerado completo quando ambos estiverem coerentes:

SPEC.md
   ↓
Implementação
   ↓
instructions.md atualizado
   ↓
Testes automatizados / E2E

Se o agente perceber que não consegue executar ou validar o sistema porque faltam instruções, dependências ou procedimentos documentados, deve primeiro registrar e, quando possível, corrigir essa lacuna no instructions.md.

11. Princípio central

A spec descreve o que deve acontecer; o código implementa; o E2E comprova o comportamento; o índice registra o estado do trabalho.

O agente deve tratar _SPECS como parte do sistema, e não como documentação descartável.