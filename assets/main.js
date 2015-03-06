"use strict";

(function() {
    function log() {
        console.log.apply(console, arguments);
    }

    function displayMessage(msg) {
        console.log(msg);
    }

    function displayError(msg) {
        console.error(msg);
        alert(msg);
    }

    // listen on shift+tab
    function tabOverride() {
        var textarea = document.querySelector("textarea");
        if (textarea == null) { return }
        textarea.onkeydown = function (e) {
            if (e.shiftKey && e.keyCode === 9) {
                e.preventDefault();
                var text = textarea.value;
                var pos = textarea.selectionStart;
                textarea.value = text.substr(0, pos) + '    ' + text.substr(pos);
                // select nothing
                console.log(textarea.selectionStart, textarea.selectionEnd, textarea);
                textarea.selectionStart = pos + 4;
                textarea.selectionEnd = pos + 4;
            }
        }
    }

    // support DELETEing resources via data-method="DELETE"
    function supportDeleteLinks() {
        var deleteLinks = document.querySelectorAll("a[data-method=DELETE]");
        for (var i = 0; i < deleteLinks.length; i++) {
            (function(deleteLink) {
                deleteLink.addEventListener("click", function(ev) {
                    ev.preventDefault();

                    var xhr = new XMLHttpRequest();
                    xhr.open("DELETE", deleteLink.href);
                    xhr.onload = function(ev) {
                        if (xhr.status == 200) {
                            location.reload();
                        } else {
                            displayError("could not delete post");
                        }
                    }
                    xhr.send();
                });
            })(deleteLinks[i]);
        }
    }

    // render the markdown preview
    function renderPreview() {
        var form = document.getElementById("edit-post");
        var titleInput = document.getElementById("edit-title");
        var contentInput = document.getElementById("edit-content");
        var previewSelect = document.getElementById("preview-select");
        var preview = document.getElementById("preview-tab");
        if (previewSelect == null) { return; }

        previewSelect.addEventListener("click", function(ev) {
            var xhr = new XMLHttpRequest();
            xhr.open('POST', '/posts/preview');
            var post = {
                "title": titleInput.value,
                "content": contentInput.value,
                "created": form.dataset.postCreated || new Date().toISOString()
            };
            xhr.send(JSON.stringify(post));
            xhr.onload = function(ev) {
                if (xhr.status == 200) {
                    preview.innerHTML = xhr.responseText;
                } else {
                    console.error(xhr);
                }
            };
        });
    }

    // open editor on doubleclick
    function editorOnDoubleClick() {
        var posts = document.querySelectorAll(".post");
        for (var i = 0; i < posts.length; i++) {
            (function(post) {
                post.addEventListener('dblclick', function(ev) {
                    location.href = '/posts/' + post.dataset.id + '/edit';
                });
            })(posts[i]);
        }
    }

    function requestFullscreen(el) {
        if (el.requestFullscreen) {
            el.requestFullscreen();
        } else if (el.mozRequestFullScreen) {
            el.mozRequestFullScreen()
        } else if (el.webkitRequestFullScreen) {
            el.webkitRequestFullScreen();
        }
    }

    // fullscreen mode
    function setupFullscreenMode() {
        var editContent = document.getElementById("edit-content");
        if (editContent == null) { return; }

        var fullscreenToggle = document.getElementById("fullscreen-toggle");
        fullscreenToggle.addEventListener("click", function(ev) {
            requestFullscreen(editContent.parentElement);
            editContent.focus();
        });
    }

    function savePost(success, error) {
        var form = document.getElementById("edit-post");
        var editTitle = document.getElementById("edit-title");
        var editContent = document.getElementById("edit-content");

        var isNew = !form.dataset.postId;
        var post = {
            "title": editTitle.value,
            "content": editContent.value
        };

        var xhr = new XMLHttpRequest();
        xhr.open('POST', isNew ? '/posts' : '/posts/' + form.dataset.postId);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.responseType = 'json'
        xhr.onload = function(ev) {
            if (xhr.status >= 200 && xhr.status < 300) {
                if (isNew) {
                    form.dataset.postId = xhr.response.id;
                }
                success(xhr, isNew);
            } else {
                error(xhr);
            }
        };
        xhr.send(JSON.stringify(post));
    }

    // save post shortcut (without stopping to write it)
    function savePostShortcut() {
        var editContent = document.getElementById("edit-content");
        if (editContent == null) { return; }

        editContent.addEventListener("keydown", function(ev) {
            if (ev.ctrlKey && ev.keyCode == 83) { // Ctrl-S
                ev.preventDefault();
                savePost(function(_, isNew) {
                    displayMessage(isNew ? "post created" : "post saved");
                }, function(xhr) { console.error(xhr.status, xhr.statusText); });
            }
        });
    }

    tabOverride();
    supportDeleteLinks();
    renderPreview();
    editorOnDoubleClick();
    setupFullscreenMode();
    savePostShortcut();
})();
