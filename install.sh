#!/bin/bash

SERVICE_NAME="input2com"
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"
BINARY_PATH="/home/lty/projects/input2com_linux/input2com"

# 检查二进制文件是否存在
if [ ! -f "$BINARY_PATH" ]; then
    echo "错误：二进制文件不存在: $BINARY_PATH"
    exit 1
fi

# 检查是否有执行权限
if [ ! -x "$BINARY_PATH" ]; then
    echo "警告：二进制文件没有执行权限，正在添加..."
    chmod +x "$BINARY_PATH"
fi

echo "正在创建服务: $SERVICE_NAME..."

# 创建systemd服务文件
sudo tee "$SERVICE_FILE" > /dev/null <<EOF
[Unit]
Description=Input2COM Service
After=network.target

[Service]
Type=simple
ExecStart=$BINARY_PATH
ExecReload=/bin/kill -HUP \$MAINPID
ExecStop=/bin/kill -SIGTERM \$MAINPID
Restart=on-failure
RestartSec=5
User=root
Group=root
WorkingDirectory=/home/lty/projects/input2com_linux
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

echo "正在重载systemd配置..."
sudo systemctl daemon-reload

echo "正在启用开机自启..."
sudo systemctl enable "$SERVICE_NAME"

echo "正在启动服务..."
sudo systemctl start "$SERVICE_NAME"

echo ""
echo "✅ 服务安装完成！"
echo ""
echo "常用命令："
echo "  查看状态: sudo systemctl status $SERVICE_NAME"
echo "  查看日志: sudo journalctl -u $SERVICE_NAME -f"
echo "  重启服务: sudo systemctl restart $SERVICE_NAME"
echo "  停止服务: sudo systemctl stop $SERVICE_NAME"