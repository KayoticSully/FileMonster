FileMonster
===========
Organizes massive amounts of files. I wrote this program for my own needs, maybe it will be useful to someone else too.

Long story short, I ended up with a directory of about 500,000 files from a HDD recovery dump.  Ever OS I could think of choked on when trying to view the directory. I tried going through the files via terminal, and that just proved frustrating and time consuming.  I needed a way to organize the files so I could easily see what I had in that directory.

This program can rip through directorys at about 100,000 files a minute and reorganizes files into folders by file extension.  The idea is, run this program then go and grab the folders you care about such as doc, jpg, png, zip, etc.

It uses a gorutine to walk the directory and several workers to process the files.  Im sure this can be made even more efficient and much faster since I hacked this together in just a few days in my free time but I think its plenty fast to be useful in its current state.  Please feel free to use this for whatever you want, but if you make improvements I ask that you issue pull requests :)

Usage
------
__Warning:__ This will flaten out directory structures.  It organizes files into extension folders and will not preserve directories, just the files within them.  Also, it moves files not copies them.

```bash
./FileMonster SourceDirectory TargetDirectory
```

*   __SourceDirectory__ - the directory with files and folders you wish to organize. (This will recurse through sub-directories)
*   __TargetDirectory__ - the directory to put the organized files.

The command will open a statistic control pannel that takes over the terminal.  It is finished when all workers say `Done!` in the filename column.  To exit at any time, press `ctrl + c`.

