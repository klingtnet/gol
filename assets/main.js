"use strict";

(function() {
    function log() {
        console.log.apply(console, arguments);
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
                    xhr.onreadystatechange = log;
                    xhr.onerror = log;
                    xhr.send();
                });
            })(deleteLinks[i]);
        }
    }

    supportDeleteLinks();
})();
