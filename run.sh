#!/bin/bash
cd "$(dirname "$0")"
while true; do
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] 启动 input2com..."
    sudo ./input2com
    code=$?
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] 进程退出 (code=$code)，1秒后重启..."
    sleep 1
done
