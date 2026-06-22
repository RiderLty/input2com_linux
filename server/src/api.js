import { useCallback, useEffect, useState } from 'react'

export function usePreConfig() {
  const [data, setData] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    fetch('/api/get/preConfig')
      .then(r => r.json())
      .then(d => { setData(d || {}); setLoading(false) })
      .catch(e => { setError(e.message); setLoading(false) })
  }, [])

  return { data: data || {}, loading, error }
}

export function useActiveConfig() {
  const [active, setActive] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch('/api/get/active')
      .then(r => r.json())
      .then(d => { setActive(d.active || ''); setLoading(false) })
      .catch(() => setLoading(false))
  }, [])

  const set = (name) => {
    setActive(name)
    fetch(`/api/set/active?name=${encodeURIComponent(name)}`).catch(() => {})
  }

  return { active: active ?? '', setActive: set, loading }
}

export function useRecoilInput() {
  const [content, setContent] = useState('')
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [saved, setSaved] = useState(false)

  useEffect(() => {
    fetch('/api/get/recoilInput')
      .then(r => r.text())
      .then(t => { setContent(t); setLoading(false) })
      .catch(() => setLoading(false))
  }, [])

  const save = useCallback(async () => {
    setSaving(true)
    setSaved(false)
    try {
      const resp = await fetch('/api/set/recoilInput', {
        method: 'POST',
        body: content,
      })
      if (resp.ok) setSaved(true)
    } finally {
      setSaving(false)
      setTimeout(() => setSaved(false), 2000)
    }
  }, [content])

  return { content, setContent, save, saving, saved, loading }
}

export function useConfigs() {
  const [configs, setConfigs] = useState([])
  const [activeContent, setActiveContent] = useState('')
  const [activeName, setActiveName] = useState('')
  const [loading, setLoading] = useState(true)

  const refresh = useCallback(() => {
    setLoading(true)
    Promise.all([
      fetch('/api/get/configs').then(r => r.json()),
      fetch('/api/get/activeConfig').then(r => r.text()),
      fetch('/api/get/activeConfigName').then(r => r.text()),
    ]).then(([names, content, name]) => {
      setConfigs(names || [])
      setActiveContent(content)
      setActiveName(name || '')
      setLoading(false)
    }).catch(() => setLoading(false))
  }, [])

  useEffect(() => { refresh() }, [refresh])

  const loadConfig = useCallback(async (name) => {
    const resp = await fetch(`/api/get/config?name=${encodeURIComponent(name)}`)
    return resp.text()
  }, [])

  const saveConfig = useCallback(async (name, content) => {
    const resp = await fetch(`/api/set/config?name=${encodeURIComponent(name)}`, {
      method: 'POST',
      body: content,
    })
    if (resp.ok) refresh()
    return resp.ok
  }, [refresh])

  const deleteConfig = useCallback(async (name) => {
    const resp = await fetch(`/api/delete/config?name=${encodeURIComponent(name)}`, {
      method: 'DELETE',
    })
    if (resp.ok) refresh()
    return resp.ok
  }, [refresh])

  const applyConfig = useCallback(async (name) => {
    const resp = await fetch(`/api/apply/config?name=${encodeURIComponent(name)}`, {
      method: 'POST',
    })
    if (resp.ok) refresh()
    return resp.ok
  }, [refresh])

  return { configs, activeContent, activeName, loading, refresh, loadConfig, saveConfig, deleteConfig, applyConfig }
}
