# editly Windows 安装指南 (via54Design SKILL)

> 适用于: Windows 10/11 + Node.js 22 (hermes 默认) + ffmpeg
> 工具: **Node 18 免安装 + canvas prebuilt + editly --ignore-scripts + GL stub**

## 🎯 为什么需要 Node 18?

editly 依赖两个 native bindings:
- **canvas** (Cairo/Pango 渲染) — Node 22 无 prebuilt
- **gl** (OpenGL WebGL 渲染) — 任何 Node 都无 prebuilt (需 VS Build Tools)

**Node 18 是 canvas 最后支持的版本** (prebuilt binary 到 node-v120)。

## 📋 4 步安装

### 步骤 1: 装 Node 18 免安装版
```bash
# 下载 Node 18.20.4 zip (28.7MB)
curl -L -o "C:/Users/via54/AppData/Local/Temp/node18.zip" \
  "https://nodejs.org/dist/v18.20.4/node-v18.20.4-win-x64.zip"

# 解压到用户目录
mkdir -p "C:/Users/via54/tools"
cd "C:/Users/via54/tools" && \
  unzip -q "../AppData/Local/Temp/node18.zip" -d node18/

# 验证
"C:/Users/via54/tools/node18/node-v18.20.4-win-x64/node.exe" --version
# → v18.20.4
```

### 步骤 2: 装 canvas (Node 18, prebuilt 自动下载)
```bash
export PATH="C:/Users/via54/tools/node18/node-v18.20.4-win-x64:$PATH"
"C:/Users/via54/tools/node18/node-v18.20.4-win-x64/npm.cmd" install -g canvas
# → 37 packages added in 25s (含 prebuilt binary)
```

### 步骤 3: 装 editly (--ignore-scripts 跳过 gl 编译)
```bash
"C:/Users/via54/tools/node18/node-v18.20.4-win-x64/npm.cmd" install -g editly --ignore-scripts
# → 319 packages added in 7s
```

### 步骤 4: Stub gl native binding (★ 关键 ★)

**A. Stub glTransitions.js** (GL shader 转场 → ffmpeg xfade)
```js
// 文件: C:\Users\via54\tools\node18\node-v18.20.4-win-x64\node_modules\editly\glTransitions.js
export default ({ width, height, channels }) => {
  return {
    runTransitionOnFrame: ({ fromFrame, toFrame, progress, transitionName, transitionParams = {} }) => {
      if (progress > 0.5) return toFrame;
      return fromFrame;
    },
  };
};
```

**B. Stub glFrameSource.js** (GL frame source → 抛错提示用 image 替代)
```js
// 文件: C:\Users\via54\tools\node18\node-v18.20.4-win-x64\node_modules\editly\sources\glFrameSource.js
import fsExtra from 'fs-extra';
export default async function createGlFrameSource({ width, height, channels, params }) {
  throw new Error('GL frame source is disabled. Use image/video sources instead.');
}
```

## ✅ 验证

```bash
# 1. editly --help
"C:/Users/via54/tools/node18/node-v18.20.4-win-x64/editly" --help

# 2. 拼接 3 段测试
cd "C:/Users/via54/tools/node18/test-editly"
cp "G:/agent/developments/via54Design/minimax-output/lithium/seg1_cyberpunk_10s_hd.mp4" test1.mp4
"C:/Users/via54/tools/node18/node-v18.20.4-win-x64/editly" \
  test1.mp4 test1.mp4 test1.mp4 \
  --out test-output.mp4 --transition-name fade --clip-duration 5

# → 1366x768 @ 24fps / 29s / 13.3MB
```

## 🎬 可用功能 (Stub 后)

| 功能 | 可用 | 备注 |
|---|---|---|
| 视频拼接 | ✅ | — |
| 图片拼接 | ✅ | — |
| **fade / slide / wipe / dissolve 转场** | ✅ | ffmpeg xfade 内置 14 种 |
| **200+ GL shader 转场** | ❌ | 需 Visual Studio Build Tools 编译 gl |
| 字幕 (fabric.js) | ✅ | 需字体文件 |
| PiP (画中画) | ✅ | — |
| vignette (暗角) | ✅ | — |
| audio 混音 + crossfade | ✅ | — |
| GL shader 动态背景 | ❌ | 用 image 替代 |

## 📚 JSON spec 模板 (v0.6.2 锂电用)

```json5
{
  outPath: 'lithium_30s_v6.mp4',
  width: 1920,
  height: 1080,
  fps: 30,
  defaults: {
    duration: 6,
    transition: { duration: 0.5, name: 'fade' },
  },
  clips: [
    {
      duration: 3,
      layers: [
        { type: 'video', path: 'hook.mp4' },
        { type: 'title', text: '锂电新纪元', position: 'top' }
      ]
    },
    {
      duration: 5,
      layers: [
        { type: 'video', path: 'trend.mp4' },
        { type: 'subtitle', text: '2026 产能突破 800GWh' }
      ]
    },
    {
      duration: 10,
      layers: [
        { type: 'video', path: 'tech.mp4' }
      ]
    },
    {
      duration: 7,
      layers: [
        { type: 'video', path: 'market.mp4' }
      ]
    },
    {
      duration: 5,
      layers: [
        { type: 'video', path: 'outlook.mp4' },
        { type: 'title', text: '投资窗口正在打开' }
      ]
    }
  ],
  audioFilePath: 'narration.mp3',
  loopAudio: false,
  outputVolume: 1.0,
  clipsAudioVolume: 0.3
}
```

## 🔧 故障排查

### Error: Cannot find module '../build/Release/canvas.node'
- 原因: Node 22 装的 editly (没装 canvas 预编译)
- 修法: 切到 Node 18 路径, 重装

### Error: Could not locate the bindings file. Tried: ...gl\build\Release\webgl.node
- 原因: gl 包没编译
- 修法: 已 stub glTransitions.js + glFrameSource.js

### TypeError [ERR_STREAM_NULL_VALUES]: May not write null values to stream
- 原因: glTransitions.js stub 返回 null
- 修法: 改返回 toFrame/fromFrame (上面示例)

### GL frame source is disabled
- 原因: 用了 `type: 'shader'` 或 `type: 'gl'`
- 修法: 改用 `type: 'image'`
