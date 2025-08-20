
// NewFeature 新功能演示
func NewFeature() string {
    return "This is a new feature added after v1.0.0"
}

// GetCommitInfo 获取commit信息
func GetCommitInfo() map[string]string {
    return map[string]string{
        "stage": "after-v1.0.0",
        "description": "This commit will demonstrate pseudo-version with v1.0.0 base",
    }
}
