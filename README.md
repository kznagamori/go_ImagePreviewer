# go_ImagePreviewer (Fyne版)

シンプルで軽量な画像表示アプリケーションです。指定した画像ファイルをウィンドウで表示します。

## 特徴

- 複数の画像形式に対応（JPG、PNG、GIF、WebP）
- アスペクト比を保持した表示
- 設定ファイルによる表示サイズ調整
- コンソールウィンドウなしのクリーンな表示
- 最前面表示モード対応
- クロスプラットフォーム対応（Windows、macOS、Linux）

## システム要件

### Windows
- Windows 10以降
- 追加の依存関係なし

### macOS
- macOS 10.12以降
- Xcode Command Line Tools

### Linux
- 現代的なLinuxディストリビューション
- X11またはWayland

## インストール

### GitHubからインストール

```bash
go install github.com/kznagamori/go_ImagePreviewer@latest
```

### ローカルビルド

```bash
git clone https://github.com/kznagamori/go_ImagePreviewer.git
cd go_ImagePreviewer
go mod tidy
```

#### Windows用ビルド（コンソールなし）
```bash
# バッチファイルを使用
build.bat

# または、コマンドラインから
go build -ldflags="-s -w -H=windowsgui" -trimpath -o go_ImagePreviewer.exe

# または、Makefileを使用
make windows
```

#### その他のプラットフォーム
```bash
# Linux用
make linux

# macOS用
make macos

# 現在のプラットフォーム用
make build
```

## 使用方法

### 基本的な使用方法

```bash
go_ImagePreviewer image.jpg
```

### オプション

- `-Q`: 最前面表示モード（任意のキーで終了）
- `--verbose`: 詳細情報を表示（デバッグ用）

```bash
# 最前面表示モード
go_ImagePreviewer -Q image.png

# デバッグ情報付き
go_ImagePreviewer --verbose image.jpg

# オプション組み合わせ
go_ImagePreviewer -Q --verbose image.gif
```

### キー操作

- **通常モード**: `ESC`キーでアプリケーション終了
- **最前面モード（-Q）**: 任意のキーでアプリケーション終了

## 設定ファイル

実行ファイルと同じディレクトリまたは現在のディレクトリに`config.toml`ファイルを配置することで、表示サイズを設定できます。

```toml
[display]
width = 800   # 最大幅（ピクセル）
height = 600  # 最大高さ（ピクセル）
```

### 設定ファイルの優先順位

1. 実行ファイルと同じディレクトリの `config.toml`
2. 現在の作業ディレクトリの `config.toml`
3. デフォルト設定（幅800px、高さ600px）

画像はアスペクト比を保持したまま、指定されたサイズに収まるように表示されます。

## 対応画像形式

- JPEG (.jpg, .jpeg)
- PNG (.png)
- GIF (.gif)
- WebP (.webp)

## ファイル構成

```
go_ImagePreviewer/
├── main.go              # メインソースコード
├── config.toml          # 設定ファイル（オプション）
├── go.mod              # Goモジュールファイル
├── README.md           # このファイル
├── build.bat           # Windows用ビルドスクリプト
├── Makefile           # クロスプラットフォーム用ビルドファイル
└── go_ImagePreviewer.exe  # 実行ファイル（ビルド後）
```

## 技術仕様

### 使用ライブラリ

- **Fyne v2**: クロスプラットフォームGUIフレームワーク
- **BurntSushi/toml**: TOML設定ファイル解析
- **golang.org/x/image**: 追加画像形式サポート

### アーキテクチャ

- **設定管理**: TOML形式の設定ファイル
- **画像処理**: Go標準ライブラリ + 拡張画像形式
- **GUI**: Fyneフレームワーク
- **プラットフォーム固有処理**: Windows API（コンソール非表示）

## デバッグ

問題が発生した場合は、`--verbose` オプションを使用して詳細情報を確認できます：

```bash
go_ImagePreviewer --verbose problematic_image.jpg
```

デバッグ情報には以下が含まれます：
- 実行ファイルのディレクトリパス
- 設定ファイルの検索パス
- 現在の作業ディレクトリ
- 読み込まれた設定値

## トラブルシューティング

### 画像が表示されない
1. 対応形式の画像ファイルか確認
2. ファイルパスが正しいか確認
3. `--verbose` オプションでエラー詳細を確認

### 設定ファイルが読み込まれない
1. `--verbose` オプションで設定ファイルのパスを確認
2. TOML形式が正しいか確認
3. ファイルの読み取り権限を確認

### Windowsでコンソールが表示される
1. `-H=windowsgui` オプション付きでビルド
2. `build.bat` または `make windows` を使用

## 開発

### 依存関係のインストール
```bash
go mod tidy
```

### テスト実行
```bash
go run main.go sample.jpg
```

### ビルド（開発用）
```bash
go build -o go_ImagePreviewer
```

## ライセンス

MIT License

## 作者

kznagamori (https://github.com/kznagamori)

## 更新履歴

- **v1.0.0**: 初回リリース
  - 基本的な画像表示機能
  - 設定ファイル対応
  - 最前面表示モード
  - Fyneフレームワーク採用

## 貢献

プルリクエストやIssueは歓迎します。

1. このリポジトリをフォーク
2. フィーチャーブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add some amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成