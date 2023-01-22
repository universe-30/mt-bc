module github.com/universe-30/mt-bc

go 1.19

replace github.com/universe-30/mt-trie => /Users/jaimin/Documents/work/meson/gitlab/chain/mt-trie

require (
	github.com/universe-30/mt-trie v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.4.0
)

require golang.org/x/sys v0.3.0 // indirect
