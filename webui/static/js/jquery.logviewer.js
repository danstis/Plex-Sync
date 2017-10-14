(function(jQuery) {
    var logViewer = (function(options) {
        (function cleanLVtags(html) {
            if (typeof html === "string") {
                return html
                    .replace(/&/g, "&amp;")
                    .replace(/</g, "&lt;")
                    .replace(/>/g, "&gt;");
            } else {
                return html;
            }
        },
            function setMyTimeOut() {
                if (logViewer.timeoutID > 0) {
                    window.clearTimeout(logViewer.timeoutID);
                }
                logViewer.timeoutID = window.setTimeout(doLVHead, logViewer.options.refreshtimeout);
            },
            function checkLVLength(newLength) {
                if (logViewer.curLength !== newLength) {
                    if (logViewer.curLength > newLength) {
                        logViewer.curLength = 0;
                        jQuery("#" + logViewer.options.targetObjectID).append("\nReseting ... \n");
                    }
                    var getBytes = logViewer.curLength;
                    var readBytes = parseInt(logViewer.options.readBytes);

                    if (logViewer.curLength === 0 && newLength > readBytes) {
                        getBytes = newLength - readBytes;
                    } else if (logViewer.curLength > 0) {
                        getBytes--;
                    }

                    jQuery.ajax({
                        type: "GET",
                        url: logViewer.options.logUrl,
                        cache: false,
                        success: function(data) {
                            data = logViewer.options.callback.call(this, data);
                            jQuery("#" + logViewer.options.targetObjectID).append(cleanLVtags(data));
                            var objDiv = document.getElementById(logViewer.options.targetObjectID);
                            objDiv.scrollTop = objDiv.scrollHeight;
                        },
                        beforeSend: function(http) {
                            http.setRequestHeader("Range", "bytes=" + getBytes + "-");
                        }
                    });
                }
                logViewer.curLength = newLength;
                setMyTimeOut();
            },
            function doLVHead(id) {
                jQuery.ajax({
                    type: "HEAD",
                    url: logViewer.options.logUrl,
                    cache: false,
                    complete: function(xhr, textStatus) {
                        if (textStatus === "success" || textStatus === "nocontent") {
                            var newLength = parseInt(xhr.getResponseHeader("Content-Length"));
                            checkLVLength(newLength);
                        }
                    }
                });
            });

        return {
            init: function(options) {
                if (!options) var options = {};
                if (options.logUrl === undefined) {
                    alert("Log URL missing");
                    return false;
                }
                if (options.refreshtimeout === undefined) {
                    options.refreshtimeout = "2000";
                }
                if (options.readBytes === undefined) {
                    options.readBytes = 10000;
                }
                if (options.callback === undefined) {
                    options.callback = function(data) {
                        return data;
                    };
                }

                options.targetObjectID = jQuery(this).attr("id");

                logViewer.options = options;
                logViewer.curLength = 0;

                doLVHead();
            }
        };
    })();
    jQuery.fn.extend({
        logViewer: logViewer.init
    });
})(jQuery);
