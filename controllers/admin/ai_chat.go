package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"glk-web-app/config"
	"glk-web-app/models"
	"io"
	"net/http"
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

	return c.Render("admin/ai/chat", contextData(c, fiber.Map{
		"Title":      "AI Chat Assistant",
		"Breadcrumb": "AI Mode / Chat",
		"MCPs":       mcps,
		"Skills":     skills,
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
	systemPrompt := "You are a helpful HR and Company assistant for Gurulesku."

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
	json.NewDecoder(resp.Body).Decode(&aiResp)

	aiMessage := ""
	if len(aiResp.Choices) > 0 {
		aiMessage = aiResp.Choices[0].Message.Content
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": aiMessage,
	})
}

// fetchFileContent tries to download content from GDrive URL (heuristic/proxy)
func fetchFileContent(driveURL string) (string, error) {
	// Heuristic: Many AI models work better if we at least provide the metadata.
	// For full content, we'd need to convert GDrive URL to export link if it's a doc.
	// For now, let's return a note that the link is attached.
	return "Document link: " + driveURL, nil
}
