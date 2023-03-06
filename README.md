# Backubrr

Backubrr is a command-line tool for backing up directories specified in a configuration file. It creates compressed tar archives of the specified directories and stores them in a destination directory.

## Installation

Clone the repository using git and compile the binary yourself:

```bash
git clone https://github.com/s0up4200/backubrr.git
cd backubrr
go build
```

## Usage

By default, Backubrr looks for a configuration file named config.yaml in the same directory as the backubrr executable. Alternatively, you can specify a custom configuration file using the --config flag. The configuration file should contain the following keys:

```yaml
output_dir: /path/to/backup/directory
source_dirs:
   - /path/to/source/directory1
   - /path/to/source/directory2
#interval: 24 #hours
#retention_days: 7 #days
#discord: https://discord.com/api/webhooks/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX

```

The output_dir key specifies the destination directory where the backup archives will be stored. The source_dirs key is a list of source directories that will be backed up. You can add as many source directories as you need to this list.

To create a backup, simply run the backubrr executable with or without the `--config` variable:

```bash
./backubrr --config /path/to/config.yaml
```

By default, Backubrr runs once and exits. If you want to run it on a schedule, add the interval key and set its value to the number of hours between backups. For example, `interval: 24` will run Backubrr once a day. You can change the interval value to any number of hours.

In addition to the `interval` key, Backubrr also provides a `retention_days` key to help manage the amount of space used by backups. Setting the retention value specifies how many days of backups you want to keep. When a new backup is created, Backubrr checks the age of the backups in the destination directory and removes those that are older than the retention period specified.


### Discord notifications
To enable Discord notifications, add a `discord` key to the configuration file and set its value to the URL of the Discord webhook that you want to use. Backubrr will send a notification to the specified channel via the webhook when a backup is created. The notification includes the name of the backup archive and the source directory that was backed up.