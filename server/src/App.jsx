import React, { useMemo, useState } from 'react'
import {
  Box, CircularProgress, CssBaseline, Typography,
  ThemeProvider, createTheme, responsiveFontSizes,
} from '@mui/material'
import DeleteSweepIcon from '@mui/icons-material/DeleteSweep'
import EditNoteIcon from '@mui/icons-material/EditNote'
import RestartAltIcon from '@mui/icons-material/RestartAlt'
import { usePreConfig, useActiveConfig, useRecoilInput } from './api'

function Tag({ label, color = 'default', sx = {} }) {
  const colors = {
    default: { bg: 'action.hover', fg: 'text.secondary' },
    primary: { bg: 'rgba(0,121,107,0.15)', fg: 'primary.main' },
    primarySolid: { bg: 'primary.main', fg: '#fff' },
    error: { bg: 'rgba(211,47,47,0.12)', fg: 'error.main' },
    errorSolid: { bg: 'error.main', fg: '#fff' },
    white: { bg: 'rgba(255,255,255,0.2)', fg: '#fff' },
  }
  const c = colors[color] || colors.default
  return (
    <Box component="span" sx={{
      display: 'inline-flex',
      alignItems: 'center',
      px: 1,
      py: 0.25,
      borderRadius: 1,
      fontSize: 12,
      fontWeight: 600,
      lineHeight: 1.5,
      bgcolor: c.bg,
      color: c.fg,
      whiteSpace: 'nowrap',
      ...sx,
    }}>
      {label}
    </Box>
  )
}

export default function App() {
  const prefersDark = usePrefersDark()
  const theme = useMemo(() => responsiveFontSizes(createTheme({
    palette: {
      mode: prefersDark ? 'dark' : 'light',
      primary: { main: '#00796B' },
      secondary: { main: '#d90051' },
      ...(prefersDark
        ? { background: { paper: '#1e1e1e', default: '#121212' } }
        : { background: { paper: '#ffffff', default: '#f5f5f5' } }
      ),
    },
    shape: { borderRadius: 16 },
    typography: {
      fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
    },
  })), [prefersDark])

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Box sx={{
        minHeight: '100vh',
        bgcolor: 'background.default',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        p: 3,
      }}>
        <ConfigPanel />
      </Box>
    </ThemeProvider>
  )
}

