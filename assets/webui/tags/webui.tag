<webui>
    <navbar app="{app}" user={user} guild={guild} logged-in={loggedIn}></navbar>
    <welcome app="{app}" if={route=="start"}></welcome>
    <guildselect guilds={guilds} if={loggedIn}></guildselect>

    <script>
        var that = this;

        this.app = {
            name: "Loading..."
        };

        this.alerts = [];

        this.user = null;
        this.guilds = []
        this.guild = null;

        this.route = "";

        this.loggedIn = false;

        route("start", function() {
            that.route = "start"
        });

        route("logout", function() {
            localStorage.clear();
            updateLoginStatus()
            updateProfile()
            route("start")
        });

        route("dashboard", function() {
            this.route = "dashboard"
        })

        updateManifest();
        updateLoginStatus();
        updateProfile();

        if (this.route == "") {
            route("start")
        }

        function setGuildContext() {

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

        function updateProfile() {
            if (!localStorage.jwt) {
                that.user = null;
                return
            }
            var headers = new Headers({
                "Authorization": "Bearer " + localStorage.jwt
            });

            fetch("/user/profile/", {
                headers: headers,
                method: "GET"
            }).then(function(res) {
                return res.json()
            }).then(function(res) {
                if (res.err != "" || res.status != 200) {
                    throw res.err;
                }
                that.user = res.data.user
                that.guilds = res.data.guilds
                that.update()
            }).catch(function() {
                route("logout")
            })
        }
    </script>
</webui>