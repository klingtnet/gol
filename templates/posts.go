package templates

var postsTemplate = `<!DOCTYPE html>
<html lang=en>
	<head>
		<title>{{ .title }}</title>

		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/materialize/0.95.3/css/materialize.min.css">

		<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/highlight.js/8.4/styles/tomorrow.min.css">

		<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no"/>
		<link rel="stylesheet" href="{{ "main.css" | assetUrl }}" />
	</head>

	<body>
		<div class="container">
			<div id="edit-button" class="fixed-action-btn">
				<a href="/posts/new" class="btn-floating btn-large waves-effect waves-light blue tooltipped" data-tooltip="Write a new post"><i class="mdi-content-add"></i></a>
			</div>

			{{ range $post := .posts }}
			<article id="post-{{ $post.Id }}" class="post">
				<div class="post-actions">
					<a href="/posts/{{ $post.Id }}/edit" class="btn-floating waves-effect waves-light blue tooltipped" data-tooltip="Edit post"><i class="mdi-editor-mode-edit"></i></a>
					<a href="/posts/{{ $post.Id }}" data-method="DELETE" class="btn-floating waves-effect waves-light red tooltipped" data-tooltip="Delete post"><i class="mdi-action-delete"></i></a>
				</div>
				<h1><a href="/posts/{{ $post.Id }}">{{ $post.Title }}</a></h1>
				<h5>Posted on <i>{{ $post.Created | formatTime }}</i></h5>

				<div class="post-content flow-text">
					{{ $post.Content | markdown }}
				</div>
			</article>
			<hr />
			{{ end }}
		</div>

		<script type="text/javascript" src="https://code.jquery.com/jquery-2.1.1.min.js"></script>
		<script src="https://cdn.rawgit.com/heyLu/materialize.css/master/dist/js/materialize.min.js"></script>

		<script src="//cdnjs.cloudflare.com/ajax/libs/highlight.js/8.4/highlight.min.js"></script>
		<script>hljs.initHighlightingOnLoad();</script>

		<script src="{{ "main.js" | assetUrl }}"></script>
	</body>
</html>`
