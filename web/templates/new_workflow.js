// ── Intent-Driven Workflow ──
let currentIntent = 'design'
let uploadedImagePath = ''
let s2pPath = ''
let vbPaths = []

// Intent switching
document.querySelectorAll('.intent-btn').forEach(btn => {
  btn.addEventListener('click', () => {
    document.querySelectorAll('.intent-btn').forEach(b => b.classList.remove('active'))
    document.querySelectorAll('.intent-pane').forEach(p => p.classList.remove('active'))
    btn.classList.add('active')
    currentIntent = btn.dataset.intent
    document.querySelector(`.intent-pane[data-intent="${currentIntent}"]`).classList.add('active')
  })
})
// Default: first intent
document.querySelector('.intent-btn').click()

// ── Intent: 做设计 ──
document.getElementById('dGenBtn').addEventListener('click', async function() {
  const title = document.getElementById('dTitle').value || 'via54Design'
  const pres = document.getElementById('dMode').value === 'true'
  const btn = this; btn.disabled = true; btn.textContent = '⏳ 生成中...'
  try {
    const r = await fetch('/api/generate', { method:'POST', headers:{'Content-Type':'application/json'},
      body: JSON.stringify({title, layout:'hero', color:'ink-wash', font:'ming-hei-editorial', presentation:pres}) })
    const d = await r.json()
    document.getElementById('dOutput').textContent = d.html ? `✅ 已生成 (${d.length} bytes)` : (d.error || 'No output')
  } catch(e) { document.getElementById('dOutput').textContent = 'Error: '+e.message }
  btn.disabled = false; btn.textContent = '🎨 生成 HTML'
})

// ── Intent: 写提示词 ──
document.getElementById('pUploadZone').addEventListener('click', () => document.getElementById('pFileInput').click())
document.getElementById('pFileInput').addEventListener('change', async () => {
  const f = document.getElementById('pFileInput').files[0]
  if (!f) return
  const fd = new FormData(); fd.append('image', f)
  try {
    const r = await fetch('/api/upload', { method:'POST', body: fd })
    const d = await r.json()
    if (d.url) { uploadedImagePath = d.path; document.getElementById('pPreview').style.display = 'block'; document.getElementById('pPreviewImg').src = d.url }
  } catch(e) {}
})

document.getElementById('pGenBtn').addEventListener('click', async function() {
  const scene = document.getElementById('pScene').value || 'test scene'
  const platform = document.getElementById('pPlatform').value
  const btn = this; btn.disabled = true; btn.textContent = '⏳ 生成中...'
  try {
    const r = await fetch('/api/prompt', { method:'POST', headers:{'Content-Type':'application/json'},
      body: JSON.stringify({scene, platform, format:'markdown'}) })
    const d = await r.json()
    document.getElementById('pOutput').textContent = d.output || d.error || 'No output'
  } catch(e) { document.getElementById('pOutput').textContent = 'Error: '+e.message }
  btn.disabled = false; btn.textContent = '✨ 生成提示词'
})

// Forge status in prompt tab
async function checkPForge() {
  try {
    const r = await fetch('http://localhost:7860/sdapi/v1/sd-models', {signal:AbortSignal.timeout(3000)})
    const d = await r.json(); const n = Array.isArray(d) ? d.length : 0
    document.getElementById('pForgeDot').className = 'dot ok'
    document.getElementById('pForgeText').textContent = n > 0 ? `✅ 已连接 (${n} 模型)` : '✅ 已连接'
  } catch(e) {
    document.getElementById('pForgeDot').className = 'dot off'
    document.getElementById('pForgeText').textContent = '❌ 未连接'
  }
}
checkPForge()

document.getElementById('pForgeBtn').addEventListener('click', async function() {
  const prompt = document.getElementById('pOutput').textContent
  if (!prompt || prompt.startsWith('Error') || prompt.startsWith('⏳')) { document.getElementById('pForgeOut').textContent = '请先生成提示词'; return }
  const btn = this; btn.disabled = true; btn.textContent = '⏳ 提交中...'
  try {
    const r = await fetch('/api/regen', { method:'POST', headers:{'Content-Type':'application/json'},
      body: JSON.stringify({prompt, workflow:'sdxl_txt2img'}) })
    const d = await r.json()
    document.getElementById('pForgeOut').textContent = d.offline_mode ? '⚠️ Forge 未运行。提示词已就绪可手动提交' : (d.submitted ? '✅ 已提交到 Forge!' : JSON.stringify(d))
  } catch(e) { document.getElementById('pForgeOut').textContent = 'Error: '+e.message }
  btn.disabled = false; btn.textContent = '⚡ 提交生成'
})

// ── Intent: 做演示 ──
document.getElementById('s2pUploadZone').addEventListener('click', () => document.getElementById('s2pFileInput').click())
document.getElementById('s2pFileInput').addEventListener('change', async () => {
  const f = document.getElementById('s2pFileInput').files[0]; if (!f) return
  const fd = new FormData(); fd.append('file', f)
  try {
    const r = await fetch('/api/upload', { method:'POST', body: fd })
    const d = await r.json()
    if (d.path) { s2pPath = d.path; document.getElementById('s2pFileInfo').style.display = 'block'; document.getElementById('s2pFileInfo').textContent = `📎 ${f.name}` }
  } catch(e) {}
})

