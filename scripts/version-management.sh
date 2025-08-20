#!/bin/bash

echo "=== Go模块版本管理脚本 ==="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 帮助函数
print_header() {
    echo -e "${BLUE}=== $1 ===${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# 函数：显示包的所有版本
show_versions() {
    local package=$1
    if [ -z "$package" ]; then
        print_error "请提供包名"
        echo "用法: $0 versions <package-name>"
        echo "示例: $0 versions github.com/gin-gonic/gin"
        return 1
    fi

    print_header "查看 $package 的所有版本"

    echo "正式发布版本:"
    go list -m -versions "$package" 2>/dev/null || print_error "无法获取版本信息"

    echo ""
    echo "获取最新commit版本:"
    echo "go get $package@main"
    echo "go get $package@master"
    echo "go get $package@develop"
}

# 函数：获取特定commit
get_commit() {
    local package=$1
    local commit=$2

    if [ -z "$package" ] || [ -z "$commit" ]; then
        print_error "请提供包名和commit哈希"
        echo "用法: $0 commit <package-name> <commit-hash>"
        echo "示例: $0 commit github.com/gin-gonic/gin abcdef123456"
        return 1
    fi

    print_header "获取 $package 的特定commit: $commit"

    echo "执行命令: go get $package@$commit"
    if go get "$package@$commit"; then
        print_success "成功获取commit版本"
        echo ""
        echo "查看go.mod中的版本:"
        grep "$package" go.mod || print_warning "在go.mod中未找到该包"
    else
        print_error "获取commit版本失败"
    fi
}

# 函数：获取最新开发版本
get_latest_dev() {
    local package=$1
    local branch=${2:-main}

    if [ -z "$package" ]; then
        print_error "请提供包名"
        echo "用法: $0 latest <package-name> [branch]"
        echo "示例: $0 latest github.com/gin-gonic/gin main"
        return 1
    fi

    print_header "获取 $package 的最新开发版本 (分支: $branch)"

    echo "执行命令: go get $package@$branch"
    if go get "$package@$branch"; then
        print_success "成功获取最新开发版本"
        echo ""
        echo "生成的伪版本:"
        grep "$package" go.mod || print_warning "在go.mod中未找到该包"
    else
        print_error "获取最新开发版本失败"
        echo "尝试其他分支: master, develop, main"
    fi
}

# 函数：解析伪版本
parse_pseudo() {
    local version=$1

    if [ -z "$version" ]; then
        print_error "请提供伪版本字符串"
        echo "用法: $0 parse <pseudo-version>"
        echo "示例: $0 parse v1.9.2-0.20240110094500-fedcba654321"
        return 1
    fi

    print_header "解析伪版本: $version"

    # 检查是否为伪版本格式
    if [[ $version =~ ^v[0-9]+\.[0-9]+\.[0-9]+-[0-9]{14}-[a-f0-9]{12}$ ]]; then
        # 提取各部分
        base_version=$(echo "$version" | cut -d'-' -f1)
        timestamp=$(echo "$version" | cut -d'-' -f2)
        commit_hash=$(echo "$version" | cut -d'-' -f3)

        echo "基础版本: $base_version"
        echo "时间戳: $timestamp"
        echo "Commit哈希: $commit_hash"

        # 转换时间戳为可读格式
        if command -v date >/dev/null 2>&1; then
            readable_time=$(date -d "${timestamp:0:8} ${timestamp:8:2}:${timestamp:10:2}:${timestamp:12:2}" 2>/dev/null || echo "时间解析失败")
            echo "时间: $readable_time"
        fi

        print_success "伪版本解析完成"
    else
        print_warning "这不是标准的伪版本格式"
        echo "标准格式: v{major}.{minor}.{patch}-{timestamp}-{commit}"
        echo "示例: v1.9.2-0.20240110094500-fedcba654321"
    fi
}

# 函数：显示依赖图
show_deps() {
    print_header "当前项目的依赖图"

    echo "直接依赖:"
    go list -m all | head -20

    echo ""
    echo "依赖关系图:"
    go mod graph | head -20

    echo ""
    echo "查看完整依赖图: go mod graph"
    echo "查看特定包的依赖: go mod why <package>"
}

# 函数：清理和更新依赖
clean_deps() {
    print_header "清理和更新依赖"

    echo "1. 下载缺失的依赖..."
    go mod download

    echo "2. 清理未使用的依赖..."
    go mod tidy

    echo "3. 验证依赖..."
    go mod verify

    print_success "依赖清理完成"
}

# 函数：显示帮助信息
show_help() {
    echo "Go模块版本管理工具"
    echo ""
    echo "用法: $0 <command> [args...]"
    echo ""
    echo "命令:"
    echo "  versions <package>           - 显示包的所有可用版本"
    echo "  commit <package> <hash>      - 获取特定commit的版本"
    echo "  latest <package> [branch]    - 获取最新开发版本"
    echo "  parse <pseudo-version>       - 解析伪版本格式"
    echo "  deps                         - 显示当前项目依赖"
    echo "  clean                        - 清理和更新依赖"
    echo "  help                         - 显示帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 versions github.com/gin-gonic/gin"
    echo "  $0 commit github.com/gin-gonic/gin abcdef123456"
    echo "  $0 latest github.com/sirupsen/logrus main"
    echo "  $0 parse v1.9.2-0.20240110094500-fedcba654321"
    echo ""
    echo "伪版本说明:"
    echo "  格式: v{base}-{timestamp}-{commit}"
    echo "  base: 最近的语义化版本标签"
    echo "  timestamp: commit时间 (YYYYMMDDHHMMSS)"
    echo "  commit: Git commit哈希前12位"
}

# 主函数
main() {
    case "$1" in
        "versions")
            show_versions "$2"
            ;;
        "commit")
            get_commit "$2" "$3"
            ;;
        "latest")
            get_latest_dev "$2" "$3"
            ;;
        "parse")
            parse_pseudo "$2"
            ;;
        "deps")
            show_deps
            ;;
        "clean")
            clean_deps
            ;;
        "help"|"")
            show_help
            ;;
        *)
            print_error "未知命令: $1"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# 检查是否在Go模块目录中
if [ ! -f "go.mod" ]; then
    print_warning "当前目录不是Go模块根目录"
    echo "请在包含go.mod文件的目录中运行此脚本"
fi

# 运行主函数
main "$@"
