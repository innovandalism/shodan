<token-status>
<div class={error: hasError}>{ status }</div>
<script>
    this.status = "pending..."
    this.hasError = false
    var that = this
    var params = parseHash()

    fetch("/oauth/testtoken/", {
        method: "POST",
        body: JSON.stringify({
            token: params.access_token
        })
    }).then(function(res) {
        if (!res.ok) {
            throw Error(response.statusText)
        }
        localStorage.setItem("token", params.access_token)
        setTimeout(function() {
            document.location = "/"
        }, 3000);
        that.status = "token verified, redirecting you to the application..."
        that.update()
    }).catch(function(err) {
        that.hasError = true;
        that.status = "an error occured"
        that.update();
    })

    function parseHash() {
        var params = {}
        var hash = window.location.hash
        hash = hash.substr(1)
        var arg = hash.split("&")
        arg.forEach(function(a) {
            var kvp = a.split("=")
            params[kvp[0]] = kvp[1]
        })
        return params
    }
</script>
</token-status>