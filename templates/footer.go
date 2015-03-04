package templates

var footerTemplate = `
		</div>

		<script type="text/javascript" src="https://code.jquery.com/jquery-2.1.1.min.js"></script>
		<script src="https://cdn.rawgit.com/heyLu/materialize.css/master/dist/js/materialize.min.js"></script>

		<script src="//cdnjs.cloudflare.com/ajax/libs/highlight.js/8.4/highlight.min.js"></script>
		<script>hljs.initHighlightingOnLoad();</script>

        <script type="text/javascript"
          src="//cdn.mathjax.org/mathjax/latest/MathJax.js?config=TeX-AMS-MML_HTMLorMML">
        </script>

		<script src="{{ "main.js" | assetUrl }}"></script>
	</body>
</html>`
