import serial
import time

ser = serial.Serial('COM4', 115200)  # 替换为你的串口号

# 发送二进制数据
data = bytes([0x01, 0x01, 0xFF, 0xFF, 0xFF, 0xFF])
ser.write(data)

ser.close()