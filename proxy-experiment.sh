#!/bin/bash

# proxy.golang.org 收录实验脚本
# 用于演示同一commit的不同伪版本现象

set -e

# 配置变量
GITHUB_USERNAME="binary-1024"
REPO_NAME="go-build-test"
REPO_URL="github.com/$GITHUB_USERNAME/$REPO_NAME"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo -e "${PURPLE}[STEP]${NC} $1"
}

# 检查proxy收录状态
check_proxy_status() {
    local repo=$1
    local version=$2

    log_info "检查 $repo@$version 的收录状态..."

    if [ -z "$version" ]; then
        # 检查版本列表
        response=$(curl -s "https://proxy.golang.org/$repo/@v/list")
        if [ $? -eq 0 ] && [ -n "$response" ]; then
            log_success "版本列表已收录，包含以下版本:"
            echo "$response" | sort
        else
            log_warning "版本列表未收录或为空"
        fi
    else
        # 检查特定版本
        response=$(curl -s -o /dev/null -w "%{http_code}" \
            "https://proxy.golang.org/$repo/@v/$version.info")

        if [ "$response" = "200" ]; then
            log_success "$version 已被收录"
            # 显示版本信息
            curl -s "https://proxy.golang.org/$repo/@v/$version.info" | jq . 2>/dev/null || \
            curl -s "https://proxy.golang.org/$repo/@v/$version.info"
        else
            log_warning "$version 未收录 (HTTP: $response)"
        fi
    fi
}

# 触发proxy收录
trigger_proxy_collection() {
    local repo=$1
    local version=$2

    log_info "触发 $repo@$version 的收录..."

    # 方法1: go get (最可靠)
    if command -v go >/dev/null 2>&1; then
        log_info "使用 go get 触发收录..."
        go get "$repo@$version" 2>/dev/null || true
        log_success "go get 请求已发送"
    fi

    # 方法2: curl API
    log_info "使用 curl API 触发收录..."
    if [ "$version" = "list" ]; then
        curl -s "https://proxy.golang.org/$repo/@v/list" >/dev/null
    elif [ "$version" = "latest" ]; then
        curl -s "https://proxy.golang.org/$repo/@latest" >/dev/null
    else
        curl -s "https://proxy.golang.org/$repo/@v/$version.info" >/dev/null
    fi

    log_success "API 请求已发送"
}

# 等待收录完成
wait_for_collection() {
    local repo=$1
    local version=$2
    local max_wait=${3:-300}  # 默认等待5分钟
    local interval=30
    local elapsed=0

    log_info "等待 $repo@$version 被收录 (最多等待 ${max_wait}s)..."

    while [ $elapsed -lt $max_wait ]; do
        response=$(curl -s -o /dev/null -w "%{http_code}" \
            "https://proxy.golang.org/$repo/@v/$version.info" 2>/dev/null || echo "000")

        if [ "$response" = "200" ]; then
            log_success "$version 已被收录 (等待时间: ${elapsed}s)"
            return 0
        fi

        log_info "等待中... (${elapsed}s/${max_wait}s)"
        sleep $interval
        elapsed=$((elapsed + interval))
    done

    log_warning "$version 在 ${max_wait}s 内未被收录"
    return 1
}

# 创建实验仓库
create_experiment_repo() {
    log_step "创建实验仓库..."

    if [ -d "$REPO_NAME" ]; then
        log_warning "目录 $REPO_NAME 已存在，将删除重建"
        rm -rf "$REPO_NAME"
    fi

    mkdir "$REPO_NAME"
    cd "$REPO_NAME"

    # 初始化Go模块
    go mod init "$REPO_URL"

    # 创建main.go
    cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("Go Proxy Experiment")
    fmt.Printf("Version: %s\n", Version())
}

// Version 返回当前版本
func Version() string {
    return "1.0.0"
}
EOF

    # 创建experiment.go
    cat > experiment.go << 'EOF'
package main

import (
    "encoding/json"
    "time"
)

