package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/yclw/go-vtracer"
)

const htmlTemplate = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>VTracer 在线图像矢量化</title>
    <style>
        body {
            font-family: 'Microsoft YaHei', Arial, sans-serif;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background: #f5f5f5;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            text-align: center;
            margin-bottom: 30px;
        }
        .form-group {
            margin-bottom: 20px;
        }
        label {
            display: block;
            margin-bottom: 5px;
            font-weight: bold;
            color: #555;
        }
        input[type="file"], select {
            width: 100%;
            padding: 10px;
            border: 2px solid #ddd;
            border-radius: 5px;
            font-size: 16px;
        }
        .config-section {
            background: #f9f9f9;
            padding: 20px;
            border-radius: 5px;
            margin: 20px 0;
        }
        .config-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
        }
        .btn {
            background: #007bff;
            color: white;
            padding: 12px 30px;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            font-size: 16px;
            width: 100%;
            margin-top: 20px;
        }
        .btn:hover {
            background: #0056b3;
        }
        .btn:disabled {
            background: #ccc;
            cursor: not-allowed;
        }
        .result {
            margin-top: 30px;
            padding: 20px;
            border: 2px solid #28a745;
            border-radius: 5px;
            background: #d4edda;
        }
        .error {
            color: #dc3545;
            background: #f8d7da;
            border-color: #dc3545;
        }
        .loading {
            text-align: center;
            color: #007bff;
        }
        .preview {
            max-width: 100%;
            border: 1px solid #ddd;
            border-radius: 5px;
            margin-top: 10px;
        }
        .download-link {
            display: inline-block;
            background: #28a745;
            color: white;
            padding: 10px 20px;
            text-decoration: none;
            border-radius: 5px;
            margin-top: 10px;
        }
        .stats {
            font-size: 14px;
            color: #666;
            margin-top: 10px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎨 VTracer 在线图像矢量化</h1>
        
        <form id="convertForm" enctype="multipart/form-data">
            <div class="form-group">
                <label for="image">选择图像文件:</label>
                <input type="file" id="image" name="image" accept="image/*" required>
            </div>
            
            <div class="form-group">
                <label for="preset">预设配置:</label>
                <select id="preset" name="preset">
                    <option value="">自定义配置</option>
                    <option value="bw">黑白模式 - 适用于简单图形</option>
                    <option value="poster" selected>海报模式 - 适用于插图海报</option>
                    <option value="photo">照片模式 - 适用于复杂照片</option>
                </select>
            </div>
            
            <div class="config-section" id="customConfig">
                <h3>自定义配置 (选择预设后自动设置)</h3>
                <div class="config-grid">
                    <div>
                        <label>颜色模式:</label>
                        <select name="colormode">
                            <option value="color">彩色</option>
                            <option value="binary">二值</option>
                        </select>
                    </div>
                    <div>
                        <label>路径模式:</label>
                        <select name="mode">
                            <option value="spline" selected>样条</option>
                            <option value="polygon">多边形</option>
                            <option value="pixel">像素</option>
                        </select>
                    </div>
                    <div>
                        <label>过滤斑点 (px):</label>
                        <input type="number" name="speckle" value="4" min="0" max="16">
                    </div>
                    <div>
                        <label>颜色精度 (1-8):</label>
                        <input type="number" name="precision" value="6" min="1" max="8">
                    </div>
                    <div>
                        <label>图层差异 (0-255):</label>
                        <input type="number" name="difference" value="16" min="0" max="255">
                    </div>
                    <div>
                        <label>角度阈值 (度):</label>
                        <input type="number" name="corner" value="60" min="0" max="180">
                    </div>
                </div>
            </div>
            
            <button type="submit" class="btn" id="submitBtn">开始转换</button>
        </form>
        
        <div id="result"></div>
    </div>

    <script>
        document.getElementById('convertForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const submitBtn = document.getElementById('submitBtn');
            const resultDiv = document.getElementById('result');
            
            submitBtn.disabled = true;
            submitBtn.textContent = '转换中...';
            resultDiv.innerHTML = '<div class="loading">正在处理图像，请稍候...</div>';
            
            const formData = new FormData(this);
            
            try {
                const response = await fetch('/convert', {
                    method: 'POST',
                    body: formData
                });
                
                if (response.ok) {
                    const svg = await response.text();
                    const blob = new Blob([svg], { type: 'image/svg+xml' });
                    const url = URL.createObjectURL(blob);
                    
                    resultDiv.innerHTML = ` + "`" + `
                        <div class="result">
                            <h3>✅ 转换成功！</h3>
                            <div class="stats">
                                SVG 大小: ${svg.length.toLocaleString()} 字符
                            </div>
                            <div>
                                <h4>预览:</h4>
                                <img src="${url}" class="preview" alt="SVG Preview">
                            </div>
                            <a href="${url}" download="converted.svg" class="download-link">
                                📥 下载 SVG 文件
                            </a>
                        </div>
                    ` + "`" + `;
                } else {
                    const error = await response.text();
                    resultDiv.innerHTML = ` + "`" + `
                        <div class="result error">
                            <h3>❌ 转换失败</h3>
                            <p>${error}</p>
                        </div>
                    ` + "`" + `;
                }
            } catch (error) {
                resultDiv.innerHTML = ` + "`" + `
                    <div class="result error">
                        <h3>❌ 网络错误</h3>
                        <p>${error.message}</p>
                    </div>
                ` + "`" + `;
            } finally {
                submitBtn.disabled = false;
                submitBtn.textContent = '开始转换';
            }
        });
    </script>
