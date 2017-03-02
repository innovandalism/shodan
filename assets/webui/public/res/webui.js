riot.tag2('webui', '<navbar app="{app}" items="{menu}" user="{user}" logged-in="{loggedIn}"></navbar> <welcome app="{app}" if="{route==⁗start⁗}"></welcome>', '', '', function(opts) {
        var that = this;

        this.app = {
            name: "Loading..."
        };

        this.user = {

        }

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
        updateProfile();

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

        function updateProfile() {
            var headers = new Headers({
                "Authorization": "Bearer " + localStorage.jwt
            })
            fetch("/user/profile/", {
                headers: headers,
                method: "GET"
            }).then(function(res) {
                return res.json()
            }).then(function(user) {
                if (user.err != "" || user.status != 200) {
                    throw user.err
                }
                that.user = user.user
            })
        }
});