# Especificações da Tela de Login - JP Mall

Este documento contém o prompt detalhado para geração da tela de login e a descrição completa de sua estrutura lógica e visual, baseada na identidade visual atual do sistema.

---

## 1. Prompt para Geração (HTML/CSS)

Copie e cole o texto abaixo em uma IA geradora de código para obter a implementação:

> "Gere o código HTML e CSS (utilizando Tailwind CSS ou CSS puro) para uma tela de login premium e profissional para o sistema 'JP Mall - Sistema de Conservação & Manutenção Predial'. 
> 
> **Identidade Visual:**
> - **Cores Primárias:** #8B1A1A (Vinho) e #C8A882 (Dourado/Bege).
> - **Fundo:** #F7F4EF (Off-white) para o modo claro; #1a1a1a para o modo escuro.
> - **Estilo:** Design limpo, corporativo, com bordas arredondadas (10px/0.625rem) e sombras suaves.
> - **Tipografia:** Sans-serif moderna (ex: Inter ou Roboto).
> 
> **Elementos da Página:**
> 1. **Cabeçalho:** Ícone de prédio (Building2) em um gradiente de #8B1A1A para #a43030, título 'JP Mall' em negrito e subtítulo descritivo.
> 2. **Card de Login:** Fundo branco (ou #252525 no escuro) com uma borda superior de 4px na cor primária (#8B1A1A).
> 3. **Alerta de Segurança:** Um box informativo (estilo warning) com ícone de alerta, indicando que o acesso é restrito a administradores autorizados.
> 4. **Formulário:** 
>    - Campo de 'E-mail ou Usuário' com ícone lateral (Mail).
>    - Campo de 'Senha' com ícone lateral (Lock) e suporte a ocultar/mostrar (opcional).
>    - Checkbox 'Lembrar acesso' e link 'Esqueceu a senha?'.
> 5. **Botão de Ação:** Botão 'Entrar no Sistema' ocupando a largura total, com efeito hover suave e cor #8B1A1A.
> 
> **Requisitos Técnicos:**
> - Layout totalmente responsivo.
> - Suporte a Dark Mode através da classe '.dark'.
> - Micro-interações (hover no botão e inputs).
> - Código semântico e acessível."

---

## 2. Estrutura Lógica (Funcionalidade)

A tela de login foi projetada para ser o portão de entrada seguro do sistema JP Mall. Seu funcionamento lógico segue os seguintes princípios:

### Fluxo de Autenticação
1. **Entrada de Dados:** O sistema aceita tanto o e-mail corporativo quanto o nome de usuário único.
2. **Validação em Tempo Real:** 
   - Verificação de formato de e-mail.
   - Verificação de campos vazios antes da submissão.
3. **Persistência (Remember Me):** Ao marcar "Lembrar acesso", o sistema deve armazenar o token de sessão ou o identificador do usuário em `localStorage` ou via `Cookies` com flag `HttpOnly` (dependendo da implementação do backend).
4. **Recuperação de Acesso:** O link "Esqueceu a senha?" redireciona para um fluxo de recuperação via e-mail cadastrado.
5. **Segurança:** O alerta visual reforça a política de uso do sistema, desencorajando tentativas de acesso não autorizado por terceiros.

---

## 3. Estrutura Visual (Design System)

A estética busca transmitir **confiabilidade, robustez e elegância**, alinhada a um sistema de gestão predial de alto padrão.

### Layout e Composição
- **Centralização:** O card de login é centralizado vertical e horizontalmente para focar a atenção do usuário no centro da tela.
- **Espaçamento (Padding/Margin):** Utilização de escalas consistentes (p-4, p-8, mt-6) para garantir respiro visual.
- **Elevação:** O card utiliza `box-shadow` profundo no modo claro e sombras sutis no modo escuro para criar profundidade.

### Cores e Contrastes
- **Hierarquia Visual:** O vinho (#8B1A1A) é usado para as ações principais (botões, bordas de destaque) e o dourado (#C8A882) para elementos decorativos ou ícones de marca, criando um contraste luxuoso sobre o fundo off-white.
- **Modo Escuro:** A paleta transita para tons de cinza carvão e preto, mantendo a legibilidade com textos em cinza claro e aumentando a vibrância do vermelho primário para #D93030.

### Elementos Gráficos
- **Ícones:** Uso da biblioteca `Lucide` para representação visual dos campos (e-mail, senha) e do alerta de segurança.
- **Bordas:** O `border-radius` de `0.625rem` é aplicado em todos os elementos (cards, inputs, botões) para suavizar a interface.
- **Gradientes:** Pequenos gradientes aplicados apenas no ícone da marca para dar um toque moderno sem sobrecarregar o design flat do restante da página.
