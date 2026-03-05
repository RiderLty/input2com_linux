// package remote
package main

// import (
// 	"encoding/binary"
// 	"input2com/internal/config"
// 	"input2com/internal/input"
// 	"input2com/internal/logger"
// 	"input2com/internal/macros"
// 	"input2com/internal/serial"
// 	"math"
// 	"math/rand"
// 	"net"
// 	"sync"
// 	"time"
// )

// // 定义数据包结构，用于线程间传递
// type UDPData struct {
// 	FirstInt, MutX, SecondInt, X2 int32
// 	//SecondInt          int64
// 	Timestamp int32
// }

// type Remote struct {
// 	Ctrl     *macros.MacroMouseKeyboard
// 	dataChan chan UDPData  // 用于传递最新数据包的通道
// 	stopChan chan struct{} // 用于停止线程的信号通道
// 	wg       sync.WaitGroup
// }

// func NewRemoteControl(ctrl *macros.MacroMouseKeyboard) *Remote {
// 	return &Remote{
// 		Ctrl:     ctrl,
// 		dataChan: make(chan UDPData), // 缓冲区大小为1，确保只保留最新数据
// 		stopChan: make(chan struct{}),
// 	}
// }

// func (r *Remote) Start() {
// 	r.wg.Add(2)
// 	go r.listenUDP()    // UDP监听线程
// 	go r.handleMotion() // 移动控制线程
// }

// // Stop 停止所有线程
// func (r *Remote) Stop() {
// 	close(r.stopChan)
// 	r.wg.Wait()
// }

// // listenUDP 负责接收UDP数据并将最新数据发送到通道
// func (r *Remote) listenUDP() {
// 	defer r.wg.Done()

// 	// 初始化UDP监听
// 	addr, err := net.ResolveUDPAddr("udp", ":8080")
// 	if err != nil {
// 		logger.Logger.Error("Error resolving address:", err)
// 		return
// 	}

// 	conn, err := net.ListenUDP("udp", addr)
// 	if err != nil {
// 		logger.Logger.Error("Error listening:", err)
// 		return
// 	}
// 	defer conn.Close()
// 	buffer := make([]byte, 28)
// 	for {
// 		select {
// 		case <-r.stopChan:
// 			return // 收到停止信号，退出线程
// 		default:
// 			// 读取UDP数据
// 			n, _, err := conn.ReadFromUDP(buffer)
// 			if err != nil {
// 				logger.Logger.Error("Error reading:", err)
// 				continue
// 			}

// 			if n != 28 {
// 				logger.Logger.Warnf("Received packet with wrong size: %d bytes (expected 28)\n", n)
// 				continue
// 			}
// 			firstInt := binary.LittleEndian.Uint32(buffer[0:4])
// 			secondInt := binary.LittleEndian.Uint32(buffer[4:8])
// 			mutx := binary.LittleEndian.Uint32(buffer[8:12])
// 			x2 := binary.LittleEndian.Uint32(buffer[16:20])
// 			y2 := binary.LittleEndian.Uint32(buffer[20:24])
// 			timeStamp := binary.LittleEndian.Uint32(buffer[24:28])
// 			r.Ctrl.SetAimData(int32(firstInt), int32(secondInt), int32(x2), int32(y2), int32(timeStamp))
// 			select {
// 			case r.dataChan <- UDPData{FirstInt: int32(firstInt), SecondInt: int32(secondInt), MutX: int32(mutx), X2: int32(x2), Timestamp: int32(timeStamp)}:
// 			default:
// 				logger.Logger.Debug("Discarding old packet, new one received")
// 			}
// 		}
// 	}
// }

// // handleMotion 负责处理鼠标移动，只处理最新数据
// func (r *Remote) handleMotion() {
// 	defer r.wg.Done()
// 	lastTriggerTime := int32(0)
// 	lastX := int32(0)
// 	for {
// 		select {
// 		case <-r.stopChan:
// 			return // 收到停止信号，退出线程
// 		case data := <-r.dataChan:
// 			if data.Timestamp-lastTriggerTime < config.GetAimDelay() {
// 				continue
// 			}
// 			firstInt, mutx, x2 := int32(float64(data.FirstInt)*1.25), int32(float64(data.FirstInt)*1.08), data.X2
// 			disRate := float64(abs(firstInt)) / float64(int32(x2)/2)
// 			if abs(abs(firstInt)-abs(lastX)) < 3 {
// 				continue
// 			}
// 			lastX = firstInt
// 			startTime := time.Now().UnixMilli()
// 			if r.Ctrl.Ctrl.IsMouseBtnPressed(input.MouseBtnForward) {
// 				if disRate < 7 && disRate > 0.8 {
// 					var remainder float64
// 					r.Ctrl.Ctrl.LockMouse(serial.MOUSE_X, 1)