function ConfigPanel() {
  const { data, loading, error } = usePreConfig()
  const { active, setActive, loading: activeLoading } = useActiveConfig()
  const [applying, setApplying] = useState(null)
  const [restarting, setRestarting] = useState(false)

  const configs = useMemo(() => {
    const keys = Object.keys(data)
    const isAscii = (s) => s.charCodeAt(0) < 128
    return keys.sort((a, b) => {
      const sa = isAscii(a) ? 0 : 1
      const sb = isAscii(b) ? 0 : 1
      if (sa !== sb) return sa - sb
      return a.localeCompare(b, 'zh')
    })
  }, [data])

  const handleApply = async (name) => {
    setApplying(name)
    try {
      const pair = data[name]
      const mouseCfg = pair?.[0] || {}
      const keyCfg = pair?.[1] || {}
      await fetch('/api/set/mouse?key=CLEAR_ALL&value=NONE')
      for (const k in mouseCfg) await fetch(`/api/set/mouse?key=${k}&value=${mouseCfg[k]}`)
      await fetch('/api/set/keyboard?key=CLEAR_ALL&value=NONE')
      for (const k in keyCfg) await fetch(`/api/set/keyboard?key=${k}&value=${keyCfg[k]}`)
      setActive(name)
    } finally {
      setApplying(null)
    }
  }

  const handleClear = async () => {
    setApplying('__clear__')
    try {
      await fetch('/api/set/mouse?key=CLEAR_ALL&value=NONE')
      await fetch('/api/set/keyboard?key=CLEAR_ALL&value=NONE')
      setActive('')
    } finally {
      setApplying(null)
    }
  }

  const handleRestart = () => {
    setRestarting(true)
    fetch('/api/restart').catch(() => {})
    // 定时刷新页面，服务器重启后自动恢复
    const timer = setInterval(() => {
      fetch('/').then(() => { clearInterval(timer); window.location.reload() }).catch(() => {})
    }, 1000)
  }

  if (loading || activeLoading) {
    return (
      <Box sx={{ textAlign: 'center' }}>
        <CircularProgress size={48} />
        <Typography sx={{ mt: 2, color: 'text.secondary' }}>加载中...</Typography>
      </Box>
    )
  }

  if (error) {
    return (
      <Typography color="error" variant="h6">
        加载失败: {error}
      </Typography>
    )
  }

  const isClearActive = active === '' && !activeLoading
  const isClearing = applying === '__clear__'

  return (
    <Box sx={{ width: '100%', maxWidth: 1000 }}>
      <Box sx={{ textAlign: 'center', mb: 5 }}>
        <Typography variant="h3" fontWeight={800} sx={{ letterSpacing: 2 }}>
          配置切换
        </Typography>
        <Typography variant="body1" color="text.secondary" sx={{ mt: 1 }}>
          选择一个配置方案，鼠标和键盘将同时切换
        </Typography>
      </Box>

      {/* 清空/重启按钮 */}
      <Box sx={{ display: 'flex', justifyContent: 'center', gap: 2, mb: 4 }}>
        <Box
          onClick={() => !isClearing && handleClear()}
          sx={{
            display: 'flex',
            alignItems: 'center',
            gap: 1.5,
            px: 4,
            py: 1.5,
            borderRadius: 2,
            cursor: isClearing ? 'wait' : 'pointer',
            userSelect: 'none',
            border: 2,
            borderColor: isClearActive ? 'error.main' : 'divider',
            bgcolor: 'transparent',
            color: isClearActive ? 'error.main' : 'text.secondary',
            transition: 'all 0.25s cubic-bezier(.4,0,.2,1)',
            opacity: isClearing ? 0.5 : 1,
            '&:hover': {
              borderColor: 'error.main',
              color: 'error.main',
              bgcolor: 'rgba(211,47,47,0.06)',
              transform: isClearing ? 'none' : 'translateY(-2px)',
            },
            '&:active': {
              transform: 'translateY(0)',
            },
          }}
        >
          <DeleteSweepIcon fontSize="small" />
          <Typography variant="body1" fontWeight={600}>
            {isClearing ? '清空中...' : '清空所有配置'}
          </Typography>
          {isClearActive && <Tag label="当前" color="error" sx={{ ml: 1 }} />}
        </Box>
        <Box
          onClick={() => !restarting && handleRestart()}
          sx={{
            display: 'flex',
            alignItems: 'center',
            gap: 1.5,
            px: 4,
            py: 1.5,
            borderRadius: 2,
            cursor: restarting ? 'wait' : 'pointer',
            userSelect: 'none',
            border: 2,
            borderColor: 'warning.main',
            bgcolor: 'transparent',
            color: 'warning.main',
            transition: 'all 0.25s cubic-bezier(.4,0,.2,1)',
            opacity: restarting ? 0.5 : 1,
            '&:hover': {
              borderColor: 'warning.dark',
              color: 'warning.dark',
              bgcolor: 'rgba(237,108,2,0.06)',
              transform: restarting ? 'none' : 'translateY(-2px)',
            },
            '&:active': {
              transform: 'translateY(0)',
            },
          }}
        >
          <RestartAltIcon fontSize="small" />
          <Typography variant="body1" fontWeight={600}>
            {restarting ? '重启中...' : '重启程序'}
          </Typography>
        </Box>
      </Box>

      {/* 配置卡片网格 */}
      <Box sx={{
        display: 'grid',
        gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, 1fr)', md: 'repeat(3, 1fr)' },
        gap: 3,
      }}>
        {configs.filter(n => n !== '清空').map(name => {
          const isSelected = active === name
          const isApplying = applying === name
          const pair = data[name]
          const mouseCfg = pair?.[0] || {}
          const keyCfg = pair?.[1] || {}
          const mouseCount = Object.keys(mouseCfg).length
          const keyCount = Object.keys(keyCfg).length

          return (
            <Box
              key={name}
              onClick={() => !isApplying && handleApply(name)}
              sx={{
                position: 'relative',
                p: 3.5,
                borderRadius: 3,
                cursor: isApplying ? 'wait' : 'pointer',
                userSelect: 'none',
                bgcolor: isSelected ? 'primary.main' : 'background.paper',
                border: 2,
                borderColor: isSelected ? 'primary.main' : 'divider',
                boxShadow: isSelected
                  ? '0 8px 32px rgba(0,121,107,0.35)'
                  : '0 2px 12px rgba(0,0,0,0.08)',
                transition: 'all 0.25s cubic-bezier(.4,0,.2,1)',
                opacity: isApplying ? 0.6 : 1,
                '&:hover': {
                  transform: isApplying ? 'none' : 'translateY(-6px)',
                  boxShadow: isApplying ? undefined
                    : isSelected ? '0 12px 40px rgba(0,121,107,0.45)'
                    : '0 12px 32px rgba(0,0,0,0.18)',
                  borderColor: 'primary.main',
                },
                '&:active': {
                  transform: 'translateY(0)',
                  boxShadow: '0 2px 8px rgba(0,0,0,0.12)',
                },
              }}
            >
              {isSelected && (
                <Tag label="当前" color="primarySolid" sx={{
                  position: 'absolute',
                  top: -10,
                  right: 14,
                }} />
              )}

              <Typography
                variant="h5"
                fontWeight={700}
                noWrap
                sx={{
                  color: isSelected ? '#fff' : 'text.primary',
                  mb: 2,
                }}
              >
                {isApplying ? '切换中...' : name}
              </Typography>

              <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                {mouseCount > 0 && (
                  <Tag label={`鼠标 ${mouseCount} 键`} color={isSelected ? 'white' : 'default'} />
                )}
                {keyCount > 0 && (
                  <Tag label={`键盘 ${keyCount} 键`} color={isSelected ? 'white' : 'default'} />
                )}
                {mouseCount === 0 && keyCount === 0 && (
                  <Typography variant="caption" sx={{
                    color: isSelected ? 'rgba(255,255,255,0.6)' : 'text.disabled',
                  }}>
                    仅鼠标
                  </Typography>
                )}
              </Box>
            </Box>
          )
        })}
      </Box>

      {/* 自定义压枪数据编辑器 */}
      <RecoilEditor />
    </Box>
  )
}

