<guildselect>
    <div class="uk-section sd-section-guildselect">
        <div class="uk-container uk-container-small">
            <h1 class="uk-heading-line uk-light uk-text-center"><span>Select Your Guild</span></h1>
            <div class="uk-flex uk-flex-center uk-flex-around uk-flex-wrap">
                <a class="guild" each={opts.guilds} title={name} uk-tooltip="pos:bottom;offset:-10">
                    <img src="https://cdn.discordapp.com/icons/{id}/{icon}.jpg">
                </a>
            </div>
        </div>
    </div>

    <style>
        .sd-section-guildselect {
            background-color: #ff9900;
        }

        .guild {
            padding: 12px;
            width: 64px;
            height: 64px;
        }

        .guild img {
            border-radius: 32px;
        }
    </style>
</guildselect>