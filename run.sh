#!/bin/bash

# 定义安装目录和缓存文件路径
INSTALL_DIR="$(pwd)"
#INSTALL_DIR="/opt/apkshifter"
CACHE_FILE="$INSTALL_DIR/apkshifter_installed"

# 确定包管理器
if command -v yum >/dev/null 2>&1; then
    PKG_MANAGER="yum"
    INSTALL_CMD="yum install -y"
elif command -v apt-get >/dev/null 2>&1; then
    PKG_MANAGER="apt-get"
    INSTALL_CMD="apt-get install -y"
else
    echo "不支持的操作系统。需要yum或apt-get包管理器。"
    exit 1
fi

# 安装 Java 和 zlib
IS_INSTALLED=$(java -version 2>&1 >/dev/null)
if [ $? -eq 0 ]; then
    echo 'Java 已安装'
else
    $INSTALL_CMD default-jre
fi

IS_INSTALLED=$(ldconfig -p | grep zlib)
if [ $? -eq 0 ]; then
    echo 'zlib 已安装'
else
    $INSTALL_CMD zlib1g zlib1g-dev
fi

# 判断是否已经通过缓存文件确认安装
if [ -f "$CACHE_FILE" ]; then
    echo "apkShifter 已经安装."
else
    # 创建安装目录
    mkdir -p "$INSTALL_DIR"

    # 下载 apkShifter.zip
    if [ ! -f "$INSTALL_DIR/apkShifter.zip" ]; then
        echo "正在下载 apkShifter.zip..."
        curl -L -o "$INSTALL_DIR/apkShifter.zip" https://github.com/YoyoCrafts/apkShifter/releases/download/1.0.0/apkShifter.zip
    else
        echo "apkShifter.zip 已存在."
    fi

    # 解压 apkShifter.zip
    if [ ! -d "$INSTALL_DIR/apkShifter" ]; then
        echo "正在解压 apkShifter.zip..."
        unzip "$INSTALL_DIR/apkShifter.zip" -d "$INSTALL_DIR/apkShifter"
    else
        echo "apkShifter 已解压."
    fi

    # 创建缓存文件，标记为已安装
    touch "$CACHE_FILE"
    echo "安装已完成."
fi

# 状态显示和操作选项
echo "当前状态："
echo "1) 启动/重启服务"
echo "2) 停止服务"
read -p "请选择操作(1或2): " choice

case $choice in
    1)
        echo "正在启动/重启服务..."
        chmod +x "$INSTALL_DIR/apkShifter/apkShifter"
        cd "$INSTALL_DIR/apkShifter"
        nohup ./apkShifter > apkshifter.log 2>&1 &
        echo "服务已启动，正在后台运行。"
        ;;
    2)
        echo "正在停止服务..."
        pkill -f ./apkShifter
        echo "服务已停止。"
        ;;
    *)
        echo "无效选项，请选择 1 或 2."
        ;;
esac

exit 0
