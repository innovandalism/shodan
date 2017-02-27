<navbar>
    <nav class="uk-navbar-container sd-navbar uk-light uk-sticky" uk-navbar>
        <div class="uk-navbar-left">
            <ul class="uk-navbar-nav">
                <a class="uk-navbar-item uk-logo" href="#">{opts.app.name}</a>
            </ul>
        </div>
        <div class="uk-navbar-right">
            <ul class="uk-navbar-nav" if={!opts.loggedIn}>
                <li><a href="/oauth/authenticate/">Sign in with Discord</a></li>
            </ul>
            <ul class="uk-navbar-nav" if={opts.loggedIn}>
                <li each={opts.items}><a onclick={navigate(target)}>{title}</a></li>
            </ul>
        </div>
    </nav>

    <script>
        this.navigate = function(r) {
            return function() {
                route(r)
            }
        }
    </script>

    <style>
        .sd-navbar:not(.uk-navbar-transparent) {
            background: linear-gradient(to left, #28a5f5, #1e87f0);
        }
    </style>
</navbar>