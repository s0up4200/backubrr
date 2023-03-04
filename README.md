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

To use Backubrr, you'll need to create a configuration file named config.yaml in the same directory as the backubrr executable. The configuration file should contain the following keys:

```yaml
output_dir: /path/to/backup/directory
source_dirs:
  - /path/to/source/directory1
  - /path/to/source/directory2
```

The output_dir key specifies the destination directory where the backup archives will be stored. The source_dirs key is a list of source directories that will be backed up. You can add as many source directories as you need to this list.

To create a backup, simply run the backubrr executable:

```bash
./backubrr
```

Backubrr will read the config.yaml file and create compressed tar archives of each source directory. The archives will be stored in the destination directory specified in the configuration file.
