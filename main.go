package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"github.com/BurntSushi/toml"
	"golang.org/x/image/webp"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// Windows用コンソール非表示
func hideConsole() {
	if runtime.GOOS == "windows" {
		kernel32 := syscall.NewLazyDLL("kernel32.dll")
		user32 := syscall.NewLazyDLL("user32.dll")
		
		getConsoleWindow := kernel32.NewProc("GetConsoleWindow")
		showWindow := user32.NewProc("ShowWindow")
		
		hwnd, _, _ := getConsoleWindow.Call()
		if hwnd != 0 {
			showWindow.Call(hwnd, 0) // SW_HIDE = 0
		}
	}
}

// ログ出力（verboseモード時のみ）
func (iv *ImageViewer) logf(format string, args ...interface{}) {
	if iv.verbose {
		fmt.Printf(format, args...)
	}
}

// ログ出力（verboseモード時のみ）
func (iv *ImageViewer) log(msg string) {
	if iv.verbose {
		fmt.Println(msg)
	}
}
type Config struct {
	Display struct {
		Width  int `toml:"width"`
		Height int `toml:"height"`
	} `toml:"display"`
}

// アプリケーション構造体
type ImageViewer struct {
	app         fyne.App
	window      fyne.Window
	config      Config
	alwaysOnTop bool
	quitOnKey   bool
	verbose     bool
}

// 設定ファイル読み込み
func (iv *ImageViewer) loadConfig() error {
	// 実行ファイルのディレクトリを取得
	execPath, err := os.Executable()
	if err != nil {
		iv.logf("実行ファイルのパスを取得できません: %v\n", err)
		execPath = "."
	}
	execDir := filepath.Dir(execPath)
	
	// 設定ファイルのパスを構築
	configPath := filepath.Join(execDir, "config.toml")
	
	// デバッグ情報を出力（verboseモード時のみ）
	iv.logf("実行ファイルのディレクトリ: %s\n", execDir)
	iv.logf("設定ファイルのパス: %s\n", configPath)
	
	// 現在の作業ディレクトリも出力
	currentDir, _ := os.Getwd()
	iv.logf("現在の作業ディレクトリ: %s\n", currentDir)
	
	// 代替パス（現在のディレクトリ）も確認
	altConfigPath := "config.toml"
	iv.logf("代替設定ファイルパス: %s\n", altConfigPath)
	
	// 実行ファイルと同じディレクトリの設定ファイルを優先
	if _, err := os.Stat(configPath); err == nil {
		iv.logf("設定ファイルが見つかりました: %s\n", configPath)
		if _, err := toml.DecodeFile(configPath, &iv.config); err != nil {
			return fmt.Errorf("設定ファイルの読み込みに失敗しました: %v", err)
		}
		return nil
	}
	
	// 現在のディレクトリの設定ファイルを確認
	if _, err := os.Stat(altConfigPath); err == nil {
		iv.logf("代替設定ファイルが見つかりました: %s\n", altConfigPath)
		if _, err := toml.DecodeFile(altConfigPath, &iv.config); err != nil {
			return fmt.Errorf("設定ファイルの読み込みに失敗しました: %v", err)
		}
		return nil
	}
	
	// 設定ファイルが見つからない場合
	iv.logf("設定ファイルが見つかりません。デフォルト設定を作成します: %s\n", altConfigPath)
	
	// デフォルト設定を設定
	iv.config.Display.Width = 800
	iv.config.Display.Height = 600
	
	// デフォルト設定ファイルを作成
	return iv.createDefaultConfig()
}

// デフォルト設定ファイル作成
func (iv *ImageViewer) createDefaultConfig() error {
	configContent := `[display]
width = 800
height = 600
`
	configPath := "config.toml"
	
	iv.logf("デフォルト設定ファイルを作成中: %s\n", configPath)
	
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		iv.logf("設定ファイルの作成に失敗しました: %v\n", err)
		return err
	}
	
	iv.logf("デフォルト設定ファイルを作成しました: %s\n", configPath)
	return nil
}

// 画像サイズ計算
func (iv *ImageViewer) calculateDisplaySize(imgWidth, imgHeight int) (int, int) {
	configW := iv.config.Display.Width
	configH := iv.config.Display.Height
	
	// 設定値が0の場合はデフォルト値を使用
	if configW <= 0 {
		configW = 800
	}
	if configH <= 0 {
		configH = 600
	}
	
	// アスペクト比を保持して表示サイズを計算
	imageAspect := float64(imgWidth) / float64(imgHeight)
	configAspect := float64(configW) / float64(configH)
	
	var displayW, displayH int
	if imageAspect > configAspect {
		// 幅に合わせる
		displayW = configW
		displayH = int(float64(configW) / imageAspect)
	} else {
		// 高さに合わせる
		displayH = configH
		displayW = int(float64(configH) * imageAspect)
	}
	
	return displayW, displayH
}

