package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"glk-web-app/config"
	"glk-web-app/models"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ShowChatPage renders the AI Chat interface
func ShowChatPage(c *fiber.Ctx) error {
	// Fetch MCPs and Skills for the dropdowns
	var mcps []models.CompanyMCP
	var skills []models.AIProfilingSkill
	config.DB.Find(&mcps)
	config.DB.Find(&skills)

	// Fetch Monthly Token Usage from Redis
	monthKey := fmt.Sprintf("ai:usage:tokens:%s", time.Now().Format("2006-01"))
	tokenUsage, _ := config.RDB.Get(config.Ctx, monthKey).Int64()

	return c.Render("admin/ai/chat", contextData(c, fiber.Map{
		"Title":       "AI Chat Assistant",
		"Breadcrumb":  "AI Mode / Chat",
		"MCPs":        mcps,
		"Skills":      skills,
		"TokensMonth": tokenUsage,
	}), "layouts/base")
}

// ChatRequest represents the incoming chat request from the frontend
type ChatRequest struct {
	Message   string `json:"message"`
	MCPID     uint   `json:"mcp_id,omitempty"`
	SkillID   uint   `json:"skill_id,omitempty"`
	SessionID string `json:"session_id"`
}

// ChatCompletionResponse represents the response from Aivene API
type ChatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
			Role    string `json:"role"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int64 `json:"prompt_tokens"`
		CompletionTokens int64 `json:"completion_tokens"`
		TotalTokens      int64 `json:"total_tokens"`
	} `json:"usage"`
}

// ProcessChat handles the communication with Aivene API
func ProcessChat(c *fiber.Ctx) error {
	var req ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid request"})
	}

	apiKey := config.GetEnv("AIVENE_API_KEY", "")
	if apiKey == "" || apiKey == "your_api_key_here" {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "AIVENE_API_KEY is not configured in .env"})
	}

	// 1. Construct the System Prompt based on MCP and Skill
	systemPrompt := `You are the "Gurulesku AI Expert", a highly specialized HR and Recruitment Intelligence system.
Your persona is professional, analytical, and direct.

RULES:
- DO NOT provide generic information about tutoring or general education unless specifically asked.
- ALWAYS prioritize using the provided COMPANY CONTEXT and PROFILING GUIDELINE.
- Keep responses concise, structured, and focused on recruitment data or company policy.
- If you don't know the answer based on the context, state it clearly rather than giving general advice.
- Avoid using long generic bullet points for common knowledge.`

	if req.MCPID > 0 {
		var mcp models.CompanyMCP
		if err := config.DB.First(&mcp, req.MCPID).Error; err == nil {
			content, _ := fetchFileContent(mcp.URL)
			systemPrompt += fmt.Sprintf("\n\nCOMPANY CONTEXT (%s):\n%s\n%s", mcp.Name, mcp.Keterangan, content)
		}
	}

	if req.SkillID > 0 {
		var skill models.AIProfilingSkill
		if err := config.DB.First(&skill, req.SkillID).Error; err == nil {
			content, _ := fetchFileContent(skill.URL)
			systemPrompt += fmt.Sprintf("\n\nPROFILING GUIDELINE (%s):\n%s\n%s", skill.Name, skill.Keterangan, content)
		}
	}

	// 2. Prepare payload for Aivene API
	// Note: For simplicity, we'll send system prompt + user message.
	// For actual multi-turn, we should fetch history from Redis.
	messages := []fiber.Map{
		{"role": "system", "content": systemPrompt},
		{"role": "user", "content": req.Message},
	}

	payload := fiber.Map{
		"model":    "gemini-2.5-flash-lite",
		"messages": messages,
	}

	payloadBytes, _ := json.Marshal(payload)

	client := &http.Client{Timeout: 60 * time.Second}
	apiReq, _ := http.NewRequest("POST", "https://api.aivene.com/v1/chat/completions", bytes.NewBuffer(payloadBytes))
	apiReq.Header.Set("Content-Type", "application/json")
	apiReq.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(apiReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to contact AI API: " + err.Error()})
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return c.Status(resp.StatusCode).JSON(fiber.Map{"success": false, "message": "AI API Error: " + string(body)})
	}

	var aiResp ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to parse AI response"})
	}

	aiMessage := ""
	if len(aiResp.Choices) > 0 {
		aiMessage = aiResp.Choices[0].Message.Content
	}

	// 3. Track Usage in Redis
	if aiResp.Usage.TotalTokens > 0 {
		monthKey := fmt.Sprintf("ai:usage:tokens:%s", time.Now().Format("2006-01"))
		config.RDB.IncrBy(config.Ctx, monthKey, aiResp.Usage.TotalTokens)
		// Set expiry to 60 days to keep history for a bit
		config.RDB.Expire(config.Ctx, monthKey, 60*24*time.Hour)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": aiMessage,
		"usage":   aiResp.Usage,
	})
}

// fetchFileContent actually downloads the content from Google Drive URL
func fetchFileContent(driveURL string) (string, error) {
	if driveURL == "" {
		return "", nil
	}

	// 1. Convert GDrive Preview Link to Direct Download Link
	// From: https://drive.google.com/file/d/FILE_ID/view?usp=drivesdk
	// To: https://drive.google.com/uc?export=download&id=FILE_ID
	fileID := ""
	if strings.Contains(driveURL, "drive.google.com") {
		parts := strings.Split(driveURL, "/")
		for i, part := range parts {
			if part == "d" && i+1 < len(parts) {
				fileID = parts[i+1]
				break
			}
		}
	}

	if fileID == "" {
		return "Document link: " + driveURL, nil
	}

	downloadURL := fmt.Sprintf("https://drive.google.com/uc?export=download&id=%s", fileID)

	// 2. Fetch the content
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(downloadURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "Note: Document content could not be fetched (Access Denied or Not Downloadable).", nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 3. Return the content (limited to 50k chars to avoid token limit)
	content := string(body)
	if len(content) > 50000 {
		content = content[:50000] + "... (truncated)"
	}

	return content, nil
}
