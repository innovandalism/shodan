<webui>
    <navbar app="{app}" items={menu} logged-in={loggedIn}></navbar>
    <welcome app="{app}" if={route=="start"}></welcome>

    <script>
        var that = this;

        this.app = {
            name: "Loading..."
        };

        this.route = "";

        this.loggedIn = false;

        this.menu = [
            {title: "Sign Out", target: "logout"}
        ];

        route("start", function() {
            that.route = "start"
        });

        route("logout", function() {
            localStorage.clear();
            updateLoginStatus()
            route("start")
        });

        route("dashboard", function() {
            this.route = "dashboard"
        })

        updateManifest();
        updateLoginStatus();

        if (this.route == "") {
            route("start")
        }

        function updateLoginStatus() {
            if (route.query().jwt) {
                localStorage.setItem("jwt", route.query().jwt)
            }
            if (localStorage && localStorage.jwt) {
                that.loggedIn = true;
            } else {
                that.loggedIn = false;
            }
            that.update()
        }

        function updateManifest() {
            fetch("/app.json").then(function(res) {
                return res.json()
            }).then(function(app) {
                that.app = app;
                that.update()
            })
        }
    </script>
</webui>