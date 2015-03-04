package templates

var headerTemplate = `<!DOCTYPE html>
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
`