// ExperimentInfo 实验信息
type ExperimentInfo struct {
    Name        string    `json:"name"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    Version     string    `json:"version"`
    CommitHash  string    `json:"commit_hash,omitempty"`
}

// GetExperimentInfo 获取实验信息
func GetExperimentInfo() *ExperimentInfo {
    return &ExperimentInfo{
        Name:        "Go Proxy Experiment",
        Description: "测试proxy.golang.org收录不同版本的实验",
        CreatedAt:   time.Now(),
        Version:     Version(),
    }
}

// ToJSON 转换为JSON字符串
func (e *ExperimentInfo) ToJSON() string {
    data, _ := json.MarshalIndent(e, "", "  ")
    return string(data)
}
EOF

    # 创建测试文件
    cat > experiment_test.go << 'EOF'
package main

import "testing"

func TestVersion(t *testing.T) {
    version := Version()
    if version == "" {
        t.Error("Version should not be empty")
    }
}

func TestGetExperimentInfo(t *testing.T) {
    info := GetExperimentInfo()
    if info.Name == "" {
        t.Error("Name should not be empty")
    }
}
EOF

    # 创建README
    cat > README.md << EOF
# Go Proxy Experiment

这是一个用于测试 proxy.golang.org 收录机制的实验项目。

## 实验目标

测试 proxy.golang.org 如何收录：
1. 正常的Go包
2. 语义化版本标签
3. 伪版本
4. 版本变化过程

## 仓库信息

- 仓库地址: https://github.com/$GITHUB_USERNAME/$REPO_NAME
- Go模块: $REPO_URL

## 使用方法

\`\`\`bash
go get $REPO_URL
\`\`\`

## 实验记录

### 阶段1: 初始提交
- 时间: $(date)
- 描述: 项目初始化

### 阶段2: 标签版本
- 待记录...

### 阶段3: 伪版本实验
- 待记录...
EOF

    # 初始化Git仓库
    git init
    git add .
    git commit -m "Initial commit: Go proxy experiment setup

This is the initial commit for testing proxy.golang.org collection behavior.
We will use this commit to demonstrate pseudo-version generation."

    # 记录初始commit
    INITIAL_COMMIT=$(git rev-parse HEAD)
    echo "INITIAL_COMMIT=$INITIAL_COMMIT" > .experiment_vars

    log_success "实验仓库创建完成"
    log_info "初始commit: $INITIAL_COMMIT"

    cd ..
}

# 执行实验阶段1: 收录伪版本
experiment_stage1() {
    log_step "阶段1: 收录初始伪版本"

    # 检查是否已经在正确的目录中
    if [ -f "go.mod" ] && [ -f ".experiment_vars" ]; then
        # 已经在正确的目录中
        log_info "检测到已在实验目录中"
    elif [ -d "$REPO_NAME" ]; then
        # 在父目录中，需要进入子目录
        log_info "进入实验目录 $REPO_NAME"
        cd "$REPO_NAME"
    else
        log_error "实验目录 $REPO_NAME 不存在，请先运行: $0 init"
        return 1
    fi

    if [ ! -f ".experiment_vars" ]; then
        log_error "实验配置文件 .experiment_vars 不存在，请先运行: $0 init"
        return 1
    fi

    source .experiment_vars

    log_info "初始commit: $INITIAL_COMMIT"

    # 触发收录
    trigger_proxy_collection "$REPO_URL" "main"
    trigger_proxy_collection "$REPO_URL" "$INITIAL_COMMIT"

    # 等待收录
    log_info "等待收录完成..."
    sleep 60

    # 检查收录状态
    check_proxy_status "$REPO_URL"

    log_success "阶段1完成"
    cd ..
}

# 执行实验阶段2: 创建标签版本
experiment_stage2() {
    log_step "阶段2: 创建v1.0.0标签"

    # if [ ! -d "$REPO_NAME" ]; then
    #     log_error "实验目录 $REPO_NAME 不存在，请先运行: $0 init"
    #     return 1
    # fi

    # cd "$REPO_NAME"

    # if [ ! -f ".experiment_vars" ]; then
    #     log_error "实验配置文件 .experiment_vars 不存在，请先运行: $0 init"
    #     return 1
    # fi

    source .experiment_vars

    # 创建v1.0.0标签
    git tag v1.0.0
    log_info "已创建标签 v1.0.0 指向 $INITIAL_COMMIT"

    # 如果有远程仓库，推送标签
    if git remote get-url origin >/dev/null 2>&1; then
        git push origin v1.0.0
        log_info "标签已推送到远程仓库"
    else
        log_warning "未配置远程仓库，请手动推送标签: git push origin v1.0.0"
    fi

    # 触发收录
    trigger_proxy_collection "$REPO_URL" "v1.0.0"

    # 等待收录
    log_info "等待v1.0.0收录..."
    sleep 60

    # 检查收录状态
    check_proxy_status "$REPO_URL" "v1.0.0"
    check_proxy_status "$REPO_URL"

    log_success "阶段2完成"
    cd ..
}

# 执行实验阶段3: 新的开发提交
experiment_stage3() {
    log_step "阶段3: 添加新的开发提交"

    # if [ ! -d "$REPO_NAME" ]; then
    #     log_error "实验目录 $REPO_NAME 不存在，请先运行: $0 init"
    #     return 1
    # fi

    # cd "$REPO_NAME"

    # if [ ! -f ".experiment_vars" ]; then
    #     log_error "实验配置文件 .experiment_vars 不存在，请先运行: $0 init"
    #     return 1
    # fi

    source .experiment_vars

    # 添加新功能
    cat >> experiment.go << 'EOF'

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
EOF

    # 更新版本号
    sed -i.bak 's/return "1.0.0"/return "1.0.1-dev"/' main.go
    rm -f main.go.bak

    # 提交更改
    git add .
    git commit -m "Add new feature after v1.0.0

This commit will be used to test pseudo-version generation
with v1.0.0 as the base version.

Features added:
- NewFeature() function
- GetCommitInfo() function
- Updated version to 1.0.1-dev"

    # 记录新commit
    NEW_COMMIT=$(git rev-parse HEAD)
    echo "NEW_COMMIT=$NEW_COMMIT" >> .experiment_vars

    log_info "新commit: $NEW_COMMIT"

    # 如果有远程仓库，推送更改
    if git remote get-url origin >/dev/null 2>&1; then
        git push origin main
        log_info "更改已推送到远程仓库"
    else
        log_warning "未配置远程仓库，请手动推送: git push origin main"
    fi

    # 触发收录新的伪版本
    trigger_proxy_collection "$REPO_URL" "main"
    trigger_proxy_collection "$REPO_URL" "$NEW_COMMIT"

    # 等待收录
    log_info "等待新伪版本收录..."
    sleep 60

    # 检查收录状态
    check_proxy_status "$REPO_URL"

    log_success "阶段3完成"
    log_info "预期的伪版本格式: v1.0.0-0.YYYYMMDDHHMMSS-${NEW_COMMIT:0:12}"

    cd ..
}

# 执行实验阶段4: 关键实验 - 同一commit的不同版本
experiment_stage4() {
    log_step "阶段4: 为同一commit创建v1.0.1标签"

    # if [ ! -d "$REPO_NAME" ]; then
    #     log_error "实验目录 $REPO_NAME 不存在，请先运行: $0 init"
    #     return 1
    # fi

    # cd "$REPO_NAME"

    # if [ ! -f ".experiment_vars" ]; then
    #     log_error "实验配置文件 .experiment_vars 不存在，请先运行: $0 init"
    #     return 1
    # fi

    source .experiment_vars

    log_info "将为commit $NEW_COMMIT 创建v1.0.1标签"

    # 创建v1.0.1标签指向新commit
    git tag v1.0.1 "$NEW_COMMIT"
    log_info "已创建标签 v1.0.1 指向 $NEW_COMMIT"

    # 推送标签
    if git remote get-url origin >/dev/null 2>&1; then
        git push origin v1.0.1
        log_info "标签已推送到远程仓库"
    else
        log_warning "未配置远程仓库，请手动推送标签: git push origin v1.0.1"
    fi

    # 触发收录
    trigger_proxy_collection "$REPO_URL" "v1.0.1"

    # 等待收录
    log_info "等待v1.0.1收录..."
    sleep 60

    # 再次触发伪版本收录 (测试是否会生成不同的伪版本)
    log_info "重新触发相同commit的收录..."
    trigger_proxy_collection "$REPO_URL" "main"
    trigger_proxy_collection "$REPO_URL" "$NEW_COMMIT"

    sleep 30

    # 检查所有版本
    check_proxy_status "$REPO_URL"
    check_proxy_status "$REPO_URL" "v1.0.1"

    log_success "阶段4完成"
    log_info "关键观察点: commit $NEW_COMMIT 现在可能有不同的伪版本表示"

    cd ..
}

# 分析实验结果
analyze_results() {
    log_step "分析实验结果"

    # if [ ! -d "$REPO_NAME" ]; then
    #     log_error "实验目录 $REPO_NAME 不存在，请先运行: $0 init"
    #     return 1
    # fi

    # cd "$REPO_NAME"

    # if [ ! -f ".experiment_vars" ]; then
    #     log_error "实验配置文件 .experiment_vars 不存在，请先运行: $0 init"
    #     return 1
    # fi

    source .experiment_vars

    log_info "实验总结:"
    echo "初始commit: $INITIAL_COMMIT"
    echo "新commit: $NEW_COMMIT"
    echo ""

    log_info "获取所有收录的版本..."
    versions=$(curl -s "https://proxy.golang.org/$REPO_URL/@v/list" 2>/dev/null || echo "")

    if [ -n "$versions" ]; then
        log_success "proxy.golang.org 已收录以下版本:"
        echo "$versions" | sort
        echo ""

        log_info "分析包含新commit的版本..."
        echo "$versions" | while read -r version; do
            if [[ "$version" == *"${NEW_COMMIT:0:12}"* ]]; then
                echo "🎯 版本 $version 包含commit ${NEW_COMMIT:0:12}"

                # 获取详细信息
                info=$(curl -s "https://proxy.golang.org/$REPO_URL/@v/$version.info" 2>/dev/null)
                if [ -n "$info" ]; then
                    echo "   详细信息: $info"
                fi
            fi
        done
    else
        log_warning "未获取到版本列表，可能需要更多时间等待收录"
    fi

    echo ""
    log_info "手动检查方法:"
    echo "curl \"https://proxy.golang.org/$REPO_URL/@v/list\""
    echo "curl \"https://proxy.golang.org/$REPO_URL/@latest\""

    cd ..
}

# 清理实验环境
cleanup_experiment() {
    log_step "清理实验环境"

    if [ -d "$REPO_NAME" ]; then
        log_warning "删除实验目录 $REPO_NAME"
        rm -rf "$REPO_NAME"
    fi

    # 清理Go模块缓存
    if command -v go >/dev/null 2>&1; then
        log_info "清理Go模块缓存..."
        go clean -modcache 2>/dev/null || true
    fi

    log_success "清理完成"
}

# 显示帮助信息
show_help() {
    echo "Go Proxy 收录实验脚本"
    echo ""
    echo "用法: $0 [command]"
    echo ""
    echo "命令:"
    echo "  init              - 创建实验仓库"
    echo "  stage1            - 执行阶段1: 收录初始伪版本"
    echo "  stage2            - 执行阶段2: 创建v1.0.0标签"
    echo "  stage3            - 执行阶段3: 添加新开发提交"
    echo "  stage4            - 执行阶段4: 同一commit的不同版本实验"
    echo "  analyze           - 分析实验结果"
    echo "  full              - 执行完整实验流程"
    echo "  check [version]   - 检查收录状态"
    echo "  trigger <version> - 触发收录"
    echo "  cleanup           - 清理实验环境"
    echo "  help              - 显示帮助信息"
    echo ""
    echo "环境变量:"
    echo "  GITHUB_USERNAME   - GitHub用户名 (默认: your-username)"
    echo ""
    echo "示例:"
    echo "  export GITHUB_USERNAME=myusername"
    echo "  $0 full"
}

# 主函数
main() {
    case "${1:-help}" in
        "triger")
            trigger_proxy_collection "$REPO_URL" "main"
            ;;
        "init")
            INITIAL_COMMIT=$(git rev-parse HEAD)
            echo "INITIAL_COMMIT=$INITIAL_COMMIT" > .experiment_vars

            log_success "实验仓库创建完成"
            log_info "初始commit: $INITIAL_COMMIT"
            # create_experiment_repo
            ;;
        "stage1")
            experiment_stage1
            ;;
        "stage2")
            experiment_stage2
            ;;
        "stage3")
            experiment_stage3
            ;;
        "stage4")
            experiment_stage4
            ;;
        "analyze")
            analyze_results
            ;;
        "full")
            log_info "开始完整的实验流程..."
            create_experiment_repo
            sleep 5
            experiment_stage1
            sleep 5
            experiment_stage2
            sleep 5
            experiment_stage3
            sleep 5
            experiment_stage4
            sleep 5
            analyze_results
            ;;
        "check")
            check_proxy_status "$REPO_URL" "$2"
            ;;
        "trigger")
            if [ -z "$2" ]; then
                log_error "请指定要触发的版本"
                exit 1
            fi
            trigger_proxy_collection "$REPO_URL" "$2"
            ;;
        "cleanup")
            cleanup_experiment
            ;;
        "help"|*)
            show_help
            ;;
    esac
}

# 检查依赖
check_dependencies() {
    local missing_deps=()

    if ! command -v curl >/dev/null 2>&1; then
        missing_deps+=("curl")
    fi

    if ! command -v jq >/dev/null 2>&1; then
        log_warning "jq 未安装，JSON输出可能不够美观"
    fi

    if ! command -v git >/dev/null 2>&1; then
        missing_deps+=("git")
    fi

    if [ ${#missing_deps[@]} -gt 0 ]; then
        log_error "缺少必要的依赖: ${missing_deps[*]}"
        log_info "请安装缺少的工具后重试"
        exit 1
    fi
}

# 脚本入口
if [ "$GITHUB_USERNAME" = "your-username" ]; then
    log_warning "请设置环境变量 GITHUB_USERNAME 为你的GitHub用户名"
    log_info "示例: export GITHUB_USERNAME=myusername"
fi

check_dependencies
main "$@"