// 					if r.Ctrl.Ctrl.IsMouseBtnPressed(input.MouseBtnRight) {
// 						// 使用调整后的mutx
// 						points := GeneratePullTrajectory(float64(mutx), int(abs(mutx)/9), 0.8, 0.9, 0)
// 						for _, v := range points {
// 							total := v + remainder
// 							moveValue := int32(total)
// 							if moveValue == 0 {
// 								continue
// 							}
// 							remainder = total - float64(moveValue)
// 							r.Ctrl.Ctrl.MouseMove(moveValue, rand.Int31n(2), 0)
// 							time.Sleep(time.Millisecond * time.Duration(rand.Intn(2)+config.GetAimSpeed()))
// 						}
// 					} else {
// 						points := GeneratePullTrajectory(float64(firstInt), int(abs(firstInt)/9), 0.8, 0.9, 0)
// 						for _, v := range points {
// 							total := v + remainder
// 							moveValue := int32(total)
// 							remainder = total - float64(moveValue)
// 							if moveValue == 0 {
// 								continue
// 							}
// 							r.Ctrl.Ctrl.MouseMove(moveValue, rand.Int31n(2), 0)
// 							time.Sleep(time.Millisecond * time.Duration(rand.Intn(2)+config.GetAimSpeed()))
// 						}
// 					}
// 					r.Ctrl.Ctrl.LockMouse(serial.MOUSE_X, 0)
// 					lastTriggerTime = data.Timestamp + int32(time.Now().UnixMilli()-startTime)
// 				}
// 			}
// 		}
// 	}
// }

// // 更精确的延迟补偿（如果需要）
// func calculateAdjustedDistance(currentDistance, lastDistance int32, timeDelta int32, delayMs int32) int32 {
// 	if lastDistance <= 0 || currentDistance >= lastDistance {
// 		return currentDistance // 距离没有缩短，使用原距离
// 	}

// 	// 计算距离变化率（单位时间内的距离变化）
// 	distanceChange := lastDistance - currentDistance

// 	// 如果时间差很小或为0，避免除零错误，使用默认估算
// 	if timeDelta <= 0 {
// 		// 使用简单的比例估算
// 		estimatedMovement := float64(distanceChange) * float64(delayMs) / 40.0 // 假设40ms为基准
// 		adjusted := float64(currentDistance) - estimatedMovement
// 		if adjusted < 0 {
// 			return 0
// 		}
// 		return int32(adjusted)
// 	}

// 	// 计算速度（距离变化/时间）
// 	speed := float64(distanceChange) / float64(timeDelta) // 单位：距离/毫秒

// 	// 预估在delayMs后的目标位置
// 	estimatedMovement := speed * float64(delayMs)

// 	// 调整后的距离 = 当前距离 - 预估的目标移动距离
// 	adjusted := float64(currentDistance) - estimatedMovement
// 	if adjusted < 0 {
// 		adjusted = 0
// 	}

// 	return int32(adjusted)
// }
// func sign(x int32) int32 {
// 	if x < 0 {
// 		return -1
// 	}
// 	return 1
// }

// // splitAndMove 步长小于15，增加随机化，Y轴特殊处理
// func (r *Remote) splitAndMove(totalX int32) {
// 	// 忽略传入的Y轴参数，仅保留X轴原始位移
// 	//remainingX := totalX
// 	// 初始化随机数种子（确保全局只初始化一次，可移至程序入口）
// 	//print("move", totalX)
// 	s := sign(totalX)
// 	remainingX := abs(totalX)
// 	// 循环处理X轴剩余位移，直到归零
// 	for remainingX != 0 {
// 		select {
// 		case <-r.stopChan:
// 			return // 收到停止信号，退出
// 		default:
// 			// 1. 计算当前X步长（最大15，根据剩余位移动态调整）
// 			// 正向移动：步长为1~15之间的随机值，不超过剩余位移
// 			maxStep := min(remainingX, 15)
// 			stepX := int32(rand.Intn(int(maxStep))) + 1 // 1~maxStep

// 			// 2. Y轴处理：10%概率移动-1，否则0
// 			stepY := int32(0)
// 			if rand.Float64() < 0.1 {
// 				stepY = 1
// 			}

// 			// 3. 执行本次移动
// 			r.Ctrl.Ctrl.MouseMove(s*stepX, stepY, 0)

// 			// 4. 扣减剩余位移（核心：移动多少扣多少）
// 			remainingX -= stepX

