# Job Manager - First version is up!
![imagem](https://github.com/Xyrsto/job_manager/assets/73367973/98d4c4cb-1430-459b-84ed-1449d882d32a)
## What is this?
So... I need a way to manager my job applications, and I also wanted to learn Go, so I've decided to build this app!
Try it out, its very easy to use, just run it and use the commands down below.
## What commands are available?
Since i haven't implemented any commands to help you use this app (I am planning on doing it in the future), here is the list of commands supported by it:
- `jm -a`: This one is used to create an entry on your job table. 
    - `-cn`: Used to set the company name.
    - `-r`: Used to set the rating of the company
    - `-n`: Used to set a short note about the job offer.
    - `-ha`: Used to set if you have received an answer about the job offer, like a job interview confirmation.
        - This one can either be `true` or `false`. The app will filter your table to show your confirmed job interviews first (the ones set to true)
    - **Example**:
        - `jm -a -cn Company Name -r 5/5 -n Very cool offer! -ha false`
- `jm -ls`: This command will list all the job offers you have registered in the database, in a pretty little table.
- `jm -u`: Using this command, you will need to provide the name of the company you want to update. What it will, is update the "Has Answered" field to `true`
- `clear`: This will clear the terminal.

## What's planned for the next version
- [ ] [Update multiple fields from the job offer](https://github.com/Xyrsto/job_manager/issues/1)
- [x] [Create `jm --help` command](https://github.com/Xyrsto/job_manager/issues/2)
- [ ] [Create `jm -d` to delete an offer](https://github.com/Xyrsto/job_manager/issues/3)
- [ ] [Create interview date field](https://github.com/Xyrsto/job_manager/issues/4)
- [x] [Implement SQLite database](https://github.com/Xyrsto/job_manager/issues/5)



