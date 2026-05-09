package htmlconv

import "testing"

func TestConvert(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", "", ""},
		{"plain text", "hello", "hello"},
		{"h1", "<h1>Title</h1>", "# Title"},
		{"h2", "<h2>Section</h2>", "## Section"},
		{"h3", "<h3>Sub</h3>", "### Sub"},
		{"h4", "<h4>Sub</h4>", "#### Sub"},
		{"h5", "<h5>Sub</h5>", "##### Sub"},
		{"h6", "<h6>Sub</h6>", "###### Sub"},
		{"p", "<p>Hello</p>", "Hello"},
		{"br", "a<br>b", "a\nb"},
		{"strong", "<strong>bold</strong>", "**bold**"},
		{"b", "<b>bold</b>", "**bold**"},
		{"em", "<em>italic</em>", "*italic*"},
		{"i", "<i>italic</i>", "*italic*"},
		{"code inline", "<code>x</code>", "`x`"},
		{"anchor", `<a href="https://example.com">link</a>`, "[link](https://example.com)"},
		{"pre", "<pre>code here</pre>", "```\ncode here\n```"},
		{"hr", "<hr>", "---"},
		{"blockquote", "<blockquote>quoted</blockquote>", "> quoted"},
		{"ul", "<ul><li>A</li><li>B</li></ul>", "- A\n- B"},
		{"ol", "<ol><li>First</li><li>Second</li></ol>", "1. First\n2. Second"},
		{
			"table simple",
			"<table><tr><th>Name</th><th>Val</th></tr><tr><td>foo</td><td>bar</td></tr></table>",
			"| Name | Val |\n|---|---|\n| foo | bar |",
		},
		{"img", `<img src="img.png" alt="desc">`, "![desc](img.png)"},
		{
			"code macro no lang",
			`<ac:structured-macro ac:name="code"><ac:plain-text-body>x</ac:plain-text-body></ac:structured-macro>`,
			"```\nx\n```",
		},
		{
			"code macro with lang",
			`<ac:structured-macro ac:name="code"><ac:parameter ac:name="language">go</ac:parameter><ac:plain-text-body>package main</ac:plain-text-body></ac:structured-macro>`,
			"```go\npackage main\n```",
		},
		{
			"info macro",
			`<ac:structured-macro ac:name="info"><ac:rich-text-body>This is info</ac:rich-text-body></ac:structured-macro>`,
			"> **Info:** This is info",
		},
		{
			"note macro",
			`<ac:structured-macro ac:name="note"><ac:rich-text-body>This is note</ac:rich-text-body></ac:structured-macro>`,
			"> **Note:** This is note",
		},
		{
			"warning macro",
			`<ac:structured-macro ac:name="warning"><ac:rich-text-body>This is warning</ac:rich-text-body></ac:structured-macro>`,
			"> **Warning:** This is warning",
		},
		{
			"tip macro",
			`<ac:structured-macro ac:name="tip"><ac:rich-text-body>This is tip</ac:rich-text-body></ac:structured-macro>`,
			"> **Tip:** This is tip",
		},
		{
			"ac:link with page title",
			`<ac:link><ri:page ri:content-title="My Page"/></ac:link>`,
			"*My Page*",
		},
		{"emoticon skipped, following text preserved", `text<ac:emoticon ac:name="smile"/>more`, "textmore"},
		{
			"h2 with emoticon and nbsp",
			"<h2><ac:emoticon ac:name=\"blue-star\" ac:emoji-shortname=\":thinking:\" /> 問題の提示</h2>",
			"## 問題の提示",
		},
		{
			"ac:placeholder text preserved",
			"<p><ac:placeholder>プレースホルダーテキスト</ac:placeholder></p>",
			"プレースホルダーテキスト",
		},
		{
			"status macro",
			`<ac:structured-macro ac:name="status"><ac:parameter ac:name="title">進行中</ac:parameter><ac:parameter ac:name="colour">Yellow</ac:parameter></ac:structured-macro>`,
			"進行中",
		},
		{"p with inline code", "<p>Use <code>go test</code> to run</p>", "Use `go test` to run"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := Convert(tt.input)
			if got != tt.want {
				t.Errorf("Convert(%q)\n  got:  %q\n  want: %q", tt.input, got, tt.want)
			}
		})
	}
}