// 			// 小延迟避免移动过快
// 			//time.Sleep(5 * time.Millisecond)
// 		}
// 	}

// 	// X轴位移归零时，仍检查一次Y轴随机移动
// 	//if rand.Float64() < 0.1 {
// 	//	r.Ctrl.Ctrl.MouseMove(0, -1, 0)
// 	//	//time.Sleep(5 * time.Millisecond)
// 	//}
// }

// // 辅助函数：计算绝对值
// func abs(x int32) int32 {
// 	if x < 0 {
// 		return -x
// 	}
// 	return x
// }

// // GeneratePullTrajectory 生成拉枪（单方向）每帧位移数组
// // totalDisp: 要移动的总位移（像素），正数代表向正方向拉（若为负则自动处理）
// // totalFrames: 输出数组长度（帧数）
// // noiseStd: 每帧高斯噪声标准差（像素），可设 0 表示无噪声
// // easeAmt: 缓动权重 0..1，0=线性，1=强 ease-in-out
// // seed: 若为0则用当前时间随机化
// func GeneratePullTrajectory(totalDisp float64, totalFrames int, noiseStd, easeAmt float64, seed int64) []float64 {
// 	if totalFrames <= 1 {
// 		return []float64{totalDisp}
// 	}
// 	// sign 保证方向不变（我们在内部用正数处理，再乘回 sign）
// 	sign := 1.0
// 	if totalDisp < 0 {
// 		sign = -1.0
// 		totalDisp = -totalDisp
// 	}

// 	// rand
// 	if seed == 0 {
// 		seed = time.Now().UnixNano()
// 	}
// 	r := rand.New(rand.NewSource(seed))

// 	// easing helper: mix linear and smoothstep (smoothstep = t^2*(3-2t))
// 	ease := func(t float64) float64 {
// 		if easeAmt <= 0 {
// 			return t
// 		}
// 		if easeAmt >= 1 {
// 			return t * t * (3 - 2*t)
// 		}
// 		linear := t
// 		smooth := t * t * (3 - 2*t)
// 		return linear*(1.0-easeAmt) + smooth*easeAmt
// 	}

// 	// 1) 生成累积位置（从0到totalDisp），等间隔 t in [0,1]
// 	pos := make([]float64, totalFrames)
// 	for i := 0; i < totalFrames; i++ {
// 		t := float64(i) / float64(totalFrames-1)
// 		pos[i] = ease(t) * totalDisp
// 	}

// 	// 2) 由位置求每帧位移（delta）
// 	deltas := make([]float64, totalFrames)
// 	prev := 0.0
// 	for i := 0; i < totalFrames; i++ {
// 		deltas[i] = pos[i] - prev
// 		prev = pos[i]
// 	}

// 	// 3) 在 delta 上加噪声，但**确保不反向**（即 delta 保持 >= minPositive）
// 	//    使用 minPositive 为 totalDisp/(totalFrames*1000) 保证极小正值边界
// 	minPositive := totalDisp / float64(totalFrames) * 1e-3 // 极小阈值，避免变为0
// 	if minPositive <= 1e-6 {
// 		minPositive = 1e-6
// 	}
// 	for i := 0; i < totalFrames; i++ {
// 		if noiseStd > 0 {
// 			// Box-Muller 高斯
// 			u1 := r.Float64()
// 			u2 := r.Float64()
// 			if u1 < 1e-12 {
// 				u1 = 1e-12
// 			}
// 			z0 := math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2*math.Pi*u2)
// 			noise := z0 * noiseStd
// 			deltas[i] += noise
// 		}
// 		// 防止方向翻转或为 0：保证 delta >= minPositive
// 		if deltas[i] < minPositive {
// 			deltas[i] = minPositive
// 		}
// 	}

// 	// 4) 修正总和：按比例缩放让 sum(deltas) == totalDisp
// 	sum := 0.0
// 	for _, v := range deltas {
// 		sum += v
// 	}
// 	if sum <= 0 {
// 		// 极端情况下退化为等分
// 		each := totalDisp / float64(totalFrames)
// 		for i := range deltas {
// 			deltas[i] = each
// 		}
// 	} else {
// 		scale := totalDisp / sum
// 		for i := range deltas {
// 			deltas[i] *= scale
// 		}
// 	}

// 	// 5) 最后微调：把可能的浮点误差修到最后一帧，确保精确相等
// 	sum = 0.0
// 	for _, v := range deltas {
// 		sum += v
// 	}
// 	diff := totalDisp - sum
// 	deltas[totalFrames-1] += diff // 这不会改变方向

// 	// 6) 恢复原始符号
// 	for i := range deltas {
// 		deltas[i] *= sign
// 	}

// 	return deltas
// }
