## go-mod-cleanup
# Version 3.

# Setup
`mkdir ~/.go-mod-cleanup`<br>
`cd ~/.go-mod-cleanup`<br>
`git clone https://github.com/andym1125/go-mod-cleanup gitfiles`<br>
`cp ./gitfiles/cmd/go-mod-cleanup modclean`<br>
`rm -rf gitfiles`<br>

# Use
Now, ensuring you're in your project directory:
`go mod graph > <Your file name>`<br>
`~/.go-mod-cleanup/modclean <Your file name>`<br>
Optionally remove HTML files: `rm -rf go_mod_graphs`
