"use strict";

(function() {
    function log() {
        console.log.apply(console, arguments);
    }

    function displayError(msg) {
        console.error(msg);
        alert(msg);
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

    supportDeleteLinks();
})();
