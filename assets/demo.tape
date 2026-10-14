# Regenerate with: vhs assets/demo.tape (from the repository root; requires gopls in PATH)

Output assets/demo.gif

Set Shell "bash"
Set FontSize 22
Set Width 1440
Set Height 760
Set Padding 24
Set Theme "Catppuccin Mocha"
Set TypingSpeed 45ms

Hide
Type `export D=/tmp/goalias-demo && rm -rf "$D" && mkdir -p "$D" && go build -o "$D/bin/goalias" ./cmd/goalias && cp -R assets/demo/. "$D" && cd "$D" && export PATH="$D/bin:$PATH" PS1='\[\e[1;35m\]$ \[\e[0m\]' && clear`
Enter
Sleep 4s
Show

Type "goalias list -p example.com/shop/pricing"
Sleep 400ms
Enter
Sleep 3s

Type "goalias set -p example.com/shop/pricing -a pricing"
Sleep 400ms
Enter
Sleep 5s

Type "goalias list -p example.com/shop/pricing"
Sleep 400ms
Enter
Sleep 3s

Type "clear"
Enter
Type "cat worker/invoice.go"
Sleep 400ms
Enter
Sleep 4s
