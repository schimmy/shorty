# Pretty Self-Hosted URL Shortener

What, you need more information than the title?

I looked around and found a bunch of toy examples of backend URL shorteners in Go, or PHP/MySQL solutions which seemed complicated and less-than-fully-featured. I wanted to make something that spins up with super-simple configuration and has a pluggable backend. Please contribute!

If you don't like what you see, you can also check out the other alternatives:
- [YOURLS](yourls.org)
- [Google Apps](TODO)

## Goals For The Project
- Optimized for simple deployment
- Designed for internal network use
- Clear ownership of URLs
- Expiry of URLs unless explicitly set (TODO: flash warning if URL is about to expire)
- Tags to allow grouping / organization

## Deployment
After a few simple steps, I promise you will be up and running and won't need to touch it again.

### Pick a backend database.
Current options are:
- PostgreSQL (AWS Redshift also supported)
- Redis

*If you have a favorite backend you don't see here, please help me by creating a pull request!*
I recommend a persistent storage - your users are not going to be happy if they lose these URLs!

### Edit the database creation script
Find the script corresponding to your database and open it in a text editor. If there are TODOs, follow the instructions.
You may modify this as you wish (for instance, perhaps you may want to use a different Postgres `schema` than '`shortener`'), however this means that you'll likely have to be more careful setting other environment variables later.
- [PostgreSQL]()
- Redis - does not require a creation script

### Run the database creation script

#### PostgreSQL
Enter your admin-user's password when prompted.
```bash
psql -p 5432 -U YOUR_ADMIN_USER -d YOUR_DATABASE -f ./scripts/postgres.sql
```
You also might want to change the default password for the `shortener` user:
```bash
psql -p 5432 -U YOUR_ADMIN_USER -d YOUR_DATABASE -c "ALTER ROLE shortener WITH PASSWORD 'TODO'"
```

### Pick a deployment method

#### Docker
The latest published image at: [schimmy/pretty-self-hosted-url-shortener](dockerhub TODO). Run this Docker container with the relevant environment variables for your database.

#### Standalone Server (Windows, Linux)
We also release binaries for Windows and Linux, available at [](TODO).
After setting the relevant environment variables, simply execute the binary.
Likely you'll want to ensure that it restarts on reboot, etc, and I've included an [init script](TODO) to help you out.

Linux example:
```bash
PG_HOST=my-postgres.colinschimmelfing.com PG_PASS=super_secret ./shortener
```

## That's it, you're done!

If you'd like other methods of deploying, please create a pull request with support and documentation.
