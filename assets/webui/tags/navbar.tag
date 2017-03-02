<navbar>
    <nav class="uk-navbar-container sd-navbar uk-light uk-sticky" uk-navbar="mode: click">
        <div class="uk-navbar-left">
            <ul class="uk-navbar-nav">
                <a class="uk-navbar-item uk-logo" href="#">{opts.app.name}</a>
            </ul>
        </div>
        <div class="uk-navbar-right">
            <ul class="uk-navbar-nav" if={!opts.loggedIn}>
                <li><a href="/oauth/authenticate/">Sign in with Discord</a></li>
            </ul>
            <ul class="uk-navbar-nav">
                <li if={opts.user}>
                    <a>
                        <img class="avatar" src="https://cdn.discordapp.com/avatars/{opts.user.id}/{opts.user.Avatar}.{(opts.user.Avatar.indexOf("a_") == 0 ? "gif" : "jpg")}">
                        &nbsp;{opts.user.username}
                        <div class="uk-navbar-dropdown">
                            <ul class="uk-nav uk-navbar-dropdown-nav">
                                <li><a onclick={navigate("logout")}>Sign out</a></li>
                            </ul>
                        </div>
                    </a>
                </li>
                <li if={opts.guild}>
                    <a>
                        <img class="avatar" src="https://cdn.discordapp.com/avatars/{opts.user.id}/{opts.user.Avatar}.{(opts.user.Avatar.indexOf("a_") == 0 ? "gif" : "jpg")}">
                        &nbsp;{opts.user.username}
                        <div class="uk-navbar-dropdown">
                            <ul class="uk-nav uk-navbar-dropdown-nav">
                                <li><a onclick={navigate("logout")}>Sign out</a></li>
                            </ul>
                        </div>
                    </a>
                </li>
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
        .avatar {
            max-width: 32px;
            max-height: 32px;
            border-radius: 16px;
            margin-right: 8px;
        }

        .sd-navbar:not(.uk-navbar-transparent) {
            background: linear-gradient(to left, #28a5f5, #1e87f0);
        }
    </style>
</navbar>