</body>
</html>
`

func main() {
	var (
		port = flag.String("port", "8080", "服务器端口")
		host = flag.String("host", "localhost", "服务器地址")
	)
	flag.Parse()

	// 主页处理器
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(htmlTemplate))
	})

	// 转换处理器
	http.HandleFunc("/convert", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "只支持 POST 方法", http.StatusMethodNotAllowed)
			return
		}

		// 解析表单数据
		err := r.ParseMultipartForm(10 << 20) // 10MB
		if err != nil {
			http.Error(w, "解析表单失败: "+err.Error(), http.StatusBadRequest)
			return
		}

		// 获取上传的文件
		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "获取图像文件失败: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// 检查文件大小
		if header.Size > 10<<20 { // 10MB
			http.Error(w, "文件太大，最大支持 10MB", http.StatusBadRequest)
			return
		}

		// 读取文件内容
		imageData, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "读取图像数据失败: "+err.Error(), http.StatusBadRequest)
			return
		}

		// 解码图像
		img, _, err := image.Decode(bytes.NewReader(imageData))
		if err != nil {
			http.Error(w, "图像格式不支持或文件损坏: "+err.Error(), http.StatusBadRequest)
			return
		}

		// 获取配置
		config := getConfigFromForm(r)

		// 记录开始时间
		startTime := time.Now()

		// 执行转换
		svg, err := vtracer.ConvertImage(img, config)
		if err != nil {
			http.Error(w, "图像转换失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 记录处理时间
		duration := time.Since(startTime)
		log.Printf("转换完成: 文件=%s, 大小=%d字节, 耗时=%v, SVG长度=%d",
			header.Filename, header.Size, duration.Round(time.Millisecond), len(svg))

		// 返回 SVG
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Content-Disposition", "attachment; filename=\"converted.svg\"")
		w.Write([]byte(svg))
	})

	// 启动服务器
	addr := *host + ":" + *port
	fmt.Printf("🚀 VTracer Web 服务启动成功！\n")
	fmt.Printf("📍 访问地址: http://%s\n", addr)
	fmt.Printf("💡 按 Ctrl+C 停止服务\n\n")

	log.Fatal(http.ListenAndServe(addr, nil))
}

func getConfigFromForm(r *http.Request) *vtracer.Config {
	preset := r.FormValue("preset")

	// 如果指定了预设，使用预设配置
	if preset != "" {
		switch preset {
		case "bw":
			return vtracer.NewConfigFromPreset(vtracer.PresetBW)
		case "poster":
			return vtracer.NewConfigFromPreset(vtracer.PresetPoster)
		case "photo":
			return vtracer.NewConfigFromPreset(vtracer.PresetPhoto)
		}
	}

	// 否则使用自定义配置
	config := vtracer.DefaultConfig()

	// 颜色模式
	if colorMode := r.FormValue("colormode"); colorMode == "binary" {
		config.ColorMode = vtracer.ColorModeBinary
	}

	// 路径模式
	switch r.FormValue("mode") {
	case "pixel":
		config.Mode = vtracer.PathModeNone
	case "polygon":
		config.Mode = vtracer.PathModePolygon
	case "spline":
		config.Mode = vtracer.PathModeSpline
	}

	// 数值参数
	if val, err := strconv.Atoi(r.FormValue("speckle")); err == nil && val >= 0 && val <= 16 {
		config.FilterSpeckle = val
	}
	if val, err := strconv.Atoi(r.FormValue("precision")); err == nil && val >= 1 && val <= 8 {
		config.ColorPrecision = val
	}
	if val, err := strconv.Atoi(r.FormValue("difference")); err == nil && val >= 0 && val <= 255 {
		config.LayerDifference = val
	}
	if val, err := strconv.Atoi(r.FormValue("corner")); err == nil && val >= 0 && val <= 180 {
		config.CornerThreshold = val
	}

	return config
}