function RecoilEditor() {
  const { content, setContent, save, saving, saved, loading } = useRecoilInput()
  const [open, setOpen] = useState(true)

  if (loading) return null

  return (
    <Box sx={{ mt: 6 }}>
      <Box
        onClick={() => setOpen(!open)}
        sx={{
          display: 'flex',
          alignItems: 'center',
          gap: 1,
          cursor: 'pointer',
          mb: open ? 2 : 0,
          '&:hover': { opacity: 0.8 },
        }}
      >
        <EditNoteIcon fontSize="small" sx={{ color: 'text.secondary' }} />
        <Typography variant="body1" fontWeight={600} sx={{ color: 'text.secondary' }}>
          自定义压枪数据编辑
        </Typography>
        <Typography variant="caption" sx={{ color: 'text.disabled', ml: 1 }}>
          {open ? '收起' : '展开'}
        </Typography>
      </Box>

      {open && (
        <Box sx={{
          bgcolor: 'background.paper',
          border: 1,
          borderColor: 'divider',
          borderRadius: 2,
          p: 3,
        }}>
          <Typography variant="caption" sx={{
            color: 'text.disabled',
            mb: 2,
            display: 'block',
            lineHeight: 1.8,
          }}>
            格式：每行3个整数，用空格分隔 — 周期数 X移动量 Y移动量（每周期10ms）
            <br />
            示例：10 0 5 表示 100ms 内每步移动(0,5)，# 开头为注释行
          </Typography>

          <textarea
            value={content}
            onChange={e => setContent(e.target.value)}
            spellCheck={false}
            style={{
              width: '100%',
              minHeight: 200,
              fontFamily: 'Consolas, "Courier New", monospace',
              fontSize: 14,
              lineHeight: 1.6,
              padding: 12,
              border: '1px solid #555',
              borderRadius: 8,
              resize: 'vertical',
              backgroundColor: 'inherit',
              color: 'inherit',
              outline: 'none',
              boxSizing: 'border-box',
            }}
          />

          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mt: 2 }}>
            <Box
              onClick={() => !saving && save()}
              sx={{
                px: 3,
                py: 1,
                borderRadius: 2,
                cursor: saving ? 'wait' : 'pointer',
                bgcolor: 'primary.main',
                color: '#fff',
                fontWeight: 600,
                fontSize: 14,
                transition: 'all 0.2s',
                opacity: saving ? 0.6 : 1,
                '&:hover': { bgcolor: saving ? undefined : 'primary.dark' },
              }}
            >
              {saving ? '保存中...' : '保存'}
            </Box>
            {saved && (
              <Typography variant="body2" sx={{ color: 'primary.main', fontWeight: 600 }}>
                已保存
              </Typography>
            )}
          </Box>
        </Box>
      )}
    </Box>
  )
}

function usePrefersDark() {
  const [dark, setDark] = useState(
    () => window.matchMedia?.('(prefers-color-scheme: dark)').matches ?? false
  )
  React.useEffect(() => {
    const mq = window.matchMedia('(prefers-color-scheme: dark)')
    const handler = (e) => setDark(e.matches)
    mq.addEventListener('change', handler)
    return () => mq.removeEventListener('change', handler)
  }, [])
  return dark
}
