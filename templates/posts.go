package templates

var postsTemplate = `
{{ template "header" .}}

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

{{ template "footer" . }}`
