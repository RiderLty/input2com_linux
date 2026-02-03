#!/bin/bash

SERVICE_NAME="input2com"
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"

echo "正在卸载服务: $SERVICE_NAME..."

# 检查服务是否存在
if [ ! -f "$SERVICE_FILE" ]; then
    echo "警告：服务文件不存在，可能已卸载"
else
    # 停止服务
    echo "正在停止服务..."
    sudo systemctl stop "$SERVICE_NAME" 2>/dev/null || true
    
    # 禁用开机自启
    echo "正在禁用开机自启..."
    sudo systemctl disable "$SERVICE_NAME" 2>/dev/null || true
    
    # 删除服务文件
    echo "正在删除服务文件..."
    sudo rm -f "$SERVICE_FILE"
    
    # 重载systemd
    echo "正在重载systemd配置..."
    sudo systemctl daemon-reload
    
    echo ""
    echo "✅ 服务卸载完成！"
fi

echo ""
echo "注意：二进制文件保留在 /home/lty/projects/input2com_linux/input2com"
echo "如需删除请手动执行: rm /home/lty/projects/input2com_linux/input2com"