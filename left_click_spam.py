"""
左键连发 UDP 发送脚本
通过 UDP 向 input2com 发送鼠标左键按下/释放事件，实现连发效果

用法:
    python left_click_spam.py                    # 默认参数
    python left_click_spam.py -H 192.168.3.3    # 指定目标主机
    python left_click_spam.py -i 50              # 每50ms一次点击
    python left_click_spam.py -d 5              # 持续5秒
"""

import socket
import struct
import time
import random
import argparse

# Linux evdev 常量
EV_KEY = 1
BTN_LEFT = 0x110  # 272


def pack_event(event_type: int, code: int, value: int) -> bytes:
    """打包一个 evdev 事件为 8 字节: type(u16 LE) + code(u16 LE) + value(i32 LE)"""
    return struct.pack('<HHi', event_type, code, value)


def make_packet(events: list[tuple[int, int, int]], dev_name: str = "python") -> bytes:
    """
    构造 UDP 数据包
    格式: event_count(u8) + events(N*8 bytes) + dev_name(str)
    """
    count = len(events)
    header = struct.pack('B', count)
    body = b''.join(pack_event(t, c, v) for t, c, v in events)
    return header + body + dev_name.encode('utf-8')


def press_packet() -> bytes:
    """左键按下"""
    return make_packet([(EV_KEY, BTN_LEFT, 1)])


def release_packet() -> bytes:
    """左键释放"""
    return make_packet([(EV_KEY, BTN_LEFT, 0)])


def main():
    parser = argparse.ArgumentParser(description='左键连发 UDP 发送脚本')
    parser.add_argument('-H', '--host', default='192.168.3.3', help='目标主机 (默认 192.168.3.3)')
    parser.add_argument('-p', '--port', type=int, default=9265, help='目标端口 (默认 9265)')
    parser.add_argument('-i', '--interval', type=float, default=500, help='点击间隔 ms (默认 500)')
    parser.add_argument('-Hd', '--hold', type=float, default=100, help='按下持续时间 ms (默认 100)')
    parser.add_argument('-d', '--duration', type=float, default=0, help='持续时间秒, 0=无限 (默认 0)')
    args = parser.parse_args()

    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    addr = (args.host, args.port)
    down_pkt = press_packet()
    up_pkt = release_packet()
    hold_s = args.hold / 1000.0
    interval_s = args.interval / 1000.0

    print(f"发送到 {addr[0]}:{addr[1]}, 间隔 {args.interval}ms, 按住 {args.hold}ms")
    print("按 Ctrl+C 停止")

    start = time.monotonic()
    count = 0
    try:
        while True:
            if args.duration > 0 and (time.monotonic() - start) >= args.duration:
                break
            sock.sendto(down_pkt, addr)
            time.sleep(hold_s)
            sock.sendto(up_pkt, addr)
            count += 1
            time.sleep(max(0, interval_s - hold_s) + random.uniform(0, 0.3))
    except KeyboardInterrupt:
        pass
    finally:
        sock.close()
        elapsed = time.monotonic() - start
        print(f"\n已停止, 共发送 {count} 次, 耗时 {elapsed:.1f}s")


if __name__ == '__main__':
    main()
