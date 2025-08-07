package ai

import (
	"encoding/json"
	"fmt"

	"org.thinkinai.com/recruit-center/pkg/config"

	"github.com/go-resty/resty/v2"
)

// 提示词模板
const resumeParsePrompt = `
你现在是一个专业的简历解析助手，需要帮我解析一份简历文件。请按照以下要求进行解析：

1. 从简历中提取关键信息，包括：基本信息、教育经历、工作经历、项目经历
2. 所有提取的信息必须严格按照JSON格式返回
3. 时间格式统一使用：YYYY-MM-DD
4. 数据必须完整且准确，不要添加不存在的信息

简历内容如下：
%s

请按照以下JSON格式返回解析结果：
{
    "basicInfo": {
        "name": "姓名",
        "phone": "手机号",
        "email": "邮箱",
        "gender": "性别(1-男,2-女)",
        "location": "所在地",
        "experience": "工作年限",
        "expectedJob": "期望职位",
        "expectedCity": "期望城市",
        "introduction": "个人简介",
        "skills": "技能描述"
    },
    "education": [{
        "school": "学校名称",
        "major": "专业",
        "degree": "学位",
        "startTime": "开始时间",
        "endTime": "结束时间"
    }],
    "workExperience": [{
        "companyName": "公司名称",
        "position": "职位",
        "department": "部门",
        "startTime": "开始时间",
        "endTime": "结束时间",
        "description": "工作描述",
        "achievement": "工作成就"
    }],
    "projects": [{
        "name": "项目名称",
        "role": "担任角色",
        "startTime": "开始时间",
        "endTime": "结束时间",
        "description": "项目描述",
        "technology": "使用技术",
        "achievement": "项目成就"
    }]
}`

// ResumeParseResult AI解析简历的返回结果
type ResumeParseResult struct {
	BasicInfo struct {
		Name         string `json:"name"`
		Phone        string `json:"phone"`
		Email        string `json:"email"`
		Gender       int    `json:"gender"`
		Location     string `json:"location"`
		Experience   int    `json:"experience"`
		ExpectedJob  string `json:"expectedJob"`
		ExpectedCity string `json:"expectedCity"`
		Introduction string `json:"introduction"`
		Skills       string `json:"skills"`
	} `json:"basicInfo"`

	Education []struct {
		School    string `json:"school"`
		Major     string `json:"major"`
		Degree    string `json:"degree"`
		StartTime string `json:"startTime"`
		EndTime   string `json:"endTime"`
	} `json:"education"`

	WorkExperience []struct {
		CompanyName string `json:"companyName"`
		Position    string `json:"position"`
		Department  string `json:"department"`
		StartTime   string `json:"startTime"`
		EndTime     string `json:"endTime"`
		Description string `json:"description"`
		Achievement string `json:"achievement"`
	} `json:"workExperience"`

	Projects []struct {
		Name        string `json:"name"`
		Role        string `json:"role"`
		StartTime   string `json:"startTime"`
		EndTime     string `json:"endTime"`
		Description string `json:"description"`
		Technology  string `json:"technology"`
		Achievement string `json:"achievement"`
	} `json:"projects"`
}

// AIResponse AI接口返回结构
type AIResponse struct {
	ID      string `json:"id"`
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

// ParseResume 解析简历文件
func ParseResume(fileContent string) (*ResumeParseResult, error) {

	config.GetConfig().ValidateAI()
	aiConfig := config.GetConfig().AI
	client := resty.New().
		SetTimeout(aiConfig.Timeout).
		SetRetryCount(aiConfig.MaxRetries)

	// 构造请求内容
	prompt := fmt.Sprintf(resumeParsePrompt, fileContent)
	reqBody := map[string]interface{}{
		"model": aiConfig.ModelName,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "你是一个专业的简历解析助手",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	// 发送请求
	var aiResp AIResponse
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+aiConfig.APIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		SetResult(&aiResp).
		Post(aiConfig.BaseURL + "/chat/completions")

	if err != nil {
		return nil, fmt.Errorf("调用AI服务失败: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("AI服务返回错误: %s", resp.String())
	}

	// 解析AI返回的JSON结果
	var result ResumeParseResult
	if err := json.Unmarshal([]byte(aiResp.Message.Content), &result); err != nil {
		return nil, fmt.Errorf("解析AI响应失败: %w", err)
	}

	return &result, nil
}

// ExtractTextFromPDF 从PDF中提取文本
func ExtractTextFromPDF(pdfURL string) (string, error) {
	// TODO: 实现PDF文本提取
	// 可以使用第三方库如 pdfcpu 或调用其他服务
	return "", nil
}

// ExtractTextFromWord 从Word文档中提取文本
func ExtractTextFromWord(docURL string) (string, error) {
	// TODO: 实现Word文本提取
	// 可以使用第三方库如 unidoc 或调用其他服务
	return "", nil
}
