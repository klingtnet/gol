package templates

var postFormTemplate = `
{{ template "header" . }}

			<h1>{{ .title }}</h1>

			<form method="POST" action="/posts{{ if .post }}/{{ .post.Id }}{{ end }}">
				<div class="input-field">
					<input class="markdown-input" name="title" type="text" value="{{ .post.Title }}"></input>
					<label for="title">Titlemania</label>
				</div>
				<div class="input-field">
					<textarea class="materialize-textarea markdown-input" name="content" rows="80" cols="100">{{ .post.Content }}</textarea>
					<label for="content">Your thoughts.</label>
				</div>


				<button class="btn waves-effect waves-light" type="submit" name="action">
					<i class="mdi-action-done left"></i>
					Submit
				</button>
			</form>

{{template "footer" . }}`
