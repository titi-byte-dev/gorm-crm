package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	anthropicAPI     = "https://api.anthropic.com/v1/messages"
	anthropicVersion = "2023-06-01"
	// claude-haiku-4-5: rápido, barato, suficiente para decisões CRM
	defaultModel = "claude-haiku-4-5-20251001"
)

// LLMClient envia prompts para a API Anthropic e devolve tool calls estruturados.
type LLMClient struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewLLMClient cria o cliente a partir de ANTHROPIC_API_KEY no ambiente.
// Devolve nil se a chave não estiver configurada — agente fica desativado.
func NewLLMClient() *LLMClient {
	key := os.Getenv("ANTHROPIC_API_KEY")
	if key == "" {
		return nil
	}
	model := os.Getenv("ANTHROPIC_MODEL")
	if model == "" {
		model = defaultModel
	}
	return &LLMClient{
		apiKey: key,
		model:  model,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// ToolCall representa uma ação que o LLM decide executar.
type ToolCall struct {
	Name  string         `json:"name"`
	Input map[string]any `json:"input"`
}

// LLMResult é o resultado de uma chamada ao LLM.
type LLMResult struct {
	ToolCalls  []ToolCall
	Summary    string // texto livre devolvido pelo modelo
	TokensUsed int
}

// --- estruturas internas para serialização da API Anthropic ---

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system,omitempty"`
	Messages  []anthropicMessage `json:"messages"`
	Tools     []anthropicTool    `json:"tools"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicTool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"input_schema"`
}

type anthropicResponse struct {
	Content []struct {
		Type  string         `json:"type"`
		Text  string         `json:"text,omitempty"`
		Name  string         `json:"name,omitempty"`
		Input map[string]any `json:"input,omitempty"`
	} `json:"content"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// crmTools define as ferramentas disponíveis ao agente no contexto CRM.
var crmTools = []anthropicTool{
	{
		Name:        "create_task",
		Description: "Cria uma nova tarefa de acompanhamento associada a um contacto ou deal.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"title":    map[string]any{"type": "string", "description": "Título claro e accionável"},
				"priority": map[string]any{"type": "string", "enum": []string{"low", "medium", "high", "urgent"}},
				"due_days": map[string]any{"type": "integer", "description": "Dias a partir de hoje para o prazo"},
			},
			"required": []string{"title", "priority"},
		},
	},
	{
		Name:        "update_lead_status",
		Description: "Muda o status de um lead (ex: new → contacted, contacted → qualified, etc.).",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"new_status": map[string]any{
					"type": "string",
					"enum": []string{"new", "contacted", "qualified", "lost"},
				},
				"reason": map[string]any{"type": "string", "description": "Justificação para a mudança"},
			},
			"required": []string{"new_status"},
		},
	},
	{
		Name:        "add_note",
		Description: "Adiciona uma nota/observação ao histórico da entidade.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"content": map[string]any{"type": "string", "description": "Conteúdo da nota"},
			},
			"required": []string{"content"},
		},
	},
	{
		Name:        "escalate_to_manager",
		Description: "Notifica o manager de que este item requer atenção urgente.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"reason": map[string]any{"type": "string"},
			},
			"required": []string{"reason"},
		},
	},
	{
		Name:        "summarize_only",
		Description: "Devolve apenas um resumo do contexto sem executar ações.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"summary": map[string]any{"type": "string"},
			},
			"required": []string{"summary"},
		},
	},
}

const systemPrompt = `És um assistente de CRM especializado em vendas B2B para equipas portuguesas.
Analisa o contexto fornecido e decide quais as ações mais adequadas para avançar o processo comercial.

Regras:
- Usa sempre linguagem profissional em Português de Portugal
- Nunca propões ações irreversíveis sem justificação clara
- Se não há ação óbvia, usa summarize_only
- Prioriza tarefas urgentes antes das de baixa prioridade
- Máximo 3 ações por execução`

// Run envia o prompt ao LLM e devolve as tool calls resultantes.
func (c *LLMClient) Run(userPrompt string) (*LLMResult, error) {
	reqBody := anthropicRequest{
		Model:     c.model,
		MaxTokens: 1024,
		System:    systemPrompt,
		Tools:     crmTools,
		Messages:  []anthropicMessage{{Role: "user", Content: userPrompt}},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, anthropicAPI, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", anthropicVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call anthropic api: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var ar anthropicResponse
	if err := json.Unmarshal(raw, &ar); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	if ar.Error != nil {
		return nil, fmt.Errorf("anthropic error: %s", ar.Error.Message)
	}

	result := &LLMResult{
		TokensUsed: ar.Usage.InputTokens + ar.Usage.OutputTokens,
	}
	for _, block := range ar.Content {
		switch block.Type {
		case "tool_use":
			result.ToolCalls = append(result.ToolCalls, ToolCall{
				Name:  block.Name,
				Input: block.Input,
			})
		case "text":
			result.Summary += block.Text
		}
	}
	return result, nil
}
