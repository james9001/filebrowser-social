# filebrowser-social

filebrowser-social is a downstream fork which intends to stay synced to the excellent upstream project, filebrowser, whilst also adding  social features. These features are those that a typical self-hoster might find useful in situations where filebrowser instances are shared among small groups of friends or family, such as: commenting, reactions, and notifications - as well as other miscellaneous relevant tweaks.

## Features

filebrowser-social adds the following social features:

- Users may now leave Comments, which can be submitted against any files which have a Preview mode - which is almost all of them.
- Users may now also submit Reactions against any files which may have Comments. The set of Reactions available is determined by the instance admin and/or users.
- Whenever Comments and Reactions are submitted by users, Notification events are generated, which will be shown to other users until acknowledged.
- A new way of viewing files, "Large Mosaic Gallery" mode. Useful for photo galleries.
- Bandwidth control, implemented with Traffic Control (tc). Useful for when you want to limit the amount of bandwidth friends and/or family may use at any given time when accessing your filebrowser-social instances.

## Quick Start

From the root repo directory:

- Run script `build.sh` to build container image
- Run `touch localdev/filebrowser.db`
- Run `docker compose up -d` and wait for the container to come up
- Open a new browser window and navigate to `http://localhost:80`
- Login with the username `admin` and password `admin`
- Create a new user with the normal filebrowser user administration system, for example, `user2`
- Open the image `1920px-MarsSunset.jpg`. You will see a Comments UI, and above that, a Reactions UI. The default localdev configuration has one reaction type configured.
- If you log in as your newly created user (i.e. `user2` as created earlier), you will see notifications in the top right of the screen leading you towards the file context where you previously left comments and/or reactions, as the admin user

## Troubleshooting

Please feel free to open an Issue in this repository if you run into any problems!
