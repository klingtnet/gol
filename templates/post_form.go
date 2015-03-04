package templates

var postFormTemplate = `<!DOCTYPE html>
<html lang=en>
	<head>
		<title>{{ .title }}</title>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/materialize/0.95.3/css/materialize.min.css">
		<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no"/>
		<link rel="stylesheet" href="{{ "main.css" | assetUrl }}" />
	</head>

	<body>
		<div class="container">
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
					Submit
				</button>
			</form>
		</div>

		<script type="text/javascript" src="https://code.jquery.com/jquery-2.1.1.min.js"></script>
		<script src="https://cdn.rawgit.com/heyLu/materialize.css/master/dist/js/materialize.min.js"></script>

		<script src="{{ "main.js" | assetUrl }}"></script>
	</body>
</html>`