// 画像読み込み
func (iv *ImageViewer) loadImage(imagePath string) (*canvas.Image, int, int, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("画像ファイルを開けません: %v", err)
	}
	defer file.Close()

	// ファイル拡張子を取得
	ext := strings.ToLower(filepath.Ext(imagePath))
	
	var img image.Image
	
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	case ".gif":
		img, err = gif.Decode(file)
	case ".webp":
		img, err = webp.Decode(file)
	default:
		// 汎用デコード
		img, _, err = image.Decode(file)
	}
	
	if err != nil {
		return nil, 0, 0, fmt.Errorf("画像をデコードできません: %v", err)
	}

	// 画像サイズを取得
	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()
	
	// 設定ファイルに基づいて表示サイズを計算
	displayW, displayH := iv.calculateDisplaySize(imgWidth, imgHeight)
	
	// ファイルからリソースを作成
	fileResource, err := fyne.LoadResourceFromPath(imagePath)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("画像リソースの作成に失敗しました: %v", err)
	}
	
	canvasImage := canvas.NewImageFromResource(fileResource)
	canvasImage.FillMode = canvas.ImageFillContain // アスペクト比を保持
	canvasImage.Resize(fyne.NewSize(float32(displayW), float32(displayH)))
	
	return canvasImage, displayW, displayH, nil
}

// ウィンドウ初期化
func (iv *ImageViewer) initWindow(imagePath string) error {
	iv.app = app.New()
	iv.app.SetIcon(nil) // アイコンなし
	
	iv.window = iv.app.NewWindow("go_ImagePreviewer")
	
	// ウィンドウ設定
	iv.window.SetPadded(false)
	
	// 最前面表示の代替実装（一部のプラットフォームでサポート）
	if iv.alwaysOnTop {
		// Fyneでは直接的な最前面表示メソッドがないため
		// ウィンドウのフルスクリーン設定やフォーカス維持で代替
		iv.window.RequestFocus()
	}
	
	// 画像読み込み（設定ファイルの値を使用）
	imageCanvas, displayW, displayH, err := iv.loadImage(imagePath)
	if err != nil {
		return err
	}
	
	// 設定ファイルで指定されたサイズでウィンドウを作成
	iv.window.Resize(fyne.NewSize(float32(displayW), float32(displayH)))
	iv.window.CenterOnScreen()
	
	// コンテンツ設定
	iv.window.SetContent(imageCanvas)
	
	// キーイベント設定
	iv.window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		if key.Name == fyne.KeyEscape && !iv.quitOnKey {
			iv.app.Quit()
		} else if iv.quitOnKey {
			iv.app.Quit()
		}
	})
	
	// ウィンドウの装飾を最小限に
	iv.window.SetFixedSize(true)
	
	return nil
}

// 実行
func (iv *ImageViewer) run() {
	iv.window.ShowAndRun()
}

// 使用方法表示
func showUsage() {
	fmt.Println("使用方法:")
	fmt.Println("  go_ImagePreviewer [オプション] <画像ファイル>")
	fmt.Println("")
	fmt.Println("オプション:")
	fmt.Println("  -Q         最前面表示モード（任意のキーで終了）")
	fmt.Println("  --verbose  詳細情報を表示")
	fmt.Println("")
	fmt.Println("キー操作:")
	fmt.Println("  ESC   アプリケーション終了（通常モード）")
	fmt.Println("  任意  アプリケーション終了（-Qモード）")
	fmt.Println("")
	fmt.Println("対応画像形式:")
	fmt.Println("  jpg, jpeg, png, gif, webp")
}

func main() {
	var viewer ImageViewer
	
	// コンソールウィンドウを非表示（Windows）
	hideConsole()
	
	// 引数解析
	args := os.Args[1:]
	var imagePath string
	
	if len(args) == 0 {
		showUsage()
		return
	}
	
	// オプション解析
	for _, arg := range args {
		if arg == "-Q" {
			viewer.alwaysOnTop = true
			viewer.quitOnKey = true
		} else if arg == "--verbose" {
			viewer.verbose = true
		} else if !strings.HasPrefix(arg, "-") {
			imagePath = arg
		}
	}
	
	// verboseモードの場合、コンソールを再表示
	if viewer.verbose && runtime.GOOS == "windows" {
		if syscall.NewLazyDLL("kernel32.dll").NewProc("AllocConsole").Find() == nil {
			syscall.NewLazyDLL("kernel32.dll").NewProc("AllocConsole").Call()
		}
	}
	
	if imagePath == "" {
		viewer.log("エラー: 画像ファイルが指定されていません")
		if viewer.verbose {
			showUsage()
		}
		os.Exit(1)
	}
	
	// ファイル存在確認
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		viewer.logf("エラー: 画像ファイルが見つかりません: %s\n", imagePath)
		os.Exit(1)
	}
	
	viewer.log("=== デバッグ情報 ===")
	
	// 設定ファイル読み込み
	if err := viewer.loadConfig(); err != nil {
		viewer.logf("設定ファイルの読み込みに失敗しました: %v\n", err)
		// エラーでも続行（デフォルト値を使用）
	}
	
	// 設定値をログ出力（verboseモード時のみ）
	viewer.logf("読み込まれた設定値: 幅=%d, 高さ=%d\n", viewer.config.Display.Width, viewer.config.Display.Height)
	viewer.log("==================")
	
	// ウィンドウ初期化
	if err := viewer.initWindow(imagePath); err != nil {
		if viewer.verbose {
			log.Printf("ウィンドウの初期化に失敗しました: %v", err)
		}
		os.Exit(1)
	}
	
	// 実行
	viewer.run()
}