
<h3 align="center">Gator - Database-driven aggregation tool</h3>
</div>

<!-- ABOUT THE PROJECT -->
## About The Project

Gator is a CLI-driven blog aggregation tool that seeks to create a streamlined process for collecting and cataloging posts from registered sources to be pulled with the press of a button. A barebones login system as been implemented to allow for multiple users to save and follow sources on the same instance of the aggregation database.

Skills shown in this project:
* Creation and automated updated/queries of Postgres databases
* Database interraction and integration using the GOLANG toolchain
* Long-running service worker that reaches out over the internet to fetch data from remote locations
* Designing CLI interface from scratch
<p align="right">(<a href="#readme-top">back to top</a>)</p>



### Built With

* [GO](https://go.dev/)
* [SQLc](https://github.com/sqlc-dev/sqlc)
* [Goose Data Migration Tool](https://github.com/pressly/goose)
* [PostgreSQL](https://www.postgresql.org/)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->

### Prerequisites

This tool needs the most recent version of PostgreSQL and GO to function

* Postgres installation
  Mac with brew:
    ```sh
    brew install postgresql@15
    ```
  Linux/WSL (Debian):
    ```sh
    sudo apt update
    sudo apt install postgresql postgresql-contrib
    ```
* Go installation
  Found on [GOLANG Website](https://go.dev/dl/)

### Installation

**1. Create local Postgres database in the background**
  Mac:
  ```sh
  brew services start postgresql
  ```
  Linux:
  ```sh
  sudo service postgresql start
  ```
**2. Create new database inside the Postgres shell, I called mine ```gator```**
  ```sh
  CREATE DATABASE gator;
  ```
  Connect to the database
  ```sh
  \c gator
  ```
  (Linux only, setting new user password):
  ```sh
  ALTER USER postgres PASSWORD 'postgres';
  ```
**3. Install the CLI using GO**
  ```sh
  go install github.com/vyynl/gator@latest
  ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE -->

### Config Setup

1. Registering database URL
  ***Get your connection Postgres connection string to then register with the CLI:***
  ```sh
  protocol://username:password@host:port/database
  ```
  Examples:
    * macOS (no password, your username): `postgres://vyynl:@localhost:5432/gator`
    * Linux (password from above, postgres user): `postgres://postgres:postgres@localhost:5432/gator`

  Test your connection string by running the below with your connection string to check (change `psql` to `sudo -u postgres psql` for Linux):
  ```sh
  psql "postgres://vyynl:@localhost:5432/gator"
  ```

  
  *Run the below command to set your database url in the CLI:*
  ```sh
  gator dburl <your connection string here w/ no quotes>
  ```
2. Registering your first user using `register`
  ```sh
  gator register <New_User_Name>
  ```

### Command List

  * `agg` - Starts scraping all registered feeds with set time between requests (format examples: 3s - 3 seconds, 2m - 2 minutes, 1h - 1 hour)
    ```sh
    gator agg <time_between_reqs>
    ```
  * `addfeed` - Adds a URL to list of feeds to scrape (also follows with active user)
    ```sh
    gator addfeed <feed_url>
    ```
  * `browse` - List most recent feeds from followed URLs up to `limit` (default: 2)
    ```sh
    gator browse <limit>
    ```
  * `feeds` - Lists all feeds flagged to be scraped
    ```sh
    gator feeds
    ```
  * `follow` - Follows feed url with active user
    ```sh
    gator follow <feed_url>
    ```
  * `following` - Lists all feeds followed by active user
    ```sh
    gator following
    ```
  * `login` - Sets active user (must be registered)
    ```sh
    gator login <User_Name>
    ```
  * `users` - Prints list of registered uses w/ active flagged
    ```sh
    gator users
    ```
  * `unfollow` - Unfollows feed url with active user
    ```sh
    gator unfollow <feed_url>
    ```
  * `reset` - Completely resets and clears database
    ```sh
    gator reset <arguement1>
    ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/vyynl/gator.svg?style=for-the-badge
[contributors-url]: https://github.com/vyynl/gator/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/vyynl/gator.svg?style=for-the-badge
[forks-url]: https://github.com/vyynl/gator/network/members
[stars-shield]: https://img.shields.io/github/stars/vyynl/gator.svg?style=for-the-badge
[stars-url]: https://github.com/vyynl/gator/stargazers