document.getElementById('s2pGenBtn').addEventListener('click', async function() {
  const desc = document.getElementById('s2pDesc').value || ''
  const btn = this; btn.disabled = true; btn.textContent = '⏳ 分析中...'
  const guidance = document.getElementById('s2pGuidance'); guidance.style.display = 'block'; guidance.textContent = '分析中...'
  const frames = document.getElementById('s2pFrames'); frames.style.display = 'block'
  
  if (s2pPath) {
    try {
      const r = await fetch('/api/story2ppt', { method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify({path:s2pPath, prompt:desc}) })
      const d = await r.json()
      if (d.error) { guidance.textContent = d.error; return }
      guidance.textContent = d.user_guidance || `✅ ${d.type||'文件'}，${d.total_slides} 页框架`
      frames.textContent = (d.slides||[]).map((s,i) => `${i+1}. ${s.title||s.type}`).join('\n')
    } catch(e) { guidance.textContent = 'Error: '+e.message }
  } else if (desc) {
    try {
      const r = await fetch('/api/narrate', { method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify({seed:desc, model:'three-act', duration:30, format:'markdown'}) })
      const d = await r.json()
      guidance.textContent = '📖 基于叙事种子生成的演示框架'
      frames.textContent = d.output || d.error || 'No output'
    } catch(e) { guidance.textContent = 'Error: '+e.message }
  } else { guidance.textContent = '请上传文件或输入描述'; frames.style.display = 'none' }
  btn.disabled = false; btn.textContent = '📊 生成框架'
})

// ── Intent: 做视频 ──
document.getElementById('vbUploadZone').addEventListener('click', () => document.getElementById('vbFileInput').click())
document.getElementById('vbFileInput').addEventListener('change', async () => {
  const files = Array.from(document.getElementById('vbFileInput').files)
  const previews = document.getElementById('vbPreviews'); previews.innerHTML = ''; vbPaths = []
  for (const f of files) {
    const fd = new FormData(); fd.append('image', f)
    try {
      const r = await fetch('/api/upload', { method:'POST', body: fd })
      const d = await r.json()
      if (d.path) { vbPaths.push(d.path); const img = document.createElement('img'); img.src = d.url; img.style.cssText = 'width:50px;height:50px;object-fit:cover;border-radius:4px'; previews.appendChild(img) }
    } catch(e) {}
  }
})

document.getElementById('vbGenBtn').addEventListener('click', async function() {
  if (!vbPaths.length) { document.getElementById('vbOutput').style.display = 'block'; document.getElementById('vbOutput').textContent = '请先上传故事板图片'; return }
  const model = document.getElementById('vbModel').value; const duration = parseInt(document.getElementById('vbDuration').value) || 30
  const btn = this; btn.disabled = true; btn.textContent = '⏳ 生成中...'
  document.getElementById('vbOutput').style.display = 'block'; document.getElementById('vbOutput').textContent = '分析图片并构建叙事...'
  try {
    const r = await fetch('/api/storyboard', { method:'POST', headers:{'Content-Type':'application/json'},
      body: JSON.stringify({paths:vbPaths, model, duration, single: vbPaths.length===1}) })
    const d = await r.json()
    let out = ''
    if (d.narrative_scaffold) {
      out += `📖 ${d.narrative_scaffold.model_name} | ${d.narrative_scaffold.total_duration}s\n`
      ;(d.narrative_scaffold.beats||[]).forEach(b => { out += `\n${b.name} (${b.start_time}s-${b.start_time+b.duration}s) [${b.mood}]` })
    }
    if (d.video_prompts) { out += '\n\n📝 Prompts:\n'; d.video_prompts.forEach(p => { out += `\nScene ${p.scene} (${p.duration}s): ${p.prompt.substring(0,80)}` }) }
    document.getElementById('vbOutput').textContent = out || JSON.stringify(d, null, 2).substring(0, 1000)
  } catch(e) { document.getElementById('vbOutput').textContent = 'Error: '+e.message }
  btn.disabled = false; btn.textContent = '🎬 生成脚本'
})

// ── Intent: 提交 Forge ──
async function checkFForge() {
  try {
    const r = await fetch('http://localhost:7860/sdapi/v1/sd-models', {signal:AbortSignal.timeout(3000)})
    const d = await r.json(); const n = Array.isArray(d) ? d.length : 0
    document.getElementById('fStatusDot').className = 'dot ok'
    document.getElementById('fStatusText').textContent = n > 0 ? `✅ 已连接 (${n} 模型)` : '✅ 已连接'
  } catch(e) {
    document.getElementById('fStatusDot').className = 'dot off'
    document.getElementById('fStatusText').textContent = '❌ 未运行'
  }
}
checkFForge()

document.getElementById('fSubmitBtn').addEventListener('click', async function() {
  const prompt = document.getElementById('fPrompt').value
  const negative = document.getElementById('fNeg').value
  const workflow = document.getElementById('fWorkflow').value
  if (!prompt) { document.getElementById('fOutput').textContent = '请先输入提示词'; return }
  const btn = this; btn.disabled = true; btn.textContent = '⏳ 提交中...'
  try {
    const r = await fetch('/api/regen', { method:'POST', headers:{'Content-Type':'application/json'},
      body: JSON.stringify({prompt, negative, workflow}) })
    const d = await r.json()
    document.getElementById('fOutput').textContent = d.submitted ? `✅ 已提交到 Forge！\n工作流: ${d.workflow}` : (d.offline_mode ? '⚠️ Forge 未运行。提示词已就绪：\n\n' + prompt : JSON.stringify(d))
  } catch(e) { document.getElementById('fOutput').textContent = 'Error: '+e.message }
  btn.disabled = false; btn.textContent = '⚡ 提交生成'
})